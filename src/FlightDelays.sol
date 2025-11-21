// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import {VotingPowers} from "./symbiotic/VotingPowers.sol";

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import {
    TimelockControllerUpgradeable
} from "@openzeppelin/contracts-upgradeable/governance/TimelockControllerUpgradeable.sol";

import {IBaseDelegator} from "@symbioticfi/core/src/interfaces/delegator/IBaseDelegator.sol";
import {
    IBaseRewards
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/extensions/IBaseRewards.sol";
import {IBaseSlasher} from "@symbioticfi/core/src/interfaces/slasher/IBaseSlasher.sol";
import {
    IBaseSlashing
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/extensions/IBaseSlashing.sol";
import {
    IDefaultStakerRewardsFactory
} from "@symbioticfi/rewards/src/interfaces/defaultStakerRewards/IDefaultStakerRewardsFactory.sol";
import {
    IDefaultStakerRewards
} from "@symbioticfi/rewards/src/interfaces/defaultStakerRewards/IDefaultStakerRewards.sol";
import {INetworkManager} from "@symbioticfi/relay-contracts/src/interfaces/modules/base/INetworkManager.sol";
import {
    IOperatorNetworkSpecificDelegator
} from "@symbioticfi/core/src/interfaces/delegator/IOperatorNetworkSpecificDelegator.sol";
import {IOperatorRegistry} from "@symbioticfi/core/src/interfaces/IOperatorRegistry.sol";
import {IOptInService} from "@symbioticfi/core/src/interfaces/service/IOptInService.sol";
import {ISettlement} from "@symbioticfi/relay-contracts/src/interfaces/modules/settlement/ISettlement.sol";
import {ISlasher} from "@symbioticfi/core/src/interfaces/slasher/ISlasher.sol";
import {IVaultConfigurator} from "@symbioticfi/core/src/interfaces/IVaultConfigurator.sol";
import {IVault} from "@symbioticfi/core/src/interfaces/vault/IVault.sol";
import {
    IVotingPowerProvider
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/IVotingPowerProvider.sol";
import {NetworkManager} from "@symbioticfi/relay-contracts/src/modules/base/NetworkManager.sol";
import {Subnetwork} from "@symbioticfi/core/src/contracts/libraries/Subnetwork.sol";

contract FlightDelays is NetworkManager {
    using SafeERC20 for IERC20;
    using Subnetwork for bytes32;

    error BuyWindowClosed();
    error FlightNotScheduled();
    error FlightNotDelayable();
    error FlightNotDelayed();
    error FlightAlreadyExists();
    error PreviousFlightIncomplete();
    error InvalidFlight();
    error InvalidTimestamp();
    error InvalidPreviousFlight();
    error InvalidMessageSignature();
    error InvalidEpoch();
    error PolicyAlreadyPurchased();
    error PolicyNotFound();
    error SlashFailed();
    error InsufficientCoverage();
    error InvalidPolicy();

    enum FlightStatus {
        NONE,
        SCHEDULED,
        DELAYED,
        DEPARTED
    }

    enum PolicyStatus {
        NONE,
        PURCHASED,
        CLAIMED
    }

    struct Flight {
        uint48 timestamp;
        FlightStatus status;
        uint128 policiesSold;
        bytes32 previousFlightId;
    }

    struct Airline {
        address vault;
        address rewards;
        uint256 covered;
        bytes32 lastFlightId;
    }

    struct InitParams {
        address votingPowers;
        address settlement;
        address collateral;
        uint48 vaultEpochDuration;
        uint32 messageExpiry;
        uint48 policyWindow;
        uint48 delayWindow;
        uint256 policyPremium;
        uint256 policyPayout;
    }

    event AirlineVaultDeployed(bytes32 indexed airlineId, address vault, address rewards);
    event FlightCreated(bytes32 indexed airlineId, bytes32 indexed flightId, uint48 scheduledTimestamp);
    event FlightDelayed(bytes32 indexed airlineId, bytes32 indexed flightId);
    event FlightDeparted(bytes32 indexed airlineId, bytes32 indexed flightId);
    event InsurancePurchased(
        bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 premium
    );
    event InsuranceClaimed(bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 payout);

    bytes32 internal constant CREATE_MESSAGE_TYPEHASH =
        keccak256("Create(bytes32 airlineId,bytes32 flightId,uint48 departure,bytes32 previousFlightId)");
    bytes32 internal constant DELAY_MESSAGE_TYPEHASH = keccak256("Delay(bytes32 airlineId,bytes32 flightId)");
    bytes32 internal constant DEPART_MESSAGE_TYPEHASH = keccak256("Depart(bytes32 airlineId,bytes32 flightId)");

    address public immutable VAULT_CONFIGURATOR;
    address public immutable DEFAULT_STAKER_REWARDS_FACTORY;
    address public immutable OPERATOR_VAULT_OPT_IN_SERVICE;
    address public immutable OPERATOR_NETWORK_OPT_IN_SERVICE;
    address public immutable OPERATOR_REGISTRY;

    address public votingPowers;
    address public settlement;
    address public collateral;
    uint48 public policyWindow;
    uint48 public delayWindow;
    uint48 public vaultEpochDuration;
    uint32 public messageExpiry;
    uint256 public policyPremium;
    uint256 public policyPayout;

    mapping(bytes32 airlineId => Airline airline) public airlines;
    mapping(bytes32 airlineId => mapping(bytes32 flightId => Flight flight)) public flights;
    mapping(bytes32 airlineId => mapping(bytes32 flightId => mapping(address buyer => PolicyStatus policyStatus)))
        public policies;

    constructor(
        address vaultConfigurator,
        address operatorVaultOptInService,
        address operatorNetworkOptInService,
        address defaultStakerRewardsFactory,
        address operatorRegistry
    ) {
        VAULT_CONFIGURATOR = vaultConfigurator;
        OPERATOR_VAULT_OPT_IN_SERVICE = operatorVaultOptInService;
        OPERATOR_NETWORK_OPT_IN_SERVICE = operatorNetworkOptInService;
        DEFAULT_STAKER_REWARDS_FACTORY = defaultStakerRewardsFactory;
        OPERATOR_REGISTRY = operatorRegistry;
    }

    function initialize(InitParams calldata initParams) external initializer {
        if (initParams.policyPremium == 0 || initParams.policyPayout < initParams.policyPremium) {
            revert InvalidPolicy();
        }

        settlement = initParams.settlement;
        collateral = initParams.collateral;
        vaultEpochDuration = initParams.vaultEpochDuration;
        messageExpiry = initParams.messageExpiry;
        policyWindow = initParams.policyWindow;
        delayWindow = initParams.delayWindow;
        policyPremium = initParams.policyPremium;
        policyPayout = initParams.policyPayout;
        votingPowers = initParams.votingPowers;

        __NetworkManager_init(
            INetworkManager.NetworkManagerInitParams({
                network: INetworkManager(votingPowers).NETWORK(),
                subnetworkId: INetworkManager(votingPowers).SUBNETWORK_IDENTIFIER()
            })
        );

        IOperatorRegistry(OPERATOR_REGISTRY).registerOperator();
        IOptInService(OPERATOR_NETWORK_OPT_IN_SERVICE).optIn(NETWORK());
    }

    function createFlight(
        bytes32 airlineId,
        bytes32 flightId,
        uint48 scheduledTimestamp,
        bytes32 previousFlightId,
        uint48 epoch,
        bytes calldata proof
    ) external {
        if (airlineId == bytes32(0) || flightId == bytes32(0)) {
            revert InvalidFlight();
        }
        _verifyFlightMessage(
            abi.encode(
                keccak256(
                    abi.encode(CREATE_MESSAGE_TYPEHASH, airlineId, flightId, scheduledTimestamp, previousFlightId)
                )
            ),
            epoch,
            proof
        );

        Airline storage airline = airlines[airlineId];
        if (airline.vault == address(0)) {
            _deployAirline(airlineId);
        }
        if (previousFlightId != airline.lastFlightId) {
            revert InvalidPreviousFlight();
        }

        Flight storage flight = flights[airlineId][flightId];
        if (flight.status != FlightStatus.NONE) {
            revert FlightAlreadyExists();
        }

        if (previousFlightId != bytes32(0) && scheduledTimestamp < flights[airlineId][previousFlightId].timestamp) {
            revert InvalidTimestamp();
        }

        flight.timestamp = scheduledTimestamp;
        flight.status = FlightStatus.SCHEDULED;
        flight.previousFlightId = previousFlightId;
        airline.lastFlightId = flightId;

        emit FlightCreated(airlineId, flightId, scheduledTimestamp);
    }

    function delayFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes calldata proof) external {
        Flight storage flight = flights[airlineId][flightId];
        if (flight.status != FlightStatus.SCHEDULED) {
            revert FlightNotDelayable();
        }

        _verifyFlightMessage(
            abi.encode(keccak256(abi.encode(DELAY_MESSAGE_TYPEHASH, airlineId, flightId))), epoch, proof
        );

        bytes32 previousFlightId = flight.previousFlightId;
        if (
            previousFlightId != bytes32(0)
                && (flights[airlineId][previousFlightId].status != FlightStatus.DEPARTED
                    && flights[airlineId][previousFlightId].status != FlightStatus.DELAYED)
        ) {
            revert PreviousFlightIncomplete();
        }

        flight.status = FlightStatus.DELAYED;

        Airline storage airline = airlines[airlineId];
        uint256 coverage = flight.policiesSold * policyPayout;

        if (coverage > 0) {
            (bool success,) = IBaseSlashing(votingPowers)
                .slashVault(flight.timestamp - policyWindow, airline.vault, address(this), coverage, new bytes(0));
            if (!success) {
                revert SlashFailed();
            }
        }

        airline.covered -= coverage;

        emit FlightDelayed(airlineId, flightId);
    }

    function departFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes calldata proof) external {
        Flight storage flight = flights[airlineId][flightId];
        if (flight.status != FlightStatus.SCHEDULED) {
            revert FlightNotScheduled();
        }

        _verifyFlightMessage(
            abi.encode(keccak256(abi.encode(DEPART_MESSAGE_TYPEHASH, airlineId, flightId))), epoch, proof
        );

        bytes32 previousFlightId = flight.previousFlightId;
        if (
            previousFlightId != bytes32(0)
                && (flights[airlineId][previousFlightId].status != FlightStatus.DEPARTED
                    && flights[airlineId][previousFlightId].status != FlightStatus.DELAYED)
        ) {
            revert PreviousFlightIncomplete();
        }

        flight.status = FlightStatus.DEPARTED;

        Airline storage airline = airlines[airlineId];
        airline.covered -= flight.policiesSold * policyPayout;

        uint256 rewardsAmount = flight.policiesSold * policyPremium;
        if (rewardsAmount > 0) {
            IERC20(collateral).safeTransfer(votingPowers, rewardsAmount);
            IBaseRewards(votingPowers)
                .distributeStakerRewards(
                    airline.rewards,
                    collateral,
                    rewardsAmount,
                    abi.encode(flight.timestamp - policyWindow, 10_000, new bytes(0), new bytes(0))
                );
        }

        emit FlightDeparted(airlineId, flightId);
    }

    function buyInsurance(bytes32 airlineId, bytes32 flightId) external {
        Flight storage flight = flights[airlineId][flightId];
        if (flight.status != FlightStatus.SCHEDULED) {
            revert FlightNotScheduled();
        }

        uint48 captureTimestamp = flight.timestamp - policyWindow;
        if (block.timestamp <= captureTimestamp || block.timestamp > flight.timestamp - delayWindow) {
            revert BuyWindowClosed();
        }

        PolicyStatus policyStatus = policies[airlineId][flightId][msg.sender];
        if (policyStatus != PolicyStatus.NONE) {
            revert PolicyAlreadyPurchased();
        }

        Airline storage airline = airlines[airlineId];
        uint256 newCovered = airline.covered + policyPayout;
        if (
            newCovered
                > IBaseSlasher(IVault(airline.vault).slasher())
                    .slashableStake(SUBNETWORK(), address(this), captureTimestamp, new bytes(0))
        ) {
            revert InsufficientCoverage();
        }

        policies[airlineId][flightId][msg.sender] = PolicyStatus.PURCHASED;
        ++flight.policiesSold;
        airline.covered = newCovered;

        IERC20(collateral).safeTransferFrom(msg.sender, address(this), policyPremium);

        emit InsurancePurchased(airlineId, flightId, msg.sender, policyPremium);
    }

    function claimInsurance(bytes32 airlineId, bytes32 flightId) external {
        Flight storage flight = flights[airlineId][flightId];
        if (flight.status != FlightStatus.DELAYED) {
            revert FlightNotDelayed();
        }

        PolicyStatus policyStatus = policies[airlineId][flightId][msg.sender];
        if (policyStatus != PolicyStatus.PURCHASED) {
            revert PolicyNotFound();
        }

        policies[airlineId][flightId][msg.sender] = PolicyStatus.CLAIMED;

        IERC20(collateral).safeTransfer(msg.sender, policyPayout);

        emit InsuranceClaimed(airlineId, flightId, msg.sender, policyPayout);
    }

    function _deployAirline(bytes32 airlineId) internal {
        Airline storage airline = airlines[airlineId];
        if (airline.vault != address(0)) {
            return;
        }

        (address vault,,) = IVaultConfigurator(VAULT_CONFIGURATOR)
            .create(
                IVaultConfigurator.InitParams({
                    version: 1,
                    owner: address(this),
                    vaultParams: abi.encode(
                        IVault.InitParams({
                            collateral: collateral,
                            burner: address(this),
                            epochDuration: vaultEpochDuration,
                            depositWhitelist: false,
                            isDepositLimit: false,
                            depositLimit: 0,
                            defaultAdminRoleHolder: address(0),
                            depositWhitelistSetRoleHolder: address(0),
                            depositorWhitelistRoleHolder: address(0),
                            isDepositLimitSetRoleHolder: address(0),
                            depositLimitSetRoleHolder: address(0)
                        })
                    ),
                    delegatorIndex: uint64(IVotingPowerProvider.DelegatorType.OPERATOR_NETWORK_SPECIFIC),
                    delegatorParams: abi.encode(
                        IOperatorNetworkSpecificDelegator.InitParams({
                            baseParams: IBaseDelegator.BaseParams({
                                defaultAdminRoleHolder: address(this),
                                hook: address(0),
                                hookSetRoleHolder: address(this)
                            }),
                            network: NETWORK(),
                            operator: address(this)
                        })
                    ),
                    withSlasher: true,
                    slasherIndex: uint64(IVotingPowerProvider.SlasherType.INSTANT),
                    slasherParams: abi.encode(
                        ISlasher.InitParams({baseParams: IBaseSlasher.BaseParams({isBurnerHook: false})})
                    )
                })
            );

        IOptInService(OPERATOR_VAULT_OPT_IN_SERVICE).optIn(vault);

        VotingPowers(votingPowers).setMaxNetworkLimit(vault);

        address rewards = IDefaultStakerRewardsFactory(DEFAULT_STAKER_REWARDS_FACTORY)
            .create(
                IDefaultStakerRewards.InitParams({
                    vault: vault,
                    adminFee: 0,
                    defaultAdminRoleHolder: address(0),
                    adminFeeClaimRoleHolder: address(0),
                    adminFeeSetRoleHolder: address(0)
                })
            );

        airline.vault = vault;
        airline.rewards = rewards;

        emit AirlineVaultDeployed(airlineId, vault, rewards);
    }

    function _verifyFlightMessage(bytes memory payload, uint48 epoch, bytes calldata proof) internal view {
        uint48 nextCaptureTimestamp = ISettlement(settlement).getCaptureTimestampFromValSetHeaderAt(epoch + 1);
        if (nextCaptureTimestamp > 0 && block.timestamp >= nextCaptureTimestamp + messageExpiry) {
            revert InvalidEpoch();
        }
        if (!ISettlement(settlement)
                .verifyQuorumSigAt(
                    payload,
                    ISettlement(settlement).getRequiredKeyTagFromValSetHeaderAt(epoch),
                    ISettlement(settlement).getQuorumThresholdFromValSetHeaderAt(epoch),
                    proof,
                    epoch,
                    new bytes(0)
                )) {
            revert InvalidMessageSignature();
        }
    }
}

// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import {Test} from "forge-std/Test.sol";

import {FlightDelays} from "../src/FlightDelays.sol";
import {SettlementMock} from "./mock/SettlementMock.sol";
import {
    MockVaultConfigurator,
    MockRewardsFactory,
    MockVault,
    MockRewards,
    OperatorRegistryMock,
    OptInServiceMock,
    VotingPowersMock
} from "./mock/FlightDelaysMocks.sol";
import {MockERC20} from "../script/mocks/MockERC20.sol";

contract FlightDelaysTest is Test {
    FlightDelays public flightDelays;
    SettlementMock public settlement;
    MockVaultConfigurator public vaultConfigurator;
    MockRewardsFactory public rewardsFactory;
    MockERC20 public stablecoin;
    OperatorRegistryMock public operatorRegistry;
    OptInServiceMock public operatorVaultOptInService;
    OptInServiceMock public operatorNetworkOptInService;
    VotingPowersMock public votingPowers;

    bytes32 internal constant AIRLINE_ID = keccak256("ALPHA");
    bytes32 internal constant FLIGHT_ID = keccak256("ALPHA-001");
    uint48 internal departure;

    function setUp() public {
        settlement = new SettlementMock();
        stablecoin = new MockERC20("MockUSD Coin", "mUSDC");
        vaultConfigurator = new MockVaultConfigurator(stablecoin);
        rewardsFactory = new MockRewardsFactory();
        operatorRegistry = new OperatorRegistryMock();
        operatorVaultOptInService = new OptInServiceMock();
        operatorNetworkOptInService = new OptInServiceMock();
        votingPowers = new VotingPowersMock(address(0xAA), 0);

        flightDelays = new FlightDelays(
            address(vaultConfigurator),
            address(operatorVaultOptInService),
            address(operatorNetworkOptInService),
            address(rewardsFactory),
            address(operatorRegistry)
        );
        votingPowers.setFlightDelays(address(flightDelays));

        FlightDelays.InitParams memory initParams = FlightDelays.InitParams({
            votingPowers: address(votingPowers),
            settlement: address(settlement),
            collateral: address(stablecoin),
            vaultEpochDuration: uint48(3 days),
            messageExpiry: 12_000,
            policyWindow: uint48(3 days),
            delayWindow: uint48(1 days),
            policyPremium: 5 ether,
            policyPayout: 50 ether
        });

        flightDelays.initialize(initParams);
        stablecoin.approve(address(flightDelays), type(uint256).max);

        departure = uint48(block.timestamp + 10 days);
    }

    function testBuyInsuranceAndRecordPolicy() public {
        _createFlight();
        _fundLatestVault(1000 ether);

        vm.warp(departure - 1 days);

        flightDelays.buyInsurance(AIRLINE_ID, FLIGHT_ID);

        (,, uint128 policiesSold,) = flightDelays.flights(AIRLINE_ID, FLIGHT_ID);
        assertEq(policiesSold, 1);

        (,, uint256 coverageReserved,) = flightDelays.airlines(AIRLINE_ID);
        assertEq(coverageReserved, 50 ether);

        uint8 policyStatus = uint8(flightDelays.policies(AIRLINE_ID, FLIGHT_ID, address(this)));
        assertEq(policyStatus, uint8(FlightDelays.PolicyStatus.PURCHASED));
    }

    function testDelayFlightAndClaim() public {
        _createFlight();
        _fundLatestVault(1000 ether);

        vm.warp(departure - 1 days);

        flightDelays.buyInsurance(AIRLINE_ID, FLIGHT_ID);
        vm.expectRevert(FlightDelays.FlightAlreadyExists.selector);
        flightDelays.createFlight(AIRLINE_ID, FLIGHT_ID, departure, 1, "");

        vm.warp(departure + 1);
        flightDelays.delayFlight(AIRLINE_ID, FLIGHT_ID, 1, "");

        uint256 balanceBefore = stablecoin.balanceOf(address(this));
        flightDelays.claimInsurance(AIRLINE_ID, FLIGHT_ID);
        assertEq(stablecoin.balanceOf(address(this)), balanceBefore + 50 ether);

        uint8 policyStatus = uint8(flightDelays.policies(AIRLINE_ID, FLIGHT_ID, address(this)));
        assertEq(policyStatus, uint8(FlightDelays.PolicyStatus.CLAIMED));
    }

    function testCompleteFlightDistributesRewards() public {
        _createFlight();
        _fundLatestVault(500 ether);

        vm.warp(departure - 1 days);

        flightDelays.buyInsurance(AIRLINE_ID, FLIGHT_ID);

        vm.warp(departure + 1);
        flightDelays.departFlight(AIRLINE_ID, FLIGHT_ID, 1, "");

        (uint48 timestamp, FlightDelays.FlightStatus status,,) = flightDelays.flights(AIRLINE_ID, FLIGHT_ID);
        assertEq(uint8(status), uint8(FlightDelays.FlightStatus.DEPARTED));
        assertEq(timestamp, departure);

        MockRewards rewards = rewardsFactory.lastRewards();
        assertEq(rewards.lastAmount(), 5 ether);
        assertEq(rewards.lastToken(), address(stablecoin));
        assertEq(rewards.lastSnapshot(), departure - uint48(3 days));
    }

    function testDelayRequiresEarlierFlightsProcessed() public {
        _createFlight();
        _fundLatestVault(1000 ether);

        bytes32 laterFlightId = keccak256("ALPHA-200");
        uint48 laterDeparture = departure + 2 days;
        flightDelays.createFlight(AIRLINE_ID, laterFlightId, laterDeparture, 1, "");

        vm.warp(laterDeparture + 1);
        vm.expectRevert(FlightDelays.PreviousFlightIncomplete.selector);
        flightDelays.delayFlight(AIRLINE_ID, laterFlightId, 1, "");

        flightDelays.departFlight(AIRLINE_ID, FLIGHT_ID, 1, "");

        flightDelays.delayFlight(AIRLINE_ID, laterFlightId, 1, "");
    }

    function testFlightsMustBeCreatedWithNonDecreasingTimestamp() public {
        _createFlight();

        bytes32 laterFlightId = keccak256("ALPHA-201");
        vm.expectRevert(FlightDelays.InvalidTimestamp.selector);
        flightDelays.createFlight(AIRLINE_ID, laterFlightId, departure - 1, 1, "");

        flightDelays.createFlight(AIRLINE_ID, laterFlightId, departure + 1, 1, "");
    }

    function testCannotBuyInsuranceWithoutFlight() public {
        vm.expectRevert(FlightDelays.FlightNotScheduled.selector);
        flightDelays.buyInsurance(AIRLINE_ID, FLIGHT_ID);
    }

    function testCannotBuyBeforePolicyWindowOpens() public {
        _createFlight();
        _fundLatestVault(1_000 ether);

        vm.warp(departure - uint256(flightDelays.policyWindow()));
        vm.expectRevert(FlightDelays.BuyWindowClosed.selector);
        flightDelays.buyInsurance(AIRLINE_ID, FLIGHT_ID);
    }

    function testCannotBuyAfterWindowCloses() public {
        _createFlight();
        _fundLatestVault(1_000 ether);

        uint256 cutoff = departure - uint256(flightDelays.delayWindow()) + 1;
        vm.warp(cutoff);
        vm.expectRevert(FlightDelays.BuyWindowClosed.selector);
        flightDelays.buyInsurance(AIRLINE_ID, FLIGHT_ID);
    }

    function testCannotBuyOnceFlightDelayed() public {
        _createFlight();
        _fundLatestVault(1_000 ether);

        vm.warp(departure - 2 days);
        flightDelays.delayFlight(AIRLINE_ID, FLIGHT_ID, 1, "");

        vm.expectRevert(FlightDelays.FlightNotScheduled.selector);
        flightDelays.buyInsurance(AIRLINE_ID, FLIGHT_ID);
    }

    function _createFlight() internal {
        flightDelays.createFlight(AIRLINE_ID, FLIGHT_ID, departure, 1, "");
    }

    function _fundLatestVault(uint256 amount) internal {
        address vault = vaultConfigurator.lastVault();
        stablecoin.mint(address(this), amount);
        stablecoin.transfer(vault, amount);
        MockVault(vault).setActiveStake(amount);
    }
}

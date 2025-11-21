// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";

import {IVaultConfigurator} from "@symbioticfi/core/src/interfaces/IVaultConfigurator.sol";
import {IVault} from "@symbioticfi/core/src/interfaces/vault/IVault.sol";
import {
    IDefaultStakerRewards
} from "@symbioticfi/rewards/src/interfaces/defaultStakerRewards/IDefaultStakerRewards.sol";
import {
    IBaseSlashing
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/extensions/IBaseSlashing.sol";
import {
    IBaseRewards
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/extensions/IBaseRewards.sol";
import {INetworkManager} from "@symbioticfi/relay-contracts/src/interfaces/modules/base/INetworkManager.sol";

contract MockVault {
    IERC20 public immutable token;
    address public immutable burner;
    uint256 public activeStakeAmount;
    address public slasher;

    constructor(IERC20 token_, address burner_) {
        token = token_;
        burner = burner_;
    }

    function setSlasher(address slasher_) external {
        slasher = slasher_;
    }

    function setActiveStake(uint256 amount) external {
        activeStakeAmount = amount;
    }

    function activeStake() external view returns (uint256) {
        return activeStakeAmount;
    }

    function onSlash(uint256 amount) external returns (uint256) {
        if (amount > activeStakeAmount) {
            amount = activeStakeAmount;
        }
        activeStakeAmount -= amount;
        token.transfer(burner, amount);
        return amount;
    }
}

contract MockSlasher {
    MockVault public immutable vault;

    constructor(address vault_) {
        vault = MockVault(vault_);
    }

    function slash(bytes32, address, uint256 amount, uint48, bytes calldata) external returns (uint256) {
        return vault.onSlash(amount);
    }

    function slashableStake(bytes32, address, uint48, bytes calldata) external view returns (uint256) {
        return vault.activeStake();
    }
}

contract MockVaultConfigurator {
    IERC20 public immutable token;
    address public lastVault;

    constructor(IERC20 token_) {
        token = token_;
    }

    function create(IVaultConfigurator.InitParams calldata params)
        external
        returns (address vault, address delegator, address slasher)
    {
        IVault.InitParams memory baseParams = abi.decode(params.vaultParams, (IVault.InitParams));
        MockVault deployedVault = new MockVault(token, baseParams.burner);
        MockSlasher deployedSlasher = new MockSlasher(address(deployedVault));
        deployedVault.setSlasher(address(deployedSlasher));
        deployedVault.setActiveStake(1_000_000 ether);
        lastVault = address(deployedVault);
        return (address(deployedVault), address(0), address(deployedSlasher));
    }
}

contract MockRewardsFactory {
    MockRewards public lastRewards;

    function create(IDefaultStakerRewards.InitParams calldata params) external returns (address) {
        lastRewards = new MockRewards(params.vault);
        return address(lastRewards);
    }
}

contract MockRewards {
    address public immutable vault;
    address public lastNetwork;
    address public lastToken;
    uint48 public lastSnapshot;
    uint256 public lastAmount;

    constructor(address vault_) {
        vault = vault_;
    }

    function distributeRewards(address network, address token, uint256 amount, bytes calldata data) external {
        (uint48 snapshotTimestamp,,,) = abi.decode(data, (uint48, uint256, bytes, bytes));
        lastNetwork = network;
        lastToken = token;
        lastSnapshot = snapshotTimestamp;
        lastAmount = amount;
        IERC20(token).transferFrom(msg.sender, address(this), amount);
    }
}

contract OperatorRegistryMock {
    function registerOperator() external {}
}

contract OptInServiceMock {
    address public lastWhere;

    function optIn(address where) external {
        lastWhere = where;
    }
}

contract VotingPowersMock is IBaseSlashing, IBaseRewards, INetworkManager {
    struct SlashCall {
        uint48 timestamp;
        address vault;
        address operator;
        uint256 amount;
    }

    address private immutable _network;
    uint96 private immutable _subnetwork;
    bool public shouldSucceed = true;
    SlashCall public lastSlash;
    address public flightDelays;
    address public rewarder;
    bytes public lastStakerRewardData;
    address public lastStakerRewards;
    address public lastRewardToken;
    uint256 public lastRewardAmount;

    constructor(address network_, uint96 subnetwork_) {
        _network = network_;
        _subnetwork = subnetwork_;
    }

    function setSuccess(bool success) external {
        shouldSucceed = success;
    }

    function getSlasher() external pure returns (address) {
        return address(0);
    }

    function setSlasher(address) external override {}

    function slashVault(uint48 timestamp, address vault, address operator, uint256 amount, bytes memory)
        external
        override
        returns (bool success, bytes memory response)
    {
        lastSlash = SlashCall({timestamp: timestamp, vault: vault, operator: operator, amount: amount});
        MockVault(vault).onSlash(amount);
        success = shouldSucceed;
        if (success) {
            response = abi.encode(amount);
        }
    }

    function executeSlashVault(address, uint256 slashIndex, bytes memory)
        external
        pure
        override
        returns (bool success, uint256 slashedAmount)
    {
        return (true, slashIndex);
    }

    function getRewarder() external view returns (address) {
        return rewarder;
    }

    function setRewarder(address rewarder_) external override {
        rewarder = rewarder_;
    }

    function distributeStakerRewards(address stakerRewards, address token, uint256 amount, bytes memory data)
        external
        override
    {
        lastStakerRewards = stakerRewards;
        lastRewardToken = token;
        lastRewardAmount = amount;
        lastStakerRewardData = data;
        IERC20(token).approve(stakerRewards, amount);
        MockRewards(stakerRewards).distributeRewards(_network, token, amount, data);
    }

    function distributeOperatorRewards(address, address, uint256, bytes32) external pure override {
        // no-op for tests
    }

    function NETWORK() external view override returns (address) {
        return _network;
    }

    function SUBNETWORK_IDENTIFIER() external view override returns (uint96) {
        return _subnetwork;
    }

    function SUBNETWORK() external view override returns (bytes32) {
        return bytes32(uint256(uint160(_network)) << 96 | _subnetwork);
    }

    function setFlightDelays(address fd) external {
        flightDelays = fd;
    }

    function setMaxNetworkLimit(address) external view {
        require(msg.sender == flightDelays, "not FD");
    }
}

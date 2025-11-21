// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import {BaseRewards} from "@symbioticfi/relay-contracts/src/modules/voting-power/extensions/BaseRewards.sol";
import {BaseSlashing} from "@symbioticfi/relay-contracts/src/modules/voting-power/extensions/BaseSlashing.sol";
import {
    EqualStakeVPCalc
} from "@symbioticfi/relay-contracts/src/modules/voting-power/common/voting-power-calc/EqualStakeVPCalc.sol";
import {
    IBaseRewards
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/extensions/IBaseRewards.sol";
import {
    IBaseSlashing
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/extensions/IBaseSlashing.sol";
import {
    IOpNetVaultAutoDeploy
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/extensions/IOpNetVaultAutoDeploy.sol";
import {IOzOwnable} from "@symbioticfi/relay-contracts/src/interfaces/modules/common/permissions/IOzOwnable.sol";
import {ISetMaxNetworkLimitHook} from "@symbioticfi/network/src/interfaces/ISetMaxNetworkLimitHook.sol";
import {IVault} from "@symbioticfi/core/src/interfaces/vault/IVault.sol";
import {
    IVotingPowerProvider
} from "@symbioticfi/relay-contracts/src/interfaces/modules/voting-power/IVotingPowerProvider.sol";
import {
    OpNetVaultAutoDeploy
} from "@symbioticfi/relay-contracts/src/modules/voting-power/extensions/OpNetVaultAutoDeploy.sol";
import {OzOwnable} from "@symbioticfi/relay-contracts/src/modules/common/permissions/OzOwnable.sol";
import {VotingPowerProvider} from "@symbioticfi/relay-contracts/src/modules/voting-power/VotingPowerProvider.sol";

contract VotingPowers is
    VotingPowerProvider,
    OzOwnable,
    EqualStakeVPCalc,
    OpNetVaultAutoDeploy,
    BaseSlashing,
    BaseRewards
{
    error NotFlightDelays();

    address public flightDelays;

    constructor(address operatorRegistry, address vaultFactory, address vaultConfigurator)
        VotingPowerProvider(operatorRegistry, vaultFactory)
        OpNetVaultAutoDeploy(vaultConfigurator)
    {}

    function initialize(
        IVotingPowerProvider.VotingPowerProviderInitParams memory votingPowerProviderInitParams,
        IOpNetVaultAutoDeploy.OpNetVaultAutoDeployInitParams memory opNetVaultAutoDeployInitParams,
        IOzOwnable.OzOwnableInitParams memory ozOwnableInitParams,
        IBaseRewards.BaseRewardsInitParams memory baseRewardsInitParams,
        IBaseSlashing.BaseSlashingInitParams memory baseSlashingInitParams
    ) public virtual initializer {
        __VotingPowerProvider_init(votingPowerProviderInitParams);
        __OpNetVaultAutoDeploy_init(opNetVaultAutoDeployInitParams);
        __OzOwnable_init(ozOwnableInitParams);
        __EqualStakeVPCalc_init();
        __BaseRewards_init(baseRewardsInitParams);
        __BaseSlashing_init(baseSlashingInitParams);
    }

    function _registerOperatorImpl(address operator) internal override(OpNetVaultAutoDeploy, VotingPowerProvider) {
        super._registerOperatorImpl(operator);
    }

    function _unregisterOperatorVaultImpl(address operator, address vault)
        internal
        override(OpNetVaultAutoDeploy, VotingPowerProvider)
    {
        super._unregisterOperatorVaultImpl(operator, vault);
    }

    /// @dev To simplify protocol management, FlightDelays contract need to deploy new Vaults for airlines automatically.
    /// Due to Network's contract implementation allowing to set max network limit only to its middleware, it is done via VotingPowers contract,
    /// because it already contains extensions for rewards and slashing so a good choice to be a "middleware".
    function setMaxNetworkLimit(address vault) public {
        if (flightDelays != msg.sender) {
            revert NotFlightDelays();
        }
        ISetMaxNetworkLimitHook(NETWORK())
            .setMaxNetworkLimit(IVault(vault).delegator(), SUBNETWORK_IDENTIFIER(), type(uint256).max);
    }

    function setFlightDelays(address flightDelays_) public checkPermission {
        flightDelays = flightDelays_;
    }
}

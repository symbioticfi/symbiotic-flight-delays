// SPDX-License-Identifier: MIT
pragma solidity ^0.8.25;

import {SymbioticRewardsConstants} from "@symbioticfi/rewards/test/integration/SymbioticRewardsConstants.sol";

contract SymbioticRewardsConstantsHelper {
    function defaultStakerRewardsFactory() public returns (address) {
        return address(SymbioticRewardsConstants.defaultStakerRewardsFactory());
    }
}

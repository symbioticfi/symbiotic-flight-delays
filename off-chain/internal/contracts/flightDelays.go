// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// FlightDelaysInitParams is an auto generated low-level Go binding around an user-defined struct.
type FlightDelaysInitParams struct {
	VotingPowers       common.Address
	Settlement         common.Address
	Collateral         common.Address
	VaultEpochDuration *big.Int
	MessageExpiry      uint32
	PolicyWindow       *big.Int
	DelayWindow        *big.Int
	PolicyPremium      *big.Int
	PolicyPayout       *big.Int
}

// FlightDelaysMetaData contains all meta data concerning the FlightDelays contract.
var FlightDelaysMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"vaultConfigurator\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorVaultOptInService\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorNetworkOptInService\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"defaultStakerRewardsFactory\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"operatorRegistry\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_STAKER_REWARDS_FACTORY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"NETWORK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"OPERATOR_NETWORK_OPT_IN_SERVICE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"OPERATOR_VAULT_OPT_IN_SERVICE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SUBNETWORK\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"SUBNETWORK_IDENTIFIER\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint96\",\"internalType\":\"uint96\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"VAULT_CONFIGURATOR\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"airlines\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"vault\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"rewards\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"covered\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"lastFlightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"buyInsurance\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claimInsurance\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"collateral\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createFlight\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"scheduledTimestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"previousFlightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delayFlight\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"delayWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"departFlight\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"epoch\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"proof\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"flights\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"timestamp\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"status\",\"type\":\"uint8\",\"internalType\":\"enumFlightDelays.FlightStatus\"},{\"name\":\"policiesSold\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"previousFlightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"initParams\",\"type\":\"tuple\",\"internalType\":\"structFlightDelays.InitParams\",\"components\":[{\"name\":\"votingPowers\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"settlement\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"collateral\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"vaultEpochDuration\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"messageExpiry\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"policyWindow\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"delayWindow\",\"type\":\"uint48\",\"internalType\":\"uint48\"},{\"name\":\"policyPremium\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"policyPayout\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"messageExpiry\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"policies\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"buyer\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"policyStatus\",\"type\":\"uint8\",\"internalType\":\"enumFlightDelays.PolicyStatus\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"policyPayout\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"policyPremium\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"policyWindow\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"settlement\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"staticDelegateCall\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"vaultEpochDuration\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint48\",\"internalType\":\"uint48\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"votingPowers\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AirlineVaultDeployed\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"vault\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"rewards\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FlightCreated\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"scheduledTimestamp\",\"type\":\"uint48\",\"indexed\":false,\"internalType\":\"uint48\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FlightDelayed\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"FlightDeparted\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InitSubnetwork\",\"inputs\":[{\"name\":\"network\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"},{\"name\":\"subnetworkId\",\"type\":\"uint96\",\"indexed\":false,\"internalType\":\"uint96\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsuranceClaimed\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"payout\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"InsurancePurchased\",\"inputs\":[{\"name\":\"airlineId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"flightId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"premium\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BuyWindowClosed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FlightAlreadyExists\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FlightNotDelayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FlightNotDelayed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FlightNotScheduled\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InsufficientCoverage\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidEpoch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidFlight\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidMessageSignature\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidPolicy\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidTimestamp\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NetworkManager_InvalidNetwork\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PolicyAlreadyPurchased\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PolicyNotFound\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PreviousFlightIncomplete\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
}

// FlightDelaysABI is the input ABI used to generate the binding from.
// Deprecated: Use FlightDelaysMetaData.ABI instead.
var FlightDelaysABI = FlightDelaysMetaData.ABI

// FlightDelays is an auto generated Go binding around an Ethereum contract.
type FlightDelays struct {
	FlightDelaysCaller     // Read-only binding to the contract
	FlightDelaysTransactor // Write-only binding to the contract
	FlightDelaysFilterer   // Log filterer for contract events
}

// FlightDelaysCaller is an auto generated read-only Go binding around an Ethereum contract.
type FlightDelaysCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlightDelaysTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FlightDelaysTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlightDelaysFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FlightDelaysFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FlightDelaysSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FlightDelaysSession struct {
	Contract     *FlightDelays     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FlightDelaysCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FlightDelaysCallerSession struct {
	Contract *FlightDelaysCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// FlightDelaysTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FlightDelaysTransactorSession struct {
	Contract     *FlightDelaysTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// FlightDelaysRaw is an auto generated low-level Go binding around an Ethereum contract.
type FlightDelaysRaw struct {
	Contract *FlightDelays // Generic contract binding to access the raw methods on
}

// FlightDelaysCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FlightDelaysCallerRaw struct {
	Contract *FlightDelaysCaller // Generic read-only contract binding to access the raw methods on
}

// FlightDelaysTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FlightDelaysTransactorRaw struct {
	Contract *FlightDelaysTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFlightDelays creates a new instance of FlightDelays, bound to a specific deployed contract.
func NewFlightDelays(address common.Address, backend bind.ContractBackend) (*FlightDelays, error) {
	contract, err := bindFlightDelays(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FlightDelays{FlightDelaysCaller: FlightDelaysCaller{contract: contract}, FlightDelaysTransactor: FlightDelaysTransactor{contract: contract}, FlightDelaysFilterer: FlightDelaysFilterer{contract: contract}}, nil
}

// NewFlightDelaysCaller creates a new read-only instance of FlightDelays, bound to a specific deployed contract.
func NewFlightDelaysCaller(address common.Address, caller bind.ContractCaller) (*FlightDelaysCaller, error) {
	contract, err := bindFlightDelays(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysCaller{contract: contract}, nil
}

// NewFlightDelaysTransactor creates a new write-only instance of FlightDelays, bound to a specific deployed contract.
func NewFlightDelaysTransactor(address common.Address, transactor bind.ContractTransactor) (*FlightDelaysTransactor, error) {
	contract, err := bindFlightDelays(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysTransactor{contract: contract}, nil
}

// NewFlightDelaysFilterer creates a new log filterer instance of FlightDelays, bound to a specific deployed contract.
func NewFlightDelaysFilterer(address common.Address, filterer bind.ContractFilterer) (*FlightDelaysFilterer, error) {
	contract, err := bindFlightDelays(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysFilterer{contract: contract}, nil
}

// bindFlightDelays binds a generic wrapper to an already deployed contract.
func bindFlightDelays(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FlightDelaysMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlightDelays *FlightDelaysRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlightDelays.Contract.FlightDelaysCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlightDelays *FlightDelaysRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlightDelays.Contract.FlightDelaysTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlightDelays *FlightDelaysRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlightDelays.Contract.FlightDelaysTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FlightDelays *FlightDelaysCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FlightDelays.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FlightDelays *FlightDelaysTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FlightDelays.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FlightDelays *FlightDelaysTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FlightDelays.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTSTAKERREWARDSFACTORY is a free data retrieval call binding the contract method 0x360d7094.
//
// Solidity: function DEFAULT_STAKER_REWARDS_FACTORY() view returns(address)
func (_FlightDelays *FlightDelaysCaller) DEFAULTSTAKERREWARDSFACTORY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "DEFAULT_STAKER_REWARDS_FACTORY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DEFAULTSTAKERREWARDSFACTORY is a free data retrieval call binding the contract method 0x360d7094.
//
// Solidity: function DEFAULT_STAKER_REWARDS_FACTORY() view returns(address)
func (_FlightDelays *FlightDelaysSession) DEFAULTSTAKERREWARDSFACTORY() (common.Address, error) {
	return _FlightDelays.Contract.DEFAULTSTAKERREWARDSFACTORY(&_FlightDelays.CallOpts)
}

// DEFAULTSTAKERREWARDSFACTORY is a free data retrieval call binding the contract method 0x360d7094.
//
// Solidity: function DEFAULT_STAKER_REWARDS_FACTORY() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) DEFAULTSTAKERREWARDSFACTORY() (common.Address, error) {
	return _FlightDelays.Contract.DEFAULTSTAKERREWARDSFACTORY(&_FlightDelays.CallOpts)
}

// NETWORK is a free data retrieval call binding the contract method 0x8759e6d1.
//
// Solidity: function NETWORK() view returns(address)
func (_FlightDelays *FlightDelaysCaller) NETWORK(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "NETWORK")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// NETWORK is a free data retrieval call binding the contract method 0x8759e6d1.
//
// Solidity: function NETWORK() view returns(address)
func (_FlightDelays *FlightDelaysSession) NETWORK() (common.Address, error) {
	return _FlightDelays.Contract.NETWORK(&_FlightDelays.CallOpts)
}

// NETWORK is a free data retrieval call binding the contract method 0x8759e6d1.
//
// Solidity: function NETWORK() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) NETWORK() (common.Address, error) {
	return _FlightDelays.Contract.NETWORK(&_FlightDelays.CallOpts)
}

// OPERATORNETWORKOPTINSERVICE is a free data retrieval call binding the contract method 0x1a80e500.
//
// Solidity: function OPERATOR_NETWORK_OPT_IN_SERVICE() view returns(address)
func (_FlightDelays *FlightDelaysCaller) OPERATORNETWORKOPTINSERVICE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "OPERATOR_NETWORK_OPT_IN_SERVICE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OPERATORNETWORKOPTINSERVICE is a free data retrieval call binding the contract method 0x1a80e500.
//
// Solidity: function OPERATOR_NETWORK_OPT_IN_SERVICE() view returns(address)
func (_FlightDelays *FlightDelaysSession) OPERATORNETWORKOPTINSERVICE() (common.Address, error) {
	return _FlightDelays.Contract.OPERATORNETWORKOPTINSERVICE(&_FlightDelays.CallOpts)
}

// OPERATORNETWORKOPTINSERVICE is a free data retrieval call binding the contract method 0x1a80e500.
//
// Solidity: function OPERATOR_NETWORK_OPT_IN_SERVICE() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) OPERATORNETWORKOPTINSERVICE() (common.Address, error) {
	return _FlightDelays.Contract.OPERATORNETWORKOPTINSERVICE(&_FlightDelays.CallOpts)
}

// OPERATORVAULTOPTINSERVICE is a free data retrieval call binding the contract method 0x128e5d82.
//
// Solidity: function OPERATOR_VAULT_OPT_IN_SERVICE() view returns(address)
func (_FlightDelays *FlightDelaysCaller) OPERATORVAULTOPTINSERVICE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "OPERATOR_VAULT_OPT_IN_SERVICE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OPERATORVAULTOPTINSERVICE is a free data retrieval call binding the contract method 0x128e5d82.
//
// Solidity: function OPERATOR_VAULT_OPT_IN_SERVICE() view returns(address)
func (_FlightDelays *FlightDelaysSession) OPERATORVAULTOPTINSERVICE() (common.Address, error) {
	return _FlightDelays.Contract.OPERATORVAULTOPTINSERVICE(&_FlightDelays.CallOpts)
}

// OPERATORVAULTOPTINSERVICE is a free data retrieval call binding the contract method 0x128e5d82.
//
// Solidity: function OPERATOR_VAULT_OPT_IN_SERVICE() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) OPERATORVAULTOPTINSERVICE() (common.Address, error) {
	return _FlightDelays.Contract.OPERATORVAULTOPTINSERVICE(&_FlightDelays.CallOpts)
}

// SUBNETWORK is a free data retrieval call binding the contract method 0x773e6b54.
//
// Solidity: function SUBNETWORK() view returns(bytes32)
func (_FlightDelays *FlightDelaysCaller) SUBNETWORK(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "SUBNETWORK")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SUBNETWORK is a free data retrieval call binding the contract method 0x773e6b54.
//
// Solidity: function SUBNETWORK() view returns(bytes32)
func (_FlightDelays *FlightDelaysSession) SUBNETWORK() ([32]byte, error) {
	return _FlightDelays.Contract.SUBNETWORK(&_FlightDelays.CallOpts)
}

// SUBNETWORK is a free data retrieval call binding the contract method 0x773e6b54.
//
// Solidity: function SUBNETWORK() view returns(bytes32)
func (_FlightDelays *FlightDelaysCallerSession) SUBNETWORK() ([32]byte, error) {
	return _FlightDelays.Contract.SUBNETWORK(&_FlightDelays.CallOpts)
}

// SUBNETWORKIDENTIFIER is a free data retrieval call binding the contract method 0xabacb807.
//
// Solidity: function SUBNETWORK_IDENTIFIER() view returns(uint96)
func (_FlightDelays *FlightDelaysCaller) SUBNETWORKIDENTIFIER(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "SUBNETWORK_IDENTIFIER")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SUBNETWORKIDENTIFIER is a free data retrieval call binding the contract method 0xabacb807.
//
// Solidity: function SUBNETWORK_IDENTIFIER() view returns(uint96)
func (_FlightDelays *FlightDelaysSession) SUBNETWORKIDENTIFIER() (*big.Int, error) {
	return _FlightDelays.Contract.SUBNETWORKIDENTIFIER(&_FlightDelays.CallOpts)
}

// SUBNETWORKIDENTIFIER is a free data retrieval call binding the contract method 0xabacb807.
//
// Solidity: function SUBNETWORK_IDENTIFIER() view returns(uint96)
func (_FlightDelays *FlightDelaysCallerSession) SUBNETWORKIDENTIFIER() (*big.Int, error) {
	return _FlightDelays.Contract.SUBNETWORKIDENTIFIER(&_FlightDelays.CallOpts)
}

// VAULTCONFIGURATOR is a free data retrieval call binding the contract method 0xb25bc0c0.
//
// Solidity: function VAULT_CONFIGURATOR() view returns(address)
func (_FlightDelays *FlightDelaysCaller) VAULTCONFIGURATOR(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "VAULT_CONFIGURATOR")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VAULTCONFIGURATOR is a free data retrieval call binding the contract method 0xb25bc0c0.
//
// Solidity: function VAULT_CONFIGURATOR() view returns(address)
func (_FlightDelays *FlightDelaysSession) VAULTCONFIGURATOR() (common.Address, error) {
	return _FlightDelays.Contract.VAULTCONFIGURATOR(&_FlightDelays.CallOpts)
}

// VAULTCONFIGURATOR is a free data retrieval call binding the contract method 0xb25bc0c0.
//
// Solidity: function VAULT_CONFIGURATOR() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) VAULTCONFIGURATOR() (common.Address, error) {
	return _FlightDelays.Contract.VAULTCONFIGURATOR(&_FlightDelays.CallOpts)
}

// Airlines is a free data retrieval call binding the contract method 0x116a272b.
//
// Solidity: function airlines(bytes32 airlineId) view returns(address vault, address rewards, uint256 covered, bytes32 lastFlightId)
func (_FlightDelays *FlightDelaysCaller) Airlines(opts *bind.CallOpts, airlineId [32]byte) (struct {
	Vault        common.Address
	Rewards      common.Address
	Covered      *big.Int
	LastFlightId [32]byte
}, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "airlines", airlineId)

	outstruct := new(struct {
		Vault        common.Address
		Rewards      common.Address
		Covered      *big.Int
		LastFlightId [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Vault = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Rewards = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Covered = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LastFlightId = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// Airlines is a free data retrieval call binding the contract method 0x116a272b.
//
// Solidity: function airlines(bytes32 airlineId) view returns(address vault, address rewards, uint256 covered, bytes32 lastFlightId)
func (_FlightDelays *FlightDelaysSession) Airlines(airlineId [32]byte) (struct {
	Vault        common.Address
	Rewards      common.Address
	Covered      *big.Int
	LastFlightId [32]byte
}, error) {
	return _FlightDelays.Contract.Airlines(&_FlightDelays.CallOpts, airlineId)
}

// Airlines is a free data retrieval call binding the contract method 0x116a272b.
//
// Solidity: function airlines(bytes32 airlineId) view returns(address vault, address rewards, uint256 covered, bytes32 lastFlightId)
func (_FlightDelays *FlightDelaysCallerSession) Airlines(airlineId [32]byte) (struct {
	Vault        common.Address
	Rewards      common.Address
	Covered      *big.Int
	LastFlightId [32]byte
}, error) {
	return _FlightDelays.Contract.Airlines(&_FlightDelays.CallOpts, airlineId)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_FlightDelays *FlightDelaysCaller) Collateral(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "collateral")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_FlightDelays *FlightDelaysSession) Collateral() (common.Address, error) {
	return _FlightDelays.Contract.Collateral(&_FlightDelays.CallOpts)
}

// Collateral is a free data retrieval call binding the contract method 0xd8dfeb45.
//
// Solidity: function collateral() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) Collateral() (common.Address, error) {
	return _FlightDelays.Contract.Collateral(&_FlightDelays.CallOpts)
}

// DelayWindow is a free data retrieval call binding the contract method 0x9e4f4b76.
//
// Solidity: function delayWindow() view returns(uint48)
func (_FlightDelays *FlightDelaysCaller) DelayWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "delayWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelayWindow is a free data retrieval call binding the contract method 0x9e4f4b76.
//
// Solidity: function delayWindow() view returns(uint48)
func (_FlightDelays *FlightDelaysSession) DelayWindow() (*big.Int, error) {
	return _FlightDelays.Contract.DelayWindow(&_FlightDelays.CallOpts)
}

// DelayWindow is a free data retrieval call binding the contract method 0x9e4f4b76.
//
// Solidity: function delayWindow() view returns(uint48)
func (_FlightDelays *FlightDelaysCallerSession) DelayWindow() (*big.Int, error) {
	return _FlightDelays.Contract.DelayWindow(&_FlightDelays.CallOpts)
}

// Flights is a free data retrieval call binding the contract method 0x45e8026d.
//
// Solidity: function flights(bytes32 airlineId, bytes32 flightId) view returns(uint48 timestamp, uint8 status, uint128 policiesSold, bytes32 previousFlightId)
func (_FlightDelays *FlightDelaysCaller) Flights(opts *bind.CallOpts, airlineId [32]byte, flightId [32]byte) (struct {
	Timestamp        *big.Int
	Status           uint8
	PoliciesSold     *big.Int
	PreviousFlightId [32]byte
}, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "flights", airlineId, flightId)

	outstruct := new(struct {
		Timestamp        *big.Int
		Status           uint8
		PoliciesSold     *big.Int
		PreviousFlightId [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Timestamp = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Status = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.PoliciesSold = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.PreviousFlightId = *abi.ConvertType(out[3], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// Flights is a free data retrieval call binding the contract method 0x45e8026d.
//
// Solidity: function flights(bytes32 airlineId, bytes32 flightId) view returns(uint48 timestamp, uint8 status, uint128 policiesSold, bytes32 previousFlightId)
func (_FlightDelays *FlightDelaysSession) Flights(airlineId [32]byte, flightId [32]byte) (struct {
	Timestamp        *big.Int
	Status           uint8
	PoliciesSold     *big.Int
	PreviousFlightId [32]byte
}, error) {
	return _FlightDelays.Contract.Flights(&_FlightDelays.CallOpts, airlineId, flightId)
}

// Flights is a free data retrieval call binding the contract method 0x45e8026d.
//
// Solidity: function flights(bytes32 airlineId, bytes32 flightId) view returns(uint48 timestamp, uint8 status, uint128 policiesSold, bytes32 previousFlightId)
func (_FlightDelays *FlightDelaysCallerSession) Flights(airlineId [32]byte, flightId [32]byte) (struct {
	Timestamp        *big.Int
	Status           uint8
	PoliciesSold     *big.Int
	PreviousFlightId [32]byte
}, error) {
	return _FlightDelays.Contract.Flights(&_FlightDelays.CallOpts, airlineId, flightId)
}

// MessageExpiry is a free data retrieval call binding the contract method 0x5c96f9b5.
//
// Solidity: function messageExpiry() view returns(uint32)
func (_FlightDelays *FlightDelaysCaller) MessageExpiry(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "messageExpiry")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// MessageExpiry is a free data retrieval call binding the contract method 0x5c96f9b5.
//
// Solidity: function messageExpiry() view returns(uint32)
func (_FlightDelays *FlightDelaysSession) MessageExpiry() (uint32, error) {
	return _FlightDelays.Contract.MessageExpiry(&_FlightDelays.CallOpts)
}

// MessageExpiry is a free data retrieval call binding the contract method 0x5c96f9b5.
//
// Solidity: function messageExpiry() view returns(uint32)
func (_FlightDelays *FlightDelaysCallerSession) MessageExpiry() (uint32, error) {
	return _FlightDelays.Contract.MessageExpiry(&_FlightDelays.CallOpts)
}

// Policies is a free data retrieval call binding the contract method 0x73714bcd.
//
// Solidity: function policies(bytes32 airlineId, bytes32 flightId, address buyer) view returns(uint8 policyStatus)
func (_FlightDelays *FlightDelaysCaller) Policies(opts *bind.CallOpts, airlineId [32]byte, flightId [32]byte, buyer common.Address) (uint8, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "policies", airlineId, flightId, buyer)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Policies is a free data retrieval call binding the contract method 0x73714bcd.
//
// Solidity: function policies(bytes32 airlineId, bytes32 flightId, address buyer) view returns(uint8 policyStatus)
func (_FlightDelays *FlightDelaysSession) Policies(airlineId [32]byte, flightId [32]byte, buyer common.Address) (uint8, error) {
	return _FlightDelays.Contract.Policies(&_FlightDelays.CallOpts, airlineId, flightId, buyer)
}

// Policies is a free data retrieval call binding the contract method 0x73714bcd.
//
// Solidity: function policies(bytes32 airlineId, bytes32 flightId, address buyer) view returns(uint8 policyStatus)
func (_FlightDelays *FlightDelaysCallerSession) Policies(airlineId [32]byte, flightId [32]byte, buyer common.Address) (uint8, error) {
	return _FlightDelays.Contract.Policies(&_FlightDelays.CallOpts, airlineId, flightId, buyer)
}

// PolicyPayout is a free data retrieval call binding the contract method 0x365ca41e.
//
// Solidity: function policyPayout() view returns(uint256)
func (_FlightDelays *FlightDelaysCaller) PolicyPayout(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "policyPayout")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PolicyPayout is a free data retrieval call binding the contract method 0x365ca41e.
//
// Solidity: function policyPayout() view returns(uint256)
func (_FlightDelays *FlightDelaysSession) PolicyPayout() (*big.Int, error) {
	return _FlightDelays.Contract.PolicyPayout(&_FlightDelays.CallOpts)
}

// PolicyPayout is a free data retrieval call binding the contract method 0x365ca41e.
//
// Solidity: function policyPayout() view returns(uint256)
func (_FlightDelays *FlightDelaysCallerSession) PolicyPayout() (*big.Int, error) {
	return _FlightDelays.Contract.PolicyPayout(&_FlightDelays.CallOpts)
}

// PolicyPremium is a free data retrieval call binding the contract method 0x22b6fefe.
//
// Solidity: function policyPremium() view returns(uint256)
func (_FlightDelays *FlightDelaysCaller) PolicyPremium(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "policyPremium")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PolicyPremium is a free data retrieval call binding the contract method 0x22b6fefe.
//
// Solidity: function policyPremium() view returns(uint256)
func (_FlightDelays *FlightDelaysSession) PolicyPremium() (*big.Int, error) {
	return _FlightDelays.Contract.PolicyPremium(&_FlightDelays.CallOpts)
}

// PolicyPremium is a free data retrieval call binding the contract method 0x22b6fefe.
//
// Solidity: function policyPremium() view returns(uint256)
func (_FlightDelays *FlightDelaysCallerSession) PolicyPremium() (*big.Int, error) {
	return _FlightDelays.Contract.PolicyPremium(&_FlightDelays.CallOpts)
}

// PolicyWindow is a free data retrieval call binding the contract method 0x0da1f7ca.
//
// Solidity: function policyWindow() view returns(uint48)
func (_FlightDelays *FlightDelaysCaller) PolicyWindow(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "policyWindow")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PolicyWindow is a free data retrieval call binding the contract method 0x0da1f7ca.
//
// Solidity: function policyWindow() view returns(uint48)
func (_FlightDelays *FlightDelaysSession) PolicyWindow() (*big.Int, error) {
	return _FlightDelays.Contract.PolicyWindow(&_FlightDelays.CallOpts)
}

// PolicyWindow is a free data retrieval call binding the contract method 0x0da1f7ca.
//
// Solidity: function policyWindow() view returns(uint48)
func (_FlightDelays *FlightDelaysCallerSession) PolicyWindow() (*big.Int, error) {
	return _FlightDelays.Contract.PolicyWindow(&_FlightDelays.CallOpts)
}

// Settlement is a free data retrieval call binding the contract method 0x51160630.
//
// Solidity: function settlement() view returns(address)
func (_FlightDelays *FlightDelaysCaller) Settlement(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "settlement")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Settlement is a free data retrieval call binding the contract method 0x51160630.
//
// Solidity: function settlement() view returns(address)
func (_FlightDelays *FlightDelaysSession) Settlement() (common.Address, error) {
	return _FlightDelays.Contract.Settlement(&_FlightDelays.CallOpts)
}

// Settlement is a free data retrieval call binding the contract method 0x51160630.
//
// Solidity: function settlement() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) Settlement() (common.Address, error) {
	return _FlightDelays.Contract.Settlement(&_FlightDelays.CallOpts)
}

// VaultEpochDuration is a free data retrieval call binding the contract method 0xee1b2207.
//
// Solidity: function vaultEpochDuration() view returns(uint48)
func (_FlightDelays *FlightDelaysCaller) VaultEpochDuration(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "vaultEpochDuration")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VaultEpochDuration is a free data retrieval call binding the contract method 0xee1b2207.
//
// Solidity: function vaultEpochDuration() view returns(uint48)
func (_FlightDelays *FlightDelaysSession) VaultEpochDuration() (*big.Int, error) {
	return _FlightDelays.Contract.VaultEpochDuration(&_FlightDelays.CallOpts)
}

// VaultEpochDuration is a free data retrieval call binding the contract method 0xee1b2207.
//
// Solidity: function vaultEpochDuration() view returns(uint48)
func (_FlightDelays *FlightDelaysCallerSession) VaultEpochDuration() (*big.Int, error) {
	return _FlightDelays.Contract.VaultEpochDuration(&_FlightDelays.CallOpts)
}

// VotingPowers is a free data retrieval call binding the contract method 0xc00f50eb.
//
// Solidity: function votingPowers() view returns(address)
func (_FlightDelays *FlightDelaysCaller) VotingPowers(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FlightDelays.contract.Call(opts, &out, "votingPowers")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// VotingPowers is a free data retrieval call binding the contract method 0xc00f50eb.
//
// Solidity: function votingPowers() view returns(address)
func (_FlightDelays *FlightDelaysSession) VotingPowers() (common.Address, error) {
	return _FlightDelays.Contract.VotingPowers(&_FlightDelays.CallOpts)
}

// VotingPowers is a free data retrieval call binding the contract method 0xc00f50eb.
//
// Solidity: function votingPowers() view returns(address)
func (_FlightDelays *FlightDelaysCallerSession) VotingPowers() (common.Address, error) {
	return _FlightDelays.Contract.VotingPowers(&_FlightDelays.CallOpts)
}

// BuyInsurance is a paid mutator transaction binding the contract method 0xb8f37ab2.
//
// Solidity: function buyInsurance(bytes32 airlineId, bytes32 flightId) returns()
func (_FlightDelays *FlightDelaysTransactor) BuyInsurance(opts *bind.TransactOpts, airlineId [32]byte, flightId [32]byte) (*types.Transaction, error) {
	return _FlightDelays.contract.Transact(opts, "buyInsurance", airlineId, flightId)
}

// BuyInsurance is a paid mutator transaction binding the contract method 0xb8f37ab2.
//
// Solidity: function buyInsurance(bytes32 airlineId, bytes32 flightId) returns()
func (_FlightDelays *FlightDelaysSession) BuyInsurance(airlineId [32]byte, flightId [32]byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.BuyInsurance(&_FlightDelays.TransactOpts, airlineId, flightId)
}

// BuyInsurance is a paid mutator transaction binding the contract method 0xb8f37ab2.
//
// Solidity: function buyInsurance(bytes32 airlineId, bytes32 flightId) returns()
func (_FlightDelays *FlightDelaysTransactorSession) BuyInsurance(airlineId [32]byte, flightId [32]byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.BuyInsurance(&_FlightDelays.TransactOpts, airlineId, flightId)
}

// ClaimInsurance is a paid mutator transaction binding the contract method 0x147dbcf1.
//
// Solidity: function claimInsurance(bytes32 airlineId, bytes32 flightId) returns()
func (_FlightDelays *FlightDelaysTransactor) ClaimInsurance(opts *bind.TransactOpts, airlineId [32]byte, flightId [32]byte) (*types.Transaction, error) {
	return _FlightDelays.contract.Transact(opts, "claimInsurance", airlineId, flightId)
}

// ClaimInsurance is a paid mutator transaction binding the contract method 0x147dbcf1.
//
// Solidity: function claimInsurance(bytes32 airlineId, bytes32 flightId) returns()
func (_FlightDelays *FlightDelaysSession) ClaimInsurance(airlineId [32]byte, flightId [32]byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.ClaimInsurance(&_FlightDelays.TransactOpts, airlineId, flightId)
}

// ClaimInsurance is a paid mutator transaction binding the contract method 0x147dbcf1.
//
// Solidity: function claimInsurance(bytes32 airlineId, bytes32 flightId) returns()
func (_FlightDelays *FlightDelaysTransactorSession) ClaimInsurance(airlineId [32]byte, flightId [32]byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.ClaimInsurance(&_FlightDelays.TransactOpts, airlineId, flightId)
}

// CreateFlight is a paid mutator transaction binding the contract method 0xd97650b0.
//
// Solidity: function createFlight(bytes32 airlineId, bytes32 flightId, uint48 scheduledTimestamp, bytes32 previousFlightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysTransactor) CreateFlight(opts *bind.TransactOpts, airlineId [32]byte, flightId [32]byte, scheduledTimestamp *big.Int, previousFlightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.contract.Transact(opts, "createFlight", airlineId, flightId, scheduledTimestamp, previousFlightId, epoch, proof)
}

// CreateFlight is a paid mutator transaction binding the contract method 0xd97650b0.
//
// Solidity: function createFlight(bytes32 airlineId, bytes32 flightId, uint48 scheduledTimestamp, bytes32 previousFlightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysSession) CreateFlight(airlineId [32]byte, flightId [32]byte, scheduledTimestamp *big.Int, previousFlightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.CreateFlight(&_FlightDelays.TransactOpts, airlineId, flightId, scheduledTimestamp, previousFlightId, epoch, proof)
}

// CreateFlight is a paid mutator transaction binding the contract method 0xd97650b0.
//
// Solidity: function createFlight(bytes32 airlineId, bytes32 flightId, uint48 scheduledTimestamp, bytes32 previousFlightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysTransactorSession) CreateFlight(airlineId [32]byte, flightId [32]byte, scheduledTimestamp *big.Int, previousFlightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.CreateFlight(&_FlightDelays.TransactOpts, airlineId, flightId, scheduledTimestamp, previousFlightId, epoch, proof)
}

// DelayFlight is a paid mutator transaction binding the contract method 0x43bfd533.
//
// Solidity: function delayFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysTransactor) DelayFlight(opts *bind.TransactOpts, airlineId [32]byte, flightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.contract.Transact(opts, "delayFlight", airlineId, flightId, epoch, proof)
}

// DelayFlight is a paid mutator transaction binding the contract method 0x43bfd533.
//
// Solidity: function delayFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysSession) DelayFlight(airlineId [32]byte, flightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.DelayFlight(&_FlightDelays.TransactOpts, airlineId, flightId, epoch, proof)
}

// DelayFlight is a paid mutator transaction binding the contract method 0x43bfd533.
//
// Solidity: function delayFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysTransactorSession) DelayFlight(airlineId [32]byte, flightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.DelayFlight(&_FlightDelays.TransactOpts, airlineId, flightId, epoch, proof)
}

// DepartFlight is a paid mutator transaction binding the contract method 0x13f8a494.
//
// Solidity: function departFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysTransactor) DepartFlight(opts *bind.TransactOpts, airlineId [32]byte, flightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.contract.Transact(opts, "departFlight", airlineId, flightId, epoch, proof)
}

// DepartFlight is a paid mutator transaction binding the contract method 0x13f8a494.
//
// Solidity: function departFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysSession) DepartFlight(airlineId [32]byte, flightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.DepartFlight(&_FlightDelays.TransactOpts, airlineId, flightId, epoch, proof)
}

// DepartFlight is a paid mutator transaction binding the contract method 0x13f8a494.
//
// Solidity: function departFlight(bytes32 airlineId, bytes32 flightId, uint48 epoch, bytes proof) returns()
func (_FlightDelays *FlightDelaysTransactorSession) DepartFlight(airlineId [32]byte, flightId [32]byte, epoch *big.Int, proof []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.DepartFlight(&_FlightDelays.TransactOpts, airlineId, flightId, epoch, proof)
}

// Initialize is a paid mutator transaction binding the contract method 0x08251708.
//
// Solidity: function initialize((address,address,address,uint48,uint32,uint48,uint48,uint256,uint256) initParams) returns()
func (_FlightDelays *FlightDelaysTransactor) Initialize(opts *bind.TransactOpts, initParams FlightDelaysInitParams) (*types.Transaction, error) {
	return _FlightDelays.contract.Transact(opts, "initialize", initParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x08251708.
//
// Solidity: function initialize((address,address,address,uint48,uint32,uint48,uint48,uint256,uint256) initParams) returns()
func (_FlightDelays *FlightDelaysSession) Initialize(initParams FlightDelaysInitParams) (*types.Transaction, error) {
	return _FlightDelays.Contract.Initialize(&_FlightDelays.TransactOpts, initParams)
}

// Initialize is a paid mutator transaction binding the contract method 0x08251708.
//
// Solidity: function initialize((address,address,address,uint48,uint32,uint48,uint48,uint256,uint256) initParams) returns()
func (_FlightDelays *FlightDelaysTransactorSession) Initialize(initParams FlightDelaysInitParams) (*types.Transaction, error) {
	return _FlightDelays.Contract.Initialize(&_FlightDelays.TransactOpts, initParams)
}

// StaticDelegateCall is a paid mutator transaction binding the contract method 0x9f86fd85.
//
// Solidity: function staticDelegateCall(address target, bytes data) returns()
func (_FlightDelays *FlightDelaysTransactor) StaticDelegateCall(opts *bind.TransactOpts, target common.Address, data []byte) (*types.Transaction, error) {
	return _FlightDelays.contract.Transact(opts, "staticDelegateCall", target, data)
}

// StaticDelegateCall is a paid mutator transaction binding the contract method 0x9f86fd85.
//
// Solidity: function staticDelegateCall(address target, bytes data) returns()
func (_FlightDelays *FlightDelaysSession) StaticDelegateCall(target common.Address, data []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.StaticDelegateCall(&_FlightDelays.TransactOpts, target, data)
}

// StaticDelegateCall is a paid mutator transaction binding the contract method 0x9f86fd85.
//
// Solidity: function staticDelegateCall(address target, bytes data) returns()
func (_FlightDelays *FlightDelaysTransactorSession) StaticDelegateCall(target common.Address, data []byte) (*types.Transaction, error) {
	return _FlightDelays.Contract.StaticDelegateCall(&_FlightDelays.TransactOpts, target, data)
}

// FlightDelaysAirlineVaultDeployedIterator is returned from FilterAirlineVaultDeployed and is used to iterate over the raw logs and unpacked data for AirlineVaultDeployed events raised by the FlightDelays contract.
type FlightDelaysAirlineVaultDeployedIterator struct {
	Event *FlightDelaysAirlineVaultDeployed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysAirlineVaultDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysAirlineVaultDeployed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysAirlineVaultDeployed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysAirlineVaultDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysAirlineVaultDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysAirlineVaultDeployed represents a AirlineVaultDeployed event raised by the FlightDelays contract.
type FlightDelaysAirlineVaultDeployed struct {
	AirlineId [32]byte
	Vault     common.Address
	Rewards   common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterAirlineVaultDeployed is a free log retrieval operation binding the contract event 0x6cb291016784959f890661d07225a645792785a83baeb2592251dcc9e4e7bf34.
//
// Solidity: event AirlineVaultDeployed(bytes32 indexed airlineId, address vault, address rewards)
func (_FlightDelays *FlightDelaysFilterer) FilterAirlineVaultDeployed(opts *bind.FilterOpts, airlineId [][32]byte) (*FlightDelaysAirlineVaultDeployedIterator, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "AirlineVaultDeployed", airlineIdRule)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysAirlineVaultDeployedIterator{contract: _FlightDelays.contract, event: "AirlineVaultDeployed", logs: logs, sub: sub}, nil
}

// WatchAirlineVaultDeployed is a free log subscription operation binding the contract event 0x6cb291016784959f890661d07225a645792785a83baeb2592251dcc9e4e7bf34.
//
// Solidity: event AirlineVaultDeployed(bytes32 indexed airlineId, address vault, address rewards)
func (_FlightDelays *FlightDelaysFilterer) WatchAirlineVaultDeployed(opts *bind.WatchOpts, sink chan<- *FlightDelaysAirlineVaultDeployed, airlineId [][32]byte) (event.Subscription, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "AirlineVaultDeployed", airlineIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysAirlineVaultDeployed)
				if err := _FlightDelays.contract.UnpackLog(event, "AirlineVaultDeployed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAirlineVaultDeployed is a log parse operation binding the contract event 0x6cb291016784959f890661d07225a645792785a83baeb2592251dcc9e4e7bf34.
//
// Solidity: event AirlineVaultDeployed(bytes32 indexed airlineId, address vault, address rewards)
func (_FlightDelays *FlightDelaysFilterer) ParseAirlineVaultDeployed(log types.Log) (*FlightDelaysAirlineVaultDeployed, error) {
	event := new(FlightDelaysAirlineVaultDeployed)
	if err := _FlightDelays.contract.UnpackLog(event, "AirlineVaultDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlightDelaysFlightCreatedIterator is returned from FilterFlightCreated and is used to iterate over the raw logs and unpacked data for FlightCreated events raised by the FlightDelays contract.
type FlightDelaysFlightCreatedIterator struct {
	Event *FlightDelaysFlightCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysFlightCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysFlightCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysFlightCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysFlightCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysFlightCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysFlightCreated represents a FlightCreated event raised by the FlightDelays contract.
type FlightDelaysFlightCreated struct {
	AirlineId          [32]byte
	FlightId           [32]byte
	ScheduledTimestamp *big.Int
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterFlightCreated is a free log retrieval operation binding the contract event 0x1e735cf33b66ed3dc00862d938ae4f8f7dbbda9493652fc483c3607ce23c0070.
//
// Solidity: event FlightCreated(bytes32 indexed airlineId, bytes32 indexed flightId, uint48 scheduledTimestamp)
func (_FlightDelays *FlightDelaysFilterer) FilterFlightCreated(opts *bind.FilterOpts, airlineId [][32]byte, flightId [][32]byte) (*FlightDelaysFlightCreatedIterator, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "FlightCreated", airlineIdRule, flightIdRule)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysFlightCreatedIterator{contract: _FlightDelays.contract, event: "FlightCreated", logs: logs, sub: sub}, nil
}

// WatchFlightCreated is a free log subscription operation binding the contract event 0x1e735cf33b66ed3dc00862d938ae4f8f7dbbda9493652fc483c3607ce23c0070.
//
// Solidity: event FlightCreated(bytes32 indexed airlineId, bytes32 indexed flightId, uint48 scheduledTimestamp)
func (_FlightDelays *FlightDelaysFilterer) WatchFlightCreated(opts *bind.WatchOpts, sink chan<- *FlightDelaysFlightCreated, airlineId [][32]byte, flightId [][32]byte) (event.Subscription, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "FlightCreated", airlineIdRule, flightIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysFlightCreated)
				if err := _FlightDelays.contract.UnpackLog(event, "FlightCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFlightCreated is a log parse operation binding the contract event 0x1e735cf33b66ed3dc00862d938ae4f8f7dbbda9493652fc483c3607ce23c0070.
//
// Solidity: event FlightCreated(bytes32 indexed airlineId, bytes32 indexed flightId, uint48 scheduledTimestamp)
func (_FlightDelays *FlightDelaysFilterer) ParseFlightCreated(log types.Log) (*FlightDelaysFlightCreated, error) {
	event := new(FlightDelaysFlightCreated)
	if err := _FlightDelays.contract.UnpackLog(event, "FlightCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlightDelaysFlightDelayedIterator is returned from FilterFlightDelayed and is used to iterate over the raw logs and unpacked data for FlightDelayed events raised by the FlightDelays contract.
type FlightDelaysFlightDelayedIterator struct {
	Event *FlightDelaysFlightDelayed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysFlightDelayedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysFlightDelayed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysFlightDelayed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysFlightDelayedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysFlightDelayedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysFlightDelayed represents a FlightDelayed event raised by the FlightDelays contract.
type FlightDelaysFlightDelayed struct {
	AirlineId [32]byte
	FlightId  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFlightDelayed is a free log retrieval operation binding the contract event 0x8e8076a6235494b484b7cd7000b0bb95229118813ebf3437799a65912956d478.
//
// Solidity: event FlightDelayed(bytes32 indexed airlineId, bytes32 indexed flightId)
func (_FlightDelays *FlightDelaysFilterer) FilterFlightDelayed(opts *bind.FilterOpts, airlineId [][32]byte, flightId [][32]byte) (*FlightDelaysFlightDelayedIterator, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "FlightDelayed", airlineIdRule, flightIdRule)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysFlightDelayedIterator{contract: _FlightDelays.contract, event: "FlightDelayed", logs: logs, sub: sub}, nil
}

// WatchFlightDelayed is a free log subscription operation binding the contract event 0x8e8076a6235494b484b7cd7000b0bb95229118813ebf3437799a65912956d478.
//
// Solidity: event FlightDelayed(bytes32 indexed airlineId, bytes32 indexed flightId)
func (_FlightDelays *FlightDelaysFilterer) WatchFlightDelayed(opts *bind.WatchOpts, sink chan<- *FlightDelaysFlightDelayed, airlineId [][32]byte, flightId [][32]byte) (event.Subscription, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "FlightDelayed", airlineIdRule, flightIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysFlightDelayed)
				if err := _FlightDelays.contract.UnpackLog(event, "FlightDelayed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFlightDelayed is a log parse operation binding the contract event 0x8e8076a6235494b484b7cd7000b0bb95229118813ebf3437799a65912956d478.
//
// Solidity: event FlightDelayed(bytes32 indexed airlineId, bytes32 indexed flightId)
func (_FlightDelays *FlightDelaysFilterer) ParseFlightDelayed(log types.Log) (*FlightDelaysFlightDelayed, error) {
	event := new(FlightDelaysFlightDelayed)
	if err := _FlightDelays.contract.UnpackLog(event, "FlightDelayed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlightDelaysFlightDepartedIterator is returned from FilterFlightDeparted and is used to iterate over the raw logs and unpacked data for FlightDeparted events raised by the FlightDelays contract.
type FlightDelaysFlightDepartedIterator struct {
	Event *FlightDelaysFlightDeparted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysFlightDepartedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysFlightDeparted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysFlightDeparted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysFlightDepartedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysFlightDepartedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysFlightDeparted represents a FlightDeparted event raised by the FlightDelays contract.
type FlightDelaysFlightDeparted struct {
	AirlineId [32]byte
	FlightId  [32]byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFlightDeparted is a free log retrieval operation binding the contract event 0x2894e01f1028238685425dd73fa9ec9fea4ea16dfd68e0fc5ff8985cd6e1a183.
//
// Solidity: event FlightDeparted(bytes32 indexed airlineId, bytes32 indexed flightId)
func (_FlightDelays *FlightDelaysFilterer) FilterFlightDeparted(opts *bind.FilterOpts, airlineId [][32]byte, flightId [][32]byte) (*FlightDelaysFlightDepartedIterator, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "FlightDeparted", airlineIdRule, flightIdRule)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysFlightDepartedIterator{contract: _FlightDelays.contract, event: "FlightDeparted", logs: logs, sub: sub}, nil
}

// WatchFlightDeparted is a free log subscription operation binding the contract event 0x2894e01f1028238685425dd73fa9ec9fea4ea16dfd68e0fc5ff8985cd6e1a183.
//
// Solidity: event FlightDeparted(bytes32 indexed airlineId, bytes32 indexed flightId)
func (_FlightDelays *FlightDelaysFilterer) WatchFlightDeparted(opts *bind.WatchOpts, sink chan<- *FlightDelaysFlightDeparted, airlineId [][32]byte, flightId [][32]byte) (event.Subscription, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "FlightDeparted", airlineIdRule, flightIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysFlightDeparted)
				if err := _FlightDelays.contract.UnpackLog(event, "FlightDeparted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFlightDeparted is a log parse operation binding the contract event 0x2894e01f1028238685425dd73fa9ec9fea4ea16dfd68e0fc5ff8985cd6e1a183.
//
// Solidity: event FlightDeparted(bytes32 indexed airlineId, bytes32 indexed flightId)
func (_FlightDelays *FlightDelaysFilterer) ParseFlightDeparted(log types.Log) (*FlightDelaysFlightDeparted, error) {
	event := new(FlightDelaysFlightDeparted)
	if err := _FlightDelays.contract.UnpackLog(event, "FlightDeparted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlightDelaysInitSubnetworkIterator is returned from FilterInitSubnetwork and is used to iterate over the raw logs and unpacked data for InitSubnetwork events raised by the FlightDelays contract.
type FlightDelaysInitSubnetworkIterator struct {
	Event *FlightDelaysInitSubnetwork // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysInitSubnetworkIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysInitSubnetwork)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysInitSubnetwork)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysInitSubnetworkIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysInitSubnetworkIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysInitSubnetwork represents a InitSubnetwork event raised by the FlightDelays contract.
type FlightDelaysInitSubnetwork struct {
	Network      common.Address
	SubnetworkId *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterInitSubnetwork is a free log retrieval operation binding the contract event 0x469c2e982e7d76d34cf5d1e72abee29749bb9971942c180e9023cea09f5f8e83.
//
// Solidity: event InitSubnetwork(address network, uint96 subnetworkId)
func (_FlightDelays *FlightDelaysFilterer) FilterInitSubnetwork(opts *bind.FilterOpts) (*FlightDelaysInitSubnetworkIterator, error) {

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "InitSubnetwork")
	if err != nil {
		return nil, err
	}
	return &FlightDelaysInitSubnetworkIterator{contract: _FlightDelays.contract, event: "InitSubnetwork", logs: logs, sub: sub}, nil
}

// WatchInitSubnetwork is a free log subscription operation binding the contract event 0x469c2e982e7d76d34cf5d1e72abee29749bb9971942c180e9023cea09f5f8e83.
//
// Solidity: event InitSubnetwork(address network, uint96 subnetworkId)
func (_FlightDelays *FlightDelaysFilterer) WatchInitSubnetwork(opts *bind.WatchOpts, sink chan<- *FlightDelaysInitSubnetwork) (event.Subscription, error) {

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "InitSubnetwork")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysInitSubnetwork)
				if err := _FlightDelays.contract.UnpackLog(event, "InitSubnetwork", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitSubnetwork is a log parse operation binding the contract event 0x469c2e982e7d76d34cf5d1e72abee29749bb9971942c180e9023cea09f5f8e83.
//
// Solidity: event InitSubnetwork(address network, uint96 subnetworkId)
func (_FlightDelays *FlightDelaysFilterer) ParseInitSubnetwork(log types.Log) (*FlightDelaysInitSubnetwork, error) {
	event := new(FlightDelaysInitSubnetwork)
	if err := _FlightDelays.contract.UnpackLog(event, "InitSubnetwork", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlightDelaysInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FlightDelays contract.
type FlightDelaysInitializedIterator struct {
	Event *FlightDelaysInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysInitialized represents a Initialized event raised by the FlightDelays contract.
type FlightDelaysInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FlightDelays *FlightDelaysFilterer) FilterInitialized(opts *bind.FilterOpts) (*FlightDelaysInitializedIterator, error) {

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FlightDelaysInitializedIterator{contract: _FlightDelays.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FlightDelays *FlightDelaysFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FlightDelaysInitialized) (event.Subscription, error) {

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysInitialized)
				if err := _FlightDelays.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FlightDelays *FlightDelaysFilterer) ParseInitialized(log types.Log) (*FlightDelaysInitialized, error) {
	event := new(FlightDelaysInitialized)
	if err := _FlightDelays.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlightDelaysInsuranceClaimedIterator is returned from FilterInsuranceClaimed and is used to iterate over the raw logs and unpacked data for InsuranceClaimed events raised by the FlightDelays contract.
type FlightDelaysInsuranceClaimedIterator struct {
	Event *FlightDelaysInsuranceClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysInsuranceClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysInsuranceClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysInsuranceClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysInsuranceClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysInsuranceClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysInsuranceClaimed represents a InsuranceClaimed event raised by the FlightDelays contract.
type FlightDelaysInsuranceClaimed struct {
	AirlineId [32]byte
	FlightId  [32]byte
	Buyer     common.Address
	Payout    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInsuranceClaimed is a free log retrieval operation binding the contract event 0x57e4f7b973a879494a4ada5b47e55405e342f437cf7280d04008aac8dde6a588.
//
// Solidity: event InsuranceClaimed(bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 payout)
func (_FlightDelays *FlightDelaysFilterer) FilterInsuranceClaimed(opts *bind.FilterOpts, airlineId [][32]byte, flightId [][32]byte, buyer []common.Address) (*FlightDelaysInsuranceClaimedIterator, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "InsuranceClaimed", airlineIdRule, flightIdRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysInsuranceClaimedIterator{contract: _FlightDelays.contract, event: "InsuranceClaimed", logs: logs, sub: sub}, nil
}

// WatchInsuranceClaimed is a free log subscription operation binding the contract event 0x57e4f7b973a879494a4ada5b47e55405e342f437cf7280d04008aac8dde6a588.
//
// Solidity: event InsuranceClaimed(bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 payout)
func (_FlightDelays *FlightDelaysFilterer) WatchInsuranceClaimed(opts *bind.WatchOpts, sink chan<- *FlightDelaysInsuranceClaimed, airlineId [][32]byte, flightId [][32]byte, buyer []common.Address) (event.Subscription, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "InsuranceClaimed", airlineIdRule, flightIdRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysInsuranceClaimed)
				if err := _FlightDelays.contract.UnpackLog(event, "InsuranceClaimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInsuranceClaimed is a log parse operation binding the contract event 0x57e4f7b973a879494a4ada5b47e55405e342f437cf7280d04008aac8dde6a588.
//
// Solidity: event InsuranceClaimed(bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 payout)
func (_FlightDelays *FlightDelaysFilterer) ParseInsuranceClaimed(log types.Log) (*FlightDelaysInsuranceClaimed, error) {
	event := new(FlightDelaysInsuranceClaimed)
	if err := _FlightDelays.contract.UnpackLog(event, "InsuranceClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FlightDelaysInsurancePurchasedIterator is returned from FilterInsurancePurchased and is used to iterate over the raw logs and unpacked data for InsurancePurchased events raised by the FlightDelays contract.
type FlightDelaysInsurancePurchasedIterator struct {
	Event *FlightDelaysInsurancePurchased // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FlightDelaysInsurancePurchasedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FlightDelaysInsurancePurchased)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FlightDelaysInsurancePurchased)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FlightDelaysInsurancePurchasedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FlightDelaysInsurancePurchasedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FlightDelaysInsurancePurchased represents a InsurancePurchased event raised by the FlightDelays contract.
type FlightDelaysInsurancePurchased struct {
	AirlineId [32]byte
	FlightId  [32]byte
	Buyer     common.Address
	Premium   *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterInsurancePurchased is a free log retrieval operation binding the contract event 0xb5741046414a47afe0df7618cfbff25d005a117a6e9aa4dde31c6c8ffb6ab4d5.
//
// Solidity: event InsurancePurchased(bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 premium)
func (_FlightDelays *FlightDelaysFilterer) FilterInsurancePurchased(opts *bind.FilterOpts, airlineId [][32]byte, flightId [][32]byte, buyer []common.Address) (*FlightDelaysInsurancePurchasedIterator, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FlightDelays.contract.FilterLogs(opts, "InsurancePurchased", airlineIdRule, flightIdRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return &FlightDelaysInsurancePurchasedIterator{contract: _FlightDelays.contract, event: "InsurancePurchased", logs: logs, sub: sub}, nil
}

// WatchInsurancePurchased is a free log subscription operation binding the contract event 0xb5741046414a47afe0df7618cfbff25d005a117a6e9aa4dde31c6c8ffb6ab4d5.
//
// Solidity: event InsurancePurchased(bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 premium)
func (_FlightDelays *FlightDelaysFilterer) WatchInsurancePurchased(opts *bind.WatchOpts, sink chan<- *FlightDelaysInsurancePurchased, airlineId [][32]byte, flightId [][32]byte, buyer []common.Address) (event.Subscription, error) {

	var airlineIdRule []interface{}
	for _, airlineIdItem := range airlineId {
		airlineIdRule = append(airlineIdRule, airlineIdItem)
	}
	var flightIdRule []interface{}
	for _, flightIdItem := range flightId {
		flightIdRule = append(flightIdRule, flightIdItem)
	}
	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _FlightDelays.contract.WatchLogs(opts, "InsurancePurchased", airlineIdRule, flightIdRule, buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FlightDelaysInsurancePurchased)
				if err := _FlightDelays.contract.UnpackLog(event, "InsurancePurchased", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInsurancePurchased is a log parse operation binding the contract event 0xb5741046414a47afe0df7618cfbff25d005a117a6e9aa4dde31c6c8ffb6ab4d5.
//
// Solidity: event InsurancePurchased(bytes32 indexed airlineId, bytes32 indexed flightId, address indexed buyer, uint256 premium)
func (_FlightDelays *FlightDelaysFilterer) ParseInsurancePurchased(log types.Log) (*FlightDelaysInsurancePurchased, error) {
	event := new(FlightDelaysInsurancePurchased)
	if err := _FlightDelays.contract.UnpackLog(event, "InsurancePurchased", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package build

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

	c4tBind "github.com/chain4travel/caminoethvm/accounts/abi/bind"
	c4tAbi "github.com/chain4travel/caminoethvm/accounts/abi"
	// c4tCommon "github.com/chain4travel/caminoethvm/common"
	c4tTypes "github.com/chain4travel/caminoethvm/core/types"
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
)

// BuildMetaData contains all meta data concerning the Build contract.
var BuildMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"role\",\"type\":\"uint256\"}],\"name\":\"DropRole\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"role\",\"type\":\"uint256\"}],\"name\":\"SetRole\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getRoles\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"role\",\"type\":\"uint256\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"role\",\"type\":\"uint256\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"role\",\"type\":\"uint256\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// BuildABI is the input ABI used to generate the binding from.
// Deprecated: Use BuildMetaData.ABI instead.
var BuildABI = BuildMetaData.ABI

// Build is an auto generated Go binding around an Ethereum contract.
type Build struct {
	BuildCaller     // Read-only binding to the contract
	BuildTransactor // Write-only binding to the contract
	BuildFilterer   // Log filterer for contract events
}

type C4TBuild struct {
	C4TBuildCaller     // Read-only binding to the contract
	C4TBuildTransactor // Write-only binding to the contract
	C4TBuildFilterer   // Log filterer for contract events
}

// BuildCaller is an auto generated read-only Go binding around an Ethereum contract.
type BuildCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BuildCaller is an auto generated read-only Go binding around an Ethereum contract.
type C4TBuildCaller struct {
	contract *c4tBind.BoundContract // Generic contract wrapper for the low level calls
}

// BuildTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BuildTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type C4TBuildTransactor struct {
	contract *c4tBind.BoundContract // Generic contract wrapper for the low level calls
}

// BuildFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BuildFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

type C4TBuildFilterer struct {
	contract *c4tBind.BoundContract // Generic contract wrapper for the low level calls
}

// BuildSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BuildSession struct {
	Contract     *Build            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

type C4TBuildSession struct {
	Contract     *C4TBuild            // Generic contract binding to set the session for
	CallOpts     c4tBind.CallOpts     // Call options to use throughout this session
	TransactOpts c4tBind.TransactOpts // Transaction auth options to use throughout this session
}

// BuildCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BuildCallerSession struct {
	Contract *BuildCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// BuildTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BuildTransactorSession struct {
	Contract     *BuildTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BuildRaw is an auto generated low-level Go binding around an Ethereum contract.
type BuildRaw struct {
	Contract *Build // Generic contract binding to access the raw methods on
}

// BuildCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BuildCallerRaw struct {
	Contract *BuildCaller // Generic read-only contract binding to access the raw methods on
}

// BuildTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BuildTransactorRaw struct {
	Contract *BuildTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBuild creates a new instance of Build, bound to a specific deployed contract.
func NewBuild(address common.Address, backend bind.ContractBackend) (*Build, error) {
	contract, err := bindBuild(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Build{BuildCaller: BuildCaller{contract: contract}, BuildTransactor: BuildTransactor{contract: contract}, BuildFilterer: BuildFilterer{contract: contract}}, nil
}

// NewBuild creates a new instance of Build, bound to a specific deployed contract.
func C4TNewBuild(address common.Address, backend c4tBind.ContractBackend) (*C4TBuild, error) {
	contract, err := c4tBindBuild(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &C4TBuild{C4TBuildCaller: C4TBuildCaller{contract: contract}, C4TBuildTransactor: C4TBuildTransactor{contract: contract}, C4TBuildFilterer: C4TBuildFilterer{contract: contract}}, nil
}

// NewBuildCaller creates a new read-only instance of Build, bound to a specific deployed contract.
func NewBuildCaller(address common.Address, caller bind.ContractCaller) (*BuildCaller, error) {
	contract, err := bindBuild(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BuildCaller{contract: contract}, nil
}

// NewBuildTransactor creates a new write-only instance of Build, bound to a specific deployed contract.
func NewBuildTransactor(address common.Address, transactor bind.ContractTransactor) (*BuildTransactor, error) {
	contract, err := bindBuild(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BuildTransactor{contract: contract}, nil
}

// NewBuildFilterer creates a new log filterer instance of Build, bound to a specific deployed contract.
func NewBuildFilterer(address common.Address, filterer bind.ContractFilterer) (*BuildFilterer, error) {
	contract, err := bindBuild(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BuildFilterer{contract: contract}, nil
}

// bindBuild binds a generic wrapper to an already deployed contract.
func bindBuild(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BuildABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// bindBuild binds a generic wrapper to an already deployed contract.
func c4tBindBuild(address common.Address, caller c4tBind.ContractCaller, transactor c4tBind.ContractTransactor, filterer c4tBind.ContractFilterer) (*c4tBind.BoundContract, error) {
	parsed, err := c4tAbi.JSON(strings.NewReader(BuildABI))
	if err != nil {
		return nil, err
	}
	return c4tBind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Build *BuildRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Build.Contract.BuildCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Build *BuildRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Build.Contract.BuildTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Build *BuildRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Build.Contract.BuildTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Build *BuildCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Build.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Build *BuildTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Build.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Build *BuildTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Build.Contract.contract.Transact(opts, method, params...)
}

// GetRoles is a free data retrieval call binding the contract method 0xce6ccfaf.
//
// Solidity: function getRoles(address addr) view returns(uint256)
func (_Build *BuildCaller) GetRoles(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Build.contract.Call(opts, &out, "getRoles", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

func (_C4TBuild *C4TBuildCaller) C4TGetRoles(opts *c4tBind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _C4TBuild.contract.Call(opts, &out, "getRoles", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetRoles is a free data retrieval call binding the contract method 0xce6ccfaf.
//
// Solidity: function getRoles(address addr) view returns(uint256)
func (_Build *BuildSession) GetRoles(addr common.Address) (*big.Int, error) {
	return _Build.Contract.GetRoles(&_Build.CallOpts, addr)
}

func (_C4TBuild *C4TBuildSession) C4TGetRoles(addr common.Address) (*big.Int, error) {
	return _C4TBuild.Contract.C4TGetRoles(&_C4TBuild.CallOpts, addr)
}

// GetRoles is a free data retrieval call binding the contract method 0xce6ccfaf.
//
// Solidity: function getRoles(address addr) view returns(uint256)
func (_Build *BuildCallerSession) GetRoles(addr common.Address) (*big.Int, error) {
	return _Build.Contract.GetRoles(&_Build.CallOpts, addr)
}

// HasRole is a free data retrieval call binding the contract method 0x5c97f4a2.
//
// Solidity: function hasRole(address addr, uint256 role) view returns(bool)
func (_Build *BuildCaller) HasRole(opts *bind.CallOpts, addr common.Address, role *big.Int) (bool, error) {
	var out []interface{}
	err := _Build.contract.Call(opts, &out, "hasRole", addr, role)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x5c97f4a2.
//
// Solidity: function hasRole(address addr, uint256 role) view returns(bool)
func (_Build *BuildSession) HasRole(addr common.Address, role *big.Int) (bool, error) {
	return _Build.Contract.HasRole(&_Build.CallOpts, addr, role)
}

// HasRole is a free data retrieval call binding the contract method 0x5c97f4a2.
//
// Solidity: function hasRole(address addr, uint256 role) view returns(bool)
func (_Build *BuildCallerSession) HasRole(addr common.Address, role *big.Int) (bool, error) {
	return _Build.Contract.HasRole(&_Build.CallOpts, addr, role)
}

// GrantRole is a paid mutator transaction binding the contract method 0x3c09e2fd.
//
// Solidity: function grantRole(address addr, uint256 role) returns()
func (_Build *BuildTransactor) GrantRole(opts *bind.TransactOpts, addr common.Address, role *big.Int) (*types.Transaction, error) {
	return _Build.contract.Transact(opts, "grantRole", addr, role)
}

// GrantRole is a paid mutator transaction binding the contract method 0x3c09e2fd.
//
// Solidity: function grantRole(address addr, uint256 role) returns()
func (_Build *C4TBuildTransactor) C4TGrantRole(opts *c4tBind.TransactOpts, addr common.Address, role *big.Int) (*c4tTypes.Transaction, error) {
	return _Build.contract.Transact(opts, "grantRole", addr, role)
}

// GrantRole is a paid mutator transaction binding the contract method 0x3c09e2fd.
//
// Solidity: function grantRole(address addr, uint256 role) returns()
func (_Build *BuildSession) GrantRole(addr common.Address, role *big.Int) (*types.Transaction, error) {
	return _Build.Contract.GrantRole(&_Build.TransactOpts, addr, role)
}

func (_Build *C4TBuildSession) C4TGrantRole(addr common.Address, role *big.Int) (*c4tTypes.Transaction, error) {
	return _Build.Contract.C4TGrantRole(&_Build.TransactOpts, addr, role)
}

// GrantRole is a paid mutator transaction binding the contract method 0x3c09e2fd.
//
// Solidity: function grantRole(address addr, uint256 role) returns()
func (_Build *BuildTransactorSession) GrantRole(addr common.Address, role *big.Int) (*types.Transaction, error) {
	return _Build.Contract.GrantRole(&_Build.TransactOpts, addr, role)
}

// RevokeRole is a paid mutator transaction binding the contract method 0x0912ed77.
//
// Solidity: function revokeRole(address addr, uint256 role) returns()
func (_Build *BuildTransactor) RevokeRole(opts *bind.TransactOpts, addr common.Address, role *big.Int) (*types.Transaction, error) {
	return _Build.contract.Transact(opts, "revokeRole", addr, role)
}

// RevokeRole is a paid mutator transaction binding the contract method 0x0912ed77.
//
// Solidity: function revokeRole(address addr, uint256 role) returns()
func (_Build *BuildSession) RevokeRole(addr common.Address, role *big.Int) (*types.Transaction, error) {
	return _Build.Contract.RevokeRole(&_Build.TransactOpts, addr, role)
}

// RevokeRole is a paid mutator transaction binding the contract method 0x0912ed77.
//
// Solidity: function revokeRole(address addr, uint256 role) returns()
func (_Build *BuildTransactorSession) RevokeRole(addr common.Address, role *big.Int) (*types.Transaction, error) {
	return _Build.Contract.RevokeRole(&_Build.TransactOpts, addr, role)
}

// BuildDropRoleIterator is returned from FilterDropRole and is used to iterate over the raw logs and unpacked data for DropRole events raised by the Build contract.
type BuildDropRoleIterator struct {
	Event *BuildDropRole // Event containing the contract specifics and raw log

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
func (it *BuildDropRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuildDropRole)
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
		it.Event = new(BuildDropRole)
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
func (it *BuildDropRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuildDropRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuildDropRole represents a DropRole event raised by the Build contract.
type BuildDropRole struct {
	Addr common.Address
	Role *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterDropRole is a free log retrieval operation binding the contract event 0xcfa5316bd1be4ceb62f363b0a162f322c33ba870641138cd8600dd4fa603fc3b.
//
// Solidity: event DropRole(address addr, uint256 role)
func (_Build *BuildFilterer) FilterDropRole(opts *bind.FilterOpts) (*BuildDropRoleIterator, error) {

	logs, sub, err := _Build.contract.FilterLogs(opts, "DropRole")
	if err != nil {
		return nil, err
	}
	return &BuildDropRoleIterator{contract: _Build.contract, event: "DropRole", logs: logs, sub: sub}, nil
}

// WatchDropRole is a free log subscription operation binding the contract event 0xcfa5316bd1be4ceb62f363b0a162f322c33ba870641138cd8600dd4fa603fc3b.
//
// Solidity: event DropRole(address addr, uint256 role)
func (_Build *BuildFilterer) WatchDropRole(opts *bind.WatchOpts, sink chan<- *BuildDropRole) (event.Subscription, error) {

	logs, sub, err := _Build.contract.WatchLogs(opts, "DropRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuildDropRole)
				if err := _Build.contract.UnpackLog(event, "DropRole", log); err != nil {
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

// ParseDropRole is a log parse operation binding the contract event 0xcfa5316bd1be4ceb62f363b0a162f322c33ba870641138cd8600dd4fa603fc3b.
//
// Solidity: event DropRole(address addr, uint256 role)
func (_Build *BuildFilterer) ParseDropRole(log types.Log) (*BuildDropRole, error) {
	event := new(BuildDropRole)
	if err := _Build.contract.UnpackLog(event, "DropRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BuildSetRoleIterator is returned from FilterSetRole and is used to iterate over the raw logs and unpacked data for SetRole events raised by the Build contract.
type BuildSetRoleIterator struct {
	Event *BuildSetRole // Event containing the contract specifics and raw log

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
func (it *BuildSetRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BuildSetRole)
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
		it.Event = new(BuildSetRole)
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
func (it *BuildSetRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BuildSetRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BuildSetRole represents a SetRole event raised by the Build contract.
type BuildSetRole struct {
	Addr common.Address
	Role *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterSetRole is a free log retrieval operation binding the contract event 0x385a9c70004a48177c93b74796d77d5ebf7e1248f9e2369624514da454cd01b0.
//
// Solidity: event SetRole(address addr, uint256 role)
func (_Build *BuildFilterer) FilterSetRole(opts *bind.FilterOpts) (*BuildSetRoleIterator, error) {

	logs, sub, err := _Build.contract.FilterLogs(opts, "SetRole")
	if err != nil {
		return nil, err
	}
	return &BuildSetRoleIterator{contract: _Build.contract, event: "SetRole", logs: logs, sub: sub}, nil
}

// WatchSetRole is a free log subscription operation binding the contract event 0x385a9c70004a48177c93b74796d77d5ebf7e1248f9e2369624514da454cd01b0.
//
// Solidity: event SetRole(address addr, uint256 role)
func (_Build *BuildFilterer) WatchSetRole(opts *bind.WatchOpts, sink chan<- *BuildSetRole) (event.Subscription, error) {

	logs, sub, err := _Build.contract.WatchLogs(opts, "SetRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BuildSetRole)
				if err := _Build.contract.UnpackLog(event, "SetRole", log); err != nil {
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

// ParseSetRole is a log parse operation binding the contract event 0x385a9c70004a48177c93b74796d77d5ebf7e1248f9e2369624514da454cd01b0.
//
// Solidity: event SetRole(address addr, uint256 role)
func (_Build *BuildFilterer) ParseSetRole(log types.Log) (*BuildSetRole, error) {
	event := new(BuildSetRole)
	if err := _Build.contract.UnpackLog(event, "SetRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

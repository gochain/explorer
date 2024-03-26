// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package tokens

import (
	"math/big"
	"strings"

	gochain "github.com/gochain/gochain/v4"
	"github.com/gochain/gochain/v4/accounts/abi"
	"github.com/gochain/gochain/v4/accounts/abi/bind"
	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/core/types"
	"github.com/gochain/gochain/v4/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = gochain.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// UpgradeableABI is the input ABI used to generate the binding from.
const UpgradeableABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"val\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"target\",\"outputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"target\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Resumed\",\"type\":\"event\"}]"

// Upgradeable is an auto generated Go binding around an GoChain contract.
type Upgradeable struct {
	UpgradeableCaller     // Read-only binding to the contract
	UpgradeableTransactor // Write-only binding to the contract
	UpgradeableFilterer   // Log filterer for contract events
}

// UpgradeableCaller is an auto generated read-only Go binding around an GoChain contract.
type UpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableTransactor is an auto generated write-only Go binding around an GoChain contract.
type UpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableFilterer is an auto generated log filtering Go binding around an GoChain contract events.
type UpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UpgradeableSession is an auto generated Go binding around an GoChain contract,
// with pre-set call and transact options.
type UpgradeableSession struct {
	Contract     *Upgradeable      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UpgradeableCallerSession is an auto generated read-only Go binding around an GoChain contract,
// with pre-set call options.
type UpgradeableCallerSession struct {
	Contract *UpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// UpgradeableTransactorSession is an auto generated write-only Go binding around an GoChain contract,
// with pre-set transact options.
type UpgradeableTransactorSession struct {
	Contract     *UpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// UpgradeableRaw is an auto generated low-level Go binding around an GoChain contract.
type UpgradeableRaw struct {
	Contract *Upgradeable // Generic contract binding to access the raw methods on
}

// UpgradeableCallerRaw is an auto generated low-level read-only Go binding around an GoChain contract.
type UpgradeableCallerRaw struct {
	Contract *UpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// UpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an GoChain contract.
type UpgradeableTransactorRaw struct {
	Contract *UpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUpgradeable creates a new instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeable(address common.Address, backend bind.ContractBackend) (*Upgradeable, error) {
	contract, err := bindUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Upgradeable{UpgradeableCaller: UpgradeableCaller{contract: contract}, UpgradeableTransactor: UpgradeableTransactor{contract: contract}, UpgradeableFilterer: UpgradeableFilterer{contract: contract}}, nil
}

// NewUpgradeableCaller creates a new read-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*UpgradeableCaller, error) {
	contract, err := bindUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableCaller{contract: contract}, nil
}

// NewUpgradeableTransactor creates a new write-only instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*UpgradeableTransactor, error) {
	contract, err := bindUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UpgradeableTransactor{contract: contract}, nil
}

// NewUpgradeableFilterer creates a new log filterer instance of Upgradeable, bound to a specific deployed contract.
func NewUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*UpgradeableFilterer, error) {
	contract, err := bindUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UpgradeableFilterer{contract: contract}, nil
}

// bindUpgradeable binds a generic wrapper to an already deployed contract.
func bindUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UpgradeableABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.UpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.UpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Upgradeable *UpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Upgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Upgradeable *UpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Upgradeable *UpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Upgradeable.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address addr)
func (_Upgradeable *UpgradeableCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Upgradeable.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address addr)
func (_Upgradeable *UpgradeableSession) Owner() (common.Address, error) {
	return _Upgradeable.Contract.Owner(&_Upgradeable.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address addr)
func (_Upgradeable *UpgradeableCallerSession) Owner() (common.Address, error) {
	return _Upgradeable.Contract.Owner(&_Upgradeable.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool val)
func (_Upgradeable *UpgradeableCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Upgradeable.contract.Call(opts, &out, "paused")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool val)
func (_Upgradeable *UpgradeableSession) Paused() (bool, error) {
	return _Upgradeable.Contract.Paused(&_Upgradeable.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() view returns(bool val)
func (_Upgradeable *UpgradeableCallerSession) Paused() (bool, error) {
	return _Upgradeable.Contract.Paused(&_Upgradeable.CallOpts)
}

// Target is a free data retrieval call binding the contract method 0xd4b83992.
//
// Solidity: function target() view returns(address addr)
func (_Upgradeable *UpgradeableCaller) Target(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Upgradeable.contract.Call(opts, &out, "target")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Target is a free data retrieval call binding the contract method 0xd4b83992.
//
// Solidity: function target() view returns(address addr)
func (_Upgradeable *UpgradeableSession) Target() (common.Address, error) {
	return _Upgradeable.Contract.Target(&_Upgradeable.CallOpts)
}

// Target is a free data retrieval call binding the contract method 0xd4b83992.
//
// Solidity: function target() view returns(address addr)
func (_Upgradeable *UpgradeableCallerSession) Target() (common.Address, error) {
	return _Upgradeable.Contract.Target(&_Upgradeable.CallOpts)
}

// UpgradeablePausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the Upgradeable contract.
type UpgradeablePausedIterator struct {
	Event *UpgradeablePaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log       // Log channel receiving the found contract events
	sub  gochain.Subscription // Subscription for errors, completion and termination
	done bool                 // Whether the subscription completed delivering logs
	fail error                // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeablePausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeablePaused)
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
		it.Event = new(UpgradeablePaused)
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
func (it *UpgradeablePausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeablePausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeablePaused represents a Paused event raised by the Upgradeable contract.
type UpgradeablePaused struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x9e87fac88ff661f02d44f95383c817fece4bce600a3dab7a54406878b965e752.
//
// Solidity: event Paused()
func (_Upgradeable *UpgradeableFilterer) FilterPaused(opts *bind.FilterOpts) (*UpgradeablePausedIterator, error) {

	logs, sub, err := _Upgradeable.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &UpgradeablePausedIterator{contract: _Upgradeable.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x9e87fac88ff661f02d44f95383c817fece4bce600a3dab7a54406878b965e752.
//
// Solidity: event Paused()
func (_Upgradeable *UpgradeableFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *UpgradeablePaused) (event.Subscription, error) {

	logs, sub, err := _Upgradeable.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeablePaused)
				if err := _Upgradeable.contract.UnpackLog(event, "Paused", log); err != nil {
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

// ParsePaused is a log parse operation binding the contract event 0x9e87fac88ff661f02d44f95383c817fece4bce600a3dab7a54406878b965e752.
//
// Solidity: event Paused()
func (_Upgradeable *UpgradeableFilterer) ParsePaused(log types.Log) (*UpgradeablePaused, error) {
	event := new(UpgradeablePaused)
	if err := _Upgradeable.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeableResumedIterator is returned from FilterResumed and is used to iterate over the raw logs and unpacked data for Resumed events raised by the Upgradeable contract.
type UpgradeableResumedIterator struct {
	Event *UpgradeableResumed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log       // Log channel receiving the found contract events
	sub  gochain.Subscription // Subscription for errors, completion and termination
	done bool                 // Whether the subscription completed delivering logs
	fail error                // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeableResumedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeableResumed)
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
		it.Event = new(UpgradeableResumed)
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
func (it *UpgradeableResumedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeableResumedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeableResumed represents a Resumed event raised by the Upgradeable contract.
type UpgradeableResumed struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterResumed is a free log retrieval operation binding the contract event 0x62451d457bc659158be6e6247f56ec1df424a5c7597f71c20c2bc44e0965c8f9.
//
// Solidity: event Resumed()
func (_Upgradeable *UpgradeableFilterer) FilterResumed(opts *bind.FilterOpts) (*UpgradeableResumedIterator, error) {

	logs, sub, err := _Upgradeable.contract.FilterLogs(opts, "Resumed")
	if err != nil {
		return nil, err
	}
	return &UpgradeableResumedIterator{contract: _Upgradeable.contract, event: "Resumed", logs: logs, sub: sub}, nil
}

// WatchResumed is a free log subscription operation binding the contract event 0x62451d457bc659158be6e6247f56ec1df424a5c7597f71c20c2bc44e0965c8f9.
//
// Solidity: event Resumed()
func (_Upgradeable *UpgradeableFilterer) WatchResumed(opts *bind.WatchOpts, sink chan<- *UpgradeableResumed) (event.Subscription, error) {

	logs, sub, err := _Upgradeable.contract.WatchLogs(opts, "Resumed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeableResumed)
				if err := _Upgradeable.contract.UnpackLog(event, "Resumed", log); err != nil {
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

// ParseResumed is a log parse operation binding the contract event 0x62451d457bc659158be6e6247f56ec1df424a5c7597f71c20c2bc44e0965c8f9.
//
// Solidity: event Resumed()
func (_Upgradeable *UpgradeableFilterer) ParseResumed(log types.Log) (*UpgradeableResumed, error) {
	event := new(UpgradeableResumed)
	if err := _Upgradeable.contract.UnpackLog(event, "Resumed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// UpgradeableUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Upgradeable contract.
type UpgradeableUpgradedIterator struct {
	Event *UpgradeableUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log       // Log channel receiving the found contract events
	sub  gochain.Subscription // Subscription for errors, completion and termination
	done bool                 // Whether the subscription completed delivering logs
	fail error                // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *UpgradeableUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(UpgradeableUpgraded)
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
		it.Event = new(UpgradeableUpgraded)
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
func (it *UpgradeableUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *UpgradeableUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// UpgradeableUpgraded represents a Upgraded event raised by the Upgradeable contract.
type UpgradeableUpgraded struct {
	Target common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed target)
func (_Upgradeable *UpgradeableFilterer) FilterUpgraded(opts *bind.FilterOpts, target []common.Address) (*UpgradeableUpgradedIterator, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _Upgradeable.contract.FilterLogs(opts, "Upgraded", targetRule)
	if err != nil {
		return nil, err
	}
	return &UpgradeableUpgradedIterator{contract: _Upgradeable.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed target)
func (_Upgradeable *UpgradeableFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *UpgradeableUpgraded, target []common.Address) (event.Subscription, error) {

	var targetRule []interface{}
	for _, targetItem := range target {
		targetRule = append(targetRule, targetItem)
	}

	logs, sub, err := _Upgradeable.contract.WatchLogs(opts, "Upgraded", targetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(UpgradeableUpgraded)
				if err := _Upgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed target)
func (_Upgradeable *UpgradeableFilterer) ParseUpgraded(log types.Log) (*UpgradeableUpgraded, error) {
	event := new(UpgradeableUpgraded)
	if err := _Upgradeable.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

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
)

// IClaimRegistryClaim is an auto generated low-level Go binding around an user-defined struct.
type IClaimRegistryClaim struct {
	PropertyType string
	PropertyId   string
	Evidence     string
	Method       string
}

// IClaimRegistryClaimRef is an auto generated low-level Go binding around an user-defined struct.
type IClaimRegistryClaimRef struct {
	Ref string
	Key string
}

// ContractsMetaData contains all meta data concerning the Contracts contract.
var ContractsMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"claimer\",\"type\":\"address\"}],\"name\":\"ClaimRefRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"claimer\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structIClaimRegistry.ClaimRef\",\"name\":\"ref\",\"type\":\"tuple\"}],\"name\":\"ClaimRefUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"claimer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"propertyType\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"propertyId\",\"type\":\"string\"}],\"name\":\"ClaimRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"claimer\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"string\",\"name\":\"propertyType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"propertyId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"evidence\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"}],\"indexed\":false,\"internalType\":\"structIClaimRegistry.Claim\",\"name\":\"claim\",\"type\":\"tuple\"}],\"name\":\"ClaimUpdated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allClaimKeys\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allClaimRefs\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"allClaims\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"propertyType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"propertyId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"evidence\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"listClaims\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"string[2]\",\"name\":\"\",\"type\":\"string[2]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"propertyType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"propertyId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"evidence\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"ref\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"key\",\"type\":\"string\"}],\"name\":\"registerRef\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"propertyType\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"propertyId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"method\",\"type\":\"string\"}],\"name\":\"remove\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"removeRef\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ContractsABI is the input ABI used to generate the binding from.
// Deprecated: Use ContractsMetaData.ABI instead.
var ContractsABI = ContractsMetaData.ABI

// Contracts is an auto generated Go binding around an Ethereum contract.
type Contracts struct {
	ContractsCaller     // Read-only binding to the contract
	ContractsTransactor // Write-only binding to the contract
	ContractsFilterer   // Log filterer for contract events
}

// ContractsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractsSession struct {
	Contract     *Contracts        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractsCallerSession struct {
	Contract *ContractsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ContractsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractsTransactorSession struct {
	Contract     *ContractsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ContractsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractsRaw struct {
	Contract *Contracts // Generic contract binding to access the raw methods on
}

// ContractsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractsCallerRaw struct {
	Contract *ContractsCaller // Generic read-only contract binding to access the raw methods on
}

// ContractsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractsTransactorRaw struct {
	Contract *ContractsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContracts creates a new instance of Contracts, bound to a specific deployed contract.
func NewContracts(address common.Address, backend bind.ContractBackend) (*Contracts, error) {
	contract, err := bindContracts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contracts{ContractsCaller: ContractsCaller{contract: contract}, ContractsTransactor: ContractsTransactor{contract: contract}, ContractsFilterer: ContractsFilterer{contract: contract}}, nil
}

// NewContractsCaller creates a new read-only instance of Contracts, bound to a specific deployed contract.
func NewContractsCaller(address common.Address, caller bind.ContractCaller) (*ContractsCaller, error) {
	contract, err := bindContracts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsCaller{contract: contract}, nil
}

// NewContractsTransactor creates a new write-only instance of Contracts, bound to a specific deployed contract.
func NewContractsTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractsTransactor, error) {
	contract, err := bindContracts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsTransactor{contract: contract}, nil
}

// NewContractsFilterer creates a new log filterer instance of Contracts, bound to a specific deployed contract.
func NewContractsFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractsFilterer, error) {
	contract, err := bindContracts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractsFilterer{contract: contract}, nil
}

// bindContracts binds a generic wrapper to an already deployed contract.
func bindContracts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.ContractsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transact(opts, method, params...)
}

// AllClaimKeys is a free data retrieval call binding the contract method 0xc3552867.
//
// Solidity: function allClaimKeys(address , uint256 ) view returns(uint256)
func (_Contracts *ContractsCaller) AllClaimKeys(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Contracts.contract.Call(opts, &out, "allClaimKeys", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AllClaimKeys is a free data retrieval call binding the contract method 0xc3552867.
//
// Solidity: function allClaimKeys(address , uint256 ) view returns(uint256)
func (_Contracts *ContractsSession) AllClaimKeys(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Contracts.Contract.AllClaimKeys(&_Contracts.CallOpts, arg0, arg1)
}

// AllClaimKeys is a free data retrieval call binding the contract method 0xc3552867.
//
// Solidity: function allClaimKeys(address , uint256 ) view returns(uint256)
func (_Contracts *ContractsCallerSession) AllClaimKeys(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _Contracts.Contract.AllClaimKeys(&_Contracts.CallOpts, arg0, arg1)
}

// AllClaimRefs is a free data retrieval call binding the contract method 0x0f937ac3.
//
// Solidity: function allClaimRefs(address ) view returns(string ref, string key)
func (_Contracts *ContractsCaller) AllClaimRefs(opts *bind.CallOpts, arg0 common.Address) (struct {
	Ref string
	Key string
}, error) {
	var out []interface{}
	err := _Contracts.contract.Call(opts, &out, "allClaimRefs", arg0)

	outstruct := new(struct {
		Ref string
		Key string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Ref = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Key = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

// AllClaimRefs is a free data retrieval call binding the contract method 0x0f937ac3.
//
// Solidity: function allClaimRefs(address ) view returns(string ref, string key)
func (_Contracts *ContractsSession) AllClaimRefs(arg0 common.Address) (struct {
	Ref string
	Key string
}, error) {
	return _Contracts.Contract.AllClaimRefs(&_Contracts.CallOpts, arg0)
}

// AllClaimRefs is a free data retrieval call binding the contract method 0x0f937ac3.
//
// Solidity: function allClaimRefs(address ) view returns(string ref, string key)
func (_Contracts *ContractsCallerSession) AllClaimRefs(arg0 common.Address) (struct {
	Ref string
	Key string
}, error) {
	return _Contracts.Contract.AllClaimRefs(&_Contracts.CallOpts, arg0)
}

// AllClaims is a free data retrieval call binding the contract method 0x1a8598a7.
//
// Solidity: function allClaims(uint256 ) view returns(string propertyType, string propertyId, string evidence, string method)
func (_Contracts *ContractsCaller) AllClaims(opts *bind.CallOpts, arg0 *big.Int) (struct {
	PropertyType string
	PropertyId   string
	Evidence     string
	Method       string
}, error) {
	var out []interface{}
	err := _Contracts.contract.Call(opts, &out, "allClaims", arg0)

	outstruct := new(struct {
		PropertyType string
		PropertyId   string
		Evidence     string
		Method       string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.PropertyType = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.PropertyId = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.Evidence = *abi.ConvertType(out[2], new(string)).(*string)
	outstruct.Method = *abi.ConvertType(out[3], new(string)).(*string)

	return *outstruct, err

}

// AllClaims is a free data retrieval call binding the contract method 0x1a8598a7.
//
// Solidity: function allClaims(uint256 ) view returns(string propertyType, string propertyId, string evidence, string method)
func (_Contracts *ContractsSession) AllClaims(arg0 *big.Int) (struct {
	PropertyType string
	PropertyId   string
	Evidence     string
	Method       string
}, error) {
	return _Contracts.Contract.AllClaims(&_Contracts.CallOpts, arg0)
}

// AllClaims is a free data retrieval call binding the contract method 0x1a8598a7.
//
// Solidity: function allClaims(uint256 ) view returns(string propertyType, string propertyId, string evidence, string method)
func (_Contracts *ContractsCallerSession) AllClaims(arg0 *big.Int) (struct {
	PropertyType string
	PropertyId   string
	Evidence     string
	Method       string
}, error) {
	return _Contracts.Contract.AllClaims(&_Contracts.CallOpts, arg0)
}

// ListClaims is a free data retrieval call binding the contract method 0x05562ae8.
//
// Solidity: function listClaims(address account) view returns(uint256[], string[2])
func (_Contracts *ContractsCaller) ListClaims(opts *bind.CallOpts, account common.Address) ([]*big.Int, [2]string, error) {
	var out []interface{}
	err := _Contracts.contract.Call(opts, &out, "listClaims", account)

	if err != nil {
		return *new([]*big.Int), *new([2]string), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)
	out1 := *abi.ConvertType(out[1], new([2]string)).(*[2]string)

	return out0, out1, err

}

// ListClaims is a free data retrieval call binding the contract method 0x05562ae8.
//
// Solidity: function listClaims(address account) view returns(uint256[], string[2])
func (_Contracts *ContractsSession) ListClaims(account common.Address) ([]*big.Int, [2]string, error) {
	return _Contracts.Contract.ListClaims(&_Contracts.CallOpts, account)
}

// ListClaims is a free data retrieval call binding the contract method 0x05562ae8.
//
// Solidity: function listClaims(address account) view returns(uint256[], string[2])
func (_Contracts *ContractsCallerSession) ListClaims(account common.Address) ([]*big.Int, [2]string, error) {
	return _Contracts.Contract.ListClaims(&_Contracts.CallOpts, account)
}

// Register is a paid mutator transaction binding the contract method 0x0e24c52c.
//
// Solidity: function register(string propertyType, string propertyId, string evidence, string method) returns()
func (_Contracts *ContractsTransactor) Register(opts *bind.TransactOpts, propertyType string, propertyId string, evidence string, method string) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "register", propertyType, propertyId, evidence, method)
}

// Register is a paid mutator transaction binding the contract method 0x0e24c52c.
//
// Solidity: function register(string propertyType, string propertyId, string evidence, string method) returns()
func (_Contracts *ContractsSession) Register(propertyType string, propertyId string, evidence string, method string) (*types.Transaction, error) {
	return _Contracts.Contract.Register(&_Contracts.TransactOpts, propertyType, propertyId, evidence, method)
}

// Register is a paid mutator transaction binding the contract method 0x0e24c52c.
//
// Solidity: function register(string propertyType, string propertyId, string evidence, string method) returns()
func (_Contracts *ContractsTransactorSession) Register(propertyType string, propertyId string, evidence string, method string) (*types.Transaction, error) {
	return _Contracts.Contract.Register(&_Contracts.TransactOpts, propertyType, propertyId, evidence, method)
}

// RegisterRef is a paid mutator transaction binding the contract method 0xc0a9f9d2.
//
// Solidity: function registerRef(string ref, string key) returns()
func (_Contracts *ContractsTransactor) RegisterRef(opts *bind.TransactOpts, ref string, key string) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "registerRef", ref, key)
}

// RegisterRef is a paid mutator transaction binding the contract method 0xc0a9f9d2.
//
// Solidity: function registerRef(string ref, string key) returns()
func (_Contracts *ContractsSession) RegisterRef(ref string, key string) (*types.Transaction, error) {
	return _Contracts.Contract.RegisterRef(&_Contracts.TransactOpts, ref, key)
}

// RegisterRef is a paid mutator transaction binding the contract method 0xc0a9f9d2.
//
// Solidity: function registerRef(string ref, string key) returns()
func (_Contracts *ContractsTransactorSession) RegisterRef(ref string, key string) (*types.Transaction, error) {
	return _Contracts.Contract.RegisterRef(&_Contracts.TransactOpts, ref, key)
}

// Remove is a paid mutator transaction binding the contract method 0xda40fa77.
//
// Solidity: function remove(string propertyType, string propertyId, string method) returns()
func (_Contracts *ContractsTransactor) Remove(opts *bind.TransactOpts, propertyType string, propertyId string, method string) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "remove", propertyType, propertyId, method)
}

// Remove is a paid mutator transaction binding the contract method 0xda40fa77.
//
// Solidity: function remove(string propertyType, string propertyId, string method) returns()
func (_Contracts *ContractsSession) Remove(propertyType string, propertyId string, method string) (*types.Transaction, error) {
	return _Contracts.Contract.Remove(&_Contracts.TransactOpts, propertyType, propertyId, method)
}

// Remove is a paid mutator transaction binding the contract method 0xda40fa77.
//
// Solidity: function remove(string propertyType, string propertyId, string method) returns()
func (_Contracts *ContractsTransactorSession) Remove(propertyType string, propertyId string, method string) (*types.Transaction, error) {
	return _Contracts.Contract.Remove(&_Contracts.TransactOpts, propertyType, propertyId, method)
}

// RemoveRef is a paid mutator transaction binding the contract method 0x387f5d3b.
//
// Solidity: function removeRef() returns()
func (_Contracts *ContractsTransactor) RemoveRef(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "removeRef")
}

// RemoveRef is a paid mutator transaction binding the contract method 0x387f5d3b.
//
// Solidity: function removeRef() returns()
func (_Contracts *ContractsSession) RemoveRef() (*types.Transaction, error) {
	return _Contracts.Contract.RemoveRef(&_Contracts.TransactOpts)
}

// RemoveRef is a paid mutator transaction binding the contract method 0x387f5d3b.
//
// Solidity: function removeRef() returns()
func (_Contracts *ContractsTransactorSession) RemoveRef() (*types.Transaction, error) {
	return _Contracts.Contract.RemoveRef(&_Contracts.TransactOpts)
}

// ContractsClaimRefRemovedIterator is returned from FilterClaimRefRemoved and is used to iterate over the raw logs and unpacked data for ClaimRefRemoved events raised by the Contracts contract.
type ContractsClaimRefRemovedIterator struct {
	Event *ContractsClaimRefRemoved // Event containing the contract specifics and raw log

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
func (it *ContractsClaimRefRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractsClaimRefRemoved)
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
		it.Event = new(ContractsClaimRefRemoved)
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
func (it *ContractsClaimRefRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractsClaimRefRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractsClaimRefRemoved represents a ClaimRefRemoved event raised by the Contracts contract.
type ContractsClaimRefRemoved struct {
	Claimer common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClaimRefRemoved is a free log retrieval operation binding the contract event 0x209c90080992a7d98eab6016675da7fd54af9f510923cea453cb66d793c66762.
//
// Solidity: event ClaimRefRemoved(address claimer)
func (_Contracts *ContractsFilterer) FilterClaimRefRemoved(opts *bind.FilterOpts) (*ContractsClaimRefRemovedIterator, error) {

	logs, sub, err := _Contracts.contract.FilterLogs(opts, "ClaimRefRemoved")
	if err != nil {
		return nil, err
	}
	return &ContractsClaimRefRemovedIterator{contract: _Contracts.contract, event: "ClaimRefRemoved", logs: logs, sub: sub}, nil
}

// WatchClaimRefRemoved is a free log subscription operation binding the contract event 0x209c90080992a7d98eab6016675da7fd54af9f510923cea453cb66d793c66762.
//
// Solidity: event ClaimRefRemoved(address claimer)
func (_Contracts *ContractsFilterer) WatchClaimRefRemoved(opts *bind.WatchOpts, sink chan<- *ContractsClaimRefRemoved) (event.Subscription, error) {

	logs, sub, err := _Contracts.contract.WatchLogs(opts, "ClaimRefRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractsClaimRefRemoved)
				if err := _Contracts.contract.UnpackLog(event, "ClaimRefRemoved", log); err != nil {
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

// ParseClaimRefRemoved is a log parse operation binding the contract event 0x209c90080992a7d98eab6016675da7fd54af9f510923cea453cb66d793c66762.
//
// Solidity: event ClaimRefRemoved(address claimer)
func (_Contracts *ContractsFilterer) ParseClaimRefRemoved(log types.Log) (*ContractsClaimRefRemoved, error) {
	event := new(ContractsClaimRefRemoved)
	if err := _Contracts.contract.UnpackLog(event, "ClaimRefRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractsClaimRefUpdatedIterator is returned from FilterClaimRefUpdated and is used to iterate over the raw logs and unpacked data for ClaimRefUpdated events raised by the Contracts contract.
type ContractsClaimRefUpdatedIterator struct {
	Event *ContractsClaimRefUpdated // Event containing the contract specifics and raw log

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
func (it *ContractsClaimRefUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractsClaimRefUpdated)
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
		it.Event = new(ContractsClaimRefUpdated)
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
func (it *ContractsClaimRefUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractsClaimRefUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractsClaimRefUpdated represents a ClaimRefUpdated event raised by the Contracts contract.
type ContractsClaimRefUpdated struct {
	Claimer common.Address
	Ref     IClaimRegistryClaimRef
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClaimRefUpdated is a free log retrieval operation binding the contract event 0x77701c707374125153cb279a5509b5ce7a416c3167f5060a07ebe9d680e7995b.
//
// Solidity: event ClaimRefUpdated(address claimer, (string,string) ref)
func (_Contracts *ContractsFilterer) FilterClaimRefUpdated(opts *bind.FilterOpts) (*ContractsClaimRefUpdatedIterator, error) {

	logs, sub, err := _Contracts.contract.FilterLogs(opts, "ClaimRefUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractsClaimRefUpdatedIterator{contract: _Contracts.contract, event: "ClaimRefUpdated", logs: logs, sub: sub}, nil
}

// WatchClaimRefUpdated is a free log subscription operation binding the contract event 0x77701c707374125153cb279a5509b5ce7a416c3167f5060a07ebe9d680e7995b.
//
// Solidity: event ClaimRefUpdated(address claimer, (string,string) ref)
func (_Contracts *ContractsFilterer) WatchClaimRefUpdated(opts *bind.WatchOpts, sink chan<- *ContractsClaimRefUpdated) (event.Subscription, error) {

	logs, sub, err := _Contracts.contract.WatchLogs(opts, "ClaimRefUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractsClaimRefUpdated)
				if err := _Contracts.contract.UnpackLog(event, "ClaimRefUpdated", log); err != nil {
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

// ParseClaimRefUpdated is a log parse operation binding the contract event 0x77701c707374125153cb279a5509b5ce7a416c3167f5060a07ebe9d680e7995b.
//
// Solidity: event ClaimRefUpdated(address claimer, (string,string) ref)
func (_Contracts *ContractsFilterer) ParseClaimRefUpdated(log types.Log) (*ContractsClaimRefUpdated, error) {
	event := new(ContractsClaimRefUpdated)
	if err := _Contracts.contract.UnpackLog(event, "ClaimRefUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractsClaimRemovedIterator is returned from FilterClaimRemoved and is used to iterate over the raw logs and unpacked data for ClaimRemoved events raised by the Contracts contract.
type ContractsClaimRemovedIterator struct {
	Event *ContractsClaimRemoved // Event containing the contract specifics and raw log

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
func (it *ContractsClaimRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractsClaimRemoved)
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
		it.Event = new(ContractsClaimRemoved)
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
func (it *ContractsClaimRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractsClaimRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractsClaimRemoved represents a ClaimRemoved event raised by the Contracts contract.
type ContractsClaimRemoved struct {
	Claimer      common.Address
	PropertyType string
	PropertyId   string
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterClaimRemoved is a free log retrieval operation binding the contract event 0x2cab426b356a1e45c43f79a1c815849191e75495e013418b790347223685117f.
//
// Solidity: event ClaimRemoved(address claimer, string propertyType, string propertyId)
func (_Contracts *ContractsFilterer) FilterClaimRemoved(opts *bind.FilterOpts) (*ContractsClaimRemovedIterator, error) {

	logs, sub, err := _Contracts.contract.FilterLogs(opts, "ClaimRemoved")
	if err != nil {
		return nil, err
	}
	return &ContractsClaimRemovedIterator{contract: _Contracts.contract, event: "ClaimRemoved", logs: logs, sub: sub}, nil
}

// WatchClaimRemoved is a free log subscription operation binding the contract event 0x2cab426b356a1e45c43f79a1c815849191e75495e013418b790347223685117f.
//
// Solidity: event ClaimRemoved(address claimer, string propertyType, string propertyId)
func (_Contracts *ContractsFilterer) WatchClaimRemoved(opts *bind.WatchOpts, sink chan<- *ContractsClaimRemoved) (event.Subscription, error) {

	logs, sub, err := _Contracts.contract.WatchLogs(opts, "ClaimRemoved")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractsClaimRemoved)
				if err := _Contracts.contract.UnpackLog(event, "ClaimRemoved", log); err != nil {
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

// ParseClaimRemoved is a log parse operation binding the contract event 0x2cab426b356a1e45c43f79a1c815849191e75495e013418b790347223685117f.
//
// Solidity: event ClaimRemoved(address claimer, string propertyType, string propertyId)
func (_Contracts *ContractsFilterer) ParseClaimRemoved(log types.Log) (*ContractsClaimRemoved, error) {
	event := new(ContractsClaimRemoved)
	if err := _Contracts.contract.UnpackLog(event, "ClaimRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ContractsClaimUpdatedIterator is returned from FilterClaimUpdated and is used to iterate over the raw logs and unpacked data for ClaimUpdated events raised by the Contracts contract.
type ContractsClaimUpdatedIterator struct {
	Event *ContractsClaimUpdated // Event containing the contract specifics and raw log

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
func (it *ContractsClaimUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractsClaimUpdated)
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
		it.Event = new(ContractsClaimUpdated)
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
func (it *ContractsClaimUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractsClaimUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractsClaimUpdated represents a ClaimUpdated event raised by the Contracts contract.
type ContractsClaimUpdated struct {
	Claimer common.Address
	Claim   IClaimRegistryClaim
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterClaimUpdated is a free log retrieval operation binding the contract event 0x38913474ea29a22c6b686bcd567285304ab34843603996723b03f0c05e9355fc.
//
// Solidity: event ClaimUpdated(address claimer, (string,string,string,string) claim)
func (_Contracts *ContractsFilterer) FilterClaimUpdated(opts *bind.FilterOpts) (*ContractsClaimUpdatedIterator, error) {

	logs, sub, err := _Contracts.contract.FilterLogs(opts, "ClaimUpdated")
	if err != nil {
		return nil, err
	}
	return &ContractsClaimUpdatedIterator{contract: _Contracts.contract, event: "ClaimUpdated", logs: logs, sub: sub}, nil
}

// WatchClaimUpdated is a free log subscription operation binding the contract event 0x38913474ea29a22c6b686bcd567285304ab34843603996723b03f0c05e9355fc.
//
// Solidity: event ClaimUpdated(address claimer, (string,string,string,string) claim)
func (_Contracts *ContractsFilterer) WatchClaimUpdated(opts *bind.WatchOpts, sink chan<- *ContractsClaimUpdated) (event.Subscription, error) {

	logs, sub, err := _Contracts.contract.WatchLogs(opts, "ClaimUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractsClaimUpdated)
				if err := _Contracts.contract.UnpackLog(event, "ClaimUpdated", log); err != nil {
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

// ParseClaimUpdated is a log parse operation binding the contract event 0x38913474ea29a22c6b686bcd567285304ab34843603996723b03f0c05e9355fc.
//
// Solidity: event ClaimUpdated(address claimer, (string,string,string,string) claim)
func (_Contracts *ContractsFilterer) ParseClaimUpdated(log types.Log) (*ContractsClaimUpdated, error) {
	event := new(ContractsClaimUpdated)
	if err := _Contracts.contract.UnpackLog(event, "ClaimUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

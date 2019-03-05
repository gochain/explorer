package backend

import (
	"math/big"
	"strings"

	"github.com/gochain-io/gochain/v3/accounts/abi"
	"github.com/gochain-io/gochain/v3/accounts/abi/bind"
	"github.com/gochain-io/gochain/v3/common"
)

// TokenABI is the input ABI used to generate the binding from.
const TokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"mintingFinished\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finishMinting\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_releaseTime\",\"type\":\"uint256\"}],\"name\":\"mintTimelocked\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"MintFinished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]"

var INTERFACE_IDENTIFIERS = map[string]string{
	"dd62ed3e": "allowance(address,address)",
	"095ea7b3": "approve(address,uint256)",
	"70a08231": "balanceOf(address)",
	"18160ddd": "totalSupply()",
	"a9059cbb": "transfer(address,uint256)",
	"23b872dd": "transferFrom(address,address,uint256)",
	"081812fc": "getApproved(uint256)",
	"e985e9c5": "isApprovedForAll(address,address)",
	"6352211e": "ownerOf(uint256)",
	"42842e0e": "safeTransferFrom(address,address,uint256)",
	"b88d4fde": "safeTransferFrom(address,address,uint256,bytes)",
	"a22cb465": "setApprovalForAll(address,bool)",
	"01ffc9a7": "supportsInterface(bytes4)",
	"249cb3fa": "canImplementInterfaceForAddress(bytes32,address)",
	"00fdd58e": "balanceOf(address,uint256)",
	"4e1273f4": "balanceOfBatch(address[],uint256[])",
	"2eb2c2d6": "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)",
	"f242432a": "safeTransferFrom(address,address,uint256,uint256,bytes)",
	"bc197c81": "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)",
	"f23a6e61": "onERC1155Received(address,address,uint256,uint256,bytes)",
	"0e89341c": "uri(uint256)",
	"150b7a02": "onERC721Received(address,address,uint256,bytes)",
	"06fdde03": "name()",
	"95d89b41": "symbol()",
	"c87b56dd": "tokenURI(uint256)",
	"4f6ccce7": "tokenByIndex(uint256)",
	"2f745c59": "tokenOfOwnerByIndex(address,uint256)",
	"313ce567": "decimals()",
	"be45fd62": "transfer(address,uint256,bytes)",
	"f6368f8a": "transfer(address,uint256,bytes,string)",
	"869e0e60": "decreaseSupply(uint256,address)",
	"124fc7e0": "increaseSupply(uint256,address)",
	"0023de29": "tokensReceived(address,address,address,uint256,bytes,bytes)",
	"75ab9782": "tokensToSend(address,address,address,uint256,bytes,bytes)",
	"959b8c3f": "authorizeOperator(address)",
	"fe9d9303": "burn(uint256,bytes)",
	"06e48538": "defaultOperators()",
	"556f0dc7": "granularity()",
	"d95b6371": "isOperatorFor(address,address)",
	"fc673c4f": "operatorBurn(address,uint256,bytes,bytes)",
	"62ad1b83": "operatorSend(address,address,uint256,bytes,bytes)",
	"fad8b32a": "revokeOperator(address)",
	"9bd9bbc6": "send(address,uint256,bytes)",
	"cae9ca51": "approveAndCall(address,uint256,bytes)",
	"d135ca1d": "decreaseAllowanceAndCall(address,uint256,bytes)",
	"5fd42775": "increaseAllowanceAndCall(address,uint256,bytes)",
	"4000aea0": "transferAndCall(address,uint256,bytes)",
	"c1d34b89": "transferFromAndCall(address,address,uint256,bytes)",
	"47089f62": "addVerified(address,bytes32)",
	"79f64720": "cancelAndReissue(address,address)",
	"cc397ed3": "getCurrentFor(address)",
	"f3221c7f": "hasHash(address,bytes32)",
	"197bc336": "holderAt(uint256)",
	"1aab9a9f": "holderCount()",
	"d4d7b19a": "isHolder(address)",
	"2da7293e": "isSuperseded(address)",
	"b9209e33": "isVerified(address)",
	"4487b392": "removeVerified(address)",
	"354b7b1d": "updateVerified(address,bytes32)",
}

var ERC_INTERFACE_IDENTIFIERS = map[string]map[string]string{
	"20": {
		"dd62ed3e": "allowance(address,address)",
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
	"165": {
		"01ffc9a7": "supportsInterface(bytes4)",
	},
	"721": {
		"095ea7b3": "approve(address,uint256)",
		"70a08231": "balanceOf(address)",
		"081812fc": "getApproved(uint256)",
		"e985e9c5": "isApprovedForAll(address,address)",
		"6352211e": "ownerOf(uint256)",
		"42842e0e": "safeTransferFrom(address,address,uint256)",
		"b88d4fde": "safeTransferFrom(address,address,uint256,bytes)",
		"a22cb465": "setApprovalForAll(address,bool)",
		"01ffc9a7": "supportsInterface(bytes4)",
		"23b872dd": "transferFrom(address,address,uint256)",
	},
	"721_receiver": {
		"150b7a02": "onERC721Received(address,address,uint256,bytes)",
	},
	"721_metadata": {
		"06fdde03": "name()",
		"95d89b41": "symbol()",
		"c87b56dd": "tokenURI(uint256)",
	},
	"721_enumerable": {
		"4f6ccce7": "tokenByIndex(uint256)",
		"2f745c59": "tokenOfOwnerByIndex(address,uint256)",
		"18160ddd": "totalSupply()",
	},
	"820": {
		"249cb3fa": "canImplementInterfaceForAddress(bytes32,address)",
	},
	"1155": {
		"00fdd58e": "balanceOf(address,uint256)",
		"4e1273f4": "balanceOfBatch(address[],uint256[])",
		"e985e9c5": "isApprovedForAll(address,address)",
		"2eb2c2d6": "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)",
		"f242432a": "safeTransferFrom(address,address,uint256,uint256,bytes)",
		"a22cb465": "setApprovalForAll(address,bool)",
	},
	"1155_receiver": {
		"bc197c81": "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)",
		"f23a6e61": "onERC1155Received(address,address,uint256,uint256,bytes)",
	},
	"1155_metadata": {
		"0e89341c": "uri(uint256)",
	},
	"223": {
		"70a08231": "balanceOf(address)",
		"313ce567": "decimals()",
		"06fdde03": "name()",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
		"a9059cbb": "transfer(address,uint256)",
		"be45fd62": "transfer(address,uint256,bytes)",
		"f6368f8a": "transfer(address,uint256,bytes,string)",
	},
	"621": {
		"869e0e60": "decreaseSupply(uint256,address)",
		"124fc7e0": "increaseSupply(uint256,address)",
	},
	"777": {
		"959b8c3f": "authorizeOperator(address)",
		"70a08231": "balanceOf(address)",
		"fe9d9303": "burn(uint256,bytes)",
		"06e48538": "defaultOperators()",
		"556f0dc7": "granularity()",
		"d95b6371": "isOperatorFor(address,address)",
		"06fdde03": "name()",
		"fc673c4f": "operatorBurn(address,uint256,bytes,bytes)",
		"62ad1b83": "operatorSend(address,address,uint256,bytes,bytes)",
		"fad8b32a": "revokeOperator(address)",
		"9bd9bbc6": "send(address,uint256,bytes)",
		"95d89b41": "symbol()",
		"18160ddd": "totalSupply()",
	},
	"777_receiver": {
		"0023de29": "tokensReceived(address,address,address,uint256,bytes,bytes)",
	},
	"777_sender": {
		"75ab9782": "tokensToSend(address,address,address,uint256,bytes,bytes)",
	},
	"827": {
		"cae9ca51": "approveAndCall(address,uint256,bytes)",
		"d135ca1d": "decreaseAllowanceAndCall(address,uint256,bytes)",
		"5fd42775": "increaseAllowanceAndCall(address,uint256,bytes)",
		"4000aea0": "transferAndCall(address,uint256,bytes)",
		"c1d34b89": "transferFromAndCall(address,address,uint256,bytes)",
	},
	"884": {
		"47089f62": "addVerified(address,bytes32)",
		"79f64720": "cancelAndReissue(address,address)",
		"cc397ed3": "getCurrentFor(address)",
		"f3221c7f": "hasHash(address,bytes32)",
		"197bc336": "holderAt(uint256)",
		"1aab9a9f": "holderCount()",
		"d4d7b19a": "isHolder(address)",
		"2da7293e": "isSuperseded(address)",
		"b9209e33": "isVerified(address)",
		"4487b392": "removeVerified(address)",
		"a9059cbb": "transfer(address,uint256)",
		"23b872dd": "transferFrom(address,address,uint256)",
		"354b7b1d": "updateVerified(address,bytes32)",
	},
}

// Token is an auto generated Go binding around an Ethereum contract.
type Token struct {
	TokenCaller     // Read-only binding to the contract
	TokenTransactor // Write-only binding to the contract
}

// TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenSession struct {
	Contract     *Token            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenCallerSession struct {
	Contract *TokenCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenTransactorSession struct {
	Contract     *TokenTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenRaw struct {
	Contract *Token // Generic contract binding to access the raw methods on
}

// TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenCallerRaw struct {
	Contract *TokenCaller // Generic read-only contract binding to access the raw methods on
}

// TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenTransactorRaw struct {
	Contract *TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenCaller creates a new read-only instance of Token, bound to a specific deployed contract.
func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, nil), nil
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_owner address) constant returns(balance uint256)
func (_Token *TokenCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "balanceOf", _owner)
	return *ret0, err
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() constant returns(uint256)
func (_Token *TokenCaller) Decimals(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "decimals")
	return *ret0, err
}

// TotalSupply is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function totalSupply() constant returns(uint256)
func (_Token *TokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "totalSupply")
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() constant returns(string)
func (_Token *TokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "name")
	return *ret0, err
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() constant returns(string)
func (_Token *TokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Token.contract.Call(opts, out, "symbol")
	return *ret0, err
}

func (_Token *TokenCaller) Types(byteCode string) []string {
	identifiers := map[string]bool{}
	for k := range INTERFACE_IDENTIFIERS {
		if strings.Contains(byteCode, k) {
			identifiers[k] = true
		}
	}
	types := map[string]bool{}
Loop:
	for k, v := range ERC_INTERFACE_IDENTIFIERS {
		for ercIdentifier := range v {
			if _, ok := identifiers[ercIdentifier]; !ok {
				continue Loop
			}
		}
		types[k] = true
	}
	v := make([]string, len(types))
	for key := range types {
		v = append(v, key)
	}
	return v
}

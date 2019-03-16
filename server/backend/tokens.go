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

type InterfaceName string

type InterfaceData struct {
	Value       string
	Description string
	Callable    bool
}

const (
	ADD_VERIFIED                        InterfaceName = "ADD_VERIFIED"
	ALLOWANCE                           InterfaceName = "ALLOWANCE"
	APPROVE                             InterfaceName = "APPROVE"
	APPROVE_AND_CALL                    InterfaceName = "APPROVE_AND_CALL"
	AUTHORIZE_OPERATOR                  InterfaceName = "AUTHORIZE_OPERATOR"
	BALANCE_OF                          InterfaceName = "BALANCE_OF"
	BALANCE_OF_1                        InterfaceName = "BALANCE_OF_1"
	BALANCE_OF_BATCH                    InterfaceName = "BALANCE_OF_BATCH"
	BURN                                InterfaceName = "BURN"
	BURN_1                              InterfaceName = "BURN_1"
	BURN_FROM                           InterfaceName = "BURN_FROM"
	CANCEL_AND_REISSUE                  InterfaceName = "CANCEL_AND_REISSUE"
	CAN_IMPLEMENT_INTERFACE_FOR_ADDRESS InterfaceName = "CAN_IMPLEMENT_INTERFACE_FOR_ADDRESS"
	CAP                                 InterfaceName = "CAP"
	DECIMALS                            InterfaceName = "DECIMALS"
	DECREASE_ALLOWANCE                  InterfaceName = "DECREASE_ALLOWANCE"
	DECREASE_ALLOWANCE_AND_CALL         InterfaceName = "DECREASE_ALLOWANCE_AND_CALL"
	DECREASE_SUPPLY                     InterfaceName = "DECREASE_SUPPLY"
	DEFAULT_OPERATORS                   InterfaceName = "DEFAULT_OPERATORS"
	GET_APPROVED                        InterfaceName = "GET_APPROVED"
	GET_CURRENT_FOR                     InterfaceName = "GET_CURRENT_FOR"
	GRANULARITY                         InterfaceName = "GRANULARITY"
	HAS_HASH                            InterfaceName = "HAS_HASH"
	HOLDER_AT                           InterfaceName = "HOLDER_AT"
	HOLDER_COUNT                        InterfaceName = "HOLDER_COUNT"
	INCREASE_ALLOWANCE                  InterfaceName = "INCREASE_ALLOWANCE"
	INCREASE_ALLOWANCE_AND_CALL         InterfaceName = "INCREASE_ALLOWANCE_AND_CALL"
	INCREASE_SUPPLY                     InterfaceName = "INCREASE_SUPPLY"
	IS_APPROVED_FOR_ALL                 InterfaceName = "IS_APPROVED_FOR_ALL"
	IS_HOLDER                           InterfaceName = "IS_HOLDER"
	IS_OPERATOR_FOR                     InterfaceName = "IS_OPERATOR_FOR"
	IS_SUPERSEDED                       InterfaceName = "IS_SUPERSEDED"
	IS_VERIFIED                         InterfaceName = "IS_VERIFIED"
	MINT                                InterfaceName = "MINT"
	NAME                                InterfaceName = "NAME"
	ON_ERC721_RECEIVED                  InterfaceName = "ON_ERC721_RECEIVED"
	ON_ERC1155_BATCH_RECEIVED           InterfaceName = "ON_ERC1155_BATCH_RECEIVED"
	ON_ERC1155_RECEIVED                 InterfaceName = "ON_ERC1155_RECEIVED"
	OPERATOR_BURN                       InterfaceName = "OPERATOR_BURN"
	OPERATOR_SEND                       InterfaceName = "OPERATOR_SEND"
	OWNER_OF                            InterfaceName = "OWNER_OF"
	REMOVE_VERIFIED                     InterfaceName = "REMOVE_VERIFIED"
	REVOKE_OPERATOR                     InterfaceName = "REVOKE_OPERATOR"
	SAFE_BATCH_TRANSFER_FROM            InterfaceName = "SAFE_BATCH_TRANSFER_FROM"
	SAFE_TRANSFER_FROM                  InterfaceName = "SAFE_TRANSFER_FROM"
	SAFE_TRANSFER_FROM_1                InterfaceName = "SAFE_TRANSFER_FROM_1"
	SEND                                InterfaceName = "SEND"
	SET_APPROVAL_FOR_ALL                InterfaceName = "SET_APPROVAL_FOR_ALL"
	SUPPORTS_INTERFACE                  InterfaceName = "SUPPORTS_INTERFACE"
	SYMBOL                              InterfaceName = "SYMBOL"
	TOKENS_RECEIVED                     InterfaceName = "TOKENS_RECEIVED"
	TOKENS_TO_SEND                      InterfaceName = "TOKENS_TO_SEND"
	TOKEN_BY_INDEX                      InterfaceName = "TOKEN_BY_INDEX"
	TOKEN_OF_OWNER_BY_INDEX             InterfaceName = "TOKEN_OF_OWNER_BY_INDEX"
	TOKEN_URI                           InterfaceName = "TOKEN_URI"
	TOTAL_SUPPLY                        InterfaceName = "TOTAL_SUPPLY"
	TRANSFER                            InterfaceName = "TRANSFER"
	TRANSFER_1                          InterfaceName = "TRANSFER_1"
	TRANSFER_2                          InterfaceName = "TRANSFER_2"
	TRANSFER_AND_CALL                   InterfaceName = "TRANSFER_AND_CALL"
	TRANSFER_FROM                       InterfaceName = "TRANSFER_FROM"
	TRANSFER_FROM_AND_CALL              InterfaceName = "TRANSFER_FROM_AND_CALL"
	UPDATE_VERIFIED                     InterfaceName = "UPDATE_VERIFIED"
	URI                                 InterfaceName = "URI"
)

//Object.keys(e).forEach(key => {
//var k = e[key].replace(/(?<![A-Z])[A-Z]/g, `_$&`).replace(/\(.*/, '').toLocaleUpperCase()
//var m = g[k] ? k+'_1' : k;
//var callable = /\(.+\)/.test(e[key])
//g[m] = {Value: key, Description: e[key], Callable: !callable}
// })

var INTERFACE_IDENTIFIERS = map[InterfaceName]InterfaceData{
	ADD_VERIFIED:                        {Value: "47089f62", Description: "addVerified(address,bytes32)", Callable: false},
	ALLOWANCE:                           {Value: "dd62ed3e", Description: "allowance(address,address)", Callable: false},
	APPROVE:                             {Value: "095ea7b3", Description: "approve(address,uint256)", Callable: false},
	APPROVE_AND_CALL:                    {Value: "cae9ca51", Description: "approveAndCall(address,uint256,bytes)", Callable: false},
	AUTHORIZE_OPERATOR:                  {Value: "959b8c3f", Description: "authorizeOperator(address)", Callable: false},
	BALANCE_OF:                          {Value: "70a08231", Description: "balanceOf(address)", Callable: false},
	BALANCE_OF_1:                        {Value: "00fdd58e", Description: "balanceOf(address,uint256)", Callable: false},
	BALANCE_OF_BATCH:                    {Value: "4e1273f4", Description: "balanceOfBatch(address[],uint256[])", Callable: false},
	BURN:                                {Value: "42966c68", Description: "burn(uint256)", Callable: false},
	BURN_1:                              {Value: "fe9d9303", Description: "burn(uint256,bytes)", Callable: false},
	BURN_FROM:                           {Value: "79cc6790", Description: "burnFrom(address,uint256)", Callable: false},
	CANCEL_AND_REISSUE:                  {Value: "79f64720", Description: "cancelAndReissue(address,address)", Callable: false},
	CAN_IMPLEMENT_INTERFACE_FOR_ADDRESS: {Value: "249cb3fa", Description: "canImplementInterfaceForAddress(bytes32,address)", Callable: false},
	CAP:                                 {Value: "355274ea", Description: "cap()", Callable: true},
	DECIMALS:                            {Value: "313ce567", Description: "decimals()", Callable: true},
	DECREASE_ALLOWANCE:                  {Value: "a457c2d7", Description: "decreaseAllowance(address,uint256)", Callable: false},
	DECREASE_ALLOWANCE_AND_CALL:         {Value: "d135ca1d", Description: "decreaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	DECREASE_SUPPLY:                     {Value: "869e0e60", Description: "decreaseSupply(uint256,address)", Callable: false},
	DEFAULT_OPERATORS:                   {Value: "06e48538", Description: "defaultOperators()", Callable: true},
	GET_APPROVED:                        {Value: "081812fc", Description: "getApproved(uint256)", Callable: false},
	GET_CURRENT_FOR:                     {Value: "cc397ed3", Description: "getCurrentFor(address)", Callable: false},
	GRANULARITY:                         {Value: "556f0dc7", Description: "granularity()", Callable: true},
	HAS_HASH:                            {Value: "f3221c7f", Description: "hasHash(address,bytes32)", Callable: false},
	HOLDER_AT:                           {Value: "197bc336", Description: "holderAt(uint256)", Callable: false},
	HOLDER_COUNT:                        {Value: "1aab9a9f", Description: "holderCount()", Callable: true},
	INCREASE_ALLOWANCE:                  {Value: "39509351", Description: "increaseAllowance(address,uint256)", Callable: false},
	INCREASE_ALLOWANCE_AND_CALL:         {Value: "5fd42775", Description: "increaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	INCREASE_SUPPLY:                     {Value: "124fc7e0", Description: "increaseSupply(uint256,address)", Callable: false},
	IS_APPROVED_FOR_ALL:                 {Value: "e985e9c5", Description: "isApprovedForAll(address,address)", Callable: false},
	IS_HOLDER:                           {Value: "d4d7b19a", Description: "isHolder(address)", Callable: false},
	IS_OPERATOR_FOR:                     {Value: "d95b6371", Description: "isOperatorFor(address,address)", Callable: false},
	IS_SUPERSEDED:                       {Value: "2da7293e", Description: "isSuperseded(address)", Callable: false},
	IS_VERIFIED:                         {Value: "b9209e33", Description: "isVerified(address)", Callable: false},
	MINT:                                {Value: "40c10f19", Description: "mint(address,uint256)", Callable: false},
	NAME:                                {Value: "06fdde03", Description: "name()", Callable: true},
	ON_ERC721_RECEIVED:                  {Value: "150b7a02", Description: "onERC721Received(address,address,uint256,bytes)", Callable: false},
	ON_ERC1155_BATCH_RECEIVED:           {Value: "bc197c81", Description: "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)", Callable: false},
	ON_ERC1155_RECEIVED:                 {Value: "f23a6e61", Description: "onERC1155Received(address,address,uint256,uint256,bytes)", Callable: false},
	OPERATOR_BURN:                       {Value: "fc673c4f", Description: "operatorBurn(address,uint256,bytes,bytes)", Callable: false},
	OPERATOR_SEND:                       {Value: "62ad1b83", Description: "operatorSend(address,address,uint256,bytes,bytes)", Callable: false},
	OWNER_OF:                            {Value: "6352211e", Description: "ownerOf(uint256)", Callable: false},
	REMOVE_VERIFIED:                     {Value: "4487b392", Description: "removeVerified(address)", Callable: false},
	REVOKE_OPERATOR:                     {Value: "fad8b32a", Description: "revokeOperator(address)", Callable: false},
	SAFE_BATCH_TRANSFER_FROM:            {Value: "2eb2c2d6", Description: "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)", Callable: false},
	SAFE_TRANSFER_FROM:                  {Value: "42842e0e", Description: "safeTransferFrom(address,address,uint256)", Callable: false},
	SAFE_TRANSFER_FROM_1:                {Value: "f242432a", Description: "safeTransferFrom(address,address,uint256,uint256,bytes)", Callable: false},
	SEND:                                {Value: "9bd9bbc6", Description: "send(address,uint256,bytes)", Callable: false},
	SET_APPROVAL_FOR_ALL:                {Value: "a22cb465", Description: "setApprovalForAll(address,bool)", Callable: false},
	SUPPORTS_INTERFACE:                  {Value: "01ffc9a7", Description: "supportsInterface(bytes4)", Callable: false},
	SYMBOL:                              {Value: "95d89b41", Description: "symbol()", Callable: true},
	TOKENS_RECEIVED:                     {Value: "0023de29", Description: "tokensReceived(address,address,address,uint256,bytes,bytes)", Callable: false},
	TOKENS_TO_SEND:                      {Value: "75ab9782", Description: "tokensToSend(address,address,address,uint256,bytes,bytes)", Callable: false},
	TOKEN_BY_INDEX:                      {Value: "4f6ccce7", Description: "tokenByIndex(uint256)", Callable: false},
	TOKEN_OF_OWNER_BY_INDEX:             {Value: "2f745c59", Description: "tokenOfOwnerByIndex(address,uint256)", Callable: false},
	TOKEN_URI:                           {Value: "c87b56dd", Description: "tokenURI(uint256)", Callable: false},
	TOTAL_SUPPLY:                        {Value: "18160ddd", Description: "totalSupply()", Callable: true},
	TRANSFER:                            {Value: "a9059cbb", Description: "transfer(address,uint256)", Callable: false},
	TRANSFER_1:                          {Value: "be45fd62", Description: "transfer(address,uint256,bytes)", Callable: false},
	TRANSFER_2:                          {Value: "f6368f8a", Description: "transfer(address,uint256,bytes,string)", Callable: false},
	TRANSFER_AND_CALL:                   {Value: "4000aea0", Description: "transferAndCall(address,uint256,bytes)", Callable: false},
	TRANSFER_FROM:                       {Value: "23b872dd", Description: "transferFrom(address,address,uint256)", Callable: false},
	TRANSFER_FROM_AND_CALL:              {Value: "c1d34b89", Description: "transferFromAndCall(address,address,uint256,bytes)", Callable: false},
	UPDATE_VERIFIED:                     {Value: "354b7b1d", Description: "updateVerified(address,bytes32)", Callable: false},
	URI:                                 {Value: "0e89341c", Description: "uri(uint256)", Callable: false},
}

type ErcName string
type ErcData []InterfaceName

const (
	ERC_20             ErcName = "ERC_20"
	ERC_20_BURNABLE    ErcName = "ERC_20_BURNABLE"
	ERC_20_CAPPED      ErcName = "ERC_20_CAPPED"
	ERC_20_DETAILED    ErcName = "ERC_20_DETAILED"
	ERC_20_MINTABLE    ErcName = "ERC_20_MINTABLE"
	ERC_20_PAUSABLE    ErcName = "ERC_20_PAUSABLE"
	ERC_165            ErcName = "ERC_165"
	ERC_721            ErcName = "ERC_721"
	ERC_721_RECEIVER   ErcName = "ERC_721_RECEIVER"
	ERC_721_METADATA   ErcName = "ERC_721_METADATA"
	ERC_721_ENUMERABLE ErcName = "ERC_721_ENUMERABLE"
	ERC_820            ErcName = "ERC_820"
	ERC_1155           ErcName = "1155"
	ERC_1155_RECEIVER  ErcName = "ERC_1155_RECEIVER"
	ERC_1155_METADATA  ErcName = "ERC_1155_METADATA"
	ERC_223            ErcName = "ERC_223"
	ERC_621            ErcName = "ERC_621"
	ERC_777            ErcName = "ERC_777"
	ERC_777_RECEIVER   ErcName = "ERC_777_RECEIVER"
	ERC_777_SENDER     ErcName = "ERC_777_SENDER"
	ERC_827            ErcName = "ERC_827"
	ERC_884            ErcName = "ERC_884"
)

var ERC_INTERFACE_IDENTIFIERS = map[ErcName]ErcData{
	ERC_20:             {ALLOWANCE, APPROVE, BALANCE_OF, TOTAL_SUPPLY, TRANSFER, TRANSFER_FROM},
	ERC_20_BURNABLE:    {BURN, BURN_FROM},
	ERC_20_CAPPED:      {CAP},
	ERC_20_DETAILED:    {DECIMALS, NAME, SYMBOL},
	ERC_20_MINTABLE:    {MINT},
	ERC_20_PAUSABLE:    {INCREASE_ALLOWANCE, APPROVE, DECREASE_ALLOWANCE, TRANSFER, TRANSFER_FROM},
	ERC_165:            {SUPPORTS_INTERFACE},
	ERC_721:            {APPROVE, BALANCE_OF, GET_APPROVED, IS_APPROVED_FOR_ALL, OWNER_OF, SAFE_TRANSFER_FROM, SAFE_TRANSFER_FROM_1, SET_APPROVAL_FOR_ALL, SUPPORTS_INTERFACE, TRANSFER_FROM},
	ERC_721_RECEIVER:   {ON_ERC721_RECEIVED},
	ERC_721_METADATA:   {NAME, SYMBOL, TOKEN_URI},
	ERC_721_ENUMERABLE: {TOKEN_BY_INDEX, TOKEN_OF_OWNER_BY_INDEX, TOTAL_SUPPLY},
	ERC_820:            {CAN_IMPLEMENT_INTERFACE_FOR_ADDRESS},
	ERC_1155:           {BALANCE_OF_1, BALANCE_OF_BATCH, IS_APPROVED_FOR_ALL, SAFE_BATCH_TRANSFER_FROM, SAFE_TRANSFER_FROM_1, SET_APPROVAL_FOR_ALL},
	ERC_1155_RECEIVER:  {ON_ERC1155_BATCH_RECEIVED, ON_ERC1155_RECEIVED},
	ERC_1155_METADATA:  {URI},
	ERC_223:            {BALANCE_OF, DECIMALS, NAME, SYMBOL, TOTAL_SUPPLY, TRANSFER, TRANSFER_1, TRANSFER_2},
	ERC_621:            {DECREASE_SUPPLY, INCREASE_SUPPLY},
	ERC_777:            {AUTHORIZE_OPERATOR, BALANCE_OF, BURN_1, DEFAULT_OPERATORS, GRANULARITY, IS_OPERATOR_FOR, NAME, OPERATOR_BURN, OPERATOR_SEND, REVOKE_OPERATOR, SEND, SYMBOL, TOTAL_SUPPLY},
	ERC_777_RECEIVER:   {TOKENS_RECEIVED},
	ERC_777_SENDER:     {TOKENS_TO_SEND},
	ERC_827:            {APPROVE_AND_CALL, DECREASE_ALLOWANCE_AND_CALL, INCREASE_ALLOWANCE_AND_CALL, TRANSFER_AND_CALL, TRANSFER_FROM_AND_CALL},
	ERC_884:            {ADD_VERIFIED, CANCEL_AND_REISSUE, GET_CURRENT_FOR, HAS_HASH, HOLDER_AT, HOLDER_COUNT, IS_HOLDER, IS_SUPERSEDED, IS_VERIFIED, REMOVE_VERIFIED, TRANSFER, TRANSFER_FROM, UPDATE_VERIFIED},
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

func (_Token *TokenCaller) GetInfo(byteCode string) ([]ErcName, []InterfaceName) {
	identifiers := map[InterfaceName]bool{}
	var interfaces []InterfaceName
	for k, v := range INTERFACE_IDENTIFIERS {
		if strings.Contains(byteCode, v.Value) {
			identifiers[k] = true
			interfaces = append(interfaces, k)
		}
	}
	types := map[ErcName]bool{}
Loop:
	for k, v := range ERC_INTERFACE_IDENTIFIERS {
		for _, ercIdentifier := range v {
			if _, ok := identifiers[ercIdentifier]; !ok {
				continue Loop
			}
		}
		types[k] = true
	}
	ercNames := make([]ErcName, len(types))
	for key := range types {
		ercNames = append(ercNames, key)
	}
	return ercNames, interfaces
}

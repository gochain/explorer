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

type InterfaceName int

type InterfaceData struct {
	Value       string
	Description string
	Callable    bool
}

const (
	AddVerified InterfaceName = iota
	Allowance
	Approve
	ApproveAndCall
	AuthorizeOperator
	BalanceOf
	BalanceOf1
	BalanceOfBatch
	Burn
	Burn1
	BurnFrom
	CancelAndReissue
	ImplementInterfaceForAddress
	Cap
	Decimals
	DecreaseAllowance
	DecreaseAllowanceAndCall
	DecreaseSupply
	DefaultOperators
	GetApproved
	GetCurrentFor
	Granularity
	HasHash
	HolderAt
	HolderCount
	IncreaseAllowance
	IncreaseAllowanceAndCall
	IncreaseSupply
	IsApprovedForAll
	IsHolder
	IsOperatorFor
	IsSuperseded
	IsVerified
	Mint
	Name
	OnErc721Received
	OnErc1155BatchReceived
	OnErc1155Received
	OperatorBurn
	OperatorSend
	OwnerOf
	RemoveVerified
	RevokeOperator
	SafeBatchTransferFrom
	SafeTransferFrom
	SafeTransferFrom1
	Send
	SetApprovalForAll
	SupportsInterface
	Symbol
	TokensReceived
	TokensToSend
	TokenByIndex
	TokenOfOwnerByIndex
	TokenUri
	TotalSupply
	Transfer
	Transfer1
	Transfer2
	TransferAndCall
	TransferFrom
	TransferFromAndCall
	UpdateVerified
	URI
)

//Object.keys(e).forEach(key => {
//var k = e[key].replace(/(?<![A-Z])[A-Z]/g, `_$&`).replace(/\(.*/, '').toLocaleUpperCase()
//var m = g[k] ? k+'_1' : k;
//var callable = /\(.+\)/.test(e[key])
//g[m] = {Value: key, Description: e[key], Callable: !callable}
// })

var InterfaceIdentifiers = map[InterfaceName]InterfaceData{
	AddVerified:                  {Value: "47089f62", Description: "addVerified(address,bytes32)", Callable: false},
	Allowance:                    {Value: "dd62ed3e", Description: "allowance(address,address)", Callable: false},
	Approve:                      {Value: "095ea7b3", Description: "approve(address,uint256)", Callable: false},
	ApproveAndCall:               {Value: "cae9ca51", Description: "approveAndCall(address,uint256,bytes)", Callable: false},
	AuthorizeOperator:            {Value: "959b8c3f", Description: "authorizeOperator(address)", Callable: false},
	BalanceOf:                    {Value: "70a08231", Description: "balanceOf(address)", Callable: false},
	BalanceOf1:                   {Value: "00fdd58e", Description: "balanceOf(address,uint256)", Callable: false},
	BalanceOfBatch:               {Value: "4e1273f4", Description: "balanceOfBatch(address[],uint256[])", Callable: false},
	Burn:                         {Value: "42966c68", Description: "burn(uint256)", Callable: false},
	Burn1:                        {Value: "fe9d9303", Description: "burn(uint256,bytes)", Callable: false},
	BurnFrom:                     {Value: "79cc6790", Description: "burnFrom(address,uint256)", Callable: false},
	CancelAndReissue:             {Value: "79f64720", Description: "cancelAndReissue(address,address)", Callable: false},
	ImplementInterfaceForAddress: {Value: "249cb3fa", Description: "canImplementInterfaceForAddress(bytes32,address)", Callable: false},
	Cap:                          {Value: "355274ea", Description: "cap()", Callable: true},
	Decimals:                     {Value: "313ce567", Description: "decimals()", Callable: true},
	DecreaseAllowance:            {Value: "a457c2d7", Description: "decreaseAllowance(address,uint256)", Callable: false},
	DecreaseAllowanceAndCall:     {Value: "d135ca1d", Description: "decreaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	DecreaseSupply:               {Value: "869e0e60", Description: "decreaseSupply(uint256,address)", Callable: false},
	DefaultOperators:             {Value: "06e48538", Description: "defaultOperators()", Callable: true},
	GetApproved:                  {Value: "081812fc", Description: "getApproved(uint256)", Callable: false},
	GetCurrentFor:                {Value: "cc397ed3", Description: "getCurrentFor(address)", Callable: false},
	Granularity:                  {Value: "556f0dc7", Description: "granularity()", Callable: true},
	HasHash:                      {Value: "f3221c7f", Description: "hasHash(address,bytes32)", Callable: false},
	HolderAt:                     {Value: "197bc336", Description: "holderAt(uint256)", Callable: false},
	HolderCount:                  {Value: "1aab9a9f", Description: "holderCount()", Callable: true},
	IncreaseAllowance:            {Value: "39509351", Description: "increaseAllowance(address,uint256)", Callable: false},
	IncreaseAllowanceAndCall:     {Value: "5fd42775", Description: "increaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	IncreaseSupply:               {Value: "124fc7e0", Description: "increaseSupply(uint256,address)", Callable: false},
	IsApprovedForAll:             {Value: "e985e9c5", Description: "isApprovedForAll(address,address)", Callable: false},
	IsHolder:                     {Value: "d4d7b19a", Description: "isHolder(address)", Callable: false},
	IsOperatorFor:                {Value: "d95b6371", Description: "isOperatorFor(address,address)", Callable: false},
	IsSuperseded:                 {Value: "2da7293e", Description: "isSuperseded(address)", Callable: false},
	IsVerified:                   {Value: "b9209e33", Description: "isVerified(address)", Callable: false},
	Mint:                         {Value: "40c10f19", Description: "mint(address,uint256)", Callable: false},
	Name:                         {Value: "06fdde03", Description: "name()", Callable: true},
	OnErc721Received:             {Value: "150b7a02", Description: "onERC721Received(address,address,uint256,bytes)", Callable: false},
	OnErc1155BatchReceived:       {Value: "bc197c81", Description: "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)", Callable: false},
	OnErc1155Received:            {Value: "f23a6e61", Description: "onERC1155Received(address,address,uint256,uint256,bytes)", Callable: false},
	OperatorBurn:                 {Value: "fc673c4f", Description: "operatorBurn(address,uint256,bytes,bytes)", Callable: false},
	OperatorSend:                 {Value: "62ad1b83", Description: "operatorSend(address,address,uint256,bytes,bytes)", Callable: false},
	OwnerOf:                      {Value: "6352211e", Description: "ownerOf(uint256)", Callable: false},
	RemoveVerified:               {Value: "4487b392", Description: "removeVerified(address)", Callable: false},
	RevokeOperator:               {Value: "fad8b32a", Description: "revokeOperator(address)", Callable: false},
	SafeBatchTransferFrom:        {Value: "2eb2c2d6", Description: "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)", Callable: false},
	SafeTransferFrom:             {Value: "42842e0e", Description: "safeTransferFrom(address,address,uint256)", Callable: false},
	SafeTransferFrom1:            {Value: "f242432a", Description: "safeTransferFrom(address,address,uint256,uint256,bytes)", Callable: false},
	Send:                         {Value: "9bd9bbc6", Description: "send(address,uint256,bytes)", Callable: false},
	SetApprovalForAll:            {Value: "a22cb465", Description: "setApprovalForAll(address,bool)", Callable: false},
	SupportsInterface:            {Value: "01ffc9a7", Description: "supportsInterface(bytes4)", Callable: false},
	Symbol:                       {Value: "95d89b41", Description: "symbol()", Callable: true},
	TokensReceived:               {Value: "0023de29", Description: "tokensReceived(address,address,address,uint256,bytes,bytes)", Callable: false},
	TokensToSend:                 {Value: "75ab9782", Description: "tokensToSend(address,address,address,uint256,bytes,bytes)", Callable: false},
	TokenByIndex:                 {Value: "4f6ccce7", Description: "tokenByIndex(uint256)", Callable: false},
	TokenOfOwnerByIndex:          {Value: "2f745c59", Description: "tokenOfOwnerByIndex(address,uint256)", Callable: false},
	TokenUri:                     {Value: "c87b56dd", Description: "tokenURI(uint256)", Callable: false},
	TotalSupply:                  {Value: "18160ddd", Description: "totalSupply()", Callable: true},
	Transfer:                     {Value: "a9059cbb", Description: "transfer(address,uint256)", Callable: false},
	Transfer1:                    {Value: "be45fd62", Description: "transfer(address,uint256,bytes)", Callable: false},
	Transfer2:                    {Value: "f6368f8a", Description: "transfer(address,uint256,bytes,string)", Callable: false},
	TransferAndCall:              {Value: "4000aea0", Description: "transferAndCall(address,uint256,bytes)", Callable: false},
	TransferFrom:                 {Value: "23b872dd", Description: "transferFrom(address,address,uint256)", Callable: false},
	TransferFromAndCall:          {Value: "c1d34b89", Description: "transferFromAndCall(address,address,uint256,bytes)", Callable: false},
	UpdateVerified:               {Value: "354b7b1d", Description: "updateVerified(address,bytes32)", Callable: false},
	URI:                          {Value: "0e89341c", Description: "uri(uint256)", Callable: false},
}

type ErcName int
type ErcData []InterfaceName

const (
	Erc20 ErcName = iota
	Erc20Burnable
	Erc20Capped
	Erc20Detailed
	Erc20Mintable
	Erc20Pausable
	Erc165
	Erc721
	Erc721Receiver
	Erc721Metadata
	Erc721Enumerable
	Erc820
	Erc1155
	Erc1155Receiver
	Erc1155Metadata
	Erc223
	Erc621
	Erc777
	Erc777Receiver
	Erc777Sender
	Erc827
	Erc884
)

var ErcInterfaceIdentifiers = map[ErcName]ErcData{
	Erc20:            {Allowance, Approve, BalanceOf, TotalSupply, Transfer, TransferFrom},
	Erc20Burnable:    {Burn, BurnFrom},
	Erc20Capped:      {Cap},
	Erc20Detailed:    {Decimals, Name, Symbol},
	Erc20Mintable:    {Mint},
	Erc20Pausable:    {IncreaseAllowance, Approve, DecreaseAllowance, Transfer, TransferFrom},
	Erc165:           {SupportsInterface},
	Erc721:           {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom},
	Erc721Receiver:   {OnErc721Received},
	Erc721Metadata:   {Name, Symbol, TokenUri},
	Erc721Enumerable: {TokenByIndex, TokenOfOwnerByIndex, TotalSupply},
	Erc820:           {ImplementInterfaceForAddress},
	Erc1155:          {BalanceOf1, BalanceOfBatch, IsApprovedForAll, SafeBatchTransferFrom, SafeTransferFrom1, SetApprovalForAll},
	Erc1155Receiver:  {OnErc1155BatchReceived, OnErc1155Received},
	Erc1155Metadata:  {URI},
	Erc223:           {BalanceOf, Decimals, Name, Symbol, TotalSupply, Transfer, Transfer1, Transfer2},
	Erc621:           {DecreaseSupply, IncreaseSupply},
	Erc777:           {AuthorizeOperator, BalanceOf, Burn1, DefaultOperators, Granularity, IsOperatorFor, Name, OperatorBurn, OperatorSend, RevokeOperator, Send, Symbol, TotalSupply},
	Erc777Receiver:   {TokensReceived},
	Erc777Sender:     {TokensToSend},
	Erc827:           {ApproveAndCall, DecreaseAllowanceAndCall, IncreaseAllowanceAndCall, TransferAndCall, TransferFromAndCall},
	Erc884:           {AddVerified, CancelAndReissue, GetCurrentFor, HasHash, HolderAt, HolderCount, IsHolder, IsSuperseded, IsVerified, RemoveVerified, Transfer, TransferFrom, UpdateVerified},
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
	for k, v := range InterfaceIdentifiers {
		if strings.Contains(byteCode, v.Value) {
			identifiers[k] = true
			interfaces = append(interfaces, k)
		}
	}
	types := map[ErcName]bool{}
Loop:
	for k, v := range ErcInterfaceIdentifiers {
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

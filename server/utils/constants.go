package utils

import "fmt"

type AbiArgument struct {
	Name    string `json:"name" bson:"name"`
	Indexed bool   `json:"indexes" bson:"indexes"`
	Type    string `json:"type" bson:"type"`
}

type AbiItem struct {
	Anonymous       bool          `json:"anonymous" bson:"anonymous"`
	Constant        bool          `json:"constant" bson:"constant"`
	Inputs          []AbiArgument `json:"inputs" bson:"inputs"`
	Name            string        `json:"name" bson:"name"`
	Outputs         []AbiArgument `json:"outputs" bson:"outputs"`
	Payable         bool          `json:"payable" bson:"payable"`
	StateMutability string        `json:"stateMutability" bson:"stateMutability"`
	Type            string        `json:"type" bson:"type"`
}

// EVMFunction is an internal enum only. Use the String() form externally.
type EVMFunction int

func (f EVMFunction) String() string {
	if f < 0 || int(f) >= len(evmFunctionNames) {
		return fmt.Sprintf("Unrecognized function: %d", f)
	}
	return evmFunctionNames[f]
}

const (
	AddPauser EVMFunction = iota
	AddVerified
	Allowance
	Approve
	ApproveAndCall
	AuthorizeOperator
	BalanceOf
	BalanceOfID
	BalanceOfBatch
	Burn
	BurnData
	BurnFrom
	CancelAndReissue
	CanImplementInterfaceForAddress
	Cap
	ChangeOwner
	CIDByHash
	Decimals
	DecreaseAllowance
	DecreaseAllowanceAndCall
	DecreaseSupply
	DefaultOperators
	Deployed
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
	IsPauser
	IsOperatorFor
	IsOwner
	IsSuperseded
	IsVerified
	Mint
	Name
	NewWallet
	OnErc721Received
	OnErc1155BatchReceived
	OnErc1155Received
	OperatorBurn
	OperatorSend
	Owner
	OwnerOf
	Rate
	RemoveVerified
	RenounceOwnership
	RenouncePauser
	RevokeOperator
	SafeBatchTransferFrom
	SafeTransferFrom
	SafeTransferFromData
	SafeTransferFromValueData
	Send
	SetApprovalForAll
	SetRate
	SupportsInterface
	Symbol
	Target
	TokensReceived
	TokensToSend
	TokenByIndex
	TokenFallback
	TokenOfOwnerByIndex
	TokenUri
	TotalSupply
	Transfer
	TransferData
	TransferDataFallback
	TransferAndCall
	TransferFrom
	TransferFromAndCall
	TransferOwnership
	UpdateVerified
	URI
	Pause
	Paused
	Pin
	Unpause
	Upgrade
	MintWithTokenURI
	Resume
	Wallet
	Withdraw
)

var evmFunctionNames = []string{
	AddPauser:                       "AddPauser",
	AddVerified:                     "AddVerified",
	Allowance:                       "Allowance",
	Approve:                         "Approve",
	ApproveAndCall:                  "ApproveAndCall",
	AuthorizeOperator:               "AuthorizeOperator",
	BalanceOf:                       "BalanceOf",
	BalanceOfID:                     "BalanceOfID",
	BalanceOfBatch:                  "BalanceOfBatch",
	Burn:                            "Burn",
	BurnData:                        "BurnData",
	BurnFrom:                        "BurnFrom",
	CancelAndReissue:                "CancelAndReissue",
	CanImplementInterfaceForAddress: "CanImplementInterfaceForAddress",
	Cap:                             "Cap",
	ChangeOwner:                     "ChangeOwner",
	CIDByHash:                       "CIDByHash",
	Decimals:                        "Decimals",
	DecreaseAllowance:               "DecreaseAllowance",
	DecreaseAllowanceAndCall:        "DecreaseAllowanceAndCall",
	DecreaseSupply:                  "DecreaseSupply",
	DefaultOperators:                "DefaultOperators",
	Deployed:                        "Deployed",
	GetApproved:                     "GetApproved",
	GetCurrentFor:                   "GetCurrentFor",
	Granularity:                     "Granularity",
	HasHash:                         "HasHash",
	HolderAt:                        "HolderAt",
	HolderCount:                     "HolderCount",
	IncreaseAllowance:               "IncreaseAllowance",
	IncreaseAllowanceAndCall:        "IncreaseAllowanceAndCall",
	IncreaseSupply:                  "IncreaseSupply",
	IsApprovedForAll:                "IsApprovedForAll",
	IsHolder:                        "IsHolder",
	IsPauser:                        "IsPauser",
	IsOperatorFor:                   "IsOperatorFor",
	IsOwner:                         "IsOwner",
	IsSuperseded:                    "IsSuperseded",
	IsVerified:                      "IsVerified",
	Mint:                            "Mint",
	Name:                            "Name",
	NewWallet:                       "NewWallet",
	OnErc721Received:                "OnErc721Received",
	OnErc1155BatchReceived:          "OnErc1155BatchReceived",
	OnErc1155Received:               "OnErc1155Received",
	OperatorBurn:                    "OperatorBurn",
	OperatorSend:                    "OperatorSend",
	Owner:                           "Owner",
	OwnerOf:                         "OwnerOf",
	Pin:                             "Pin",
	Rate:                            "Rate",
	RemoveVerified:                  "RemoveVerified",
	RenounceOwnership:               "RenounceOwnership",
	RenouncePauser:                  "RenouncePauser",
	RevokeOperator:                  "RevokeOperator",
	SafeBatchTransferFrom:           "SafeBatchTransferFrom",
	SafeTransferFrom:                "SafeTransferFrom",
	SafeTransferFromData:            "SafeTransferFromData",
	SafeTransferFromValueData:       "SafeTransferFromValueData",
	Send:                            "Send",
	SetRate:                         "SetRate",
	SetApprovalForAll:               "SetApprovalForAll",
	SupportsInterface:               "SupportsInterface",
	Symbol:                          "Symbol",
	Target:                          "Target",
	TokensReceived:                  "TokensReceived",
	TokensToSend:                    "TokensToSend",
	TokenByIndex:                    "TokenByIndex",
	TokenFallback:                   "TokenFallback",
	TokenOfOwnerByIndex:             "TokenOfOwnerByIndex",
	TokenUri:                        "TokenUri",
	TotalSupply:                     "TotalSupply",
	Transfer:                        "Transfer",
	TransferData:                    "TransferData",
	TransferDataFallback:            "TransferDataFallback",
	TransferAndCall:                 "TransferAndCall",
	TransferFrom:                    "TransferFrom",
	TransferFromAndCall:             "TransferFromAndCall",
	TransferOwnership:               "TransferOwnership",
	UpdateVerified:                  "UpdateVerified",
	URI:                             "URI",
	Pause:                           "Pause",
	Paused:                          "Paused",
	Unpause:                         "Unpause",
	Upgrade:                         "Upgrade",
	MintWithTokenURI:                "MintWithTokenURI",
	Resume:                          "Resume",
	Wallet:                          "Wallet",
	Withdraw:                        "Withdraw",
}

var (
	EVMFunctionsByName map[string]EVMFunction
	EVMFunctionsByID   map[string]EVMFunction
)

func init() {
	EVMFunctionsByName = make(map[string]EVMFunction, len(evmFunctionNames))
	for i, nm := range evmFunctionNames {
		EVMFunctionsByName[nm] = EVMFunction(i)
	}
	EVMFunctionsByID = make(map[string]EVMFunction, len(EVMFunctions))
	for i, data := range EVMFunctions {
		EVMFunctionsByID[data.ID] = i
	}
}

type EVMFunctionData struct {
	ID        string
	Signature string
	Callable  bool
}

var EVMFunctions = map[EVMFunction]EVMFunctionData{
	AddPauser:                       {ID: "82dc1ec4", Signature: "addPauser(address)", Callable: false},
	AddVerified:                     {ID: "47089f62", Signature: "addVerified(address,bytes32)", Callable: false},
	Allowance:                       {ID: "dd62ed3e", Signature: "allowance(address,address)", Callable: false},
	Approve:                         {ID: "095ea7b3", Signature: "approve(address,uint256)", Callable: false},
	ApproveAndCall:                  {ID: "cae9ca51", Signature: "approveAndCall(address,uint256,bytes)", Callable: false},
	AuthorizeOperator:               {ID: "959b8c3f", Signature: "authorizeOperator(address)", Callable: false},
	BalanceOf:                       {ID: "70a08231", Signature: "balanceOf(address)", Callable: false},
	BalanceOfID:                     {ID: "00fdd58e", Signature: "balanceOf(address,uint256)", Callable: false},
	BalanceOfBatch:                  {ID: "4e1273f4", Signature: "balanceOfBatch(address[],uint256[])", Callable: false},
	Burn:                            {ID: "42966c68", Signature: "burn(uint256)", Callable: false},
	BurnData:                        {ID: "fe9d9303", Signature: "burn(uint256,bytes)", Callable: false},
	BurnFrom:                        {ID: "79cc6790", Signature: "burnFrom(address,uint256)", Callable: false},
	CancelAndReissue:                {ID: "79f64720", Signature: "cancelAndReissue(address,address)", Callable: false},
	CanImplementInterfaceForAddress: {ID: "249cb3fa", Signature: "canImplementInterfaceForAddress(bytes32,address)", Callable: false},
	Cap:                             {ID: "355274ea", Signature: "cap()", Callable: true},
	ChangeOwner:                     {ID: "a6f9dae1", Signature: "changeOwner(address)", Callable: false},
	CIDByHash:                       {ID: "e16cf225", Signature: "cidByHash(bytes32)", Callable: true},
	Decimals:                        {ID: "313ce567", Signature: "decimals()", Callable: true},
	DecreaseAllowance:               {ID: "a457c2d7", Signature: "decreaseAllowance(address,uint256)", Callable: false},
	DecreaseAllowanceAndCall:        {ID: "d135ca1d", Signature: "decreaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	DecreaseSupply:                  {ID: "869e0e60", Signature: "decreaseSupply(uint256,address)", Callable: false},
	DefaultOperators:                {ID: "06e48538", Signature: "defaultOperators()", Callable: true},
	Deployed:                        {ID: "f905c15a", Signature: "deployed()", Callable: true},
	GetApproved:                     {ID: "081812fc", Signature: "getApproved(uint256)", Callable: false},
	GetCurrentFor:                   {ID: "cc397ed3", Signature: "getCurrentFor(address)", Callable: false},
	Granularity:                     {ID: "556f0dc7", Signature: "granularity()", Callable: true},
	HasHash:                         {ID: "f3221c7f", Signature: "hasHash(address,bytes32)", Callable: false},
	HolderAt:                        {ID: "197bc336", Signature: "holderAt(uint256)", Callable: false},
	HolderCount:                     {ID: "1aab9a9f", Signature: "holderCount()", Callable: true},
	IncreaseAllowance:               {ID: "39509351", Signature: "increaseAllowance(address,uint256)", Callable: false},
	IncreaseAllowanceAndCall:        {ID: "5fd42775", Signature: "increaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	IncreaseSupply:                  {ID: "124fc7e0", Signature: "increaseSupply(uint256,address)", Callable: false},
	IsApprovedForAll:                {ID: "e985e9c5", Signature: "isApprovedForAll(address,address)", Callable: false},
	IsHolder:                        {ID: "d4d7b19a", Signature: "isHolder(address)", Callable: false},
	IsPauser:                        {ID: "46fbf68e", Signature: "isPauser(address)", Callable: true},
	IsOperatorFor:                   {ID: "d95b6371", Signature: "isOperatorFor(address,address)", Callable: false},
	IsOwner:                         {ID: "8f32d59b", Signature: "isOwner()", Callable: false},
	IsSuperseded:                    {ID: "2da7293e", Signature: "isSuperseded(address)", Callable: false},
	IsVerified:                      {ID: "b9209e33", Signature: "isVerified(address)", Callable: false},
	Mint:                            {ID: "40c10f19", Signature: "mint(address,uint256)", Callable: false},
	Name:                            {ID: "06fdde03", Signature: "name()", Callable: true},
	NewWallet:                       {ID: "28c6fa6f", Signature: "newWallet(bytes)", Callable: false},
	OnErc721Received:                {ID: "150b7a02", Signature: "onERC721Received(address,address,uint256,bytes)", Callable: false},
	OnErc1155BatchReceived:          {ID: "bc197c81", Signature: "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)", Callable: false},
	OnErc1155Received:               {ID: "f23a6e61", Signature: "onERC1155Received(address,address,uint256,uint256,bytes)", Callable: false},
	OperatorBurn:                    {ID: "fc673c4f", Signature: "operatorBurn(address,uint256,bytes,bytes)", Callable: false},
	OperatorSend:                    {ID: "62ad1b83", Signature: "operatorSend(address,address,uint256,bytes,bytes)", Callable: false},
	Owner:                           {ID: "8da5cb5b", Signature: "owner()", Callable: false},
	OwnerOf:                         {ID: "6352211e", Signature: "ownerOf(uint256)", Callable: false},
	Pin:                             {ID: "7d1962f8", Signature: "pin(bytes)", Callable: false},
	Rate:                            {ID: "2c4e722e", Signature: "rate()", Callable: true},
	RemoveVerified:                  {ID: "4487b392", Signature: "removeVerified(address)", Callable: false},
	RenounceOwnership:               {ID: "715018a6", Signature: "renounceOwnership()", Callable: false},
	RenouncePauser:                  {ID: "6ef8d66d", Signature: "renouncePauser()", Callable: false},
	RevokeOperator:                  {ID: "fad8b32a", Signature: "revokeOperator(address)", Callable: false},
	SafeBatchTransferFrom:           {ID: "2eb2c2d6", Signature: "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)", Callable: false},
	SafeTransferFrom:                {ID: "42842e0e", Signature: "safeTransferFrom(address,address,uint256)", Callable: false},
	SafeTransferFromData:            {ID: "b88d4fde", Signature: "safeTransferFrom(address,address,uint256,bytes)", Callable: false},
	SafeTransferFromValueData:       {ID: "f242432a", Signature: "safeTransferFrom(address,address,uint256,uint256,bytes)", Callable: false},
	Send:                            {ID: "9bd9bbc6", Signature: "send(address,uint256,bytes)", Callable: false},
	SetApprovalForAll:               {ID: "a22cb465", Signature: "setApprovalForAll(address,bool)", Callable: false},
	SetRate:                         {ID: "34fcf437", Signature: "setRate(uint256)", Callable: false},
	SupportsInterface:               {ID: "01ffc9a7", Signature: "supportsInterface(bytes4)", Callable: false},
	Symbol:                          {ID: "95d89b41", Signature: "symbol()", Callable: true},
	Target:                          {ID: "d4b83992", Signature: "target()", Callable: true},
	TokensReceived:                  {ID: "0023de29", Signature: "tokensReceived(address,address,address,uint256,bytes,bytes)", Callable: false},
	TokensToSend:                    {ID: "75ab9782", Signature: "tokensToSend(address,address,address,uint256,bytes,bytes)", Callable: false},
	TokenByIndex:                    {ID: "4f6ccce7", Signature: "tokenByIndex(uint256)", Callable: false},
	TokenFallback:                   {ID: "c0ee0b8a", Signature: "tokenFallback(address,uint256,bytes)", Callable: false},
	TokenOfOwnerByIndex:             {ID: "2f745c59", Signature: "tokenOfOwnerByIndex(address,uint256)", Callable: false},
	TokenUri:                        {ID: "c87b56dd", Signature: "tokenURI(uint256)", Callable: false},
	TotalSupply:                     {ID: "18160ddd", Signature: "totalSupply()", Callable: true},
	Transfer:                        {ID: "a9059cbb", Signature: "transfer(address,uint256)", Callable: false},
	TransferData:                    {ID: "be45fd62", Signature: "transfer(address,uint256,bytes)", Callable: false},
	TransferDataFallback:            {ID: "f6368f8a", Signature: "transfer(address,uint256,bytes,string)", Callable: false},
	TransferAndCall:                 {ID: "4000aea0", Signature: "transferAndCall(address,uint256,bytes)", Callable: false},
	TransferFrom:                    {ID: "23b872dd", Signature: "transferFrom(address,address,uint256)", Callable: false},
	TransferFromAndCall:             {ID: "c1d34b89", Signature: "transferFromAndCall(address,address,uint256,bytes)", Callable: false},
	TransferOwnership:               {ID: "f2fde38b", Signature: "transferOwnership(address)", Callable: false},
	UpdateVerified:                  {ID: "354b7b1d", Signature: "updateVerified(address,bytes32)", Callable: false},
	URI:                             {ID: "0e89341c", Signature: "uri(uint256)", Callable: false},
	Pause:                           {ID: "8456cb59", Signature: "pause()", Callable: false},
	Paused:                          {ID: "5c975abb", Signature: "paused()", Callable: true},
	Unpause:                         {ID: "3f4ba83a", Signature: "unpause()", Callable: false},
	Upgrade:                         {ID: "0900f010", Signature: "upgrade(address)", Callable: false},
	MintWithTokenURI:                {ID: "50bb4e7f", Signature: "mintWithTokenURI(address,uint256,string)", Callable: false},
	Resume:                          {ID: "046f7da2", Signature: "resume()", Callable: false},
	Wallet:                          {ID: "521eb273", Signature: "wallet()", Callable: true},
	Withdraw:                        {ID: "51cff8d9", Signature: "withdraw(address)", Callable: false},
}

// EVMInterface is an internal enum only. Use the String() form externally.
type EVMInterface int

func (e EVMInterface) String() string {
	if e < 0 || int(e) >= len(evmInterfaceNames) {
		return fmt.Sprintf("Unrecognized interface: %d", e)
	}
	return evmInterfaceNames[e]
}

const (
	Go20 EVMInterface = iota
	Go20Burnable
	Go20Capped
	Go20Detailed
	Go20Mintable
	Go20Pausable
	Go165
	Go721
	Go721Burnable
	Go721Receiver
	Go721Metadata
	Go721Enumerable
	Go721Pausable
	Go721Mintable
	Go721MetadataMintable
	Go721Full
	Go820
	Go1155
	Go1155Receiver
	Go1155Metadata
	Go223
	Go223Receiver
	Go621
	Go777
	Go777Receiver
	Go777Sender
	Go827
	Go884
	Upgradeable
	Ownable
	PauserRole
	GoFS
)

var evmInterfaceNames = []string{
	Go20:                  "Go20",
	Go20Burnable:          "Go20Burnable",
	Go20Capped:            "Go20Capped",
	Go20Detailed:          "Go20Detailed",
	Go20Mintable:          "Go20Mintable",
	Go20Pausable:          "Go20Pausable",
	Go165:                 "Go165",
	Go721:                 "Go721",
	Go721Burnable:         "Go721Burnable",
	Go721Receiver:         "Go721Receiver",
	Go721Metadata:         "Go721Metadata",
	Go721Enumerable:       "Go721Enumerable",
	Go721Pausable:         "Go721Pausable",
	Go721Mintable:         "Go721Mintable",
	Go721MetadataMintable: "Go721MetadataMintable",
	Go721Full:             "Go721Full",
	Go820:                 "Go820",
	Go1155:                "Go1155",
	Go1155Receiver:        "Go1155Receiver",
	Go1155Metadata:        "Go1155Metadata",
	Go223:                 "Go223",
	Go223Receiver:         "Go223Receiver",
	Go621:                 "Go621",
	Go777:                 "Go777",
	Go777Receiver:         "Go777Receiver",
	Go777Sender:           "Go777Sender",
	Go827:                 "Go827",
	Go884:                 "Go884",
	Upgradeable:           "Upgradeable",
	Ownable:               "Ownable",
	PauserRole:            "PauserRole",
	GoFS:                  "GoFS",
}

var EVMInterfacesByName map[string]EVMInterface

func init() {
	EVMInterfacesByName = make(map[string]EVMInterface, len(evmInterfaceNames))
	for i, nm := range evmInterfaceNames {
		EVMInterfacesByName[nm] = EVMInterface(i)
	}
}

var (
	go20Functions  = []EVMFunction{Allowance, Approve, BalanceOf, TotalSupply, Transfer, TransferFrom}
	go721Functions = []EVMFunction{Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom,
		SafeTransferFromData, SetApprovalForAll, SupportsInterface, TransferFrom}
)

// EVMFunctionsByInterface maps each EVMInterface to its set of EVMFunctions.
var EVMFunctionsByInterface = [][]EVMFunction{
	Go20:         go20Functions,
	Go20Burnable: append(go20Functions, Burn, BurnFrom),
	Go20Capped:   append(go20Functions, Mint, Cap),
	Go20Detailed: append(go20Functions, Decimals, Name, Symbol),
	Go20Mintable: append(go20Functions, Mint),
	Go20Pausable: append(go20Functions, IncreaseAllowance, DecreaseAllowance, Transfer, TransferFrom, Pause, Paused,
		Unpause, AddPauser, IsPauser, RenouncePauser),

	Go165: {SupportsInterface},

	Go721:                 go721Functions,
	Go721Burnable:         append(go721Functions, Burn),
	Go721Receiver:         {OnErc721Received},
	Go721Metadata:         append(go721Functions, Name, Symbol, TokenUri),
	Go721Enumerable:       append(go721Functions, TokenByIndex, TokenOfOwnerByIndex, TotalSupply),
	Go721Pausable:         append(go721Functions, Pause, Paused, Unpause, AddPauser, IsPauser, RenouncePauser),
	Go721Mintable:         append(go721Functions, Mint),
	Go721MetadataMintable: append(go721Functions, Name, Symbol, TokenUri, MintWithTokenURI),
	Go721Full:             append(go721Functions, TokenByIndex, TokenOfOwnerByIndex, TotalSupply, Name, Symbol, TokenUri),

	Go820: {CanImplementInterfaceForAddress},

	Go1155:         {BalanceOfID, BalanceOfBatch, IsApprovedForAll, SafeBatchTransferFrom, SafeTransferFromValueData, SetApprovalForAll},
	Go1155Receiver: {OnErc1155BatchReceived, OnErc1155Received},
	Go1155Metadata: {URI},

	Go223:         {BalanceOf, Decimals, Name, Symbol, TotalSupply, Transfer, TransferData, TransferDataFallback},
	Go223Receiver: {TokenFallback},

	Go621: {DecreaseSupply, IncreaseSupply},

	Go777: {AuthorizeOperator, BalanceOf, BurnData, DefaultOperators, Granularity, IsOperatorFor, Name,
		OperatorBurn, OperatorSend, RevokeOperator, Send, Symbol, TotalSupply},
	Go777Receiver: {TokensReceived},
	Go777Sender:   {TokensToSend},

	Go827: {ApproveAndCall, DecreaseAllowanceAndCall, IncreaseAllowanceAndCall, TransferAndCall, TransferFromAndCall},

	Go884: {AddVerified, CancelAndReissue, GetCurrentFor, HasHash, HolderAt, HolderCount, IsHolder, IsSuperseded,
		IsVerified, RemoveVerified, Transfer, TransferFrom, UpdateVerified},

	Upgradeable: {Target, Upgrade, Pause, Paused, Resume, Owner},

	Ownable: {IsOwner, RenounceOwnership, TransferOwnership},

	PauserRole: {IsPauser, AddPauser, RenouncePauser},

	GoFS: {SetRate, NewWallet, Rate, Withdraw, Pin, Owner, ChangeOwner, Wallet, CIDByHash, Deployed},
}

package utils

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

type FunctionName string

type FunctionData struct {
	ID        string
	Signature string
	Callable  bool
}

const (
	AddPauser                       FunctionName = "AddPauser"
	AddVerified                     FunctionName = "AddVerified"
	Allowance                       FunctionName = "Allowance"
	Approve                         FunctionName = "Approve"
	ApproveAndCall                  FunctionName = "ApproveAndCall"
	AuthorizeOperator               FunctionName = "AuthorizeOperator"
	BalanceOf                       FunctionName = "BalanceOf"
	BalanceOfID                     FunctionName = "BalanceOfID"
	BalanceOfBatch                  FunctionName = "BalanceOfBatch"
	Burn                            FunctionName = "Burn"
	BurnData                        FunctionName = "BurnData"
	BurnFrom                        FunctionName = "BurnFrom"
	CancelAndReissue                FunctionName = "CancelAndReissue"
	CanImplementInterfaceForAddress FunctionName = "CanImplementInterfaceForAddress"
	Cap                             FunctionName = "Cap"
	Decimals                        FunctionName = "Decimals"
	DecreaseAllowance               FunctionName = "DecreaseAllowance"
	DecreaseAllowanceAndCall        FunctionName = "DecreaseAllowanceAndCall"
	DecreaseSupply                  FunctionName = "DecreaseSupply"
	DefaultOperators                FunctionName = "DefaultOperators"
	GetApproved                     FunctionName = "GetApproved"
	GetCurrentFor                   FunctionName = "GetCurrentFor"
	Granularity                     FunctionName = "Granularity"
	HasHash                         FunctionName = "HasHash"
	HolderAt                        FunctionName = "HolderAt"
	HolderCount                     FunctionName = "HolderCount"
	IncreaseAllowance               FunctionName = "IncreaseAllowance"
	IncreaseAllowanceAndCall        FunctionName = "IncreaseAllowanceAndCall"
	IncreaseSupply                  FunctionName = "IncreaseSupply"
	IsApprovedForAll                FunctionName = "IsApprovedForAll"
	IsHolder                        FunctionName = "IsHolder"
	IsPauser                        FunctionName = "IsPauser"
	IsOperatorFor                   FunctionName = "IsOperatorFor"
	IsSuperseded                    FunctionName = "IsSuperseded"
	IsVerified                      FunctionName = "IsVerified"
	Mint                            FunctionName = "Mint"
	Name                            FunctionName = "Name"
	OnErc721Received                FunctionName = "OnErc721Received"
	OnErc1155BatchReceived          FunctionName = "OnErc1155BatchReceived"
	OnErc1155Received               FunctionName = "OnErc1155Received"
	OperatorBurn                    FunctionName = "OperatorBurn"
	OperatorSend                    FunctionName = "OperatorSend"
	OwnerOf                         FunctionName = "OwnerOf"
	RemoveVerified                  FunctionName = "RemoveVerified"
	RenouncePauser                  FunctionName = "RenouncePauser"
	RevokeOperator                  FunctionName = "RevokeOperator"
	SafeBatchTransferFrom           FunctionName = "SafeBatchTransferFrom"
	SafeTransferFrom                FunctionName = "SafeTransferFrom"
	SafeTransferFromData            FunctionName = "SafeTransferFromData"
	SafeTransferFromValueData       FunctionName = "SafeTransferFromValueData"
	Send                            FunctionName = "Send"
	SetApprovalForAll               FunctionName = "SetApprovalForAll"
	SupportsInterface               FunctionName = "SupportsInterface"
	Symbol                          FunctionName = "Symbol"
	TokensReceived                  FunctionName = "TokensReceived"
	TokensToSend                    FunctionName = "TokensToSend"
	TokenByIndex                    FunctionName = "TokenByIndex"
	TokenFallback                   FunctionName = "TokenFallback"
	TokenOfOwnerByIndex             FunctionName = "TokenOfOwnerByIndex"
	TokenUri                        FunctionName = "TokenUri"
	TotalSupply                     FunctionName = "TotalSupply"
	Transfer                        FunctionName = "Transfer"
	TransferData                    FunctionName = "TransferData"
	TransferDataFallback            FunctionName = "TransferDataFallback"
	TransferAndCall                 FunctionName = "TransferAndCall"
	TransferFrom                    FunctionName = "TransferFrom"
	TransferFromAndCall             FunctionName = "TransferFromAndCall"
	UpdateVerified                  FunctionName = "UpdateVerified"
	URI                             FunctionName = "URI"
	Pause                           FunctionName = "Pause"
	Paused                          FunctionName = "Paused"
	Unpause                         FunctionName = "Unpause"
	MintWithTokenURI                FunctionName = "MintWithTokenURI"
)

var Functions = map[FunctionName]FunctionData{
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
	Decimals:                        {ID: "313ce567", Signature: "decimals()", Callable: true},
	DecreaseAllowance:               {ID: "a457c2d7", Signature: "decreaseAllowance(address,uint256)", Callable: false},
	DecreaseAllowanceAndCall:        {ID: "d135ca1d", Signature: "decreaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	DecreaseSupply:                  {ID: "869e0e60", Signature: "decreaseSupply(uint256,address)", Callable: false},
	DefaultOperators:                {ID: "06e48538", Signature: "defaultOperators()", Callable: true},
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
	IsSuperseded:                    {ID: "2da7293e", Signature: "isSuperseded(address)", Callable: false},
	IsVerified:                      {ID: "b9209e33", Signature: "isVerified(address)", Callable: false},
	Mint:                            {ID: "40c10f19", Signature: "mint(address,uint256)", Callable: false},
	Name:                            {ID: "06fdde03", Signature: "name()", Callable: true},
	OnErc721Received:                {ID: "150b7a02", Signature: "onERC721Received(address,address,uint256,bytes)", Callable: false},
	OnErc1155BatchReceived:          {ID: "bc197c81", Signature: "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)", Callable: false},
	OnErc1155Received:               {ID: "f23a6e61", Signature: "onERC1155Received(address,address,uint256,uint256,bytes)", Callable: false},
	OperatorBurn:                    {ID: "fc673c4f", Signature: "operatorBurn(address,uint256,bytes,bytes)", Callable: false},
	OperatorSend:                    {ID: "62ad1b83", Signature: "operatorSend(address,address,uint256,bytes,bytes)", Callable: false},
	OwnerOf:                         {ID: "6352211e", Signature: "ownerOf(uint256)", Callable: false},
	RemoveVerified:                  {ID: "4487b392", Signature: "removeVerified(address)", Callable: false},
	RenouncePauser:                  {ID: "6ef8d66d", Signature: "renouncePauser()", Callable: false},
	RevokeOperator:                  {ID: "fad8b32a", Signature: "revokeOperator(address)", Callable: false},
	SafeBatchTransferFrom:           {ID: "2eb2c2d6", Signature: "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)", Callable: false},
	SafeTransferFrom:                {ID: "42842e0e", Signature: "safeTransferFrom(address,address,uint256)", Callable: false},
	SafeTransferFromData:            {ID: "b88d4fde", Signature: "safeTransferFrom(address,address,uint256,bytes)", Callable: false},
	SafeTransferFromValueData:       {ID: "f242432a", Signature: "safeTransferFrom(address,address,uint256,uint256,bytes)", Callable: false},
	Send:                            {ID: "9bd9bbc6", Signature: "send(address,uint256,bytes)", Callable: false},
	SetApprovalForAll:               {ID: "a22cb465", Signature: "setApprovalForAll(address,bool)", Callable: false},
	SupportsInterface:               {ID: "01ffc9a7", Signature: "supportsInterface(bytes4)", Callable: false},
	Symbol:                          {ID: "95d89b41", Signature: "symbol()", Callable: true},
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
	UpdateVerified:                  {ID: "354b7b1d", Signature: "updateVerified(address,bytes32)", Callable: false},
	URI:                             {ID: "0e89341c", Signature: "uri(uint256)", Callable: false},
	Pause:                           {ID: "8456cb59", Signature: "pause()", Callable: false},
	Paused:                          {ID: "5c975abb", Signature: "paused()", Callable: true},
	Unpause:                         {ID: "3f4ba83a", Signature: "unpause()", Callable: false},
	MintWithTokenURI:                {ID: "50bb4e7f", Signature: "mintWithTokenURI(address,uint256,string)", Callable: false},
}

type ErcName string

const (
	Go20                  ErcName = "Go20"
	Go20Burnable          ErcName = "Go20Burnable"
	Go20Capped            ErcName = "Go20Capped"
	Go20Detailed          ErcName = "Go20Detailed"
	Go20Mintable          ErcName = "Go20Mintable"
	Go20Pausable          ErcName = "Go20Pausable"
	Go165                 ErcName = "Go165"
	Go721                 ErcName = "Go721"
	Go721Burnable         ErcName = "Go721Burnable"
	Go721Receiver         ErcName = "Go721Receiver"
	Go721Metadata         ErcName = "Go721Metadata"
	Go721Enumerable       ErcName = "Go721Enumerable"
	Go721Pausable         ErcName = "Go721Pausable"
	Go721Mintable         ErcName = "Go721Mintable"
	Go721MetadataMintable ErcName = "Go721MetadataMintable"
	Go721Full             ErcName = "Go721Full"
	Go820                 ErcName = "Go820"
	Go1155                ErcName = "Go1155"
	Go1155Receiver        ErcName = "Go1155Receiver"
	Go1155Metadata        ErcName = "Go1155Metadata"
	Go223                 ErcName = "Go223"
	Go223Receiver         ErcName = "Go223Receiver"
	Go621                 ErcName = "Go621"
	Go777                 ErcName = "Go777"
	Go777Receiver         ErcName = "Go777Receiver"
	Go777Sender           ErcName = "Go777Sender"
	Go827                 ErcName = "Go827"
	Go884                 ErcName = "Go884"
)

var (
	go20Functions  = []FunctionName{Allowance, Approve, BalanceOf, TotalSupply, Transfer, TransferFrom}
	go721Functions = []FunctionName{Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom,
		SafeTransferFromData, SetApprovalForAll, SupportsInterface, TransferFrom}
)

var Interfaces = map[ErcName][]FunctionName{
	Go20:         go20Functions,
	Go20Burnable: append(go20Functions, Burn, BurnFrom),
	Go20Capped:   append(go20Functions, Mint, Cap),
	Go20Detailed: append(go20Functions, Decimals, Name, Symbol),
	Go20Mintable: append(go20Functions, Mint),
	Go20Pausable: append(go20Functions, IncreaseAllowance, Approve, DecreaseAllowance, Transfer, TransferFrom, Pause,
		Paused, Unpause, AddPauser, IsPauser, RenouncePauser),

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
}

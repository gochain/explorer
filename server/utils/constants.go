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
	Value       string
	Description string
	Callable    bool
}

const (
	AddPauser                       FunctionName = "AddPauser"
	AddVerified                     FunctionName = "AddVerified"
	Allowance                       FunctionName = "Allowance"
	Approve                         FunctionName = "Approve"
	ApproveAndCall                  FunctionName = "ApproveAndCall"
	AuthorizeOperator               FunctionName = "AuthorizeOperator"
	BalanceOf                       FunctionName = "BalanceOf"
	BalanceOf1                      FunctionName = "BalanceOf1"
	BalanceOfBatch                  FunctionName = "BalanceOfBatch"
	Burn                            FunctionName = "Burn"
	Burn1                           FunctionName = "Burn1"
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
	SafeTransferFrom1               FunctionName = "SafeTransferFrom1"
	SafeTransferFrom2               FunctionName = "SafeTransferFrom2"
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
	Transfer1                       FunctionName = "Transfer1"
	Transfer2                       FunctionName = "Transfer2"
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

var InterfaceIdentifiers = map[FunctionName]FunctionData{
	AddPauser:                       {Value: "82dc1ec4", Description: "addPauser(address)", Callable: false},
	AddVerified:                     {Value: "47089f62", Description: "addVerified(address,bytes32)", Callable: false},
	Allowance:                       {Value: "dd62ed3e", Description: "allowance(address,address)", Callable: false},
	Approve:                         {Value: "095ea7b3", Description: "approve(address,uint256)", Callable: false},
	ApproveAndCall:                  {Value: "cae9ca51", Description: "approveAndCall(address,uint256,bytes)", Callable: false},
	AuthorizeOperator:               {Value: "959b8c3f", Description: "authorizeOperator(address)", Callable: false},
	BalanceOf:                       {Value: "70a08231", Description: "balanceOf(address)", Callable: false},
	BalanceOf1:                      {Value: "00fdd58e", Description: "balanceOf(address,uint256)", Callable: false},
	BalanceOfBatch:                  {Value: "4e1273f4", Description: "balanceOfBatch(address[],uint256[])", Callable: false},
	Burn:                            {Value: "42966c68", Description: "burn(uint256)", Callable: false},
	Burn1:                           {Value: "fe9d9303", Description: "burn(uint256,bytes)", Callable: false},
	BurnFrom:                        {Value: "79cc6790", Description: "burnFrom(address,uint256)", Callable: false},
	CancelAndReissue:                {Value: "79f64720", Description: "cancelAndReissue(address,address)", Callable: false},
	CanImplementInterfaceForAddress: {Value: "249cb3fa", Description: "canImplementInterfaceForAddress(bytes32,address)", Callable: false},
	Cap:                             {Value: "355274ea", Description: "cap()", Callable: true},
	Decimals:                        {Value: "313ce567", Description: "decimals()", Callable: true},
	DecreaseAllowance:               {Value: "a457c2d7", Description: "decreaseAllowance(address,uint256)", Callable: false},
	DecreaseAllowanceAndCall:        {Value: "d135ca1d", Description: "decreaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	DecreaseSupply:                  {Value: "869e0e60", Description: "decreaseSupply(uint256,address)", Callable: false},
	DefaultOperators:                {Value: "06e48538", Description: "defaultOperators()", Callable: true},
	GetApproved:                     {Value: "081812fc", Description: "getApproved(uint256)", Callable: false},
	GetCurrentFor:                   {Value: "cc397ed3", Description: "getCurrentFor(address)", Callable: false},
	Granularity:                     {Value: "556f0dc7", Description: "granularity()", Callable: true},
	HasHash:                         {Value: "f3221c7f", Description: "hasHash(address,bytes32)", Callable: false},
	HolderAt:                        {Value: "197bc336", Description: "holderAt(uint256)", Callable: false},
	HolderCount:                     {Value: "1aab9a9f", Description: "holderCount()", Callable: true},
	IncreaseAllowance:               {Value: "39509351", Description: "increaseAllowance(address,uint256)", Callable: false},
	IncreaseAllowanceAndCall:        {Value: "5fd42775", Description: "increaseAllowanceAndCall(address,uint256,bytes)", Callable: false},
	IncreaseSupply:                  {Value: "124fc7e0", Description: "increaseSupply(uint256,address)", Callable: false},
	IsApprovedForAll:                {Value: "e985e9c5", Description: "isApprovedForAll(address,address)", Callable: false},
	IsHolder:                        {Value: "d4d7b19a", Description: "isHolder(address)", Callable: false},
	IsPauser:                        {Value: "46fbf68e", Description: "isPauser(address)", Callable: true},
	IsOperatorFor:                   {Value: "d95b6371", Description: "isOperatorFor(address,address)", Callable: false},
	IsSuperseded:                    {Value: "2da7293e", Description: "isSuperseded(address)", Callable: false},
	IsVerified:                      {Value: "b9209e33", Description: "isVerified(address)", Callable: false},
	Mint:                            {Value: "40c10f19", Description: "mint(address,uint256)", Callable: false},
	Name:                            {Value: "06fdde03", Description: "name()", Callable: true},
	OnErc721Received:                {Value: "150b7a02", Description: "onERC721Received(address,address,uint256,bytes)", Callable: false},
	OnErc1155BatchReceived:          {Value: "bc197c81", Description: "onERC1155BatchReceived(address,address,uint256[],uint256[],bytes)", Callable: false},
	OnErc1155Received:               {Value: "f23a6e61", Description: "onERC1155Received(address,address,uint256,uint256,bytes)", Callable: false},
	OperatorBurn:                    {Value: "fc673c4f", Description: "operatorBurn(address,uint256,bytes,bytes)", Callable: false},
	OperatorSend:                    {Value: "62ad1b83", Description: "operatorSend(address,address,uint256,bytes,bytes)", Callable: false},
	OwnerOf:                         {Value: "6352211e", Description: "ownerOf(uint256)", Callable: false},
	RemoveVerified:                  {Value: "4487b392", Description: "removeVerified(address)", Callable: false},
	RenouncePauser:                  {Value: "6ef8d66d", Description: "renouncePauser()", Callable: false},
	RevokeOperator:                  {Value: "fad8b32a", Description: "revokeOperator(address)", Callable: false},
	SafeBatchTransferFrom:           {Value: "2eb2c2d6", Description: "safeBatchTransferFrom(address,address,uint256[],uint256[],bytes)", Callable: false},
	SafeTransferFrom:                {Value: "42842e0e", Description: "safeTransferFrom(address,address,uint256)", Callable: false},
	SafeTransferFrom1:               {Value: "f242432a", Description: "safeTransferFrom(address,address,uint256,uint256,bytes)", Callable: false},
	Send:                            {Value: "9bd9bbc6", Description: "send(address,uint256,bytes)", Callable: false},
	SetApprovalForAll:               {Value: "a22cb465", Description: "setApprovalForAll(address,bool)", Callable: false},
	SupportsInterface:               {Value: "01ffc9a7", Description: "supportsInterface(bytes4)", Callable: false},
	Symbol:                          {Value: "95d89b41", Description: "symbol()", Callable: true},
	TokensReceived:                  {Value: "0023de29", Description: "tokensReceived(address,address,address,uint256,bytes,bytes)", Callable: false},
	TokensToSend:                    {Value: "75ab9782", Description: "tokensToSend(address,address,address,uint256,bytes,bytes)", Callable: false},
	TokenByIndex:                    {Value: "4f6ccce7", Description: "tokenByIndex(uint256)", Callable: false},
	TokenFallback:                   {Value: "c0ee0b8a", Description: "tokenFallback(address,uint256,bytes)", Callable: false},
	TokenOfOwnerByIndex:             {Value: "2f745c59", Description: "tokenOfOwnerByIndex(address,uint256)", Callable: false},
	TokenUri:                        {Value: "c87b56dd", Description: "tokenURI(uint256)", Callable: false},
	TotalSupply:                     {Value: "18160ddd", Description: "totalSupply()", Callable: true},
	Transfer:                        {Value: "a9059cbb", Description: "transfer(address,uint256)", Callable: false},
	Transfer1:                       {Value: "be45fd62", Description: "transfer(address,uint256,bytes)", Callable: false},
	Transfer2:                       {Value: "f6368f8a", Description: "transfer(address,uint256,bytes,string)", Callable: false},
	TransferAndCall:                 {Value: "4000aea0", Description: "transferAndCall(address,uint256,bytes)", Callable: false},
	TransferFrom:                    {Value: "23b872dd", Description: "transferFrom(address,address,uint256)", Callable: false},
	TransferFromAndCall:             {Value: "c1d34b89", Description: "transferFromAndCall(address,address,uint256,bytes)", Callable: false},
	UpdateVerified:                  {Value: "354b7b1d", Description: "updateVerified(address,bytes32)", Callable: false},
	URI:                             {Value: "0e89341c", Description: "uri(uint256)", Callable: false},
	Pause:                           {Value: "8456cb59", Description: "pause()", Callable: false},
	Paused:                          {Value: "5c975abb", Description: "paused()", Callable: true},
	Unpause:                         {Value: "3f4ba83a", Description: "unpause()", Callable: false},
	MintWithTokenURI:                {Value: "50bb4e7f", Description: "mintWithTokenURI(address,uint256,string)", Callable: false},
}

type ErcName string
type ErcData []FunctionName

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

var ErcInterfaceIdentifiers = map[ErcName]ErcData{
	Go20:                  {Allowance, Approve, BalanceOf, TotalSupply, Transfer, TransferFrom},
	Go20Burnable:          {Burn, BurnFrom},
	Go20Capped:            {Mint, Cap},
	Go20Detailed:          {Decimals, Name, Symbol},
	Go20Mintable:          {Mint},
	Go20Pausable:          {IncreaseAllowance, Approve, DecreaseAllowance, Transfer, TransferFrom, Pause, Paused, Unpause, AddPauser, IsPauser, RenouncePauser},
	Go165:                 {SupportsInterface},
	Go721:                 {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom},
	Go721Burnable:         {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom, Burn},
	Go721Receiver:         {OnErc721Received},
	Go721Metadata:         {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom, Name, Symbol, TokenUri},
	Go721Enumerable:       {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom, TokenByIndex, TokenOfOwnerByIndex, TotalSupply},
	Go721Pausable:         {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom, Pause, Paused, Unpause, AddPauser, IsPauser, RenouncePauser},
	Go721Mintable:         {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom, Mint},
	Go721MetadataMintable: {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom, Name, Symbol, TokenUri, MintWithTokenURI},
	Go721Full:             {Approve, BalanceOf, GetApproved, IsApprovedForAll, OwnerOf, SafeTransferFrom, SafeTransferFrom1, SetApprovalForAll, SupportsInterface, TransferFrom, TokenByIndex, TokenOfOwnerByIndex, TotalSupply, Name, Symbol, TokenUri},
	Go820:                 {CanImplementInterfaceForAddress},
	Go1155:                {BalanceOf1, BalanceOfBatch, IsApprovedForAll, SafeBatchTransferFrom, SafeTransferFrom2, SetApprovalForAll},
	Go1155Receiver:        {OnErc1155BatchReceived, OnErc1155Received},
	Go1155Metadata:        {URI},
	Go223:                 {BalanceOf, Decimals, Name, Symbol, TotalSupply, Transfer, Transfer1, Transfer1},
	Go223Receiver:         {TokenFallback},
	Go621:                 {DecreaseSupply, IncreaseSupply},
	Go777:                 {AuthorizeOperator, BalanceOf, Burn1, DefaultOperators, Granularity, IsOperatorFor, Name, OperatorBurn, OperatorSend, RevokeOperator, Send, Symbol, TotalSupply},
	Go777Receiver:         {TokensReceived},
	Go777Sender:           {TokensToSend},
	Go827:                 {ApproveAndCall, DecreaseAllowanceAndCall, IncreaseAllowanceAndCall, TransferAndCall, TransferFromAndCall},
	Go884:                 {AddVerified, CancelAndReissue, GetCurrentFor, HasHash, HolderAt, HolderCount, IsHolder, IsSuperseded, IsVerified, RemoveVerified, Transfer, TransferFrom, UpdateVerified},
}

package backend

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/explorer/server/tokens"
	"github.com/gochain-io/explorer/server/utils"

	"github.com/gochain/gochain/v3/common"
	"github.com/gochain/gochain/v3/common/hexutil"
	"github.com/gochain/gochain/v3/consensus/clique"
	"github.com/gochain/gochain/v3/core/types"
	"github.com/gochain/gochain/v3/goclient"
	"github.com/gochain/gochain/v3/rpc"
	"go.uber.org/zap"
)

type Backend struct {
	mongo           *MongoBackend
	goRPC           *rpc.Client
	goClient        *goclient.Client
	tokenClient     *tokens.TokenClient
	dockerhubAPI    *DockerHubAPI
	reCaptchaSecret string
	lockedAccounts  []string
	signers         map[common.Address]models.Signer
	Lgr             *zap.Logger
}

func NewBackend(ctx context.Context, mongoUrl, rpcUrl, dbName string, lockedAccounts []string, signers map[common.Address]models.Signer, lgr *zap.Logger) (*Backend, error) {
	rpcClient, err := rpc.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc %q: %v", rpcUrl, err)
	}
	client := goclient.NewClient(rpcClient)
	mongoBackend, err := NewMongoClient(client, mongoUrl, dbName, lgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %v", err)
	}
	importer := new(Backend)
	importer.goRPC = rpcClient
	importer.goClient = client
	importer.mongo = mongoBackend
	importer.tokenClient, err = tokens.NewERC20Balance(ctx, client, lgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create erc20 balance client: %v", err)
	}
	importer.dockerhubAPI = new(DockerHubAPI)
	importer.lockedAccounts = lockedAccounts
	importer.signers = signers
	importer.Lgr = lgr
	return importer, nil
}

//METHODS USED IN API
func (self *Backend) PingDB() error {
	return self.mongo.PingDB()
}

// Balance returns the latest balance for the address.
func (self *Backend) Balance(ctx context.Context, address common.Address) (*big.Int, error) {
	var value *big.Int
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = self.goClient.BalanceAt(ctx, address, nil /* latest */)
		return err
	})
	return value, err
}

func (self *Backend) CodeAt(ctx context.Context, address string) ([]byte, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid hex address: %s", address)
	}
	var value []byte
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = self.goClient.CodeAt(ctx, common.HexToAddress(address), nil)
		return err
	})
	return value, err
}

func (self *Backend) TotalSupply(ctx context.Context) (*big.Int, error) {
	var value *big.Int
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		var result hexutil.Big
		err = self.goRPC.CallContext(ctx, &result, "eth_totalSupply", "latest")
		if err != nil {
			return err
		}
		value = result.ToInt()
		return nil
	})
	return value, err
}
func (self *Backend) CirculatingSupply(ctx context.Context) (*big.Int, error) {
	var value *big.Int
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		var result hexutil.Big
		err = self.goRPC.CallContext(ctx, &result, "eth_totalSupply", "latest")
		if err != nil {
			return err
		}
		total := result.ToInt()
		locked := new(big.Int)
		for _, l := range self.lockedAccounts {
			bal, err := self.Balance(ctx, common.HexToAddress(l))
			if err != nil {
				return err
			}
			locked = locked.Add(locked, bal)
		}
		value = new(big.Int).Sub(total, locked)
		return nil
	})
	return value, err
}
func (self *Backend) GetStats() (*models.Stats, error) {
	return self.mongo.getStats()
}

func (self *Backend) GetSignersStats() ([]models.SignersStats, error) {
	return self.mongo.getSignersStats()
}

func (self *Backend) GetSignersList() map[common.Address]models.Signer {
	return self.signers
}

func (self *Backend) GetRichlist(filter *models.PaginationFilter) ([]*models.Address, error) {
	return self.mongo.getRichlist(filter, self.lockedAccounts)

}
func (self *Backend) GetAddressByHash(ctx context.Context, hash string) (*models.Address, error) {
	if !common.IsHexAddress(hash) {
		return nil, errors.New("wrong address format")
	}
	addr := common.HexToAddress(hash)
	addressHash := addr.Hex()
	address, err := self.mongo.getAddressByHash(addressHash)
	if err != nil {
		return nil, err
	}
	balance, err := self.Balance(ctx, addr)
	if err != nil {
		return nil, err
	}
	if address == nil { //edge case if the balance for the address found but we haven't imported the address yet
		address = &models.Address{Address: addressHash, UpdatedAt: time.Now()}
		if err := self.mongo.UpdateActiveAddress(addressHash); err != nil {
			return nil, fmt.Errorf("failed to update active address: %s", err)
		}
	}
	address.BalanceWei = balance.String() //to make sure that we are showing most recent balance even if db is outdated
	address.BalanceString = new(big.Rat).SetFrac(balance, wei).FloatString(18)
	return address, nil

}
func (self *Backend) GetContracts(filter *models.ContractsFilter) ([]*models.Address, error) {
	return self.mongo.getContracts(filter)
}
func (self *Backend) GetTransactionByHash(ctx context.Context, hash string) (*models.Transaction, error) {
	return self.mongo.getTransactionByHash(ctx, hash)
}
func (self *Backend) GetTxByAddressAndNonce(ctx context.Context, addr string, nonce int64) (*models.Transaction, error) {
	return self.mongo.getTxByAddressAndNonce(ctx, addr, nonce)
}
func (self *Backend) GetTransactionList(address string, filter *models.TxsFilter) ([]*models.Transaction, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid hex address: %s", address)
	}
	return self.mongo.getTransactionList(common.HexToAddress(address).Hex(), filter)
}
func (self *Backend) GetTokenHoldersList(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenHolder, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	return self.mongo.getTokenHoldersList(common.HexToAddress(contractAddress).Hex(), filter)
}
func (self *Backend) GetOwnedTokensList(ownerAddress string, filter *models.PaginationFilter) ([]*models.TokenHolder, error) {
	if !common.IsHexAddress(ownerAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", ownerAddress)
	}
	return self.mongo.getOwnedTokensList(common.HexToAddress(ownerAddress).Hex(), filter)
}

// GetInternalTokenTransfers gets token transfer events emitted by an ERC20 or ERC721 contract.
func (self *Backend) GetInternalTokenTransfers(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenTransfer, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	return self.mongo.getInternalTokenTransfers(common.HexToAddress(contractAddress).Hex(), filter)
}

// GetHeldTokenTransfers gets token transfer events to or from this contract, emitted by any ERC20 or ERC721 contract.
func (self *Backend) GetHeldTokenTransfers(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenTransfer, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	return self.mongo.getHeldTokenTransfers(common.HexToAddress(contractAddress).Hex(), filter)
}
func (self *Backend) GetContract(contractAddress string) (*models.Contract, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	normalizedAddress := common.HexToAddress(contractAddress).Hex()
	contract, err := self.mongo.getContract(normalizedAddress)
	if contract != nil || err != nil {
		return contract, err
	}
	contractDataArray, err := self.CodeAt(context.Background(), normalizedAddress)
	if err != nil {
		return nil, err
	}
	contractData := string(contractDataArray[:])
	if contractData == "" {
		return nil, fmt.Errorf("invalid contract: %s", contractAddress)
	}
	byteCode := hex.EncodeToString(contractDataArray)
	err = self.ImportContract(normalizedAddress, byteCode)
	if err != nil {
		return nil, err
	}
	contract, err = self.mongo.getContract(normalizedAddress)
	return contract, err

}
func (self *Backend) GetLatestsBlocks(filter *models.PaginationFilter) ([]*models.LightBlock, error) {
	var lightBlocks []*models.LightBlock
	blocks, err := self.mongo.getLatestsBlocks(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest blocks: %v", err)
	}
	for _, block := range blocks {
		lightBlocks = append(lightBlocks, fillExtraLight(block))
	}
	return lightBlocks, nil
}
func (self *Backend) GetBlockTransactionsByNumber(blockNumber int64, filter *models.PaginationFilter) ([]*models.Transaction, error) {
	return self.mongo.getBlockTransactionsByNumber(blockNumber, filter)
}

func (self *Backend) GetBlockByNumber(ctx context.Context, number int64) (*models.Block, error) {
	block, err := self.mongo.getBlockByNumber(number)
	if err != nil {
		return nil, err
	}
	if block == nil || block.NonceBool == nil { //redownload block if it has no NonceBool filled, sort of lazy load
		self.Lgr.Info("Cannot get block from db or block is not up to date, importing it", zap.Int64("block", number))
		blockEth, err := self.goClient.BlockByNumber(ctx, big.NewInt(number))
		if err != nil {
			return nil, fmt.Errorf("failed to get block from rpc: %v", err)
		}
		block, err = self.ImportBlock(ctx, blockEth)
		if err != nil {
			return nil, fmt.Errorf("failed to import block: %v", err)
		}
	}
	return fillExtra(block), nil
}

func (self *Backend) GetBlockByHash(ctx context.Context, hash string) (*models.Block, error) {
	b, err := self.mongo.getBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	if b == nil { //redownload block if it has no NonceBool filled, sort of lazy load
		self.Lgr.Info("Cannot get block from db or block is not up to date, importing it", zap.String("block", hash))
		blockEth, err := self.goClient.BlockByHash(ctx, common.HexToHash(hash))
		if err != nil {
			return nil, fmt.Errorf("failed to get block from rpc: %v", err)
		}
		b, err = self.ImportBlock(ctx, blockEth)
		if err != nil {
			return nil, fmt.Errorf("failed to import block: %v", err)
		}
	}
	return fillExtra(b), nil
}

func (self *Backend) GetCompilerVersion() ([]string, error) {
	return self.dockerhubAPI.GetSolcImageTags()
}

func (self *Backend) VerifyContract(ctx context.Context, contractData *models.Contract) (*models.Contract, error) {
	contract, err := self.GetContract(contractData.Address)
	if err != nil {
		return nil, err
	}
	if contract == nil {
		return nil, errors.New("contract with given address not found")
	}
	if contract.Valid == true {
		return nil, errors.New("contract with given address is already verified")
	}
	compileData, err := CompileSolidityString(ctx, contractData.CompilerVersion, contractData.SourceCode, contractData.Optimization, contractData.EVMVersion)
	if err != nil {
		self.Lgr.Error("error while compilation", zap.Error(err))
		return nil, errors.New("error occurred while compiling source code")
	}
	// compiler gives map with keys starting with <stdin>:
	key := "<stdin>:" + contractData.ContractName
	if _, ok := compileData[key]; !ok {
		return nil, errors.New("invalid contract name")
	}
	if compileData[key].RuntimeCode == "" {
		return nil, errors.New("contract binary is empty")
	}
	if SolcBinEqual(compileData[key].RuntimeCode, contract.Bytecode) {
		contract.Valid = true
		contract.Optimization = contractData.Optimization
		contract.EVMVersion = contractData.EVMVersion
		contract.ContractName = contractData.ContractName
		contract.SourceCode = compileData[key].Info.Source
		contract.CompilerVersion = compileData[key].Info.CompilerVersion
		contract.Abi = compileData[key].Info.AbiDefinition
		contract.UpdatedAt = time.Now()
		if err := self.mongo.updateContract(contract); err != nil {
			return nil, err
		}
		return contract, nil
	}
	return nil, fmt.Errorf("the compiled result does not match the input creation bytecode located at %s", contractData.Address)
}

//METHODS USED IN GRABBER

func (self *Backend) UpdateStats() (*models.Stats, error) {
	return self.mongo.updateStats()
}

func (self *Backend) GetTokenBalance(contract, wallet string) (*tokens.TokenHolderDetails, error) {
	return self.tokenClient.GetTokenHolderDetails(contract, wallet)
}

func (self *Backend) GetTokenDetails(contractAddress string, byteCode string) (*tokens.TokenDetails, error) {
	return self.tokenClient.GetTokenDetails(contractAddress, byteCode)
}

func (self *Backend) GetTransferEvents(ctx context.Context, tokenDetails *tokens.TokenDetails, contractBlock int64, blockRangeLimit uint64) ([]*tokens.TransferEvent, error) {
	return self.tokenClient.GetTransferEvents(ctx, tokenDetails, contractBlock, blockRangeLimit)
}

func (self *Backend) CountTokenTransfers(address string) (int, error) {
	addr, err := self.mongo.getAddressByHash(address)
	if err != nil {
		return 0, fmt.Errorf("failed to get address: %v", err)
	}
	if addr != nil {
		return addr.NumberOfInternalTransactions, nil
	}
	return 0, nil

}

func (self *Backend) ImportBlock(ctx context.Context, block *types.Block) (*models.Block, error) {
	return self.mongo.importBlock(ctx, block)
}

// NeedReloadParent returns true if the parent block is missing or does not match the hash from this block number.
func (self *Backend) NeedReloadParent(blockNumber int64) (bool, error) {
	return self.mongo.needReloadParent(blockNumber)
}
func (self *Backend) TransactionsConsistent(blockNumber int64) (bool, error) {
	return self.mongo.transactionsConsistent(blockNumber)
}

//return false if a number of transactions in DB is different from a number of transactions in the blockchain
func (self *Backend) TransactionCountConsistent(ctx context.Context, blockNumber int64) (bool, error) {
	lgr := self.Lgr.With(zap.Int64("number", blockNumber))
	lgr.Debug("Checking transaction count for the block")
	block, err := self.GetBlockByNumber(ctx, blockNumber)
	if err != nil {
		return false, err
	}
	if block == nil {
		return false, errors.New("block not found")
	}
	txCount, err := self.goClient.TransactionCount(ctx, common.HexToHash(block.BlockHash))
	if err != nil {
		return false, fmt.Errorf("failed to get rpc transaction count: %v", err)
	}
	lgr.Debug("Got transaction count for the block", zap.Uint("rpc", txCount), zap.Int("db", block.TxCount))
	return txCount == uint(block.TxCount), nil
}

func (self *Backend) GetActiveAdresses(fromDate time.Time, onlyContracts bool) ([]*models.ActiveAddress, error) {
	addrs, err := self.mongo.getActiveAddresses(fromDate)
	if err != nil {
		return nil, err
	}
	var selectedAddresses []*models.ActiveAddress
	for _, address := range addrs {
		isContract, err := self.mongo.isContract(address.Address)
		if err != nil {
			return nil, fmt.Errorf("active address %s: %v", address.Address, err)
		}
		if onlyContracts == isContract {
			selectedAddresses = append(selectedAddresses, address)
		}
	}
	return selectedAddresses, nil
}
func (self *Backend) ImportAddress(address string, balance *big.Int, token *tokens.TokenDetails, contract bool, updatedAtBlock int64) (*models.Address, error) {
	return self.mongo.importAddress(address, balance, token, contract, updatedAtBlock)
}
func (self *Backend) ImportTokenHolder(contractAddress, tokenHolderAddress string, token *tokens.TokenHolderDetails, address *models.Address) (*models.TokenHolder, error) {
	return self.mongo.importTokenHolder(contractAddress, tokenHolderAddress, token, address)
}
func (self *Backend) ImportTransferEvent(ctx context.Context, contractAddress string, transferEvent *tokens.TransferEvent) (*models.TokenTransfer, error) {
	createdAt := time.Now()
	block, err := self.GetBlockByNumber(ctx, transferEvent.BlockNumber)
	if err != nil {
		return nil, err
	}
	if block != nil {
		createdAt = block.CreatedAt
	}
	return self.mongo.importTransferEvent(contractAddress, transferEvent, createdAt)
}
func (self *Backend) ImportContract(contractAddress string, byteCode string) error {
	return self.mongo.importContract(contractAddress, byteCode)
}

func (self *Backend) GetContractBlock(contractAddress string) (int64, error) {
	return self.mongo.getContractBlock(contractAddress)
}

func (self *Backend) BlockByNumber(ctx context.Context, blockNumber int64) (*types.Block, error) {
	var value *types.Block
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = self.goClient.BlockByNumber(ctx, big.NewInt(blockNumber))
		return err
	})
	return value, err
}
func (self *Backend) GetLatestBlockNumber(ctx context.Context) (int64, error) {
	var value *big.Int
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = self.goClient.LatestBlockNumber(ctx)
		return err
	})
	return value.Int64(), err
}

func fillExtra(block *models.Block) *models.Block {
	if block == nil {
		return block
	}
	extra := []byte(block.ExtraData)
	block.Extra.Auth = (block.NonceBool != nil && *block.NonceBool) //workaround for get old block by hash
	block.Extra.Vanity = strings.TrimRight(string(clique.ExtraVanity(extra)), "\u0000")
	block.Extra.HasVote = clique.ExtraHasVote(extra)
	block.Extra.Candidate = clique.ExtraCandidate(extra).String()
	block.Extra.IsVoterElection = clique.ExtraIsVoterElection(extra)
	return block
}
func fillExtraLight(block *models.LightBlock) *models.LightBlock {
	extra := []byte(block.ExtraData)
	block.Extra.Vanity = strings.TrimRight(string(clique.ExtraVanity(extra)), "\u0000")
	return block
}
func (self *Backend) DeleteBlockByHash(hash string) error {
	return self.mongo.deleteBlockByHash(hash)
}
func (self *Backend) DeleteBlockByNumber(bnum int64) error {
	return self.mongo.deleteBlockByNumber(bnum)
}

func (self *Backend) DeleteContract(contractAddress string) error {
	return self.mongo.deleteContract(contractAddress)
}

func (self *Backend) MigrateDB(ctx context.Context, lgr *zap.Logger) (int, error) {
	return self.mongo.migrate(ctx, lgr)
}

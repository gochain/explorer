package backend

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gochain/gochain/v4"
	"gopkg.in/mgo.v2/bson"

	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/explorer/server/tokens"
	"github.com/gochain-io/explorer/server/utils"

	"github.com/dgraph-io/ristretto"
	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/consensus/clique"
	"github.com/gochain/gochain/v4/core"
	"github.com/gochain/gochain/v4/core/types"
	"github.com/gochain/gochain/v4/goclient"
	"github.com/gochain/gochain/v4/params"
	"github.com/gochain/gochain/v4/rpc"
	"go.uber.org/zap"
)

const rpcClientTimeout = time.Minute

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
	Cache           *ristretto.Cache
	chainID         *big.Int
	Config          *params.ChainConfig
	alloc           atomic.Value // Lazy-loaded with first TotalSupply call.
	totalBurned     atomic.Value // *TotalBurned // Latest known. Eventually consistent.
}

func NewBackend(ctx context.Context, mongoUrl, rpcUrl, dbName string, lockedAccounts []string, signers map[common.Address]models.Signer, initialAllocation *big.Int,
	lgr *zap.Logger, cache *ristretto.Cache) (*Backend, error) {
	rpcClient, err := rpc.DialHTTPWithClient(rpcUrl, &http.Client{Timeout: rpcClientTimeout})
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc %q: %v", rpcUrl, err)
	}
	client := goclient.NewClient(rpcClient)
	mongoBackend, err := NewMongoClient(rpcClient, client, mongoUrl, dbName, lgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create mongo client: %v", err)
	}
	b := new(Backend)
	b.goRPC = rpcClient
	b.goClient = client
	b.mongo = mongoBackend
	b.tokenClient, err = tokens.NewERC20Balance(ctx, client, lgr)
	if err != nil {
		return nil, fmt.Errorf("failed to create erc20 balance client: %v", err)
	}
	b.dockerhubAPI = new(DockerHubAPI)
	b.lockedAccounts = lockedAccounts
	b.signers = signers
	b.alloc.Store(initialAllocation)
	b.Lgr = lgr

	if cache == nil {
		// should this be a no-op cache if it's not passed in?
		cache, err = ristretto.NewCache(&ristretto.Config{
			NumCounters: 1e6,   // number of keys to track frequency of (1M).
			MaxCost:     10000, // maximum cost of cache (1GB).
			BufferItems: 64,    // number of keys per Get buffer.
		})
		if err != nil {
			panic(err)
		}
	}
	b.chainID, err = b.goClient.ChainID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get chain ID: %v", err)
	}
	chainID := b.chainID.Uint64()
	switch chainID {
	case params.MainnetChainID:
		b.Config = params.MainnetChainConfig
	case params.TestnetChainID:
		b.Config = params.TestnetChainConfig
	}
	if b.Config == nil {
		b.Lgr.Info("Backend configured", zap.Uint64("chainID", chainID))
	} else {
		b.Lgr.Info("Backend configured", zap.String("config", b.Config.String()))
	}
	return b, nil
}

//METHODS USED IN API
func (b *Backend) PingDB() error {
	return b.mongo.PingDB()
}

// Balance returns the latest balance for the address.
func (b *Backend) Balance(ctx context.Context, address common.Address) (*big.Int, error) {
	var value *big.Int
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = b.goClient.BalanceAt(ctx, address, nil /* latest */)
		return err
	})
	return value, err
}

func (b *Backend) Nonce(ctx context.Context, address common.Address) (uint64, error) {
	var value uint64
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = b.goClient.PendingNonceAt(ctx, address)
		return err
	})
	return value, err
}

func (b *Backend) CodeAt(ctx context.Context, address string) ([]byte, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid hex address: %s", address)
	}
	var value []byte
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = b.goClient.CodeAt(ctx, common.HexToAddress(address), nil)
		return err
	})
	return value, err
}

// TotalSupply returns the total supply and the fees burned (already subtracted from total).
func (b *Backend) TotalSupply(ctx context.Context) (*big.Int, *big.Int, error) {
	var alloc *big.Int
	if l := b.alloc.Load(); l != nil {
		alloc = l.(*big.Int)
	}
	if alloc == nil {
		if err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
			var result *core.GenesisAlloc
			err = b.goRPC.CallContext(ctx, &result, "eth_genesisAlloc")
			if err != nil {
				return err
			}
			alloc = result.Total()
			b.alloc.Store(alloc)
			return nil
		}); err != nil {
			return nil, nil, err
		}
	}
	n, err := b.GetLatestBlockNumber(ctx)
	if err != nil {
		return nil, nil, err
	}
	rewards := new(big.Int).Mul(clique.BlockReward, n)
	total := new(big.Int).Add(alloc, rewards)
	totalFeesBurned := big.NewInt(0)
	if b.isDarvazaFork(n) {
		if totalFeesBurned, err = b.totalFeesBurned(ctx, n.Int64()); err != nil {
			return nil, nil, err
		} else if totalFeesBurned != nil {
			total = total.Sub(total, totalFeesBurned)
		}
	}
	return total, totalFeesBurned, nil
}

// totalFeesBurned may return a cached value.
func (b *Backend) totalFeesBurned(ctx context.Context, n int64) (*big.Int, error) {
	v := b.totalBurned.Load()
	if v != nil {
		if tb := v.(*TotalBurned); tb != nil {
			if n == tb.Number || time.Since(tb.CachedAt) < 5*time.Second {
				return tb.TotalFeesBurned, nil
			}
		}
	}

	tb, err := b.mongo.getLatestTotalFeesBurned()
	if err != nil {
		return nil, err
	}
	if tb == nil {
		return nil, nil
	}
	b.totalBurned.Store(tb)
	return tb.TotalFeesBurned, nil
}

func (b *Backend) CirculatingSupply(ctx context.Context) (*big.Int, error) {
	supplyStats, err := b.SupplyStats(ctx)
	if err != nil {
		return nil, err
	}
	return supplyStats.Circulating, nil
}

func (b *Backend) SupplyStats(ctx context.Context) (*models.SupplyStats, error) {
	total, feesBurned, err := b.TotalSupply(ctx)
	if err != nil {
		return nil, err
	}
	locked := new(big.Int)
	for _, l := range b.lockedAccounts {
		bal, err := b.Balance(ctx, common.HexToAddress(l))
		if err != nil {
			return nil, err
		}
		locked = locked.Add(locked, bal)
	}
	return &models.SupplyStats{
		Total: total, FeesBurned: feesBurned, Locked: locked,
		Circulating: new(big.Int).Sub(total, locked),
	}, err
}

func (b *Backend) isDarvazaFork(n *big.Int) bool {
	return b.Config != nil && b.Config.IsDarvaza(n)
}

func (b *Backend) GetStats() (*models.Stats, error) {
	return b.mongo.getStats()
}

func (b *Backend) GetSignersStats() ([]models.SignersStats, error) {
	return b.mongo.getSignersStats()
}

func (b *Backend) GetSignersList() map[common.Address]models.Signer {
	return b.signers
}

func (b *Backend) GetRichlist(filter *models.PaginationFilter) ([]*models.Address, error) {
	return b.mongo.getRichlist(filter, b.lockedAccounts)

}
func (b *Backend) GetAddressByHash(ctx context.Context, hash string) (*models.Address, error) {
	if !common.IsHexAddress(hash) {
		return nil, errors.New("wrong address format")
	}
	// check cache first
	a, found := b.Cache.Get(hash)
	if found {
		return a.(*models.Address), nil
	}
	addr := common.HexToAddress(hash)
	addressHash := addr.Hex()
	address, err := b.mongo.getAddressByHash(addressHash)
	if err != nil {
		return nil, err
	}
	balance, err := b.Balance(ctx, addr)
	if err != nil {
		return nil, err
	}
	if address == nil { //edge case if the balance for the address found but we haven't imported the address yet
		address = &models.Address{Address: addressHash, UpdatedAt: time.Now()}
		if err := b.mongo.UpdateActiveAddress(addressHash); err != nil {
			return nil, fmt.Errorf("failed to update active address: %s", err)
		}
	}
	transactionCounter, err := b.mongo.mongo.C("TransactionsByAddress").Find(bson.M{"address": address}).Count()
	if err != nil {
		return nil, fmt.Errorf("failed to get txs from TransactionsByAddress: %v", err)
	}

	address.NumberOfTransactions = transactionCounter
	address.BalanceWei = balance.String() //to make sure that we are showing most recent balance even if db is outdated
	address.BalanceString = new(big.Rat).SetFrac(balance, wei).FloatString(18)
	// todo: only store a subset of this data in the cache, just the metadata like token name, decimals, etc. Things that won't change.
	b.Cache.Set(hash, address, 5)
	return address, nil

}
func (b *Backend) GetContracts(filter *models.ContractsFilter) ([]*models.Address, error) {
	return b.mongo.getContracts(filter)
}
func (b *Backend) GetTransactionByHash(ctx context.Context, hash string) (*models.Transaction, error) {
	return b.mongo.getTransactionByHash(ctx, hash)
}
func (b *Backend) GetTxByAddressAndNonce(ctx context.Context, addr string, nonce int64) (*models.Transaction, error) {
	return b.mongo.getTxByAddressAndNonce(ctx, addr, nonce)
}
func (b *Backend) GetTransactionList(address string, filter *models.TxsFilter) ([]*models.Transaction, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid hex address: %s", address)
	}
	return b.mongo.getTransactionList(common.HexToAddress(address).Hex(), filter)
}
func (b *Backend) GetTokenHoldersList(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenHolder, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	return b.mongo.getTokenHoldersList(common.HexToAddress(contractAddress).Hex(), filter)
}
func (b *Backend) GetOwnedTokensList(ctx context.Context, ownerAddress string, filter *models.PaginationFilter) ([]*models.TokenHolder, error) {
	if !common.IsHexAddress(ownerAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", ownerAddress)
	}
	held, err := b.mongo.getOwnedTokensList(common.HexToAddress(ownerAddress).Hex(), filter)
	if err != nil {
		return nil, err
	}
	// loop through and add decimals
	for _, h := range held {
		addrInfo, err := b.GetAddressByHash(ctx, h.ContractAddress)
		if err != nil {
			return nil, fmt.Errorf("error finding decimals for %v: %v", h.ContractAddress, err)
		}
		h.TokenDecimals = addrInfo.Decimals
	}
	return held, nil
}

// GetInternalTokenTransfers gets token transfer events emitted by an ERC20 or ERC721 contract.
func (b *Backend) GetInternalTokenTransfers(contractAddress string, filter *models.InternalTxFilter) ([]*models.TokenTransfer, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	return b.mongo.getInternalTokenTransfers(common.HexToAddress(contractAddress).Hex(), filter)
}

// GetHeldTokenTransfers gets token transfer events to or from this contract, emitted by any ERC20 or ERC721 contract.
func (b *Backend) GetHeldTokenTransfers(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenTransfer, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	return b.mongo.getHeldTokenTransfers(common.HexToAddress(contractAddress).Hex(), filter)
}
func (b *Backend) GetContract(contractAddress string) (*models.Contract, error) {
	if !common.IsHexAddress(contractAddress) {
		return nil, fmt.Errorf("invalid hex address: %s", contractAddress)
	}
	normalizedAddress := common.HexToAddress(contractAddress).Hex()
	contract, err := b.mongo.getContract(normalizedAddress)
	if contract != nil || err != nil {
		return contract, err
	}
	contractDataArray, err := b.CodeAt(context.Background(), normalizedAddress)
	if err != nil {
		return nil, err
	}
	contractData := string(contractDataArray[:])
	if contractData == "" {
		return nil, nil
	}
	byteCode := hex.EncodeToString(contractDataArray)
	err = b.ImportContract(normalizedAddress, byteCode)
	if err != nil {
		return nil, err
	}
	contract, err = b.mongo.getContract(normalizedAddress)
	return contract, err

}
func (b *Backend) GetLatestsBlocks(filter *models.PaginationFilter) ([]*models.LightBlock, error) {
	var lightBlocks []*models.LightBlock
	blocks, err := b.mongo.getLatestsBlocks(filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest blocks: %v", err)
	}
	for _, block := range blocks {
		lightBlocks = append(lightBlocks, fillExtraLight(block))
	}
	return lightBlocks, nil
}
func (b *Backend) GetBlockTransactionsByNumber(blockNumber int64, filter *models.PaginationFilter) ([]*models.Transaction, error) {
	return b.mongo.getBlockTransactionsByNumber(blockNumber, filter)
}

func (b *Backend) GetBlockByNumber(ctx context.Context, number int64, asIs bool) (*models.Block, error) {
	block, err := b.mongo.getBlockByNumber(number)
	if err != nil {
		return nil, err
	}
	reload := false
	if !asIs {
		if block == nil {
			reload = true
			b.Lgr.Info("Block not found in DB, importing", zap.Int64("block", number))
		} else if block.NonceBool == nil {
			reload = true
			b.Lgr.Info("Block not up to date, reimporting", zap.Int64("block", number))
		}
		if reload {
			blockEth, err := b.goClient.BlockByNumber(ctx, big.NewInt(number))
			if err != nil {
				return nil, fmt.Errorf("failed to get block from rpc: %v", err)
			}
			block, err = b.ImportBlock(ctx, blockEth)
			if err != nil {
				return nil, fmt.Errorf("failed to import block: %v", err)
			}
		}
	}
	return fillExtra(block), nil
}

func (b *Backend) GetBlockByHash(ctx context.Context, hash string) (*models.Block, error) {
	block, err := b.mongo.getBlockByHash(hash)
	if err != nil {
		return nil, err
	}
	reload := false
	if block == nil {
		reload = true
		b.Lgr.Info("Block not found in DB, importing", zap.String("block", hash))
	} else if block.NonceBool == nil {
		reload = true
		b.Lgr.Info("Block not up to date, reimporting", zap.String("block", hash))
	}
	if reload {
		blockEth, err := b.goClient.BlockByHash(ctx, common.HexToHash(hash))
		if err != nil {
			if err == gochain.NotFound {
				return nil, nil
			}
			return nil, fmt.Errorf("failed to get block from rpc: %v", err)
		}
		block, err = b.ImportBlock(ctx, blockEth)
		if err != nil {
			return nil, fmt.Errorf("failed to import block: %v", err)
		}
	}
	return fillExtra(block), nil
}

func (b *Backend) GetCompilerVersion() ([]string, error) {
	return b.dockerhubAPI.GetSolcImageTags()
}

func (b *Backend) VerifyContract(ctx context.Context, contractData *models.Contract) (*models.Contract, error) {
	contract, err := b.GetContract(contractData.Address)
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
		b.Lgr.Error("error while compilation", zap.Error(err))
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
		if err := b.mongo.updateContract(contract); err != nil {
			return nil, err
		}
		return contract, nil
	}
	return nil, fmt.Errorf("the compiled result does not match the input creation bytecode located at %s", contractData.Address)
}

//METHODS USED IN GRABBER

func (b *Backend) UpdateStats() (*models.Stats, error) {
	return b.mongo.updateStats()
}

func (b *Backend) GetTokenBalance(contract, wallet string) (*tokens.TokenHolderDetails, error) {
	return b.tokenClient.GetTokenHolderDetails(contract, wallet)
}

func (b *Backend) GetTokenDetails(contractAddress string, byteCode string) (*tokens.TokenDetails, error) {
	return b.tokenClient.GetTokenDetails(contractAddress, byteCode)
}

func (b *Backend) GetTransferEvents(ctx context.Context, tokenDetails *tokens.TokenDetails, contractBlock int64, blockRangeLimit uint64) ([]*tokens.TransferEvent, error) {
	return b.tokenClient.GetTransferEvents(ctx, tokenDetails, contractBlock, blockRangeLimit)
}

func (b *Backend) CountTokenTransfers(address string) (int, error) {
	addr, err := b.mongo.getAddressByHash(address)
	if err != nil {
		return 0, fmt.Errorf("failed to get address: %v", err)
	}
	if addr != nil {
		return addr.NumberOfInternalTransactions, nil
	}
	return 0, nil

}

func (b *Backend) ImportBlock(ctx context.Context, block *types.Block) (*models.Block, error) {
	return b.mongo.importBlock(ctx, block, b.isDarvazaFork)
}

// NeedReloadParent returns true if the parent block is missing or does not match the hash from this block number.
func (b *Backend) NeedReloadParent(blockNumber int64) (bool, error) {
	return b.mongo.needReloadParent(blockNumber)
}

func (b *Backend) InternalTxsConsistent(blockNumber int64) (*models.Block, bool, error) {
	return b.mongo.internalTxsConsistent(blockNumber)
}

//return false if a number of transactions in DB is different from a number of transactions in the blockchain
func (b *Backend) ExternalTxsConsistent(ctx context.Context, block *models.Block) (bool, error) {
	lgr := b.Lgr.With(zap.String("hash", block.BlockHash))
	lgr.Debug("Checking transaction count for the block")
	txCount, err := b.goClient.TransactionCount(ctx, common.HexToHash(block.BlockHash))
	if err != nil {
		return false, fmt.Errorf("failed to get rpc transaction count: %v", err)
	}
	lgr.Debug("Got transaction count for the block", zap.Uint("rpc", txCount), zap.Int("db", block.TxCount))
	return txCount == uint(block.TxCount), nil
}

func (b *Backend) GetActiveAdresses(fromDate time.Time, onlyContracts bool) ([]*models.ActiveAddress, error) {
	addrs, err := b.mongo.getActiveAddresses(fromDate)
	if err != nil {
		return nil, err
	}
	var selectedAddresses []*models.ActiveAddress
	for _, address := range addrs {
		isContract, err := b.mongo.isContract(address.Address)
		if err != nil {
			return nil, fmt.Errorf("active address %s: %v", address.Address, err)
		}
		if onlyContracts == isContract {
			selectedAddresses = append(selectedAddresses, address)
		}
	}
	return selectedAddresses, nil
}
func (b *Backend) ImportAddress(address string, balance *big.Int, token *tokens.TokenDetails, contract bool, updatedAtBlock int64) (*models.Address, error) {
	return b.mongo.importAddress(address, balance, token, contract, updatedAtBlock)
}
func (b *Backend) ImportTokenHolder(contractAddress, tokenHolderAddress string, token *tokens.TokenHolderDetails, address *models.Address) (*models.TokenHolder, error) {
	return b.mongo.importTokenHolder(contractAddress, tokenHolderAddress, token, address)
}
func (b *Backend) ImportTransferEvent(ctx context.Context, contractAddress string, transferEvent *tokens.TransferEvent) (*models.TokenTransfer, error) {
	createdAt := time.Now()
	block, err := b.GetBlockByNumber(ctx, transferEvent.BlockNumber, false)
	if err != nil {
		return nil, err
	}
	if block != nil {
		createdAt = block.CreatedAt
	}
	return b.mongo.importTransferEvent(contractAddress, transferEvent, createdAt)
}
func (b *Backend) ImportContract(contractAddress string, byteCode string) error {
	return b.mongo.importContract(contractAddress, byteCode)
}

func (b *Backend) GetContractBlock(contractAddress string) (int64, error) {
	return b.mongo.getContractBlock(contractAddress)
}

func (b *Backend) BlockByNumber(ctx context.Context, blockNumber int64) (*types.Block, error) {
	var value *types.Block
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = b.goClient.BlockByNumber(ctx, big.NewInt(blockNumber))
		return err
	})
	return value, err
}

func (b *Backend) GetLatestBlockNumber(ctx context.Context) (*big.Int, error) {
	var value *big.Int
	err := utils.Retry(ctx, 5, 2*time.Second, func() (err error) {
		value, err = b.goClient.LatestBlockNumber(ctx)
		return err
	})
	return value, err
}

func (b *Backend) UpdateTotalFees(hash string, totalFees string) error {
	return b.mongo.updateTotalFees(hash, totalFees)
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
func (b *Backend) DeleteBlockByHash(hash string) error {
	return b.mongo.deleteBlockByHash(hash)
}
func (b *Backend) DeleteBlockByNumber(bnum int64) error {
	return b.mongo.deleteBlockByNumber(bnum)
}

func (b *Backend) DeleteContract(contractAddress string) error {
	return b.mongo.deleteContract(contractAddress)
}

func (b *Backend) MigrateDB(ctx context.Context, lgr *zap.Logger) (int, error) {
	return b.mongo.migrate(ctx, lgr)
}

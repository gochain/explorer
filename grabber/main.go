package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/explorer/server/tokens"
	"github.com/gochain-io/explorer/server/utils"

	"github.com/blendle/zapdriver"
	"github.com/gochain/gochain/v3/common"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	cfg := zapdriver.NewProductionConfig()
	logger, err := cfg.Build()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()
	defer func() {
		if rerr := recover(); rerr != nil {
			logger.Error("Fatal panic", zap.String("panic", fmt.Sprintf("%+v", rerr)))
		}
	}()

	var rpcUrl string
	var checkTxCount bool
	var mongoUrl string
	var dbName string
	var startFrom int64
	var blockRangeLimit uint64
	var workersCount uint

	app := cli.NewApp()
	app.Usage = "Grabber populates a mongo database with explorer data."

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "rpc-url, u",
			Value:       "https://rpc.gochain.io",
			Usage:       "rpc api url",
			Destination: &rpcUrl,
		},
		cli.StringFlag{
			Name:        "mongo-url, m",
			Value:       "127.0.0.1:27017",
			Usage:       "mongo connection url",
			Destination: &mongoUrl,
		},
		cli.BoolFlag{
			Name:        "tx-count, tx",
			Usage:       "check a transactions count for every block(a heavy operation)",
			Destination: &checkTxCount,
		},
		cli.StringFlag{
			Name:        "mongo-dbname, db",
			Value:       "blocks",
			Usage:       "mongo database name",
			Destination: &dbName,
		},
		cli.Int64Flag{
			Name:        "start-from, s",
			Value:       0, //1365000
			Usage:       "refill from this block",
			Destination: &startFrom,
		},
		cli.Uint64Flag{
			Name:        "block-range-limit, b",
			Value:       10000,
			Usage:       "block range limit",
			Destination: &blockRangeLimit,
		},
		cli.UintFlag{
			Name:        "workers-amount, w",
			Value:       10,
			Usage:       "parallel workers amount",
			Destination: &workersCount,
		},
		cli.StringSliceFlag{
			Name:  "locked-accounts",
			Usage: "accounts with locked funds to exclude from rich list and circulating supply",
		},
		cli.StringFlag{
			Name:  "log-level",
			Usage: "Minimum log level to include. Lower levels will be discarded. (debug, info, warn, error, dpanic, panic, fatal)",
		},
	}

	ctx, cancelFn := context.WithCancel(context.Background())
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range sigCh {
			cancelFn()
		}
	}()

	app.Action = func(c *cli.Context) error {
		if c.IsSet("log-level") {
			var lvl zapcore.Level
			s := c.String("log-level")
			if err := lvl.Set(s); err != nil {
				return fmt.Errorf("invalid log-level %q: %v", s, err)
			}
			cfg.Level.SetLevel(lvl)
		}
		lockedAccounts := c.StringSlice("locked-accounts")
		for i, l := range lockedAccounts {
			if !common.IsHexAddress(l) {
				return fmt.Errorf("invalid hex address: %s", l)
			}
			// Ensure canonical form, since queries are case-sensitive.
			lockedAccounts[i] = common.HexToAddress(l).Hex()
		}
		importer, err := backend.NewBackend(ctx, mongoUrl, rpcUrl, dbName, lockedAccounts, nil, logger)
		if err != nil {
			return fmt.Errorf("failed to create backend: %v", err)
		}
		go listener(ctx, importer)
		go updateStats(ctx, importer)
		go backfill(ctx, importer, startFrom, checkTxCount)
		go updateAddresses(ctx, 3*time.Minute, false, blockRangeLimit, workersCount, importer) // update only addresses
		updateAddresses(ctx, 5*time.Second, true, blockRangeLimit, workersCount, importer)     // update contracts
		return nil
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal("Fatal error", zap.Error(err))
	}
	logger.Info("Stopping")
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func listener(ctx context.Context, importer *backend.Backend) {
	var prevHeader int64
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			latest, err := importer.GetLatestBlockNumber(ctx)
			if err != nil {
				importer.Lgr.Error("Listener: Failed to get latest block number", zap.Error(err))
				continue
			}
			lgr := importer.Lgr.With(zap.Int64("block", latest))
			lgr.Debug("Listener: Getting block")
			if prevHeader != latest {
				lgr.Info("Listener: Getting block")
				block, err := importer.BlockByNumber(ctx, latest)
				if err != nil {
					lgr.Error("Listener: Failed to get block", zap.Error(err))
					continue
				}
				if block == nil {
					lgr.Error("Listener: Block not found")
					continue
				}
				if _, err := importer.ImportBlock(ctx, block); err != nil {
					lgr.Error("Listener: Failed to import block", zap.Error(err))
					continue
				}
				if err := checkAncestors(ctx, importer, block.Number().Int64(), 100); err != nil {
					lgr.Warn("Listener: Failed to check ancestors", zap.Error(err))
				}
				prevHeader = latest

			}
		}
	}
}

// backfill continuously loops over all blocks in reverse order, verifying that all
// data is loaded and consistent, and reloading as necessary.
func backfill(ctx context.Context, importer *backend.Backend, blockNumber int64, checkTxCount bool) {
	var hash common.Hash // Expected hash.
	for {
		if blockNumber < 0 {
			// Start over from latest.
			hash = common.Hash{}
			var err error
			blockNumber, err = importer.GetLatestBlockNumber(ctx)
			if err != nil {
				importer.Lgr.Error("Backfill: Failed to get latest block number", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			}
		}

		logger := importer.Lgr.With(zap.Int64("block", blockNumber))
		if (blockNumber % 1000) == 0 {
			logger.Info("Backfill: Progress")
		}
		var backfill bool
		dbBlock, err := importer.GetBlockByNumber(ctx, blockNumber)
		if err != nil {
			logger.Error("Backfill: Failed to get block", zap.Error(err))
			if utils.SleepCtx(ctx, 5*time.Second) != nil {
				return
			}
			continue
		} else if dbBlock == nil {
			// Missing.
			backfill = true
			logger = logger.With(zap.String("reason", "missing"))
		} else if !common.EmptyHash(hash) && dbBlock.BlockHash != hash.Hex() {
			// Mismatch with expected hash from parent of previous block.
			backfill = true
			logger = logger.With(zap.String("reason", "hash"))
		} else {
			// Note parent as next expected hash.
			hash = common.HexToHash(dbBlock.ParentHash)
		}
		if backfill {
			logger.Info("Backfill: Backfilling block")
			rpcBlock, err := importer.BlockByNumber(ctx, blockNumber)
			if err != nil {
				logger.Error("Backfill: Failed to get block", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			} else if rpcBlock == nil {
				logger.Error("Backfill: Block not found", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			}
			if _, err := importer.ImportBlock(ctx, rpcBlock); err != nil {
				logger.Error("Backfill: Failed to import block", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			}
			// Note parent as next expected hash.
			hash = rpcBlock.ParentHash()
		} else {
			newParent, err := ensureTxsConsistent(ctx, importer, blockNumber, checkTxCount)
			if err != nil {
				logger.Error("Backfill: Failed to check tx consistency", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			}
			if !common.EmptyHash(newParent) {
				hash = newParent
			}
		}
		blockNumber--
	}
}

// checkAncestors scans the ancestors of a block to verify that they exists and their hashes match.
func checkAncestors(ctx context.Context, importer *backend.Backend, blockNumber int64, generations int) error {
	if blockNumber == 0 || generations <= 0 {
		return nil
	}
	oldest := blockNumber - int64(generations)
	if oldest < 0 {
		oldest = 0
	}
	for blockNumber > oldest {
		lgr := importer.Lgr.With(zap.Int64("block", blockNumber))
		if ok, err := importer.NeedReloadParent(blockNumber); err != nil {
			lgr.Error("Failed to check parent", zap.Error(err))
			if err := utils.SleepCtx(ctx, 5*time.Second); err != nil {
				return err
			}
			continue
		} else if !ok {
			// Parent matches, so we're done.
			return nil
		}
		lgr.Info("Redownloading corrupted or missing ancestor")
		block, err := importer.BlockByNumber(ctx, blockNumber)
		if err != nil {
			lgr.Error("Failed to get ancestor", zap.Error(err))
			if err := utils.SleepCtx(ctx, 5*time.Second); err != nil {
				return err
			}
			continue
		} else if block == nil {
			lgr.Error("Ancestor not found")
			if err := utils.SleepCtx(ctx, 5*time.Second); err != nil {
				return err
			}
			continue
		}
		if _, err := importer.ImportBlock(ctx, block); err != nil {
			if err := utils.SleepCtx(ctx, 5*time.Second); err != nil {
				return err
			}
			continue
		}
		blockNumber--
	}
	return nil
}

// ensureTxsConsistent checks if the block tx count matches the number of txs in the db, and reimports the block if not.
// If a block is reimported, then the parent hash is returned.
func ensureTxsConsistent(ctx context.Context, importer *backend.Backend, blockNumber int64, checkTxCount bool) (common.Hash, error) {
	if ok, err := importer.TransactionsConsistent(blockNumber); err != nil {
		return common.Hash{}, fmt.Errorf("failed to check tx consistentcy: %v", err)
	} else if ok {
		if !checkTxCount {
			return common.Hash{}, nil
		}
		if ok, err := importer.TransactionCountConsistent(ctx, blockNumber); err != nil {
			return common.Hash{}, fmt.Errorf("failed to check tx count: %v", err)
		} else if ok {
			return common.Hash{}, nil
		}
	}

	importer.Lgr.Warn("Reimporting block due to inconsistent tx count", zap.Int64("block", blockNumber))

	block, err := importer.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get block: %v", err)
	}
	if block == nil {
		return common.Hash{}, errors.New("block not found")
	}
	if _, err := importer.ImportBlock(ctx, block); err != nil {
		return common.Hash{}, err
	}
	return block.ParentHash(), nil
}

func updateAddresses(ctx context.Context, sleep time.Duration, updateContracts bool, blockRangeLimit uint64, workersCount uint, importer *backend.Backend) {
	lastUpdatedAt := time.Unix(0, 0)
	lastBlockUpdatedAt := int64(0)
	t := time.NewTicker(sleep)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			start := time.Now()
			currentTime := time.Now()
			currentBlock, err := importer.GetLatestBlockNumber(ctx)
			if err != nil {
				importer.Lgr.Error("Update Addresses: Failed to get latest block number", zap.Error(err))
				continue
			}
			addresses, err := importer.GetActiveAdresses(lastUpdatedAt, updateContracts)
			if err != nil {
				importer.Lgr.Error("Update Addresses: Failed to get active addresses", zap.Error(err))
				continue
			}
			importer.Lgr.Info("Update Addresses: Starting", zap.Time("lastUpdatedAt", lastUpdatedAt), zap.Int("count", len(addresses)))
			var jobs = make(chan *models.ActiveAddress, workersCount)
			var wg sync.WaitGroup
			for i := 0; i < int(workersCount); i++ {
				wg.Add(1)
				lgr := importer.Lgr.With(zap.Int("worker", i))
				go worker(ctx, wg.Done, lgr, jobs, currentBlock, blockRangeLimit, importer)
			}
		forAddrs:
			for _, address := range addresses {
				select {
				case <-ctx.Done():
					break forAddrs
				case jobs <- address:
				}
			}
			close(jobs)
			wg.Wait()
			elapsed := time.Since(start)
			importer.Lgr.Info("Update Addresses: Complete", zap.Bool("updateContracts", updateContracts),
				zap.Duration("elapsed", elapsed), zap.Int64("block", lastBlockUpdatedAt))
			lastBlockUpdatedAt = currentBlock
			lastUpdatedAt = currentTime
		}
	}
}

func worker(ctx context.Context, done func(), lgr *zap.Logger, jobs chan *models.ActiveAddress, currentBlock int64, blockRangeLimit uint64, importer *backend.Backend) {
	defer done()
	var errs []error
	var updated int
loop:
	for {
		select {
		case <-ctx.Done():
			lgr.Error("Update Addresses: worker cancelled", zap.NamedError("ctxErr", ctx.Err()),
				zap.Int("updated", updated), zap.Int("failed", len(errs)), zap.Errors("errors", errs))
			return
		case address, ok := <-jobs:
			if !ok {
				break loop
			}
			if address != nil {
				err := updateAddress(ctx, address, currentBlock, blockRangeLimit, importer)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to update address %s: %v", address.Address, err))
				} else {
					updated++
				}
			}
		}
	}
	if len(errs) > 0 {
		lgr.Error("Update Addresses: worker done", zap.Int("updated", updated), zap.Int("failed", len(errs)), zap.Errors("errors", errs))
		return
	}
	lgr.Info("Update Addresses: worker done", zap.Int("updated", updated))
}

func updateAddress(ctx context.Context, address *models.ActiveAddress, currentBlock int64, blockRangeLimit uint64, importer *backend.Backend) error {
	if !common.IsHexAddress(address.Address) {
		return fmt.Errorf("invalid hex address: %s", address.Address)
	}
	addr := common.HexToAddress(address.Address)
	normalizedAddress := addr.Hex()
	lgr := importer.Lgr.With(zap.String("address", normalizedAddress))
	balance, err := importer.Balance(ctx, addr)
	if err != nil {
		return fmt.Errorf("failed to get balance")
	}
	contractDataArray, err := importer.CodeAt(ctx, normalizedAddress)
	contractData := string(contractDataArray[:])
	var tokenDetails = &tokens.TokenDetails{TotalSupply: big.NewInt(0)}
	contract := false
	if contractData != "" {
		contract = true
		byteCode := hex.EncodeToString(contractDataArray)
		if err := importer.ImportContract(normalizedAddress, byteCode); err != nil {
			return fmt.Errorf("failed to import contract: %v", err)
		}
		var fromBlock int64
		contractFromDB, err := importer.GetAddressByHash(ctx, normalizedAddress)
		if err != nil {
			return fmt.Errorf("failed to get contract from DB: %v", err)
		}
		if contractFromDB != nil && contractFromDB.UpdatedAtBlock > 0 {
			fromBlock = contractFromDB.UpdatedAtBlock
		} else {
			fromBlock, err = importer.GetContractBlock(normalizedAddress)
			if err != nil {
				fmt.Errorf("failed to get contract block: %v", err)
			}
		}
		tokenDetails, err = importer.GetTokenDetails(normalizedAddress, byteCode)
		if err != nil {
			return fmt.Errorf("failed to get token details: %v", err)
		}
		if contractFromDB.TokenName == "" || contractFromDB.TokenSymbol == "" {
			lgr.Info("Updating token details", zap.Int64("block", fromBlock),
				zap.String("symbol", tokenDetails.Symbol), zap.String("name", tokenDetails.Name))
			contractFromDB.TokenName = tokenDetails.Name
			contractFromDB.TokenSymbol = tokenDetails.Symbol
		}
		tokenTransfers, err := importer.GetTransferEvents(ctx, tokenDetails, fromBlock, blockRangeLimit)
		if err != nil {
			return fmt.Errorf("failed to get internal txs: %v", err)
		}
		tokenTransfersFromDB, err := importer.CountTokenTransfers(normalizedAddress)
		if err != nil {
			return fmt.Errorf("failed to count internal txs: %v", err)
		}
		lgr.Info("Comparing internal tx count from DB against RPC", zap.Int("db", tokenTransfersFromDB),
			zap.Int("rpc", len(tokenTransfers)))
		if len(tokenTransfers) != tokenTransfersFromDB {
			var tokenHoldersList []string
			for _, itx := range tokenTransfers {
				lgr.Debug("Internal Transaction", zap.Stringer("from", itx.From),
					zap.Stringer("to", itx.To), zap.Stringer("value", itx.Value))
				if _, err := importer.ImportTransferEvent(ctx, normalizedAddress, itx); err != nil {
					return fmt.Errorf("failed to import internal tx: %v", err)
				}
				// if itx.BlockNumber > lastBlockUpdatedAt.Int64() {
				lgr.Debug("Updating following token holder addresses", zap.Stringer("from", itx.From),
					zap.Stringer("to", itx.To), zap.Stringer("value", itx.Value))
				tokenHoldersList = appendIfMissing(tokenHoldersList, itx.To.String())
				tokenHoldersList = appendIfMissing(tokenHoldersList, itx.From.String())
				// }
			}
			for index, tokenHolderAddress := range tokenHoldersList {
				if tokenHolderAddress == "0x0000000000000000000000000000000000000000" {
					continue
				}
				lgr.Info("Importing token holder", zap.String("holder", tokenHolderAddress),
					zap.Int("index", index), zap.Int("total", len(tokenHoldersList)))
				tokenHolder, err := importer.GetTokenBalance(normalizedAddress, tokenHolderAddress)
				if err != nil {
					lgr.Error("Failed to get token balance", zap.Error(err), zap.String("holder", tokenHolderAddress))
					continue
				}
				if contractFromDB == nil {
					lgr.Error("Cannot find contract in DB", zap.Error(err), zap.String("holder", tokenHolderAddress))
					continue
				}
				if _, err := importer.ImportTokenHolder(normalizedAddress, tokenHolderAddress, tokenHolder, contractFromDB); err != nil {
					return fmt.Errorf("failed to import token holder: %v", err)
				}
			}
		}
	}
	lgr.Info("Update Addresses: updated address", zap.Stringer("balance", balance))
	_, err = importer.ImportAddress(normalizedAddress, balance, tokenDetails, contract, currentBlock)
	return err
}

func updateStats(ctx context.Context, importer *backend.Backend) {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			stats, err := importer.UpdateStats()
			if err != nil {
				importer.Lgr.Error("Failed to update stats", zap.Error(err), zap.Reflect("stats", stats))
				continue
			}
			importer.Lgr.Info("Updated stats", zap.Reflect("stats", stats))
		}
	}
}

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
	"github.com/gochain/gochain/v4/common"
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
	var checkExternal bool
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
			Usage:       "verify transaction count against RPC for every block(a heavy operation)",
			Destination: &checkExternal,
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
		b, err := backend.NewBackend(ctx, mongoUrl, rpcUrl, dbName, lockedAccounts, nil, nil, logger, nil)
		if err != nil {
			return fmt.Errorf("failed to create backend: %v", err)
		}
		go migrator(ctx, b, logger)
		go listener(ctx, b)
		go updateStats(ctx, b)
		go backfill(ctx, b, startFrom, checkExternal)
		go updateAddresses(ctx, 3*time.Minute, false, blockRangeLimit, workersCount, b) // update only addresses
		updateAddresses(ctx, 5*time.Second, true, blockRangeLimit, workersCount, b)     // update contracts
		return nil
	}
	err = app.Run(os.Args)
	if err != nil {
		logger.Fatal("Fatal error", zap.Error(err))
	}
	logger.Info("Stopping")
}
func migrator(ctx context.Context, b *backend.Backend, lgr *zap.Logger) {
	version, err := b.MigrateDB(ctx, lgr)
	if err != nil {
		lgr.Error("Migration failed", zap.Error(err))
		return
	}
	lgr.Info("Migrations successfully complete", zap.Int("version", version))
}

func listener(ctx context.Context, b *backend.Backend) {
	var prevHeader int64
	t := time.NewTicker(time.Second * 1)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			latestBig, err := b.GetLatestBlockNumber(ctx)
			if err != nil {
				b.Lgr.Error("Listener: Failed to get latest block number", zap.Error(err))
				continue
			}
			latest := latestBig.Int64()
			lgr := b.Lgr.With(zap.Int64("block", latest))
			lgr.Debug("Listener: Getting block")
			if prevHeader != latest {
				lgr.Debug("Listener: Getting block", zap.Int64("block", latest))
				block, err := b.BlockByNumber(ctx, latest)
				if err != nil {
					lgr.Error("Listener: Failed to get block", zap.Error(err))
					continue
				}
				if block == nil {
					lgr.Error("Listener: Block not found")
					continue
				}
				if _, err := b.ImportBlock(ctx, block); err != nil {
					lgr.Error("Listener: Failed to import block", zap.Error(err))
					continue
				}
				if err := checkAncestors(ctx, b, block.Number().Int64(), 100); err != nil {
					lgr.Warn("Listener: Failed to check ancestors", zap.Error(err))
				}
				prevHeader = latest

			}
		}
	}
}

// backfill continuously loops over all blocks in reverse order, verifying that all
// data is loaded and consistent, and reloading as necessary.
func backfill(ctx context.Context, b *backend.Backend, blockNumber int64, checkExternal bool) {
	var hash common.Hash // Expected hash.
	for {
		if blockNumber < 0 {
			// Start over from latest.
			hash = common.Hash{}
			blockNumberBig, err := b.GetLatestBlockNumber(ctx)
			if err != nil {
				b.Lgr.Error("Backfill: Failed to get latest block number", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			}
			blockNumber = blockNumberBig.Int64()
		}

		logger := b.Lgr.With(zap.Int64("block", blockNumber))
		if (blockNumber % 1000) == 0 {
			logger.Debug("Backfill: Progress", zapcore.Field{Key: "block", Type: zapcore.Int64Type, Integer: blockNumber})
		}
		var backfill bool
		dbBlock, err := b.GetBlockByNumber(ctx, blockNumber, true)
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
		} else if dbBlock.GasFees == "" {
			backfill = true
			logger = logger.With(zap.String("reason", "missing gas fees"))
		} else {
			// Note parent as next expected hash.
			hash = common.HexToHash(dbBlock.ParentHash)
		}
		if backfill {
			logger.Info("Backfill: Backfilling block")
			rpcBlock, err := b.BlockByNumber(ctx, blockNumber)
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
			if bl, err := b.ImportBlock(ctx, rpcBlock); err != nil {
				if ctx.Err() != nil {
					return
				}
				logger.Error("Backfill: Failed to import block", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			} else if bl.TotalFeesBurned != "" {
				go fillTotalFees(ctx, b, bl)
			}

			// Note parent as next expected hash.
			hash = rpcBlock.ParentHash()
		} else {
			updated, newParent, err := ensureTxsConsistent(ctx, b, blockNumber, checkExternal)
			if err != nil {
				logger.Error("Backfill: Failed to check tx consistency", zap.Error(err))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			} else if updated != nil {
				if total := updated.TotalFeesBurned; total != "" && total != dbBlock.TotalFeesBurned {
					go fillTotalFees(ctx, b, updated)
				}
			}
			if !common.EmptyHash(newParent) {
				hash = newParent
			}
		}
		blockNumber--
	}
}

// checkAncestors scans the ancestors of a block to verify that they exists and their hashes match.
func checkAncestors(ctx context.Context, b *backend.Backend, blockNumber int64, generations int) error {
	if blockNumber == 0 || generations <= 0 {
		return nil
	}
	oldest := blockNumber - int64(generations)
	if oldest < 0 {
		oldest = 0
	}
	for blockNumber > oldest {
		lgr := b.Lgr.With(zap.Int64("block", blockNumber))
		if ok, err := b.NeedReloadParent(blockNumber); err != nil {
			lgr.Error("Failed to check parent", zap.Error(err))
			if err := utils.SleepCtx(ctx, 5*time.Second); err != nil {
				return err
			}
			continue
		} else if !ok {
			// Parent matches, so we're done.
			return nil
		}
		blockNumber--
		lgr.Info("Redownloading corrupted or missing ancestor")
		block, err := b.BlockByNumber(ctx, blockNumber)
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
		if bl, err := b.ImportBlock(ctx, block); err != nil {
			if err := utils.SleepCtx(ctx, 5*time.Second); err != nil {
				return err
			}
			continue
		} else if bl.TotalFeesBurned != "" {
			go fillTotalFees(ctx, b, bl)
		}
	}
	return nil
}

// ensureTxsConsistent checks if the block tx count matches the number of txs in the db, and reimports the block if not.
// If a block is reimported and the parent changed, then the parent hash is returned.
func ensureTxsConsistent(ctx context.Context, b *backend.Backend, blockNumber int64, checkExternal bool) (*models.Block, common.Hash, error) {
	block, ok, err := b.InternalTxsConsistent(blockNumber)
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to check tx consistentcy: %v", err)
	} else if ok {
		if !checkExternal {
			return nil, common.Hash{}, nil
		}
		if ok, err := b.ExternalTxsConsistent(ctx, block); err != nil {
			return nil, common.Hash{}, fmt.Errorf("failed to check tx count: %v", err)
		} else if ok {
			return nil, common.Hash{}, nil
		}
	}

	b.Lgr.Warn("Reimporting block due to inconsistent tx count", zap.Int64("block", blockNumber))

	rpcBlock, err := b.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, common.Hash{}, fmt.Errorf("failed to get block: %v", err)
	}
	if rpcBlock == nil {
		return nil, common.Hash{}, errors.New("block not found")
	}
	var newParent common.Hash
	updated, err := b.ImportBlock(ctx, rpcBlock)
	if err != nil {
		return nil, common.Hash{}, err
	} else if updated.ParentHash != block.BlockHash {
		newParent = common.HexToHash(updated.ParentHash)
	}
	return updated, newParent, nil
}

func updateAddresses(ctx context.Context, sleep time.Duration, updateContracts bool, blockRangeLimit uint64, workersCount uint, b *backend.Backend) {
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
			currentBlockBig, err := b.GetLatestBlockNumber(ctx)
			if err != nil {
				b.Lgr.Error("Update Addresses: Failed to get latest block number", zap.Error(err))
				continue
			}
			currentBlock := currentBlockBig.Int64()
			addresses, err := b.GetActiveAdresses(lastUpdatedAt, updateContracts)
			if err != nil {
				b.Lgr.Error("Update Addresses: Failed to get active addresses", zap.Error(err))
				continue
			}
			b.Lgr.Info("Update Addresses: Starting", zap.Time("lastUpdatedAt", lastUpdatedAt), zap.Int("count", len(addresses)))
			var jobs = make(chan *models.ActiveAddress, workersCount)
			var wg sync.WaitGroup
			for i := 0; i < int(workersCount); i++ {
				wg.Add(1)
				lgr := b.Lgr.With(zap.Int("worker", i))
				go worker(ctx, wg.Done, lgr, jobs, currentBlock, blockRangeLimit, b)
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
			b.Lgr.Info("Update Addresses: Complete", zap.Bool("updateContracts", updateContracts),
				zap.Duration("elapsed", elapsed), zap.Int64("block", lastBlockUpdatedAt))
			lastBlockUpdatedAt = currentBlock
			lastUpdatedAt = currentTime
		}
	}
}

func worker(ctx context.Context, done func(), lgr *zap.Logger, jobs chan *models.ActiveAddress, currentBlock int64, blockRangeLimit uint64, b *backend.Backend) {
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
				err := updateAddress(ctx, address, currentBlock, blockRangeLimit, b)
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to update address %s: %v", address.Address, err))
				} else {
					updated++
				}
			}
		}
	}
	if len(errs) > 0 {
		lgr.Error("Update Addresses: worker failed", zap.Int("updated", updated), zap.Int("failed", len(errs)), zap.Errors("errors", errs))
		return
	}
	lgr.Debug("Update Addresses: worker done", zap.Int("updated", updated))
}

func updateAddress(ctx context.Context, address *models.ActiveAddress, currentBlock int64, blockRangeLimit uint64, b *backend.Backend) error {
	if !common.IsHexAddress(address.Address) {
		return fmt.Errorf("invalid hex address: %s", address.Address)
	}
	addr := common.HexToAddress(address.Address)
	normalizedAddress := addr.Hex()
	lgr := b.Lgr.With(zap.String("address", normalizedAddress))
	balance, err := b.Balance(ctx, addr)
	if err != nil {
		return fmt.Errorf("failed to get balance")
	}
	contractDataArray, err := b.CodeAt(ctx, normalizedAddress)
	contractData := string(contractDataArray[:])
	var tokenDetails = &tokens.TokenDetails{TotalSupply: big.NewInt(0)}
	contract := false
	if contractData != "" {
		contract = true
		byteCode := hex.EncodeToString(contractDataArray)
		if err := b.ImportContract(normalizedAddress, byteCode); err != nil {
			return fmt.Errorf("failed to import contract: %v", err)
		}
		var fromBlock int64
		contractFromDB, err := b.GetAddressByHash(ctx, normalizedAddress)
		if err != nil {
			return fmt.Errorf("failed to get contract from DB: %v", err)
		}
		if contractFromDB != nil && contractFromDB.UpdatedAtBlock > 0 {
			fromBlock = contractFromDB.UpdatedAtBlock
		} else {
			fromBlock, err = b.GetContractBlock(normalizedAddress)
			if err != nil {
				fmt.Errorf("failed to get contract block: %v", err)
			}
		}
		tokenDetails, err = b.GetTokenDetails(normalizedAddress, byteCode)
		if err != nil {
			return fmt.Errorf("failed to get token details: %v", err)
		}
		if contractFromDB.TokenName == "" || contractFromDB.TokenSymbol == "" {
			lgr.Info("Updating token details", zap.Int64("block", fromBlock),
				zap.String("symbol", tokenDetails.Symbol), zap.String("name", tokenDetails.Name))
			contractFromDB.TokenName = tokenDetails.Name
			contractFromDB.TokenSymbol = tokenDetails.Symbol
		}
		tokenTransfers, err := b.GetTransferEvents(ctx, tokenDetails, fromBlock, blockRangeLimit)
		if err != nil {
			return fmt.Errorf("failed to get internal txs: %v", err)
		}
		tokenTransfersFromDB, err := b.CountTokenTransfers(normalizedAddress)
		if err != nil {
			return fmt.Errorf("failed to count internal txs: %v", err)
		}
		lgr.Info("Comparing internal tx count from DB against RPC", zap.Int("db", tokenTransfersFromDB),
			zap.Int("rpc", len(tokenTransfers)))
		if len(tokenTransfers) != tokenTransfersFromDB {
			tokenHoldersList := make(map[string]struct{})
			for _, itx := range tokenTransfers {
				lgr.Debug("Internal Transaction", zap.Stringer("from", itx.From),
					zap.Stringer("to", itx.To), zap.Stringer("value", itx.Value))
				if _, err := b.ImportTransferEvent(ctx, normalizedAddress, itx); err != nil {
					return fmt.Errorf("failed to import internal tx: %v", err)
				}
				// if itx.BlockNumber > lastBlockUpdatedAt.Int64() {
				lgr.Debug("Updating following token holder addresses", zap.Stringer("from", itx.From),
					zap.Stringer("to", itx.To), zap.Stringer("value", itx.Value))
				tokenHoldersList[itx.To.String()] = struct{}{}
				tokenHoldersList[itx.From.String()] = struct{}{}
				// }
			}
			for tokenHolderAddress := range tokenHoldersList {
				if tokenHolderAddress == "0x0000000000000000000000000000000000000000" {
					continue
				}
				lgr.Info("Importing token holder", zap.String("holder", tokenHolderAddress),
					zap.Int("total", len(tokenHoldersList)))
				tokenHolder, err := b.GetTokenBalance(normalizedAddress, tokenHolderAddress)
				if err != nil {
					lgr.Error("Failed to get token balance", zap.Error(err), zap.String("holder", tokenHolderAddress))
					continue
				}
				if contractFromDB == nil {
					lgr.Error("Cannot find contract in DB", zap.Error(err), zap.String("holder", tokenHolderAddress))
					continue
				}
				if _, err := b.ImportTokenHolder(normalizedAddress, tokenHolderAddress, tokenHolder, contractFromDB); err != nil {
					return fmt.Errorf("failed to import token holder: %v", err)
				}
			}
		}
	}
	lgr.Info("Update Addresses: updated address", zap.Stringer("balance", balance))
	_, err = b.ImportAddress(normalizedAddress, balance, tokenDetails, contract, currentBlock)
	return err
}

func updateStats(ctx context.Context, b *backend.Backend) {
	t := time.NewTicker(5 * time.Minute)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			stats, err := b.UpdateStats()
			if err != nil {
				b.Lgr.Error("Failed to update stats", zap.Error(err), zap.Reflect("stats", stats))
				continue
			}
			b.Lgr.Info("Updated stats", zap.Reflect("stats", stats))
		}
	}
}

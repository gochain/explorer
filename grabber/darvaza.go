package main

import (
	"context"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/explorer/server/utils"
)

// fillTotalFees sets the TotalFeesBurned field of parent's descendents on the canonical chain.
func fillTotalFees(ctx context.Context, b *backend.Backend, parent *models.Block) {
	b.Lgr.Info("Fees: Filling total fees", zap.Int64("start", parent.Number), zap.String("hash", parent.BlockHash))
	defer func() {
		b.Lgr.Info("Fees: Done filling total fees", zap.Int64("end", parent.Number), zap.String("hash", parent.BlockHash))
	}()
	for {
		if time.Since(parent.CreatedAt) < 5*time.Second {
			b.Lgr.Info("Fees: Filled total fees up to latest")
			return
		}
		bl, err := b.GetBlockByNumber(ctx, parent.Number+1, false)
		if err != nil {
			b.Lgr.Warn("Fees: Failed to get block", zap.Int64("number", parent.Number+1))
			if utils.SleepCtx(ctx, 5*time.Second) != nil {
				return
			}
			continue
		}
		if bl.ParentHash != parent.BlockHash {
			// reorg makes this irrelevant
			return
		}
		if bl.GasFees == "" {
			b.Lgr.Warn("Fees: Missing gas fee. Reimporting block", zap.Int64("number", bl.Number))
			rpcBlock, err := b.BlockByNumber(ctx, bl.Number)
			if err != nil {
				b.Lgr.Warn("Fees: Failed to get rpc block", zap.Int64("number", bl.Number))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			}
			bl, err = b.ImportBlock(ctx, rpcBlock)
			if err != nil {
				b.Lgr.Warn("Fees: Failed to import block", zap.Int64("number", bl.Number))
				if utils.SleepCtx(ctx, 5*time.Second) != nil {
					return
				}
				continue
			}
			if bl.ParentHash != parent.BlockHash {
				// reorg makes this irrelevant
				return
			}
			if bl.GasFees == "" {
				// give up
				b.Lgr.Error("Fees: Missing gas fee. Unable to calculate total", zap.Int64("number", bl.Number))
				return
			}
		}
		parentTotal, ok := new(big.Int).SetString(parent.TotalFeesBurned, 10)
		if !ok {
			b.Lgr.Error("Fees: Failed to parse total fee", zap.Int64("number", parent.Number), zap.String("total", parent.TotalFeesBurned))
			return
		}
		gasFees, ok := new(big.Int).SetString(bl.GasFees, 10)
		if !ok {
			b.Lgr.Error("Fees: Failed to parse block total gas fees", zap.Int64("number", bl.Number), zap.String("total", bl.GasFees))
			return
		}
		bl.TotalFeesBurned = new(big.Int).Add(parentTotal, gasFees).String()
		if err := b.UpdateTotalFees(bl.BlockHash, bl.TotalFeesBurned); err != nil {
			if utils.SleepCtx(ctx, 5*time.Second) != nil {
				return
			}
			continue
		}
		parent = bl
	}
}

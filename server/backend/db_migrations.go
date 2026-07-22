package backend

import (
	"context"
	"math/big"
	"time"

	"github.com/gochain-io/explorer/internal/migrate"
	"github.com/gochain-io/explorer/server/models"

	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var migrationCollection = "Migrations"

func (mb *MongoBackend) getDatabaseVersion() (int, error) {
	mb.databaseVersionMutex.RLock()
	version := mb.databaseVersion
	mb.databaseVersionMutex.RUnlock()
	if version != 0 {
		return version, nil
	}
	var result struct {
		ID int `bson:"ID"`
	}
	err := mb.mongo.C(migrationCollection).Find(nil).Sort("-ID").One(&result)
	if err != nil {
		return 0, err
	}
	mb.databaseVersionMutex.Lock()
	mb.databaseVersion = result.ID
	mb.databaseVersionMutex.Unlock()
	return result.ID, nil
}

var migrationTransactionsByAddress = migrate.Migration{
	ID:      1,
	Comment: "Creating TransactionsByAddress collection",
	Migrate: func(ctx context.Context, d *mgo.Database, lgr *zap.Logger) error {
		var tx *models.Transaction
		find := d.C("Transactions").Find(bson.M{})
		txs := find.Iter()
		defer txs.Close()

		const batchUpsertSize = 1000
		bulk := d.C("TransactionsByAddress").Bulk()
		bulk.Unordered()

		var txsCnt, bulkCnt int
		for txs.Next(&tx) {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			txsCnt++
			if (txsCnt % 1000) == 0 {
				lgr.Info("Processed:", zap.Int("records:", txsCnt))
			}
			bulkCnt++
			bulk.Upsert(bson.M{"address": tx.From, "tx_hash": tx.TxHash},
				bson.M{"address": tx.From, "tx_hash": tx.TxHash, "created_at": tx.CreatedAt})
			if tx.To != tx.From {
				bulkCnt++
				bulk.Upsert(bson.M{"address": tx.To, "tx_hash": tx.TxHash},
					bson.M{"address": tx.To, "tx_hash": tx.TxHash, "created_at": tx.CreatedAt})
			}
			if bulkCnt > batchUpsertSize {
				bulkCnt = 0
				// Execute batch of upserts.
				if _, err := bulk.Run(); err != nil {
					return err
				}
				bulk = d.C("TransactionsByAddress").Bulk()
				bulk.Unordered()

			}
		}
		if bulkCnt > 0 {
			// Flush remaining upserts.
			if _, err := bulk.Run(); err != nil {
				return err
			}
		}
		return txs.Err()
	},
	Rollback: func(ctx context.Context, s *mgo.Database, lgr *zap.Logger) error {
		return nil
	},
}

var migrationTotalFeesBurned = migrate.Migration{
	ID:      2,
	Comment: "Backfilling total_fees_burned on Blocks collection",
	Migrate: func(ctx context.Context, d *mgo.Database, lgr *zap.Logger) error {
		var startBlock struct {
			Number          int64  `bson:"number"`
			TotalFeesBurned string `bson:"total_fees_burned"`
		}
		err := d.C("Blocks").
			Find(bson.M{"total_fees_burned": bson.M{"$gt": ""}}).
			Hint("-number", "total_fees_burned").
			Sort("-number").
			Select(bson.M{"number": 1, "total_fees_burned": 1}).One(&startBlock)

		var startNum int64 = 0
		accum := new(big.Int)
		if err == nil && startBlock.TotalFeesBurned != "" {
			startNum = startBlock.Number
			accum.SetString(startBlock.TotalFeesBurned, 10)
			lgr.Info("Starting fee burn backfill", zap.Int64("from_block", startNum), zap.String("start_burned", startBlock.TotalFeesBurned))
		} else {
			lgr.Info("No existing fee burn found, starting from block 0")
		}

		go func() {
			bgCtx := context.Background()
			iter := d.C("Blocks").
				Find(bson.M{"number": bson.M{"$gt": startNum}}).
				Select(bson.M{"number": 1, "gas_fees": 1}).
				Sort("number").
				Iter()
			defer iter.Close()

			var b struct {
				Number  int64  `bson:"number"`
				GasFees string `bson:"gas_fees"`
			}

			const batchSize = 5000
			bulk := d.C("Blocks").Bulk()
			bulk.Unordered()
			var count, batchCount int

			for iter.Next(&b) {
				if bgCtx.Err() != nil {
					return
				}
				if b.GasFees != "" && b.GasFees != "0" {
					if fee, ok := new(big.Int).SetString(b.GasFees, 10); ok {
						accum.Add(accum, fee)
					}
				}
				bulk.Update(
					bson.M{"number": b.Number},
					bson.M{"$set": bson.M{"total_fees_burned": accum.String()}},
				)
				count++
				batchCount++

				if count%50000 == 0 {
					lgr.Info("Backfilled total_fees_burned progress", zap.Int("blocks_processed", count), zap.Int64("current_block", b.Number), zap.String("total_burned", accum.String()))
				}

				if batchCount >= batchSize {
					if _, err := bulk.Run(); err != nil {
						lgr.Error("Failed bulk update during backfill", zap.Error(err))
						return
					}
					bulk = d.C("Blocks").Bulk()
					bulk.Unordered()
					batchCount = 0
					time.Sleep(10 * time.Millisecond)
				}
			}
			if batchCount > 0 {
				if _, err := bulk.Run(); err != nil {
					lgr.Error("Failed final bulk update during backfill", zap.Error(err))
					return
				}
			}
			lgr.Info("Completed total_fees_burned backfill", zap.Int("total_blocks", count), zap.String("final_total_burned", accum.String()))
		}()
		return nil
	},
	Rollback: func(ctx context.Context, s *mgo.Database, lgr *zap.Logger) error {
		return nil
	},
}

func (mb *MongoBackend) migrate(ctx context.Context, lgr *zap.Logger) (int, error) {
	m := migrate.New(ctx, mb.mongo, lgr, migrationCollection, []*migrate.Migration{&migrationTransactionsByAddress, &migrationTotalFeesBurned})
	return m.Migrate()
}

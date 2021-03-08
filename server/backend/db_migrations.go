package backend

import (
	"context"

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

func (mb *MongoBackend) migrate(ctx context.Context, lgr *zap.Logger) (int, error) {
	m := migrate.New(ctx, mb.mongo, lgr, migrationCollection, []*migrate.Migration{&migrationTransactionsByAddress})
	return m.Migrate()
}

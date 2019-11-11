package backend

import (
	"context"

	"github.com/gochain-io/explorer/internal/migrate"
	"github.com/gochain-io/explorer/server/models"

	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (self *MongoBackend) Migrate(ctx context.Context, lgr *zap.Logger) error {
	m := migrate.New(ctx, self.mongo, lgr, migrate.DefaultOptions, []*migrate.Migration{{
		ID:      1,
		Comment: "Creating TransactionsByAddress collection",
		Migrate: func(ctx context.Context, d *mgo.Database, lgr *zap.Logger) error {
			counter := 0
			var tx *models.Transaction
			find := d.C("Transactions").Find(bson.M{})
			txs := find.Iter()
			defer txs.Close()
			for txs.Next(&tx) {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				counter++
				if (counter % 1000) == 0 {
					lgr.Info("Processed:", zap.Int("records:", counter))
				}
				_, err := self.mongo.C("TransactionsByAddress").Upsert(bson.M{"address": tx.From, "tx_hash": tx.TxHash},
					bson.M{"address": tx.From, "tx_hash": tx.TxHash, "created_at": tx.CreatedAt})
				if err != nil {
					return err
				}
				_, err = self.mongo.C("TransactionsByAddress").Upsert(bson.M{"address": tx.To, "tx_hash": tx.TxHash},
					bson.M{"address": tx.To, "tx_hash": tx.TxHash, "created_at": tx.CreatedAt})
				if err != nil {
					return err
				}
			}
			return txs.Err()
		},
		Rollback: func(ctx context.Context, s *mgo.Database, lgr *zap.Logger) error {
			return nil
		},
	}})

	return m.Migrate()
}

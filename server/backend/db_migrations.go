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

func (self *MongoBackend) getDatabaseVersion() (int, error) {
	self.databaseVersionMutex.RLock()
	version := self.databaseVersion
	self.databaseVersionMutex.RUnlock()
	if version != 0 {
		return version, nil
	}
	var result struct {
		ID int `bson:"ID"`
	}
	err := self.mongo.C(migrationCollection).Find(nil).Sort("-ID").One(&result)
	if err != nil {
		return 0, err
	}
	self.databaseVersionMutex.Lock()
	self.databaseVersion = result.ID
	self.databaseVersionMutex.Unlock()
	return result.ID, nil
}

var migrationTransactionsByAddress = migrate.Migration{
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
			_, err := d.C("TransactionsByAddress").Upsert(bson.M{"address": tx.From, "tx_hash": tx.TxHash},
				bson.M{"address": tx.From, "tx_hash": tx.TxHash, "created_at": tx.CreatedAt})
			if err != nil {
				return err
			}
			_, err = d.C("TransactionsByAddress").Upsert(bson.M{"address": tx.To, "tx_hash": tx.TxHash},
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
}

func (self *MongoBackend) migrate(ctx context.Context, lgr *zap.Logger) (int, error) {
	m := migrate.New(ctx, self.mongo, lgr, migrationCollection, []*migrate.Migration{&migrationTransactionsByAddress})
	return m.Migrate()
}

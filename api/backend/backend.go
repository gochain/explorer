package backend

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gochain-io/explorer/api/models"
	"github.com/gochain-io/gochain/core/types"
	"github.com/gochain-io/gochain/ethclient"
)

type MongoBackend struct {
	mongo     *mgo.Database
	ethclient *ethclient.Client
}

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

func (self *MongoBackend) parseTx(tx *types.Transaction, block *types.Block) *models.Transaction {
	from, err := self.ethclient.TransactionSender(context.Background(), tx, block.Header().Hash(), 0)
	if err != nil {
		log.Fatal().Err(err).Msg("parseTx")
	}
	txx := &models.Transaction{TxHash: tx.Hash().Hex(),
		To:          tx.To().Hex(),
		From:        from.Hex(),
		Amount:      tx.Value().Int64(),
		Price:       tx.GasPrice().String(),
		GasLimit:    strconv.Itoa(int(tx.Gas())),
		BlockNumber: block.Number().Int64(),
		Nonce:       string(tx.Nonce()),
		BlockHash:   block.Hash().Hex(),
		CreatedAt:   time.Unix(block.Time().Int64(), 0)}
	return txx
}
func (self *MongoBackend) createIndexes() {
	err := self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"from"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"to"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"tx_hash"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("Blocks").EnsureIndex(mgo.Index{Key: []string{"number"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("Blocks").EnsureIndex(mgo.Index{Key: []string{"miner"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("ActiveAddress").EnsureIndex(mgo.Index{Key: []string{"updated_at"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("ActiveAddress").EnsureIndex(mgo.Index{Key: []string{"address"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("Address").EnsureIndex(mgo.Index{Key: []string{"address"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
}
func NewBackend(ethclient *ethclient.Client) *MongoBackend {

	Host := []string{
		"127.0.0.1:27017",
	}
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: Host,
	})
	if err != nil {
		panic(err)
	}
	// defer session.Close()

	importer := new(MongoBackend)

	importer.mongo = session.DB("blocks")
	importer.ethclient = ethclient
	importer.createIndexes()
	return importer
}

func (self *MongoBackend) ImportBlock(block *types.Block) {
	blockNumber := block.Header().Number
	txAmount := uint64(len(block.Transactions()))
	log.Debug().Str("BlockNumber", blockNumber.String()).Str("Hash", block.Hash().Hex()).Str("ParentHash", block.ParentHash().Hex()).Msg("Importing block")
	b := &models.Block{Number: blockNumber.Int64(),
		GasLimit:   int(block.Header().GasLimit),
		BlockHash:  block.Hash().Hex(),
		CreatedAt:  time.Unix(block.Time().Int64(), 0),
		ParentHash: block.ParentHash().Hex(),
		TxHash:     block.Header().TxHash.Hex(),
		GasUsed:    strconv.Itoa(int(block.Header().GasUsed)),
		Nonce:      string(block.Nonce()),
		Miner:      block.Coinbase().Hex(),
		TxCount:    int(txAmount)}
	for _, tx := range block.Transactions() {
		self.ImportTx(tx, block)
	}
	log.Debug().Interface("Block", b)
	_, err := self.mongo.C("Blocks").Upsert(bson.M{"number": b.Number}, b)
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

	_, err = self.mongo.C("ActiveAddress").Upsert(bson.M{"address": block.Coinbase().Hex()}, &models.ActiveAddress{Address: block.Coinbase().Hex(), UpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

}
func (self *MongoBackend) ImportTx(tx *types.Transaction, block *types.Block) {
	log.Debug().Msg("Importing tx" + tx.Hash().Hex())
	transaction := self.parseTx(tx, block)
	_, err := self.mongo.C("Transactions").Upsert(bson.M{"tx_hash": tx.Hash().String()}, transaction)
	if err != nil {
		log.Fatal().Err(err).Msg("importTx")
	}
	_, err = self.mongo.C("ActiveAddress").Upsert(bson.M{"address": transaction.From}, &models.ActiveAddress{Address: transaction.From, UpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

	_, err = self.mongo.C("ActiveAddress").Upsert(bson.M{"address": transaction.To}, &models.ActiveAddress{Address: transaction.To, UpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}
}
func (self *MongoBackend) NeedReloadBlock(blockNumber int64) bool {
	block := self.GetBlockByNumber(blockNumber)
	if block == nil {
		log.Debug().Msg("Checking parent - main block not found")
		return true
	}
	parentBlockNumber := (block.Number - 1)
	parentBlock := self.GetBlockByNumber(parentBlockNumber)
	if parentBlock != nil {
		log.Debug().Str("ParentHash", block.ParentHash).Str("Hash from parent", parentBlock.BlockHash).Int64("BlockNumber", block.Number).Int64("ParentNumber", parentBlock.Number).Msg("Checking parent")
	}
	return parentBlock == nil || (block.ParentHash != parentBlock.BlockHash)

}

func (self *MongoBackend) TransactionsConsistent(blockNumber int64) bool {
	block := self.GetBlockByNumber(blockNumber)
	if block != nil {
		transactionCounter, err := self.mongo.C("Transactions").Find(bson.M{"block_number": blockNumber}).Count()
		log.Debug().Int("Transactions in block", block.TxCount).Int("Num of transactions in db", transactionCounter).Msg("TransactionsConsistent")
		if err != nil {
			log.Fatal().Err(err).Msg("TransactionsConsistent")
		}

		if transactionCounter != block.TxCount {
			log.Fatal().Err(err).Msg("TransactionsConsistent")
		}

		return transactionCounter == block.TxCount
	}
	return true
}

func (self *MongoBackend) ImportAddress(address string, balance *big.Int) {
	log.Debug().Str("address", address).Str("balance", balance.String()).Msg("Updating address")
	_, err := self.mongo.C("Addresses").Upsert(bson.M{"address": address}, &models.Address{Address: address, Balance: balance.String(), LastUpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importAddress")
	}

}

func (self *MongoBackend) GetBlockByNumber(blockNumber int64) *models.Block {
	var c models.Block
	err := self.mongo.C("Blocks").Find(bson.M{"number": blockNumber}).One(&c)
	if err != nil {
		log.Debug().Int64("Block", blockNumber).Err(err).Msg("GetBlockByNumber")
		return nil
	}
	return &c
}

func (self *MongoBackend) GetLatestsBlocks(numOfBlocks int) []*models.Block {
	var blocks []*models.Block
	err := self.mongo.C("Blocks").Find(nil).Sort("-number").Limit(numOfBlocks).All(&blocks)
	if err != nil {
		log.Debug().Int("Block", numOfBlocks).Err(err).Msg("GetLatestsBlocks")
		return nil
	}
	return blocks
}

func (self *MongoBackend) GetActiveAdresses(fromDate time.Time) *[]models.ActiveAddress {
	var addresses []models.ActiveAddress

	err := self.mongo.C("ActiveAddress").Find(bson.M{"updated_at": bson.M{"$gte": fromDate}}).All(&addresses)
	if err != nil {
		log.Debug().Err(err).Msg("GetActiveAdresses")
	}
	return &addresses
}

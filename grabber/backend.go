package main

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

type ImportMaster struct {
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

func (self *ImportMaster) parseTx(tx *types.Transaction, block *types.Block) *models.Transaction {
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
func (self *ImportMaster) createIndexes() {
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
func NewImporter(ethclient *ethclient.Client) *ImportMaster {

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

	importer := new(ImportMaster)

	importer.mongo = session.DB("blocks")
	importer.ethclient = ethclient
	importer.createIndexes()
	return importer
}

func (self *ImportMaster) importBlock(block *types.Block) {
	blockNumber := block.Header().Number
	txAmount := uint64(len(block.Transactions()))
	log.Info().Str("BlockNumber", blockNumber.String()).Str("Hash", block.Hash().Hex()).Str("ParentHash", block.ParentHash().Hex()).Msg("Importing block")
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
		self.importTx(tx, block)
	}
	log.Info().Interface("Block", b)
	// _, err := self.fs.Collection("Blocks").Doc(strconv.Itoa(blockNumber)).Set(self.ctx, b)
	_, err := self.mongo.C("Blocks").Upsert(bson.M{"number": b.Number}, b)
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

	// _, err = self.fs.Collection("ActiveAddress").Doc(block.Coinbase().Hex()).Set(self.ctx, &models.ActiveAddress{time.Now()})
	_, err = self.mongo.C("ActiveAddress").Upsert(bson.M{"address": block.Coinbase().Hex()}, &models.ActiveAddress{Address: block.Coinbase().Hex(), UpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

}
func (self *ImportMaster) importTx(tx *types.Transaction, block *types.Block) {
	log.Info().Msg("Importing tx" + tx.Hash().Hex())
	transaction := self.parseTx(tx, block)
	// _, err := self.fs.Collection("Transactions").Doc(tx.Hash().String()).Set(self.ctx, transaction)
	_, err := self.mongo.C("Transactions").Upsert(bson.M{"tx_hash": tx.Hash().String()}, transaction)
	if err != nil {
		log.Fatal().Err(err).Msg("importTx")
	}

	// _, err = self.fs.Collection("ActiveAddress").Doc(transaction.From).Set(self.ctx, &models.ActiveAddress{time.Now()})
	_, err = self.mongo.C("ActiveAddress").Upsert(bson.M{"address": transaction.From}, &models.ActiveAddress{Address: transaction.From, UpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

	// _, err = self.fs.Collection("ActiveAddress").Doc(transaction.To).Set(self.ctx, &models.ActiveAddress{time.Now()})
	_, err = self.mongo.C("ActiveAddress").Upsert(bson.M{"address": transaction.To}, &models.ActiveAddress{Address: transaction.To, UpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}
}
func (self *ImportMaster) needReloadBlock(blockNumber int64) bool {
	block := self.GetBlocksByNumber(blockNumber)
	if block == nil {
		log.Info().Msg("Checking parent - main block not found")
		return true
	}
	parentBlockNumber := (block.Number - 1)
	parentBlock := self.GetBlocksByNumber(parentBlockNumber)
	if parentBlock != nil {
		log.Info().Str("ParentHash", block.ParentHash).Str("Hash from parent", parentBlock.BlockHash).Int64("BlockNumber", block.Number).Int64("ParentNumber", parentBlock.Number).Msg("Checking parent")
	}
	return parentBlock == nil || (block.ParentHash != parentBlock.BlockHash)

}

func (self *ImportMaster) TransactionsConsistent(blockNumber int64) bool {
	block := self.GetBlocksByNumber(blockNumber)
	if block != nil {
		transactionCounter, err := self.mongo.C("Transactions").Find(bson.M{"block_number": blockNumber}).Count()
		log.Info().Int("Transactions in block", block.TxCount).Int("Num of transactions in db", transactionCounter).Msg("TransactionsConsistent")
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

func (self *ImportMaster) importAddress(address string, balance *big.Int) {
	log.Info().Str("address", address).Str("balance", balance.String()).Msg("Updating address")
	// _, err := self.fs.Collection("Addresses").Doc(address).Set(self.ctx, &models.Address{address, "", balance.String(), time.Now()})
	_, err := self.mongo.C("Addresses").Upsert(bson.M{"address": address}, &models.Address{Address: address, Balance: balance.String(), LastUpdatedAt: time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importAddress")
	}

}

func (self *ImportMaster) GetBlocksByNumber(blockNumber int64) *models.Block {
	var c models.Block
	err := self.mongo.C("Blocks").Find(bson.M{"number": blockNumber}).One(&c)
	if err != nil {
		log.Info().Int64("Block", blockNumber).Err(err).Msg("GetBlocksByNumber")
		return nil
	}
	// dsnap, err := self.fs.Collection("Blocks").Doc(blockNumber).Get(self.ctx)
	// if err != nil {
	// 	// log.Info().Err(err).Msg("GetBlocksByNumber")
	// 	return nil
	// }
	// var c models.Block
	// dsnap.DataTo(&c)
	return &c
}

func (self *ImportMaster) GetActiveAdresses(fromDate time.Time) *[]models.ActiveAddress {
	var addresses []models.ActiveAddress
	// iter := self.fs.Collection("ActiveAddress").Where("updated_at", ">", fromDate).Documents(self.ctx)

	err := self.mongo.C("ActiveAddress").Find(bson.M{"updated_at": bson.M{"$gte": fromDate}}).All(&addresses)
	if err != nil {
		log.Info().Err(err).Msg("GetActiveAdresses")
	}
	return &addresses
}

package main

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"

	"github.com/gochain-io/explorer/api/models"
	"github.com/gochain-io/gochain/core/types"
)

type ImportMaster struct {
	fs  *firestore.Client
	ctx context.Context
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
	from, err := types.Sender(self.ctx, types.HomesteadSigner{}, tx)
	if err != nil {
		log.Fatal().Err(err).Msg("parseTx")
	}
	// TxHash,	To,	From,	Amount,	Price,	GasLimit,	BlockNumber,	Nonce,	BlockHash,	CreatedAt,
	txx := &models.Transaction{tx.Hash().Hex(), tx.To().Hex(), from.Hex(),
		tx.Value().String(), tx.GasPrice().String(), strconv.Itoa(int(tx.Gas())),
		block.Number().String(), string(tx.Nonce()),
		block.Hash().Hex(), time.Unix(block.Time().Int64(), 0)}
	return txx
}

func NewImporter() *ImportMaster {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "gochain-explorer"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatal().Err(err).Msg("NewImporter")
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("NewImporter")
	}
	importer := new(ImportMaster)
	importer.fs = client
	importer.ctx = ctx
	return importer
}

func (self *ImportMaster) importBlock(block *types.Block) {
	blockNumber := int(block.Header().Number.Int64())
	txAmount := uint64(len(block.Transactions()))
	log.Info().Int("BlockNumber", blockNumber).Str("Hash", block.Hash().Hex()).Str("ParentHash", block.ParentHash().Hex()).Msg("Importing block")
	// Number,	GasLimit,	BlockHash,	CreatedAt,	ParentHash,	TxHash,	GasUsed,	Nonce,	Miner,
	b := &models.Block{blockNumber,
		int(block.Header().GasLimit), block.Hash().Hex(),
		time.Unix(block.Time().Int64(), 0), block.ParentHash().Hex(), block.Header().TxHash.Hex(),
		strconv.Itoa(int(block.Header().GasUsed)), string(block.Nonce()),
		block.Coinbase().Hex(), int(txAmount)}
	for _, tx := range block.Transactions() {
		self.importTx(tx, block)
	}
	log.Info().Interface("Block", b)
	_, err := self.fs.Collection("Blocks").Doc(strconv.Itoa(blockNumber)).Set(self.ctx, b)
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

	_, err = self.fs.Collection("ActiveAddress").Doc(block.Coinbase().Hex()).Set(self.ctx, &models.ActiveAddress{time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

}
func (self *ImportMaster) importTx(tx *types.Transaction, block *types.Block) {
	log.Info().Msg("Importing tx" + tx.Hash().Hex())
	transaction := self.parseTx(tx, block)
	_, err := self.fs.Collection("Transactions").Doc(tx.Hash().String()).Set(self.ctx, transaction)
	if err != nil {
		log.Fatal().Err(err).Msg("importTx")
	}

	_, err = self.fs.Collection("ActiveAddress").Doc(transaction.From).Set(self.ctx, &models.ActiveAddress{time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}

	_, err = self.fs.Collection("ActiveAddress").Doc(transaction.To).Set(self.ctx, &models.ActiveAddress{time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importBlock")
	}
}
func (self *ImportMaster) needReloadBlock(block *types.Block) bool {
	parentBlockNumber := strconv.Itoa(int(block.Header().Number.Int64()) - 1)
	parentBlock := self.GetBlocksByNumber(parentBlockNumber)
	if parentBlock != nil {
		log.Info().Str("ParentHash", block.ParentHash().Hex()).Str("Hash from parent", parentBlock.BlockHash).Str("BlockNumber", block.Header().Number.String()).Int("ParentNumber", parentBlock.Number).Msg("Checking parent")
	}
	return parentBlock == nil || (block.ParentHash().Hex() != parentBlock.BlockHash)

}
func (self *ImportMaster) importAddress(address string, balance *big.Int) {
	log.Info().Str("address", address).Str("balance", balance.String()).Msg("Updating address")
	_, err := self.fs.Collection("Addresses").Doc(address).Set(self.ctx, &models.Address{address, "", balance.String(), time.Now()})
	if err != nil {
		log.Fatal().Err(err).Msg("importAddress")
	}

}

func (self *ImportMaster) GetBlocksByNumber(blockNumber string) *models.Block {
	dsnap, err := self.fs.Collection("Blocks").Doc(blockNumber).Get(self.ctx)
	if err != nil {
		// log.Info().Err(err).Msg("GetBlocksByNumber")
		return nil
	}
	var c models.Block
	dsnap.DataTo(&c)
	return &c
}

func (self *ImportMaster) GetActiveAdresses(fromDate time.Time) *[]string {
	var addresses []string
	iter := self.fs.Collection("ActiveAddress").Where("updated_at", ">", fromDate).Documents(self.ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Info().Err(err).Msg("GetActiveAdresses")
		}
		// log.Info().Str("DOC:", doc.Ref.ID).Msg("GetActiveAdresses")
		addresses = append(addresses, doc.Ref.ID)
	}
	return &addresses
}

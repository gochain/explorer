package main

import (
	"context"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"

	"cloud.google.com/go/datastore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gochain-io/explorer/api/models"
	"golang.org/x/oauth2/google"
)

type ImportMaster struct {
	ds  *datastore.Client
	ctx context.Context
}

func (self *ImportMaster) parseTx(tx *types.Transaction, block *types.Block) *models.Transaction {
	var from = ""
	if msg, err := tx.AsMessage(types.HomesteadSigner{}); err != nil {
		from = msg.From().Hex()
	} else {
		utils.Fatalf("Could not parse from address: %v", err)
	}
	// TxHash,	To,	From,	Amount,	Price,	GasLimit,	BlockNumber,	Nonce,	BlockHash,	CreatedAt,
	txx := &models.Transaction{tx.Hash().Hex(), tx.To().Hex(), from,
		tx.Value().String(), tx.GasPrice().String(), strconv.Itoa(int(tx.Gas())),
		block.Number().String(), string(tx.Nonce()),
		block.Header().Hash().Hex(), time.Unix(block.Time().Int64(), 0)}
	return txx
}

func NewImporter() *ImportMaster {
	var ds *datastore.Client
	ctx := context.Background()
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/datastore")
	if err != nil {
		log.Fatal().Err(err).Msg("oops")
	}

	ds, err = datastore.NewClient(ctx, creds.ProjectID)
	if err != nil {
		log.Fatal().Err(err).Msg("oops")
	}

	importer := new(ImportMaster)
	importer.ds = ds
	importer.ctx = ctx
	return importer
}

func (self *ImportMaster) importBlock(block *types.Block) {
	blockHash := block.Header().Hash().Hex()
	txAmount := uint64(len(block.Transactions()))
	log.Info().Msg("Importing block " + blockHash + "Hash with " + string(txAmount) + "transactions")
	// Number,	GasLimit,	BlockHash,	CreatedAt,	ParentHash,	TxHash,	GasUsed,	Nonce,	Miner,
	b := &models.Block{int(block.Header().Number.Int64()),
		int(block.Header().GasLimit),
		blockHash, time.Unix(block.Time().Int64(), 0), block.ParentHash().Hex(), block.Header().TxHash.Hex(),
		strconv.Itoa(int(block.Header().GasUsed)), string(block.Nonce()),
		block.Coinbase().Hex(), int(txAmount)}

	for _, tx := range block.Transactions() {
		self.importTx(tx, block)
	}
	log.Info().Interface("Block", b)
	blockKey := datastore.NameKey("Blocks", blockHash, nil)
	if _, err := self.ds.Put(self.ctx, blockKey, b); err != nil {
		log.Fatal().Err(err).Msg("oops")
	}

}
func (self *ImportMaster) importTx(tx *types.Transaction, block *types.Block) {
	log.Info().Msg("Importing tx" + tx.Hash().Hex())
	txKey := datastore.NameKey("Transactions", tx.Hash().String(), nil)
	if _, err := self.ds.Put(self.ctx, txKey, self.parseTx(tx, block)); err != nil {
		log.Fatal().Err(err).Msg("oops")
	}
}

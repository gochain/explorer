package main

import (
	"context"
	"math/big"
	"strconv"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"

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

func appendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
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

func (self *ImportMaster) importBlockIfNotExists(block *types.Block) {
	blockNumber := int(block.Header().Number.Int64())
	txAmount := uint64(len(block.Transactions()))
	log.Info().Msg("Importing block " + strconv.Itoa(blockNumber) + "Hash with " + string(txAmount) + "transactions")
	// Number,	GasLimit,	BlockHash,	CreatedAt,	ParentHash,	TxHash,	GasUsed,	Nonce,	Miner,
	b := &models.Block{blockNumber,
		int(block.Header().GasLimit), block.Header().Hash().Hex(),
		time.Unix(block.Time().Int64(), 0), block.ParentHash().Hex(), block.Header().TxHash.Hex(),
		strconv.Itoa(int(block.Header().GasUsed)), string(block.Nonce()),
		block.Coinbase().Hex(), int(txAmount)}

	for _, tx := range block.Transactions() {
		self.importTx(tx, block)
	}
	log.Info().Interface("Block", b)
	blockKey := datastore.NameKey("Blocks", strconv.Itoa(blockNumber), nil)
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
func (self *ImportMaster) needReloadBlock(block *types.Block) bool {
	parentBlockNumber := strconv.Itoa(int(block.Header().Number.Int64()) - 1)
	parentBlock := *self.GetBlocksByNumber(parentBlockNumber)
	log.Info().Int("Number of blocks", len(parentBlock)).Msg("Checking parent")
	if len(parentBlock) == 1 {
		log.Info().Str("ParentHash", block.ParentHash().Hex()).Str("Hash from parent", parentBlock[0].BlockHash).Msg("Checking parent")
		log.Info().Str("BlockNumber", block.Header().Number.String()).Int("ParentNumber", parentBlock[0].Number).Msg("Checking parent")
	}
	return len(parentBlock) < 1 || (len(parentBlock) == 1 && block.ParentHash().Hex() != parentBlock[0].BlockHash)

}
func (self *ImportMaster) GetBlocksByNumber(blockNumber string) *[]models.Block {
	var blocks []models.Block
	key := datastore.NameKey("Blocks", blockNumber, nil)
	query := datastore.NewQuery("Blocks").Filter("__key__ =", key)
	it := self.ds.Run(self.ctx, query)
	for {
		var block models.Block
		_, err := it.Next(&block)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("oops")
		}
		blocks = append(blocks, block)
		log.Info().Str("blockNumber", blockNumber).Msg("Got block")
	}
	return &blocks
}

func (self *ImportMaster) importAddress(address string, balance *big.Int) {
	log.Info().Str("address", address).Str("balance", balance.String()).Msg("Updating address")
	adKey := datastore.NameKey("Addresses", address, nil)
	if _, err := self.ds.Put(self.ctx, adKey, &models.Address{address, "", balance.String(), time.Now()}); err != nil {
		log.Fatal().Err(err).Msg("oops")
	}
}

func (self *ImportMaster) GetActiveAdresses(fromDate time.Time) *[]string {
	var addresses []string
	query := datastore.NewQuery("Blocks").
		// Filter("num >", 0)
		DistinctOn("miner")
	it := self.ds.Run(self.ctx, query)
	for {
		var block models.Block
		_, err := it.Next(&block)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("distinct on blocks")
		}
		addresses = appendIfMissing(addresses, block.Miner)
		log.Info().Str("miner", block.Miner).Msg("Got miner")
	}

	query = datastore.NewQuery("Transactions").
		// Filter("created_at >", fromDate)
		DistinctOn("from")
	it = self.ds.Run(self.ctx, query)
	for {
		var tx models.Transaction
		_, err := it.Next(&tx)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("oops")
		}
		addresses = appendIfMissing(addresses, tx.From)
		log.Info().Str("tx_from", tx.From).Msg("Got from tx")
	}
	query = datastore.NewQuery("Transactions").
		// Filter("created_at >", fromDate)
		DistinctOn("to")
	it = self.ds.Run(self.ctx, query)
	for {
		var tx models.Transaction
		_, err := it.Next(&tx)
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("oops")
		}
		addresses = appendIfMissing(addresses, tx.To)
		log.Info().Str("tx_to", tx.To).Msg("Got to tx")
	}
	return &addresses
}

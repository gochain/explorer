package backend

import (
	"math/big"
	"time"

	"github.com/gochain-io/explorer/api/models"
	"github.com/gochain-io/gochain/core/types"
	"github.com/gochain-io/gochain/ethclient"
	"github.com/rs/zerolog/log"
)

type Backend struct {
	mongo             *MongoBackend
	ethClient         *ethclient.Client
	extendedEthClient *EthRPC
}

func NewBackend(mongoUrl, rpcUrl string) *Backend {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create eth client")
	}
	exClient := NewEthClient(rpcUrl)
	mongoBackend := NewMongoClient(mongoUrl, rpcUrl)
	importer := new(Backend)
	importer.ethClient = client
	importer.extendedEthClient = exClient
	importer.mongo = mongoBackend
	return importer
}

//METHODS USED IN API
func (self *Backend) BalanceAt(address, block string) (*big.Int, error) {
	return self.extendedEthClient.ethGetBalance(address, block)
}
func (self *Backend) TotalSupply() (*big.Int, error) {
	return self.extendedEthClient.ethTotalSupply()
}
func (self *Backend) CirculatingSupply() (*big.Int, error) {
	return self.extendedEthClient.circulatingSupply()
}
func (self *Backend) GetStats() *models.Stats {
	return self.mongo.getStats()
}
func (self *Backend) GetRichlist(skip, limit int) []*models.Address {
	return self.mongo.getRichlist(skip, limit)
}
func (self *Backend) GetAddressByHash(hash string) *models.Address {
	return self.mongo.getAddressByHash(hash)
}
func (self *Backend) GetTransactionByHash(hash string) *models.Transaction {
	return self.mongo.getTransactionByHash(hash)
}
func (self *Backend) GetTransactionList(address string) []*models.Transaction {
	return self.mongo.getTransactionList(address)
}
func (self *Backend) GetLatestsBlocks(skip, limit int) []*models.Block {
	return self.mongo.getLatestsBlocks(skip, limit)
}
func (self *Backend) GetBlockByNumber(number int64) *models.Block {
	return self.mongo.getBlockByNumber(number)
}

//METHODS USED IN GRABBER
//
func (self *Backend) ImportBlock(block *types.Block) {
	self.mongo.importBlock(block)
}
func (self *Backend) NeedReloadBlock(blockNumber int64) bool {
	return self.mongo.needReloadBlock(blockNumber)
}
func (self *Backend) TransactionsConsistent(blockNumber int64) bool {
	return self.mongo.transactionsConsistent(blockNumber)
}
func (self *Backend) GetActiveAdresses(fromDate time.Time) []*models.ActiveAddress {
	return self.mongo.getActiveAdresses(fromDate)
}
func (self *Backend) ImportAddress(address string, balance *big.Int) {
	self.mongo.importAddress(address, balance)
}

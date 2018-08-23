package backend

import (
	"context"
	"math/big"
	"time"

	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/gochain/common"
	"github.com/gochain-io/gochain/core/types"
	"github.com/gochain-io/gochain/ethclient"
	"github.com/rs/zerolog/log"
)

type Backend struct {
	mongo             *MongoBackend
	ethClient         *ethclient.Client
	extendedEthClient *EthRPC
	tokenBalance      *TokenBalance
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
	importer.tokenBalance = NewTokenBalanceClient(rpcUrl)
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
	address := self.mongo.getAddressByHash(hash)
	if address == nil {
		balance, err := self.ethClient.BalanceAt(context.Background(), common.HexToAddress(hash), nil)
		if err != nil {
			log.Info().Err(err).Str("address", hash).Msg("cannot get address information neither from eth or db")
			return nil
		}
		address = self.mongo.importAddress(hash, balance, "", "", false, false)
	}
	return address
}
func (self *Backend) GetTransactionByHash(hash string) *models.Transaction {
	return self.mongo.getTransactionByHash(hash)
}
func (self *Backend) GetTransactionList(address string) []*models.Transaction {
	return self.mongo.getTransactionList(address)
}
func (self *Backend) GetTokenHoldersList(contractAddress string) []*models.TokenHolder {
	return self.mongo.getTokenHoldersList(contractAddress)
}
func (self *Backend) GetLatestsBlocks(skip, limit int) []*models.Block {
	return self.mongo.getLatestsBlocks(skip, limit)
}
func (self *Backend) GetBlockByNumber(number int64) *models.Block {
	block := self.mongo.getBlockByNumber(number)
	if block == nil {
		log.Info().Int64("blockNumber", number).Msg("cannot get block from db, importing it")
		blockEth, err := self.ethClient.BlockByNumber(context.Background(), big.NewInt(number))
		if err != nil {
			log.Info().Err(err).Int64("blockNumber", number).Msg("cannot get block from eth and db")
			return nil
		}
		block = self.ImportBlock(blockEth)
	}
	return block
}

//METHODS USED IN GRABBER
func (self *Backend) GetTokenBalance(contract, wallet string) (*tokenBalance, error) {
	return self.tokenBalance.GetTokenBalance(contract, wallet)
}
func (self *Backend) ImportBlock(block *types.Block) *models.Block {
	return self.mongo.importBlock(block)
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
func (self *Backend) ImportAddress(address string, balance *big.Int, tokenName, tokenSymbol string, contract, go20 bool) *models.Address {
	return self.mongo.importAddress(address, balance, tokenName, tokenSymbol, contract, go20)
}
func (self *Backend) ImportTokenHolder(contractAddress, tokenHolderAddress string, balance *big.Int, tokenName, tokenSymbol string) *models.TokenHolder {
	return self.mongo.importTokenHolder(contractAddress, tokenHolderAddress, balance, tokenName, tokenSymbol)
}

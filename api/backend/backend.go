package backend

import (
	"math/big"

	"github.com/gochain-io/gochain/ethclient"
	"github.com/rs/zerolog/log"
	mgo "gopkg.in/mgo.v2"
)

type Backend struct {
	mongo             *mgo.Database
	ethClient         *ethclient.Client
	extendedEthClient *EthRPC
}

func NewBackend(rpcUrl string) *Backend {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("main")
	}
	exClient := NewEthClient(rpcUrl)
	Host := []string{
		"127.0.0.1:27017",
	}
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: Host,
	})
	if err != nil {
		panic(err)
	}

	importer := new(Backend)

	importer.mongo = session.DB("blocks")
	importer.ethClient = client
	importer.extendedEthClient = exClient
	importer.createIndexes()

	return importer
}

func (self *Backend) BalanceAt(address, block string) (*big.Int, error) {
	return self.extendedEthClient.EthGetBalance(address, block)
}

func (self *Backend) TotalSupply() (*big.Int, error) {
	return self.extendedEthClient.EthTotalSupply()
}

func (self *Backend) CirculatingSupply() (*big.Int, error) {
	return self.extendedEthClient.CirculatingSupply()
}

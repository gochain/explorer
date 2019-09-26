package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"time"

	"go.uber.org/zap"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/gochain/v3/common"
	"github.com/gochain-io/gochain/v3/core/types"
	"github.com/gochain-io/gochain/v3/goclient"
)

var wei = big.NewInt(1000000000000000000)

const defaultFetchLimit = 100
const defaultSkip = 0

type MongoBackend struct {
	host         string
	mongo        *mgo.Database
	mongoSession *mgo.Session
	goClient     *goclient.Client
	Lgr          *zap.Logger
}

// New create new rpc client with given url
func NewMongoClient(host, rpcUrl, dbName string, lgr *zap.Logger) *MongoBackend {
	client, err := goclient.Dial(rpcUrl)
	if err != nil {
		lgr.Fatal("main", zap.Error(err))
	}
	Host := []string{
		host,
	}
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs: Host,
	})
	if err != nil {
		panic(err)
	}
	importer := new(MongoBackend)
	importer.Lgr = lgr
	importer.mongoSession = session
	importer.mongo = session.DB(dbName)
	importer.goClient = client
	importer.createIndexes()

	return importer

}
func (self *MongoBackend) PingDB() error {
	return self.mongoSession.Ping()
}
func (self *MongoBackend) parseTx(tx *types.Transaction, block *types.Block) *models.Transaction {
	from, err := self.goClient.TransactionSender(context.Background(), tx, block.Header().Hash(), 0)
	if err != nil {
		self.Lgr.Fatal("parseTx", zap.Error(err))
	}
	gas := tx.Gas()
	to := ""
	if tx.To() != nil {
		to = tx.To().Hex()
	}
	self.Lgr.Debug("parseTx", zap.String("TX:", tx.Hash().Hex()))
	InputDataEmpty := hex.EncodeToString(tx.Data()[:]) == ""
	return &models.Transaction{TxHash: tx.Hash().Hex(),
		To:              to,
		From:            from.Hex(),
		Value:           tx.Value().String(),
		GasPrice:        tx.GasPrice().String(),
		ReceiptReceived: false,
		GasLimit:        tx.Gas(),
		BlockNumber:     block.Number().Int64(),
		GasFee:          new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(gas))).String(),
		Nonce:           tx.Nonce(),
		BlockHash:       block.Hash().Hex(),
		CreatedAt:       time.Unix(block.Time().Int64(), 0),
		InputData:       hex.EncodeToString(tx.Data()[:]),
		InputDataEmpty:  InputDataEmpty,
	}
}
func (self *MongoBackend) parseBlock(block *types.Block) *models.Block {
	var transactions []string
	for _, tx := range block.Transactions() {
		transactions = append(transactions, tx.Hash().Hex())
	}
	nonceBool := false
	if block.Nonce() == 0xffffffffffffffff {
		nonceBool = true
	}
	return &models.Block{Number: block.Header().Number.Int64(),
		GasLimit:   int(block.Header().GasLimit),
		BlockHash:  block.Hash().Hex(),
		CreatedAt:  time.Unix(block.Time().Int64(), 0),
		ParentHash: block.ParentHash().Hex(),
		TxHash:     block.Header().TxHash.Hex(),
		GasUsed:    strconv.Itoa(int(block.Header().GasUsed)),
		NonceBool:  &nonceBool,
		Miner:      block.Coinbase().Hex(),
		TxCount:    int(uint64(len(block.Transactions()))),
		Difficulty: block.Difficulty().Int64(),
		// TotalDifficulty: block.DeprecatedTd().Int64(), # deprecated https://github.com/ethereum/go-ethereum/blob/master/core/types/block.go#L154
		Sha3Uncles: block.UncleHash().Hex(),
		ExtraData:  string(block.Extra()[:]),
		// Transactions: transactions,
	}
}

func (self *MongoBackend) createIndexes() {
	err := self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"tx_hash"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"block_number"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"from", "created_at", "input_data_empty"}, Background: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"to", "created_at", "input_data_empty"}, Background: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"-created_at"}, Background: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Transactions").EnsureIndex(mgo.Index{Key: []string{"contract_address"}, Background: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Blocks").EnsureIndex(mgo.Index{Key: []string{"number"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Blocks").EnsureIndex(mgo.Index{Key: []string{"-number"}, Background: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Blocks").EnsureIndex(mgo.Index{Key: []string{"miner"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Blocks").EnsureIndex(mgo.Index{Key: []string{"created_at", "miner"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Blocks").EnsureIndex(mgo.Index{Key: []string{"hash"}, Background: true, Sparse: true})
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
	err = self.mongo.C("Addresses").EnsureIndex(mgo.Index{Key: []string{"address"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Addresses").EnsureIndex(mgo.Index{Key: []string{"contract"}, Background: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Addresses").EnsureIndex(mgo.Index{Key: []string{"-balance_float", "address"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("TokensHolders").EnsureIndex(mgo.Index{Key: []string{"contract_address", "token_holder_address"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("TokensHolders").EnsureIndex(mgo.Index{Key: []string{"token_holder_address"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("TokensHolders").EnsureIndex(mgo.Index{Key: []string{"balance_int"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("InternalTransactions").EnsureIndex(mgo.Index{Key: []string{"contract_address", "from_address", "to_address"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("InternalTransactions").EnsureIndex(mgo.Index{Key: []string{"from_address", "block_number"}, Background: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("InternalTransactions").EnsureIndex(mgo.Index{Key: []string{"to_address", "block_number"}, Background: true})
	if err != nil {
		panic(err)
	}
	err = self.mongo.C("InternalTransactions").EnsureIndex(mgo.Index{Key: []string{"transaction_hash"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("InternalTransactions").EnsureIndex(mgo.Index{Key: []string{"block_number"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Stats").EnsureIndex(mgo.Index{Key: []string{"-updated_at"}, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}

	err = self.mongo.C("Contracts").EnsureIndex(mgo.Index{Key: []string{"address"}, Unique: true, DropDups: true, Background: true, Sparse: true})
	if err != nil {
		panic(err)
	}
}

func (self *MongoBackend) importBlock(block *types.Block) *models.Block {
	self.Lgr.Debug("Importing block", zap.String("BlockNumber", block.Header().Number.String()), zap.String("Hash", block.Hash().Hex()), zap.String("ParentHash", block.ParentHash().Hex()))
	b := self.parseBlock(block)
	_, err := self.mongo.C("Blocks").Upsert(bson.M{"number": b.Number}, b)
	if err != nil {
		self.Lgr.Fatal("importBlock", zap.Error(err))
	}
	_, err = self.mongo.C("Transactions").RemoveAll(bson.M{"block_number": b.Number}) //deleting all txs belong to this block if any exist
	if err != nil {
		self.Lgr.Fatal("importBlock", zap.Error(err))
	}
	for _, tx := range block.Transactions() {
		self.importTx(tx, block)
	}
	self.UpdateActiveAddress(block.Coinbase().Hex())
	return b

}

func (self *MongoBackend) UpdateActiveAddress(address string) {
	_, err := self.mongo.C("ActiveAddress").Upsert(bson.M{"address": address}, &models.ActiveAddress{Address: address, UpdatedAt: time.Now()})
	if err != nil {
		self.Lgr.Fatal("UpdateActiveAddress", zap.Error(err))
	}
}

func (self *MongoBackend) importTx(tx *types.Transaction, block *types.Block) {
	self.Lgr.Debug("Importing", zap.String("tx", tx.Hash().Hex()))
	transaction := self.parseTx(tx, block)

	toAddress := transaction.To
	if transaction.To == "" {
		self.Lgr.Info("Hash doesn't have an address", zap.String("hash", transaction.TxHash))
		receipt, err := self.goClient.TransactionReceipt(context.Background(), tx.Hash())
		if err == nil {
			contractAddress := receipt.ContractAddress.String()
			if contractAddress != "0x0000000000000000000000000000000000000000" {
				transaction.ContractAddress = contractAddress
			}
			transaction.Status = false
			if receipt.Status == 1 {
				transaction.Status = true
			}
			toAddress = transaction.ContractAddress
		} else {
			self.Lgr.Error("Cannot get a receipt in importTX", zap.Error(err), zap.String("hash", transaction.TxHash))
		}
	}

	_, err := self.mongo.C("Transactions").Upsert(bson.M{"tx_hash": tx.Hash().String()}, transaction)
	if err != nil {
		self.Lgr.Fatal("importTx", zap.Error(err))
	}

	self.UpdateActiveAddress(toAddress)
	self.UpdateActiveAddress(transaction.From)
}

// needReloadParent returns true if the parent block is missing or does not match the hash from this block number.
func (self *MongoBackend) needReloadParent(blockNumber int64) bool {
	block := self.getBlockByNumber(blockNumber)
	if block == nil {
		self.Lgr.Debug("Checking parent - main block not found")
		return true
	}
	parentBlockNumber := (block.Number - 1)
	parentBlock := self.getBlockByNumber(parentBlockNumber)
	if parentBlock != nil {
		self.Lgr.Debug("Checking parent", zap.String("ParentHash", block.ParentHash), zap.String("Hash from parent", parentBlock.BlockHash), zap.Int64("BlockNumber", block.Number), zap.Int64("ParentNumber", parentBlock.Number))
	}
	return parentBlock == nil || (block.ParentHash != parentBlock.BlockHash)

}

func (self *MongoBackend) transactionsConsistent(blockNumber int64) bool {
	block := self.getBlockByNumber(blockNumber)
	if block != nil {
		transactionCounter, err := self.mongo.C("Transactions").Find(bson.M{"block_number": blockNumber}).Count()
		self.Lgr.Debug("TransactionsConsistent", zap.Int("Transactions in block", block.TxCount), zap.Int("Num of transactions in db", transactionCounter))
		if err != nil {
			self.Lgr.Fatal("TransactionsConsistent", zap.Error(err))
		}
		return transactionCounter == block.TxCount
	}
	return true
}

func (self *MongoBackend) importAddress(address string, balance *big.Int, token *TokenDetails, contract bool, updatedAtBlock int64) *models.Address {
	balanceGoFloat, _ := new(big.Float).SetPrec(100).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(wei)).Float64() //converting to GO from wei
	balanceGoString := new(big.Rat).SetFrac(balance, wei).FloatString(18)
	self.Lgr.Debug("Updating address", zap.String("address", address), zap.String("precise balance", balanceGoString), zap.Float64("balance float", balanceGoFloat))
	tokenHoldersCounter, err := self.mongo.C("TokensHolders").Find(bson.M{"contract_address": address}).Count()
	if err != nil {
		self.Lgr.Fatal("importAddress", zap.Error(err))
	}

	internalTransactionsCounter, err := self.mongo.C("InternalTransactions").Find(bson.M{"contract_address": address}).Count()

	if err != nil {
		self.Lgr.Fatal("importAddress", zap.Error(err))
	}

	tokenTransactionsCounter, err := self.mongo.C("InternalTransactions").Find(bson.M{"$or": []bson.M{bson.M{"from_address": address}, bson.M{"to_address": address}}}).Count()
	if err != nil {
		self.Lgr.Fatal("importAddress", zap.Error(err))
	}

	addressM := &models.Address{Address: address,
		BalanceWei:     balance.String(),
		UpdatedAt:      time.Now(),
		UpdatedAtBlock: updatedAtBlock,
		TokenName:      token.Name,
		TokenSymbol:    token.Symbol,
		Decimals:       token.Decimals,
		TotalSupply:    token.TotalSupply.String(),
		Contract:       contract,
		ErcTypes:       token.Types,
		Interfaces:     token.Interfaces,
		BalanceFloat:   balanceGoFloat,
		BalanceString:  balanceGoString,
		// NumberOfTransactions:         transactionCounter,
		NumberOfTokenHolders:         tokenHoldersCounter,
		NumberOfInternalTransactions: internalTransactionsCounter,
		NumberOfTokenTransactions:    tokenTransactionsCounter,
	}
	_, err = self.mongo.C("Addresses").Upsert(bson.M{"address": address}, addressM)
	if err != nil {
		self.Lgr.Fatal("importAddress", zap.Error(err))
	}
	return addressM
}

func (self *MongoBackend) importTokenHolder(contractAddress, tokenHolderAddress string, token *TokenHolderDetails, address *models.Address) *models.TokenHolder {
	balanceInt := new(big.Int).Div(token.Balance, wei) //converting to GO from wei
	self.Lgr.Info("Updating token holder", zap.String("contractAddress", contractAddress), zap.String("tokenAddress", tokenHolderAddress), zap.String("balance", token.Balance.String()), zap.String("Balance int", balanceInt.String()))
	tokenHolder := &models.TokenHolder{
		TokenName:          address.TokenName,
		TokenSymbol:        address.TokenSymbol,
		ContractAddress:    contractAddress,
		TokenHolderAddress: tokenHolderAddress,
		Balance:            token.Balance.String(),
		UpdatedAt:          time.Now(),
		BalanceInt:         balanceInt.Int64()}
	_, err := self.mongo.C("TokensHolders").Upsert(bson.M{"contract_address": contractAddress, "token_holder_address": tokenHolderAddress}, tokenHolder)
	if err != nil {
		self.Lgr.Fatal("importTokenHolder", zap.Error(err))
	}
	return tokenHolder

}

func (self *MongoBackend) importInternalTransaction(contractAddress string, transferEvent TransferEvent, createdAt time.Time) *models.InternalTransaction {

	internalTransaction := &models.InternalTransaction{
		Contract:        contractAddress,
		From:            transferEvent.From.String(),
		To:              transferEvent.To.String(),
		Value:           transferEvent.Value.String(),
		BlockNumber:     transferEvent.BlockNumber,
		TransactionHash: transferEvent.TransactionHash,
		CreatedAt:       createdAt,
		UpdatedAt:       time.Now(),
	}
	_, err := self.mongo.C("InternalTransactions").Upsert(bson.M{"transaction_hash": transferEvent.TransactionHash}, internalTransaction)
	if err != nil {
		self.Lgr.Fatal("importInternalTransaction", zap.Error(err))
	}
	return internalTransaction
}

func (self *MongoBackend) importContract(contractAddress string, byteCode string) {
	//https://stackoverflow.com/questions/43278696/golang-mgo-insert-or-update-not-working-as-expected/43278832
	_, err := self.mongo.C("Contracts").Upsert(bson.M{"address": contractAddress}, bson.M{"$set": bson.M{"address": contractAddress, "byte_code": byteCode, "created_at": time.Now()}})
	if err != nil {
		self.Lgr.Fatal("importContract", zap.Error(err))
	}
}

func (self *MongoBackend) getBlockByNumber(blockNumber int64) *models.Block {
	var c models.Block
	err := self.mongo.C("Blocks").Find(bson.M{"number": blockNumber}).Select(bson.M{"transactions": 0}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		self.Lgr.Error("Failed to get block by number", zap.Int64("block", blockNumber), zap.Error(err))
		return nil
	}
	return &c
}

func (self *MongoBackend) getBlockByHash(blockHash string) *models.Block {
	var c models.Block
	err := self.mongo.C("Blocks").Find(bson.M{"hash": blockHash}).Select(bson.M{"transactions": 0}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		self.Lgr.Error("Failed to get block by hash", zap.String("block", blockHash), zap.Error(err))
		return nil
	}
	return &c
}

func (self *MongoBackend) getBlockTransactionsByNumber(blockNumber int64, skip, limit int) []*models.Transaction {
	var transactions []*models.Transaction
	err := self.mongo.C("Transactions").Find(bson.M{"block_number": blockNumber}).Skip(skip).Limit(limit).All(&transactions)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		self.Lgr.Error("Failed to get txs for block", zap.Int64("block", blockNumber), zap.Error(err))
	}
	return transactions
}

func (self *MongoBackend) getLatestsBlocks(skip, limit int) []*models.LightBlock {
	var blocks []*models.LightBlock
	err := self.mongo.C("Blocks").Find(nil).Sort("-number").Select(bson.M{"number": 1, "created_at": 1, "miner": 1, "tx_count": 1, "extra_data": 1}).Skip(skip).Limit(limit).All(&blocks)
	if err != nil {
		self.Lgr.Error("Failed to get latest blocks", zap.Int("skip", skip), zap.Int("limit", limit), zap.Error(err))
		return nil
	}
	return blocks
}

func (self *MongoBackend) getActiveAddresses(fromDate time.Time) []*models.ActiveAddress {
	var addresses []*models.ActiveAddress
	err := self.mongo.C("ActiveAddress").Find(bson.M{"updated_at": bson.M{"$gte": fromDate}}).Select(bson.M{"address": 1}).Sort("-updated_at").All(&addresses)
	if err != nil {
		self.Lgr.Error("Failed to get active addresses", zap.Time("from", fromDate), zap.Error(err))
	}
	return addresses
}

func (self *MongoBackend) isContract(address string) bool {
	var c models.Address
	err := self.mongo.C("Addresses").Find(bson.M{"address": address}).Select(bson.M{"contract": 1}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false
		}
		self.Lgr.Error("Failed to check if contract", zap.String("address", address), zap.Error(err))
		return false
	}
	return c.Contract
}

func (self *MongoBackend) getAddressByHash(address string) *models.Address {
	var c models.Address
	err := self.mongo.C("Addresses").Find(bson.M{"address": address}).One(&c)
	if err != nil {
		self.Lgr.Error("Failed to get address", zap.String("address", address), zap.Error(err))
		return nil
	}
	//lazy calculation for number of transactions
	transactionCounter, err := self.mongo.C("Transactions").Find(bson.M{"$or": []bson.M{bson.M{"from": address}, bson.M{"to": address}}}).Count()
	if err != nil {
		self.Lgr.Fatal("importAddress", zap.Error(err))
	}
	c.NumberOfTransactions = transactionCounter
	return &c
}

func (self *MongoBackend) getTransactionByHash(transactionHash string) *models.Transaction {
	lgr := self.Lgr.With(zap.String("tx", transactionHash))
	var c models.Transaction
	err := self.mongo.C("Transactions").Find(bson.M{"tx_hash": transactionHash}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		lgr.Error("Failed to get tx", zap.Error(err))
		return nil
	}
	// lazy calculation for receipt
	if !c.ReceiptReceived {
		receipt, err := self.goClient.TransactionReceipt(context.Background(), common.HexToHash(transactionHash))
		if err != nil {
			lgr.Warn("Failed to get transaction receipt", zap.Error(err))
		} else {
			gasPrice, ok := new(big.Int).SetString(c.GasPrice, 0)
			if !ok {
				lgr.Error("Failed to parse gas price", zap.String("gasPrice", c.GasPrice))
			}
			c.GasFee = new(big.Int).Mul(gasPrice, big.NewInt(int64(receipt.GasUsed))).String()
			c.ContractAddress = receipt.ContractAddress.String()
			c.Status = false
			if receipt.Status == 1 {
				c.Status = true
			}
			c.ReceiptReceived = true
			jsonValue, err := json.Marshal(receipt.Logs)
			if err != nil {
				lgr.Error("Failed to marshal JSON receipt logs", zap.Error(err))
			}
			c.Logs = string(jsonValue)
			_, err = self.mongo.C("Transactions").Upsert(bson.M{"tx_hash": c.TxHash}, c)
			if err != nil {
				lgr.Error("Failed to upsert tx", zap.Error(err))
			}
		}
	}
	return &c
}

func (self *MongoBackend) getTransactionList(address string, skip, limit int, fromTime, toTime time.Time, inputDataEmpty *bool) []*models.Transaction {
	var transactions []*models.Transaction
	var err error
	if inputDataEmpty != nil {
		err = self.mongo.C("Transactions").Find(bson.M{"$or": []bson.M{bson.M{"from": address}, bson.M{"to": address}}, "created_at": bson.M{"$gte": fromTime, "$lte": toTime}, "input_data_empty": *inputDataEmpty}).Sort("-created_at").Skip(skip).Limit(limit).All(&transactions)
	} else {
		err = self.mongo.C("Transactions").Find(bson.M{"$or": []bson.M{bson.M{"from": address}, bson.M{"to": address}}, "created_at": bson.M{"$gte": fromTime, "$lte": toTime}}).Sort("-created_at").Skip(skip).Limit(limit).All(&transactions)
	}
	if err != nil {
		self.Lgr.Error("Failed to get transaction list", zap.String("address", address),
			zap.Time("from", fromTime), zap.Time("to", toTime), zap.Error(err))
	}
	return transactions
}

func (self *MongoBackend) getTokenHoldersList(contractAddress string, skip, limit int) []*models.TokenHolder {
	var tokenHoldersList []*models.TokenHolder
	err := self.mongo.C("TokensHolders").Find(bson.M{"contract_address": contractAddress}).Sort("-balance_int").Skip(skip).Limit(limit).All(&tokenHoldersList)
	if err != nil {
		self.Lgr.Error("Failed to get token holders list", zap.String("address", contractAddress), zap.Error(err))
	}
	return tokenHoldersList
}
func (self *MongoBackend) getOwnedTokensList(ownerAddress string, skip, limit int) []*models.TokenHolder {
	var tokenHoldersList []*models.TokenHolder
	err := self.mongo.C("TokensHolders").Find(bson.M{"token_holder_address": ownerAddress}).Sort("-balance_int").Skip(skip).Limit(limit).All(&tokenHoldersList)
	if err != nil {
		self.Lgr.Error("Failed to get owned tokens list", zap.String("address", ownerAddress), zap.Error(err))
	}
	return tokenHoldersList
}

func (self *MongoBackend) getInternalTransactionsList(contractAddress string, tokenTransactions bool, skip, limit int) []*models.InternalTransaction {
	var internalTransactionsList []*models.InternalTransaction
	var query bson.M
	if tokenTransactions {
		query = bson.M{"$or": []bson.M{bson.M{"from_address": contractAddress}, bson.M{"to_address": contractAddress}}}
	} else {
		query = bson.M{"contract_address": contractAddress}
	}
	err := self.mongo.C("InternalTransactions").Find(query).Sort("-block_number").Skip(skip).Limit(limit).All(&internalTransactionsList)
	if err != nil {
		self.Lgr.Error("Failed to get internal txs list", zap.String("address", contractAddress), zap.Error(err))
	}
	return internalTransactionsList
}

func (self *MongoBackend) getContract(contractAddress string) *models.Contract {
	var contract *models.Contract
	err := self.mongo.C("Contracts").Find(bson.M{"address": contractAddress}).One(&contract)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil
		}
		self.Lgr.Error("Failed to get contract", zap.String("contractAddress", contractAddress), zap.Error(err))
	}
	return contract
}

func (self *MongoBackend) getContractBlock(contractAddress string) int64 {
	var transaction *models.Transaction
	err := self.mongo.C("Transactions").Find(bson.M{"contract_address": contractAddress}).One(&transaction)
	if err != nil {
		self.Lgr.Error("Failed to get tx by contract address", zap.String("address", contractAddress), zap.Error(err))
	}
	if transaction != nil {
		return transaction.BlockNumber
	} else {
		return 0
	}

}

func (self *MongoBackend) updateContract(contract *models.Contract) error {
	_, err := self.mongo.C("Contracts").Upsert(bson.M{"address": contract.Address}, contract)
	if err != nil {
		return fmt.Errorf("failed to update contract: %v", err)
	}
	return nil
}

func (self *MongoBackend) getContracts(filter *models.ContractsFilter) []*models.Address {
	var addresses []*models.Address
	var sortQuery string
	findQuery := bson.M{"contract": true}
	if filter.TokenName != "" {
		findQuery["token_name"] = bson.RegEx{regexp.QuoteMeta(filter.TokenName), "i"}
	}
	if filter.TokenSymbol != "" {
		findQuery["token_symbol"] = bson.RegEx{regexp.QuoteMeta(filter.TokenSymbol), "i"}
	}
	if filter.ErcType != "" {
		findQuery["erc_types"] = filter.ErcType
	}
	if filter.SortBy != "" {
		sortQuery = filter.SortBy
		if filter.Asc == false {
			sortQuery = "-" + sortQuery
		}
	} else {
		sortQuery = "-number_of_token_holders"
	}
	if filter.Skip < 0 {
		filter.Skip = defaultSkip
	}
	if filter.Limit < 0 || filter.Limit > defaultFetchLimit {
		filter.Limit = defaultFetchLimit
	}
	err := self.mongo.
		C("Addresses").
		Find(findQuery).
		Sort(sortQuery).
		Skip(filter.Skip).
		Limit(filter.Limit).
		All(&addresses)
	if err != nil {
		self.Lgr.Error("Failed to query contracts", zap.Error(err))
	}
	return addresses
}

func (self *MongoBackend) getRichlist(skip, limit int, lockedAddresses []string) []*models.Address {
	var addresses []*models.Address
	err := self.mongo.C("Addresses").Find(bson.M{"balance_float": bson.M{"$gt": 0}, "address": bson.M{"$nin": lockedAddresses}}).Sort("-balance_float").Skip(skip).Limit(limit).All(&addresses)
	if err != nil {
		self.Lgr.Error("Failed to get rich list", zap.Error(err))
	}
	return addresses
}
func (self *MongoBackend) updateStats() {
	numOfTotalTransactions, err := self.mongo.C("Transactions").Find(nil).Count()
	if err != nil {
		self.Lgr.Error("GetStats: Failed to get Total Transactions", zap.Error(err))
	}
	numOfLastWeekTransactions, err := self.mongo.C("Transactions").Find(bson.M{"created_at": bson.M{"$gte": time.Now().AddDate(0, 0, -7)}}).Count()
	if err != nil {
		self.Lgr.Error("GetStats: Failed to get Last week Transactions", zap.Error(err))
	}
	numOfLastDayTransactions, err := self.mongo.C("Transactions").Find(bson.M{"created_at": bson.M{"$gte": time.Now().AddDate(0, 0, -1)}}).Count()
	if err != nil {
		self.Lgr.Error("GetStats: Failed to get 24H Transactions", zap.Error(err))
	}
	stats := &models.Stats{
		NumberOfTotalTransactions:    int64(numOfTotalTransactions),
		NumberOfLastWeekTransactions: int64(numOfLastWeekTransactions),
		NumberOfLastDayTransactions:  int64(numOfLastDayTransactions),
		UpdatedAt:                    time.Now(),
	}
	err = self.mongo.C("Stats").Insert(stats)
	if err != nil {
		self.Lgr.Error("Failed to update stats", zap.Error(err), zap.Reflect("stats", stats))
	}
}
func (self *MongoBackend) getStats() *models.Stats {
	var s *models.Stats
	err := self.mongo.C("Stats").Find(nil).Sort("-updated_at").One(&s)
	if err != nil {
		self.Lgr.Error("Failed to get stats", zap.Error(err))
		s = new(models.Stats)
	}
	return s
}

func (self *MongoBackend) getSignerStatsForRange(endTime time.Time, dur time.Duration) []models.SignerStats {
	var resp []bson.M
	stats := []models.SignerStats{}
	queryDayStats := []bson.M{bson.M{"$match": bson.M{"created_at": bson.M{"$gte": endTime.Add(dur)}}}, bson.M{"$group": bson.M{"_id": "$miner", "count": bson.M{"$sum": 1}}}}
	pipe := self.mongo.C("Blocks").Pipe(queryDayStats)
	err := pipe.All(&resp)
	if err != nil {
		self.Lgr.Info("Cannot run pipe", zap.Error(err))
	}
	for _, el := range resp {
		signerStats := models.SignerStats{SignerAddress: common.HexToAddress(el["_id"].(string)), BlocksCount: el["count"].(int)}
		stats = append(stats, signerStats)
	}
	return stats
}

func (self *MongoBackend) getBlockRange(endTime time.Time, dur time.Duration) models.BlockRange {
	var startBlock, endBlock models.Block
	var resp models.BlockRange
	err := self.mongo.C("Blocks").Find(bson.M{"created_at": bson.M{"$gte": endTime.Add(dur)}}).Select(bson.M{"number": 1}).Sort("created_at").One(&startBlock)
	if err != nil {
		self.Lgr.Error("Failed to get start block number", zap.Error(err))
	} else {
		resp.StartBlock = startBlock.Number
	}
	err = self.mongo.C("Blocks").Find(bson.M{"created_at": bson.M{"$gte": endTime.Add(dur)}}).Select(bson.M{"number": 1}).Sort("-created_at").One(&endBlock)
	if err != nil {
		self.Lgr.Error("Failed to get end block number", zap.Error(err))
	} else {
		resp.EndBlock = endBlock.Number
	}
	return resp
}

func (self *MongoBackend) getSignersStats() []models.SignersStats {
	var stats []models.SignersStats
	const day = -24 * time.Hour
	kvs := map[string]time.Duration{"daily": day, "weekly": 7 * day, "monthly": 30 * day}
	endTime := time.Now()
	for k, v := range kvs {
		stats = append(stats, models.SignersStats{BlockRange: self.getBlockRange(endTime, v), SignerStats: self.getSignerStatsForRange(endTime, v), Range: k})
	}
	return stats
}

func (self *MongoBackend) cleanUp() {
	collectionNames, err := self.mongo.CollectionNames()
	if err != nil {
		self.Lgr.Info("Cannot get list of collections", zap.Error(err))
		return
	}
	for _, collectionName := range collectionNames {
		self.Lgr.Info("cleanUp", zap.String("collection name", collectionName))
		self.mongo.C(collectionName).RemoveAll(nil)
	}
}

package backend

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"time"

	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/explorer/server/tokens"

	"github.com/gochain/gochain/v3/common"
	"github.com/gochain/gochain/v3/core/types"
	"github.com/gochain/gochain/v3/goclient"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var wei = big.NewInt(1000000000000000000)

type MongoBackend struct {
	host         string
	mongo        *mgo.Database
	mongoSession *mgo.Session
	goClient     *goclient.Client
	Lgr          *zap.Logger
}

// New create new rpc client with given url
func NewMongoClient(client *goclient.Client, host, dbName string, lgr *zap.Logger) (*MongoBackend, error) {
	session, err := mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:   []string{host},
		Timeout: 120 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to dial mongo: %v", err)
	}
	importer := new(MongoBackend)
	importer.Lgr = lgr
	importer.mongoSession = session
	importer.mongo = session.DB(dbName)
	importer.goClient = client
	importer.createIndexes()

	return importer, nil

}
func (self *MongoBackend) PingDB() error {
	return self.mongoSession.Ping()
}
func (self *MongoBackend) parseTx(ctx context.Context, tx *types.Transaction, block *types.Block) (*models.Transaction, error) {
	from, err := self.goClient.TransactionSender(ctx, tx, block.Header().Hash(), 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx sender: %v", err)
	}
	gas := tx.Gas()
	to := ""
	if tx.To() != nil {
		to = tx.To().Hex()
	}
	self.Lgr.Debug("Parsing tx", zap.Stringer("tx", tx.Hash()))
	txInputData := hex.EncodeToString(tx.Data()[:])
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
		InputData:       txInputData,
		InputDataEmpty:  txInputData == "",
	}, nil
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

func (self *MongoBackend) createIndexes() error {
	type CIndex struct {
		c     string
		index mgo.Index
	}
	for i, cIdx := range []CIndex{
		{c: "Transactions", index: mgo.Index{Key: []string{"tx_hash"}, Unique: true, DropDups: true, Background: true, Sparse: true}},
		{c: "Transactions", index: mgo.Index{Key: []string{"block_number"}, Background: true, Sparse: true}},
		{c: "Transactions", index: mgo.Index{Key: []string{"from", "created_at", "input_data_empty"}, Background: true}},
		{c: "Transactions", index: mgo.Index{Key: []string{"to", "created_at", "input_data_empty"}, Background: true}},
		{c: "Transactions", index: mgo.Index{Key: []string{"-created_at"}, Background: true}},
		{c: "Transactions", index: mgo.Index{Key: []string{"contract_address"}, Background: true}},
		{c: "Blocks", index: mgo.Index{Key: []string{"number"}, Unique: true, DropDups: true, Background: true, Sparse: true}},
		{c: "Blocks", index: mgo.Index{Key: []string{"-number"}, Background: true}},
		{c: "Blocks", index: mgo.Index{Key: []string{"miner"}, Background: true, Sparse: true}},
		{c: "Blocks", index: mgo.Index{Key: []string{"created_at", "miner"}, Background: true, Sparse: true}},
		{c: "Blocks", index: mgo.Index{Key: []string{"hash"}, Background: true, Sparse: true}},
		{c: "ActiveAddress", index: mgo.Index{Key: []string{"updated_at"}, Background: true, Sparse: true}},
		{c: "ActiveAddress", index: mgo.Index{Key: []string{"address"}, Unique: true, DropDups: true, Background: true, Sparse: true}},
		{c: "Address", index: mgo.Index{Key: []string{"address"}, Unique: true, DropDups: true, Background: true, Sparse: true}},
		{c: "Address", index: mgo.Index{Key: []string{"contract"}, Background: true}},
		{c: "Address", index: mgo.Index{Key: []string{"-balance_float", "address"}, Background: true, Sparse: true}},
		{c: "TokenHolders", index: mgo.Index{Key: []string{"contract_address", "token_holder_address"}, Background: true, Sparse: true}},
		{c: "TokenHolders", index: mgo.Index{Key: []string{"token_holder_address"}, Background: true, Sparse: true}},
		{c: "TokenHolders", index: mgo.Index{Key: []string{"balance_int"}, Background: true, Sparse: true}},
		{c: "InternalTransactions", index: mgo.Index{Key: []string{"contract_address", "from_address", "to_address"}, Background: true, Sparse: true}},
		{c: "InternalTransactions", index: mgo.Index{Key: []string{"from_address", "block_number"}, Background: true}},
		{c: "InternalTransactions", index: mgo.Index{Key: []string{"to_address", "block_number"}, Background: true}},
		{c: "InternalTransactions", index: mgo.Index{Key: []string{"transaction_hash"}, Background: true, Sparse: true}},
		{c: "InternalTransactions", index: mgo.Index{Key: []string{"block_number"}, Background: true, Sparse: true}},
		{c: "Stats", index: mgo.Index{Key: []string{"-updated_at"}, Background: true, Sparse: true}},
		{c: "Contracts", index: mgo.Index{Key: []string{"address"}, Unique: true, DropDups: true, Background: true, Sparse: true}},
		{c: "TransactionsByAddress", index: mgo.Index{Key: []string{"address", "tx_hash"}, Unique: true, DropDups: true, Background: true, Sparse: true}},
		{c: "TransactionsByAddress", index: mgo.Index{Key: []string{"created_at"}, Background: true, Sparse: true}},
	} {
		if err := self.mongo.C(cIdx.c).EnsureIndex(cIdx.index); err != nil {
			return fmt.Errorf("failed to create index %d for collection %q: %v", i, cIdx.c, err)
		}
	}
	return nil
}

func (self *MongoBackend) importBlock(ctx context.Context, block *types.Block) (*models.Block, error) {
	lgr := self.Lgr.With(zap.Int64("number", block.Header().Number.Int64()),
		zap.Stringer("hash", block.Hash()), zap.Stringer("parentHash", block.ParentHash()))
	lgr.Debug("Importing block")
	b := self.parseBlock(block)
	_, err := self.mongo.C("Blocks").Upsert(bson.M{"number": b.Number}, b)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert block: %v", err)
	}
	// deleting all txs belong to this block if any exist
	_, err = self.mongo.C("Transactions").RemoveAll(bson.M{"block_number": b.Number})
	if err != nil {
		return nil, fmt.Errorf("failed to remove old txs: %v", err)
	}
	for _, tx := range block.Transactions() {
		if err := self.importTx(ctx, tx, block); err != nil {
			return nil, fmt.Errorf("failed to import tx: %v", err)
		}
	}
	if err := self.UpdateActiveAddress(block.Coinbase().Hex()); err != nil {
		return nil, fmt.Errorf("failed to update active signer address: %s", err)
	}
	return b, nil

}

func (self *MongoBackend) deleteBlockByNumber(bnum int64) error {
	//delete block
	_, err := self.mongo.C("Blocks").RemoveAll(bson.M{"number": bnum})
	if err != nil {
		return fmt.Errorf("failed to remove block: %v", err)
	}
	// deleting all txs belong to this block if any exist
	_, err = self.mongo.C("Transactions").RemoveAll(bson.M{"block_number": bnum})
	if err != nil {
		return fmt.Errorf("failed to remove old txs: %v", err)
	}
	return nil
}

func (self *MongoBackend) deleteBlockByHash(hash string) error {
	//delete block
	_, err := self.mongo.C("Blocks").RemoveAll(bson.M{"BlockHash": hash})
	if err != nil {
		return fmt.Errorf("failed to remove block: %v", err)
	}
	// deleting all txs belong to this block if any exist
	_, err = self.mongo.C("Transactions").RemoveAll(bson.M{"block_hash": hash})
	if err != nil {
		return fmt.Errorf("failed to remove old txs: %v", err)
	}
	return nil
}

func (self *MongoBackend) UpdateActiveAddress(address string) error {
	_, err := self.mongo.C("ActiveAddress").Upsert(bson.M{"address": address}, &models.ActiveAddress{Address: address, UpdatedAt: time.Now()})
	return err
}

func (self *MongoBackend) importTx(ctx context.Context, tx *types.Transaction, block *types.Block) error {
	lgr := self.Lgr.With(zap.Stringer("tx", tx.Hash()))
	lgr.Debug("Importing tx")
	transaction, err := self.parseTx(ctx, tx, block)
	if err != nil {
		return fmt.Errorf("failed to parse tx: %v", err)
	}

	toAddress := transaction.To
	if transaction.To == "" {
		lgr.Debug("Importing contract creation")
		receipt, err := self.goClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return fmt.Errorf("failed to get tx receipt: %v", err)
		}
		contractAddress := receipt.ContractAddress.String()
		if contractAddress != "0x0000000000000000000000000000000000000000" {
			transaction.ContractAddress = contractAddress
		}
		transaction.Status = false
		if receipt.Status == 1 {
			transaction.Status = true
		}
		toAddress = transaction.ContractAddress
	}

	if _, err := self.mongo.C("Transactions").Upsert(bson.M{"tx_hash": tx.Hash().String()}, transaction); err != nil {
		return fmt.Errorf("failed to upsert tx: %v", err)
	}

	if err := self.UpdateActiveAddress(toAddress); err != nil {
		return fmt.Errorf("failed to update active to address: %s", err)
	}
	if err := self.UpdateActiveAddress(transaction.From); err != nil {
		return fmt.Errorf("failed to update active from address: %s", err)
	}
	return nil
}

// needReloadParent returns true if the parent block is missing or does not match the hash from this block number.
func (self *MongoBackend) needReloadParent(blockNumber int64) (bool, error) {
	block, err := self.getBlockByNumber(blockNumber)
	if err != nil {
		return false, err
	}
	if block == nil {
		self.Lgr.Debug("Checking parent - main block not found", zap.Int64("block", blockNumber))
		return true, nil
	}
	parentBlockNumber := (block.Number - 1)
	parentBlock, err := self.getBlockByNumber(parentBlockNumber)
	if err != nil {
		return false, fmt.Errorf("failed to get parent: %v", err)
	}
	if parentBlock != nil {
		self.Lgr.Debug("Checking parent", zap.Int64("block.number", block.Number), zap.String("block.parentHash", block.ParentHash),
			zap.Int64("parent.number", parentBlock.Number), zap.String("parent.hash", parentBlock.BlockHash))
	}
	return parentBlock == nil || (block.ParentHash != parentBlock.BlockHash), nil

}

// transactionsConsistent returns true if the block count matches the number of transactions with that block number.
func (self *MongoBackend) transactionsConsistent(blockNumber int64) (bool, error) {
	block, err := self.getBlockByNumber(blockNumber)
	if err != nil {
		return false, fmt.Errorf("failed to get block: %v", err)
	}
	if block == nil {
		return false, errors.New("block not found")
	}
	txCount, err := self.mongo.C("Transactions").Find(bson.M{"block_number": blockNumber}).Count()
	if err != nil {
		return false, fmt.Errorf("failed to count txs in db: %v", err)
	}
	self.Lgr.Debug("Checking tx count", zap.Int64("blockNumber", blockNumber),
		zap.Int("block.count", block.TxCount), zap.Int("db.count", txCount))
	return txCount == block.TxCount, nil
}

func (self *MongoBackend) importAddress(address string, balance *big.Int, token *tokens.TokenDetails, contract bool, updatedAtBlock int64) (*models.Address, error) {
	balanceGoFloat, _ := new(big.Float).SetPrec(100).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(wei)).Float64() //converting to GO from wei
	balanceGoString := new(big.Rat).SetFrac(balance, wei).FloatString(18)
	lgr := self.Lgr.With(zap.String("address", address))
	lgr.Debug("Updating address", zap.String("balance", balanceGoString), zap.Float64("balanceFloat", balanceGoFloat))
	tokenHoldersCounter, err := self.mongo.C("TokensHolders").Find(bson.M{"contract_address": address}).Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count token holders: %v", err)
	}

	internalTransactionsCounter, err := self.mongo.C("InternalTransactions").Find(bson.M{"contract_address": address}).Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count internal txs: %v", err)
	}

	tokenTransactionsCounter, err := self.mongo.C("InternalTransactions").Find(bson.M{"$or": []bson.M{{"from_address": address}, {"to_address": address}}}).Count()
	if err != nil {
		return nil, fmt.Errorf("failed to count held token txs: %v", err)
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
		ErcTypes:       token.ERCTypesSlice(),
		Interfaces:     token.FunctionsSlice(),
		BalanceFloat:   balanceGoFloat,
		BalanceString:  balanceGoString,
		// NumberOfTransactions:         transactionCounter,
		NumberOfTokenHolders:         tokenHoldersCounter,
		NumberOfInternalTransactions: internalTransactionsCounter,
		NumberOfTokenTransactions:    tokenTransactionsCounter,
	}
	_, err = self.mongo.C("Addresses").Upsert(bson.M{"address": address}, addressM)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert address: %v", err)
	}
	return addressM, nil
}

func (self *MongoBackend) importTokenHolder(contractAddress, tokenHolderAddress string, token *tokens.TokenHolderDetails, address *models.Address) (*models.TokenHolder, error) {
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
		return nil, fmt.Errorf("failed to upsert token holders: %v", err)
	}
	return tokenHolder, nil

}

func (self *MongoBackend) importTransferEvent(contractAddress string, transferEvent *tokens.TransferEvent, createdAt time.Time) (*models.TokenTransfer, error) {
	internalTransaction := &models.TokenTransfer{
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
		return nil, fmt.Errorf("failed to upsert internal txs: %v", zap.Error(err))
	}
	return internalTransaction, nil
}

func (self *MongoBackend) importContract(contractAddress string, byteCode string) error {
	//https://stackoverflow.com/questions/43278696/golang-mgo-insert-or-update-not-working-as-expected/43278832
	_, err := self.mongo.C("Contracts").Upsert(bson.M{"address": contractAddress}, bson.M{"$set": bson.M{"address": contractAddress, "byte_code": byteCode, "created_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("failed to upsert contract: %v", err)
	}
	return nil
}

func (self *MongoBackend) deleteContract(contractAddress string) error {
	//delete internal transactions
	_, err := self.mongo.C("InternalTransactions").RemoveAll(bson.M{"contract_address": contractAddress})
	if err != nil {
		return fmt.Errorf("failed to remove internal transactions: %v", err)
	}
	// deleting all token holders
	_, err = self.mongo.C("TokensHolders").RemoveAll(bson.M{"contract_address": contractAddress})
	if err != nil {
		return fmt.Errorf("failed to remove token holders: %v", err)
	}
	// deleting contract
	_, err = self.mongo.C("Contracts").RemoveAll(bson.M{"address": contractAddress})
	if err != nil {
		return fmt.Errorf("failed to remove contract: %v", err)
	}
	return nil
}

func (self *MongoBackend) getBlockByNumber(blockNumber int64) (*models.Block, error) {
	var c models.Block
	err := self.mongo.C("Blocks").Find(bson.M{"number": blockNumber}).Select(bson.M{"transactions": 0}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get block: %v", err)
	}
	return &c, nil
}

func (self *MongoBackend) getBlockByHash(blockHash string) (*models.Block, error) {
	var c models.Block
	err := self.mongo.C("Blocks").Find(bson.M{"hash": blockHash}).Select(bson.M{"transactions": 0}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get block by hash: %v", err)
	}
	return &c, nil
}

func (self *MongoBackend) getBlockTransactionsByNumber(blockNumber int64, filter *models.PaginationFilter) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	err := self.mongo.C("Transactions").
		Find(bson.M{"block_number": blockNumber}).
		Skip(filter.Skip).
		Limit(filter.Limit).
		All(&transactions)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get txs for block: %v", err)
	}
	return transactions, nil
}

func (self *MongoBackend) getLatestsBlocks(filter *models.PaginationFilter) ([]*models.LightBlock, error) {
	var blocks []*models.LightBlock
	err := self.mongo.C("Blocks").
		Find(nil).
		Sort("-number").
		Select(bson.M{"number": 1, "created_at": 1, "miner": 1, "tx_count": 1, "extra_data": 1}).
		Skip(filter.Skip).
		Limit(filter.Limit).
		All(&blocks)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest blocks: %v", err)
	}
	return blocks, nil
}

func (self *MongoBackend) getActiveAddresses(fromDate time.Time) ([]*models.ActiveAddress, error) {
	var addresses []*models.ActiveAddress
	err := self.mongo.C("ActiveAddress").Find(bson.M{"updated_at": bson.M{"$gte": fromDate}}).Select(bson.M{"address": 1}).Sort("-updated_at").All(&addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get active addresses: %v", err)
	}
	return addresses, nil
}

func (self *MongoBackend) isContract(address string) (bool, error) {
	var c models.Address
	err := self.mongo.C("Addresses").Find(bson.M{"address": address}).Select(bson.M{"contract": 1}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if contract: %v", err)
	}
	return c.Contract, nil
}

func (self *MongoBackend) getAddressByHash(address string) (*models.Address, error) {
	var c models.Address
	err := self.mongo.C("Addresses").Find(bson.M{"address": address}).One(&c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get address: %v", err)
	}
	//lazy calculation for number of transactions
	transactionCounter, err := self.mongo.C("Transactions").Find(bson.M{"$or": []bson.M{{"from": address}, {"to": address}}}).Count()
	if err != nil {
		return nil, fmt.Errorf("failed to get txs: %v", err)
	}
	c.NumberOfTransactions = transactionCounter
	return &c, nil
}

func (self *MongoBackend) getTxByAddressAndNonce(ctx context.Context, address string, nonce int64) (*models.Transaction, error) {
	var tx models.Transaction
	err := self.mongo.C("Transactions").Find(bson.M{"from": address, "nonce": nonce}).One(&tx)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tx: %v", err)
	}
	return self.ensureReceipt(ctx, &tx)
}

func (self *MongoBackend) getTransactionByHash(ctx context.Context, hash string) (*models.Transaction, error) {
	var tx models.Transaction
	err := self.mongo.C("Transactions").Find(bson.M{"tx_hash": hash}).One(&tx)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get tx: %v", err)
	}
	return self.ensureReceipt(ctx, &tx)
}

// ensureReceipt does lazy loads receipt info if necessary.
func (self *MongoBackend) ensureReceipt(ctx context.Context, tx *models.Transaction) (*models.Transaction, error) {
	if tx.ReceiptReceived {
		return tx, nil
	}
	lgr := self.Lgr.With(zap.String("tx", tx.TxHash))
	receipt, err := self.goClient.TransactionReceipt(ctx, common.HexToHash(tx.TxHash))
	if err != nil {
		lgr.Warn("Failed to get transaction receipt", zap.Error(err))
	} else {
		gasPrice, ok := new(big.Int).SetString(tx.GasPrice, 0)
		if !ok {
			lgr.Error("Failed to parse gas price", zap.String("gasPrice", tx.GasPrice))
		}
		tx.GasFee = new(big.Int).Mul(gasPrice, big.NewInt(int64(receipt.GasUsed))).String()
		tx.ContractAddress = receipt.ContractAddress.String()
		tx.Status = false
		if receipt.Status == 1 {
			tx.Status = true
		}
		tx.ReceiptReceived = true
		jsonValue, err := json.Marshal(receipt.Logs)
		if err != nil {
			lgr.Error("Failed to marshal JSON receipt logs", zap.Error(err))
		}
		tx.Logs = string(jsonValue)
		_, err = self.mongo.C("Transactions").Upsert(bson.M{"tx_hash": tx.TxHash}, tx)
		if err != nil {
			return nil, fmt.Errorf("failed to upsert tx: %v", err)
		}
	}
	return tx, nil
}

func (self *MongoBackend) getTransactionList(
	address string,
	filter *models.TxsFilter,
) ([]*models.Transaction, error) {
	var transactions []*models.Transaction
	findQuery := bson.M{
		"$or": []bson.M{
			{"from": address},
			{"to": address},
		},
		"created_at": bson.M{
			"$gte": filter.FromTime,
			"$lte": filter.ToTime,
		},
	}
	if filter.InputDataEmpty != nil {
		findQuery["input_data_empty"] = *filter.InputDataEmpty
	}
	err := self.mongo.C("Transactions").
		Find(findQuery).
		Sort("-created_at").
		Skip(filter.Skip).
		Limit(filter.Limit).
		All(&transactions)
	if err != nil {
		return nil, fmt.Errorf("failed to get tx list: %v", err)
	}
	return transactions, nil
}

func (self *MongoBackend) getTokenHoldersList(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenHolder, error) {
	var tokenHoldersList []*models.TokenHolder
	err := self.mongo.C("TokensHolders").
		Find(bson.M{"contract_address": contractAddress}).
		Sort("-balance_int").
		Skip(filter.Skip).
		Limit(filter.Limit).
		All(&tokenHoldersList)
	if err != nil {
		return nil, fmt.Errorf("failed to get token holders list: %v", err)
	}
	return tokenHoldersList, nil
}
func (self *MongoBackend) getOwnedTokensList(ownerAddress string, filter *models.PaginationFilter) ([]*models.TokenHolder, error) {
	var tokenHoldersList []*models.TokenHolder
	err := self.mongo.C("TokensHolders").
		Find(bson.M{"token_holder_address": ownerAddress}).
		Sort("-balance_int").
		Skip(filter.Skip).
		Limit(filter.Limit).
		All(&tokenHoldersList)
	if err != nil {
		return nil, fmt.Errorf("failed to get owned tokens list: %v", err)
	}
	return tokenHoldersList, nil
}

// getInternalTokenTransfers gets token transfer events emitted by this contract.
func (self *MongoBackend) getInternalTokenTransfers(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenTransfer, error) {
	var internalTransactionsList []*models.TokenTransfer
	err := self.mongo.C("InternalTransactions").
		Find(bson.M{"contract_address": contractAddress}).
		Sort("-block_number").Skip(filter.Skip).Limit(filter.Limit).All(&internalTransactionsList)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal txs list: %v", err)
	}
	return internalTransactionsList, nil
}

// getHeldTokenTransfers gets token transfer events to or from this contract, for any token.
func (self *MongoBackend) getHeldTokenTransfers(contractAddress string, filter *models.PaginationFilter) ([]*models.TokenTransfer, error) {
	var internalTransactionsList []*models.TokenTransfer
	err := self.mongo.C("InternalTransactions").
		Find(bson.M{"$or": []bson.M{{"from_address": contractAddress}, {"to_address": contractAddress}}}).
		Sort("-block_number").Skip(filter.Skip).Limit(filter.Limit).All(&internalTransactionsList)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal txs list: %v", err)
	}
	return internalTransactionsList, nil
}

func (self *MongoBackend) getContract(contractAddress string) (*models.Contract, error) {
	var contract *models.Contract
	err := self.mongo.C("Contracts").Find(bson.M{"address": contractAddress}).One(&contract)
	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get contract: %v", err)
	}
	return contract, nil
}

func (self *MongoBackend) getContractBlock(contractAddress string) (int64, error) {
	var transaction *models.Transaction
	err := self.mongo.C("Transactions").Find(bson.M{"contract_address": contractAddress}).One(&transaction)
	if err != nil {
		if err == mgo.ErrNotFound {
			return 0, errors.New("tx that deployed contract not found")
		}
		return 0, fmt.Errorf("failed to get tx that deployed contract: %v", err)
	}
	if transaction == nil {
		return 0, errors.New("tx that deployed contract not found")
	}
	return transaction.BlockNumber, nil
}

func (self *MongoBackend) updateContract(contract *models.Contract) error {
	_, err := self.mongo.C("Contracts").Upsert(bson.M{"address": contract.Address}, contract)
	if err != nil {
		return fmt.Errorf("failed to update contract: %v", err)
	}
	return nil
}

func (self *MongoBackend) getContracts(filter *models.ContractsFilter) ([]*models.Address, error) {
	var addresses []*models.Address
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
	if filter.SortBy == "" {
		filter.SortBy = "number_of_token_holders"
		filter.Asc = false
	}

	contractQuery := bson.M{
		"attached_contract.valid": true,
	}
	if filter.ContractName != "" {
		contractQuery["attached_contract.contract_name"] = bson.RegEx{regexp.QuoteMeta(filter.ContractName), "i"}
	}

	sortDir := -1
	if filter.Asc {
		sortDir = 1
	}
	sortQuery := bson.M{filter.SortBy: sortDir}
	query := []bson.M{
		{"$match": findQuery},
		{"$lookup": bson.M{
			"from":         "Contracts",
			"localField":   "address",
			"foreignField": "address",
			"as":           "attached_contract",
		}},
		{"$match": contractQuery},
		{"$unwind": bson.M{
			"path": "$attached_contract",
		}},
		{"$sort": sortQuery},
		{"$skip": filter.Skip},
		{"$limit": filter.Limit},
	}
	err := self.mongo.
		C("Addresses").
		Pipe(query).
		All(&addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to query contracts: %v", err)
	}
	return addresses, nil
}

func (self *MongoBackend) getRichlist(filter *models.PaginationFilter, lockedAddresses []string) ([]*models.Address, error) {
	var addresses []*models.Address
	err := self.mongo.C("Addresses").Find(bson.M{"balance_float": bson.M{"$gt": 0}, "address": bson.M{"$nin": lockedAddresses}}).Sort("-balance_float").Skip(filter.Skip).Limit(filter.Limit).All(&addresses)
	if err != nil {
		return nil, fmt.Errorf("failed to get rich list: %v", err)
	}
	return addresses, nil
}
func (self *MongoBackend) updateStats() (*models.Stats, error) {
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
	return stats, self.mongo.C("Stats").Insert(stats)
}
func (self *MongoBackend) getStats() (*models.Stats, error) {
	var s *models.Stats
	err := self.mongo.C("Stats").Find(nil).Sort("-updated_at").One(&s)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %v", err)
	}
	return s, nil
}

func (self *MongoBackend) getSignerStatsForRange(endTime time.Time, dur time.Duration) ([]models.SignerStats, error) {
	var resp []bson.M
	stats := []models.SignerStats{}
	queryDayStats := []bson.M{bson.M{"$match": bson.M{"created_at": bson.M{"$gte": endTime.Add(dur)}}}, bson.M{"$group": bson.M{"_id": "$miner", "count": bson.M{"$sum": 1}}}}
	err := self.mongo.C("Blocks").Pipe(queryDayStats).All(&resp)
	if err != nil {
		return nil, fmt.Errorf("failed to query signers stats: %v", err)
	}
	for _, el := range resp {
		addr := el["_id"].(string)
		if !common.IsHexAddress(addr) {
			return nil, fmt.Errorf("invalid hex address: %s", addr)
		}
		signerStats := models.SignerStats{SignerAddress: common.HexToAddress(addr), BlocksCount: el["count"].(int)}
		stats = append(stats, signerStats)
	}
	return stats, nil
}

func (self *MongoBackend) getBlockRange(endTime time.Time, dur time.Duration) (models.BlockRange, error) {
	var startBlock, endBlock models.Block
	err := self.mongo.C("Blocks").Find(bson.M{"created_at": bson.M{"$gte": endTime.Add(dur)}}).Select(bson.M{"number": 1}).Sort("created_at").One(&startBlock)
	if err != nil {
		return models.BlockRange{}, fmt.Errorf("failed to get start block number: %v", err)
	}
	err = self.mongo.C("Blocks").Find(bson.M{"created_at": bson.M{"$gte": endTime.Add(dur)}}).Select(bson.M{"number": 1}).Sort("-created_at").One(&endBlock)
	if err != nil {
		return models.BlockRange{}, fmt.Errorf("failed to get end block number: %v", err)
	}
	return models.BlockRange{StartBlock: startBlock.Number, EndBlock: endBlock.Number}, nil
}

func (self *MongoBackend) getSignersStats() ([]models.SignersStats, error) {
	var stats []models.SignersStats
	const day = -24 * time.Hour
	kvs := map[string]time.Duration{"daily": day, "weekly": 7 * day, "monthly": 30 * day}
	endTime := time.Now()
	for k, v := range kvs {
		blockRange, err := self.getBlockRange(endTime, v)
		if err != nil {
			return nil, fmt.Errorf("failed to get block range: %v", err)
		}
		signerStats, err := self.getSignerStatsForRange(endTime, v)
		if err != nil {
			return nil, fmt.Errorf("failed to get signer stats: %v", err)
		}
		stats = append(stats, models.SignersStats{BlockRange: blockRange, SignerStats: signerStats, Range: k})
	}
	return stats, nil
}

func (self *MongoBackend) cleanUp() {
	collectionNames, err := self.mongo.CollectionNames()
	if err != nil {
		self.Lgr.Error("Cannot get list of collections", zap.Error(err))
		return
	}
	for _, collectionName := range collectionNames {
		_, err := self.mongo.C(collectionName).RemoveAll(nil)
		if err != nil {
			self.Lgr.Error("Failed to clean collection", zap.String("collection", collectionName), zap.Error(err))
			continue
		}
		self.Lgr.Info("Cleaned collection", zap.String("collection", collectionName))
	}
}

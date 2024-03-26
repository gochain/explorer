//go:build integration
// +build integration

package backend

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/gochain-io/explorer/server/models"
	"github.com/gochain-io/explorer/server/tokens"

	"github.com/gochain/gochain/v4/common"
	"github.com/gochain/gochain/v4/core/types"
	"github.com/gochain/gochain/v4/rlp"
	"go.uber.org/zap/zaptest"
)

func getBackend(t *testing.T) *Backend {
	t.Helper()
	b, err := NewBackend(context.Background(), "127.0.0.1:27017", "https://rpc.gochain.io", "testdb", nil, nil, nil, zaptest.NewLogger(t), nil)
	if err != nil {
		t.Fatalf("failed to create backend: %v", err)
	}
	return b
}

func createImportBlock(t *testing.T, testBackend *Backend) types.Block {
	t.Helper()
	blockEnc := common.FromHex("0xf90425f90260a0ca12d831416f1fd29336c535ee814d228f638b10798b5cf437bef29415e63762a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479448c67d87cd7d716ec044dbe33a0152557bf86062c0c0b84105faed8a008cdb41757cfca7462e440f13f1369a9dc14aafbc1eca77575b6f4753d34458da828f3028fc075e3d6e7b944c3b6b8327dbbd702a1a3ee0c6f8663301a0df5891f347b7525b4369d81238f399a046fb94194e0785d42ba93446b1e27c19a013d863c4c708ff270b60de30faa63e49f073b5ac0e90cdcddb1b10c0736cc923a004deb4be6955e1a300123be48007597f67e4229f8ce70f4f10388de6fd3fa267b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e8314be5c840822d32083014820845b634071a0322e312e33362f6c696e75782d616d6436342f676f312e31302e330000000000a00000000000000000000000000000000000000000000000000000000000000000880000000000000000f901bef86d80847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf8911efc99828edb2600080819ca08f2579f5372bfadca519dfa96b9a7e16689632b65a469c8f96262034d306e104a051692d3de03154e5506afe4616304e6e70bcd7ba7c367b59226045eacd21a8ddf86d01847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf894add3970f31b5f600080819ca0e95246628a5c110356c23a176a97d9dc8a62069b75ef59b090b9074a69eedb63a07882befc79d5e2a5214e33fcc7f2beb45d85551bc4c63588817af2cb83df00b1f86e09847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf8a0b69448bfc2ed424600080819ca08e0645ff002bb43239c22d64cab5dd799b2bb9c56f92d90eb50be7860d4aadd4a00bf2701bb7e1df17921c598954621c724c6a7c50cd1583e117c10708d3c1a537f86e05847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf8a05b7ac452dab97cb600080819ba00a0cacc5f5c556c2435f3e6e43635eda857f71e2f95d73c450ea01a481018011a02cb5ec486060a6f53fd5afdb4f87c57842b7600f3790d66a7bda90fdf0135e93c0")
	var block types.Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatalf("Failed to decode bytes: %v", err)
	}
	if _, err := testBackend.ImportBlock(context.Background(), &block); err != nil {
		t.Fatalf("Failed to import block: %v", err)
	}
	return block
}

func TestImportAddress(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	var token = &tokens.TokenDetails{TotalSupply: big.NewInt(0)}

	addrHash := "0x0000000000000000000000000000000000000000"

	_, err := testBackend.ImportAddress(addrHash, big.NewInt(1000), token, false, 0)
	if err != nil {
		t.Fatal(err)
	}
	address, err := testBackend.mongo.getAddressByHash(addrHash)
	if err != nil {
		t.Fatal(err)
	}

	if address.BalanceWei != "1000" {
		t.Errorf("Balance was incorrect, got: %s, want: %d.", address.BalanceWei, 1000)
	}
	if address.BalanceString != "0.000000000000001000" {
		t.Errorf("Balance was incorrect, got: %s, want: %d.", address.BalanceString, 1000)
	}
	wrongAddressHash := "0x000"

	address, err = testBackend.GetAddressByHash(context.Background(), wrongAddressHash)
	if err == nil {
		t.Errorf("Address was incorrect, expected error, got: %v", address)
	}
}

func TestImportBlockTransaction(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	block := createImportBlock(t, testBackend)
	blockFromDb, err := testBackend.GetBlockByNumber(context.Background(), block.Header().Number.Int64(), false)
	if err != nil {
		t.Fatalf("failed to get block: %v", err)
	}

	if block.Header().Number.Int64() != blockFromDb.Number || block.Header().Number.Int64() != 1359452 {
		t.Errorf("Block number was incorrect, got: %d, want: %d.", block.Header().Number.Int64(), 1359452)
	}

	if block.Header().Hash().Hex() != blockFromDb.BlockHash {
		t.Errorf("Block hash was incorrect, got: %s, want: %s.", block.Header().Hash().Hex(), blockFromDb.BlockHash)
	}

	if len(block.Transactions()) != blockFromDb.TxCount {
		t.Errorf("Block transactions were incorrect, got: %d, want: %d.", len(block.Transactions()), blockFromDb.TxCount)
	}
}
func TestTransactions(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	block := createImportBlock(t, testBackend)
	filter1 := &models.PaginationFilter{
		Skip:  0,
		Limit: 4,
	}
	transactionsFromDb, err := testBackend.GetBlockTransactionsByNumber(block.Header().Number.Int64(), filter1)
	if err != nil {
		t.Fatal(err)
	}
	if len(block.Transactions()) != len(transactionsFromDb) {
		t.Errorf("Block transactions were incorrect, got: %d, want: %d.", len(block.Transactions()), len(transactionsFromDb))
	}

	if block.Transactions()[0].Hash().Hex() != transactionsFromDb[0].TxHash {
		t.Errorf("Block transaction was incorrect, got: %s, want: %s.", block.Transactions()[0].Hash().Hex(), transactionsFromDb[0].TxHash)
	}
	transactionFromDB, err := testBackend.GetTransactionByHash(context.Background(), block.Transactions()[0].Hash().Hex())
	if err != nil {
		t.Fatal(err)
	}
	if block.Transactions()[0].Hash().Hex() != transactionFromDB.TxHash {
		t.Errorf("Block transaction was incorrect, got: %s, want: %s.", block.Transactions()[0].Hash().Hex(), transactionFromDB.TxHash)
	}
	transactionFromDB, err = testBackend.GetTxByAddressAndNonce(context.Background(), transactionFromDB.From, int64(block.Transactions()[0].Nonce()))
	if err != nil {
		t.Fatal(err)
	}
	if block.Transactions()[0].Hash().Hex() != transactionFromDB.TxHash {
		t.Errorf("Block transaction was incorrect, got: %s, want: %s.", block.Transactions()[0].Hash().Hex(), transactionFromDB.TxHash)
	}
	filter2 := &models.TxsFilter{
		PaginationFilter: models.PaginationFilter{
			Skip:  0,
			Limit: 3,
		},
		TimeFilter: models.TimeFilter{
			FromTime: time.Unix(0, 0),
			ToTime:   time.Now(),
		},
	}
	transactionsFromAddress, err := testBackend.GetTransactionList(transactionFromDB.From, filter2)
	if err != nil {
		t.Fatal(err)
	}
	if len(transactionsFromAddress) != 3 {
		t.Errorf("Wrong number of the transactions for address, got: %d, want: %d.", len(transactionsFromAddress), 4)
	}

	transactionsToAddress, err := testBackend.GetTransactionList(transactionFromDB.To, filter2)
	if err != nil {
		t.Fatal(err)
	}
	if len(transactionsToAddress) != 3 {
		t.Errorf("Wrong number of the transactions for address, got: %d, want: %d.", len(transactionsToAddress), 4)
	}
	filter2.Skip = 2
	transactionsToAddress, err = testBackend.GetTransactionList(transactionFromDB.To, filter2)
	if err != nil {
		t.Fatal(err)
	}
	if len(transactionsToAddress) != 2 {
		t.Errorf("Wrong number of the transactions for address, got: %d, want: %d.", len(transactionsToAddress), 2)
	}

}
func TestBlockByHash(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	block := createImportBlock(t, testBackend)

	blockFromDbByHash, err := testBackend.GetBlockByHash(context.Background(), block.Header().Hash().Hex())
	if err != nil {
		t.Fatal(err)
	} else if blockFromDbByHash == nil {
		t.Fatal("missing")
	}

	if block.Header().Number.Int64() != blockFromDbByHash.Number {
		t.Errorf("Block hash was incorrect, got: %d, want: %d.", block.Header().Number.Int64(), blockFromDbByHash.Number)
	}
}
func TestLatestBlocks(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	block := createImportBlock(t, testBackend)
	filter := &models.PaginationFilter{
		Limit: 100,
		Skip:  0,
	}
	latestBlocks, err := testBackend.GetLatestsBlocks(filter)
	if err != nil {
		t.Fatal(err)
	}

	if len(latestBlocks) != 1 {
		t.Errorf("Wrong number of the latest blocks , got: %d, want: %d.", len(latestBlocks), 1)
	}

	if latestBlocks[0].Number != block.Header().Number.Int64() {
		t.Errorf("Wrong the latest block number , got: %d, want: %d.", latestBlocks[0].Number, block.Header().Number.Int64())
	}
}
func TestActiveAddresses(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	block := createImportBlock(t, testBackend)

	activeNonContracts, err := testBackend.GetActiveAdresses(time.Unix(0, 0), false)
	if err != nil {
		t.Fatal(err)
	}

	activeContracts, err := testBackend.GetActiveAdresses(time.Unix(0, 0), true)
	if err != nil {
		t.Fatal(err)
	}

	if len(activeNonContracts) != 3 {
		t.Errorf("activeNonContracts was incorrect, got: %d, want: %d.", len(activeNonContracts), 3)
	}

	if activeNonContracts[0].Address != block.Coinbase().Hex() {
		t.Errorf("activeContracts  was incorrect, got: %s, want: %s.", activeNonContracts[len(activeNonContracts)-1].Address, block.Coinbase().Hex())
	}

	if len(activeContracts) != 0 {
		t.Errorf("activeContracts was incorrect, got: %d, want: %d.", len(activeContracts), 0)
	}

	activeNonContracts, err = testBackend.GetActiveAdresses(time.Now(), false)
	if err != nil {
		t.Fatal(err)
	}
	if len(activeNonContracts) != 0 {
		t.Errorf("activeNonContracts was incorrect, got: %d, want: %d.", len(activeNonContracts), 0)
	}

}

func TestRichList(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	var token = &tokens.TokenDetails{TotalSupply: big.NewInt(0)}

	addrHash := "0x0000000000000000000000000000000000000000"

	_, err := testBackend.ImportAddress(addrHash, big.NewInt(1000), token, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = testBackend.ImportAddress("0x0000000000000000000000000000000000000001", big.NewInt(999), token, false, 0)
	if err != nil {
		t.Fatal(err)
	}

	filter := &models.PaginationFilter{
		Skip:  0,
		Limit: 100,
	}

	richList, err := testBackend.GetRichlist(filter)
	if err != nil {
		t.Fatal(err)
	}

	if len(richList) != 2 {
		t.Errorf("Richlist  was incorrect, got: %d, want: %d.", len(richList), 2)
	}

	if richList[0].Address != addrHash {
		t.Errorf("Richlist  was incorrect, got: %s, want: %s.", richList[0].Address, addrHash)
	}
	filter.Skip = 1
	richList, err = testBackend.GetRichlist(filter)
	if err != nil {
		t.Fatal(err)
	}

	if len(richList) != 1 {
		t.Errorf("Richlist  was incorrect, got: %d, want: %d.", len(richList), 1)
	}

}

func TestStats(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	_ = createImportBlock(t, testBackend)

	testBackend.UpdateStats()
	stats, err := testBackend.GetStats()
	if err != nil {
		t.Fatal(err)
	}

	if stats.NumberOfLastDayTransactions != 0 {
		t.Errorf("Wrong number of transactions for 24 hours , got: %d, want: %d.", stats.NumberOfLastDayTransactions, 0)
	}
	if stats.NumberOfLastWeekTransactions != 0 {
		t.Errorf("Wrong number of transactions for last week , got: %d, want: %d.", stats.NumberOfLastWeekTransactions, 0)
	}
	if stats.NumberOfTotalTransactions != 4 {
		t.Errorf("Wrong number of total transactions , got: %d, want: %d.", stats.NumberOfTotalTransactions, 4)
	}

}

func TestTokenHolder(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()

	var token = &tokens.TokenHolderDetails{Balance: big.NewInt(1000000000000000000)}

	addr := models.Address{Address: "addrHash ", TokenSymbol: "GoGo", TokenName: "GoGoTOken"}
	addrHash := "0x0000000000000000000000000000000000000000"
	tokenHolderHash1 := "0x0000000000000000000000000000000000000001"
	tokenHolderHash2 := "0x0000000000000000000000000000000000000002"
	_, err := testBackend.ImportTokenHolder(addrHash, tokenHolderHash1, token, &addr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = testBackend.ImportTokenHolder(addrHash, tokenHolderHash2, token, &addr)
	if err != nil {
		t.Fatal(err)
	}
	filter := &models.PaginationFilter{
		Skip:  0,
		Limit: 100,
	}
	holders, err := testBackend.GetTokenHoldersList(addrHash, filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(holders) != 2 {
		t.Fatalf("HolderList  was incorrect, got: %d, want: %d.", len(holders), 2)
	}

	if holders[0].TokenHolderAddress != tokenHolderHash1 {
		t.Errorf("HolderList  was incorrect, got: %s, want: %s.", holders[0].TokenHolderAddress, tokenHolderHash1)
	}

	if holders[0].Balance != "1000000000000000000" {
		t.Errorf("HolderList  was incorrect, got: %s, want: %s.", holders[0].Balance, "1000000000000000000")
	}

	if holders[0].BalanceInt != 1 { // 1 GO
		t.Errorf("HolderList  was incorrect, got: %d, want: %s.", holders[0].BalanceInt, "1")
	}

	if holders[0].TokenName != addr.TokenName {
		t.Errorf("HolderList  was incorrect, got: %s, want: %s.", holders[0].TokenName, addr.TokenName)
	}

	if holders[0].TokenSymbol != addr.TokenSymbol {
		t.Errorf("HolderList  was incorrect, got: %s, want: %s.", holders[0].TokenSymbol, addr.TokenSymbol)
	}

	if holders[1].TokenHolderAddress != tokenHolderHash2 {
		t.Errorf("HolderList  was incorrect, got: %s, want: %s.", holders[1].TokenHolderAddress, tokenHolderHash2)
	}
	filter.Skip = 1
	holders, err = testBackend.GetTokenHoldersList(addrHash, filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(holders) != 1 {
		t.Errorf("HolderList  was incorrect, got: %d, want: %d.", len(holders), 1)
	}

}

func TestInternalTransactions(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()

	tokenHolderHash1 := "0x0000000000000000000000000000000000000001"
	tokenHolderHash2 := "0x0000000000000000000000000000000000000002"

	var transaction1 = &tokens.TransferEvent{BlockNumber: 10, TransactionHash: "hash1", From: common.HexToAddress(tokenHolderHash1), To: common.HexToAddress(tokenHolderHash2), Value: big.NewInt(10)}
	var transaction2 = &tokens.TransferEvent{BlockNumber: 20, TransactionHash: "hash2", From: common.HexToAddress(tokenHolderHash2), To: common.HexToAddress(tokenHolderHash1), Value: big.NewInt(100)}

	addrHash := "0x0000000000000000000000000000000000000000"

	if _, err := testBackend.ImportTransferEvent(context.Background(), addrHash, transaction1); err != nil {
		t.Fatalf("failed to import internal transaction: %v", err)
	}
	if _, err := testBackend.ImportTransferEvent(context.Background(), addrHash, transaction2); err != nil {
		t.Fatalf("failed to import internal transaction: %v", err)
	}

	filter := &models.PaginationFilter{
		Skip:  0,
		Limit: 100,
	}
	itxFilter := &models.InternalTxFilter{
		InternalAddress: "",
	}
	internalTokenTransfers, err := testBackend.GetInternalTokenTransfers(addrHash, itxFilter)
	if err != nil {
		t.Fatal(err)
	}
	heldTokenTransfers, err := testBackend.GetHeldTokenTransfers(tokenHolderHash1, filter)
	if err != nil {
		t.Fatal(err)
	}
	if len(internalTokenTransfers) != 2 {
		t.Fatalf("InternalTokenTransfers  was incorrect, got: %d, want: %d.", len(internalTokenTransfers), 2)
	}
	if len(heldTokenTransfers) != 2 {
		t.Fatalf("HeldTokenTransfers  was incorrect, got: %d, want: %d.", len(heldTokenTransfers), 2)
	}
	if internalTokenTransfers[0].BlockNumber != transaction2.BlockNumber {
		t.Errorf("InternalTokenTransfers was incorrect, got: %d, want: %d.", internalTokenTransfers[0].BlockNumber, transaction1.BlockNumber)
	}
	if internalTokenTransfers[1].Value != transaction1.Value.String() {
		t.Errorf("InternalTokenTransfers was incorrect, got: %s, want: %s.", internalTokenTransfers[1].Value, transaction1.Value.String())
	}
}

func TestReloadBlock(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	block := createImportBlock(t, testBackend)
	if ok, err := testBackend.NeedReloadParent(block.Header().Number.Int64()); err != nil {
		t.Fatalf("Failed to check block: %v", err)
	} else if !ok {
		t.Error("Block is wrong and should be reloaded")
	}
}

func TestTransactionsConsistent(t *testing.T) {
	testBackend := getBackend(t)
	defer testBackend.mongo.cleanUp()
	block := createImportBlock(t, testBackend)
	if _, ok, err := testBackend.InternalTxsConsistent(block.Header().Number.Int64()); err != nil {
		t.Fatalf("Failed to check block: %v", err)
	} else if !ok {
		t.Error("Number of transactions is not consistent")
	}
}

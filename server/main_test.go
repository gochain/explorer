package main

import (
	"testing"

	"github.com/gochain-io/explorer/server/backend"
	"github.com/gochain-io/gochain/common"
	"github.com/gochain-io/gochain/core/types"
	"github.com/gochain-io/gochain/rlp"
)

var testBackend = backend.NewBackend("127.0.0.1:27017", "https://rpc.gochain.io")

func TestImportAddress(t *testing.T) {
	addrHash := "0x0000000000000000000000000000000000000000"
	address := testBackend.GetAddressByHash(addrHash)

	if address.BalanceWei != "1000" {
		t.Errorf("Balance was incorrect, got: %s, want: %d.", address.BalanceWei, 1000)
	}
	if address.BalanceString != "0.000000000000001000" {
		t.Errorf("Balance was incorrect, got: %s, want: %d.", address.BalanceString, 1000)
	}
}

func TestConvertAddress(t *testing.T) {
	addrHash := "0x0000000000000000000000000000000000000000"
	wrongAddressHash := "0x000"
	address := testBackend.GetAddressByHash(wrongAddressHash)

	if address.Address != addrHash {
		t.Errorf("Address was incorrect, got: %s, want: %s.", address.Address, addrHash)
	}
}

func TestImportBlock(t *testing.T) {
	blockEnc := common.FromHex("0xf90425f90260a0ca12d831416f1fd29336c535ee814d228f638b10798b5cf437bef29415e63762a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d493479448c67d87cd7d716ec044dbe33a0152557bf86062c0c0b84105faed8a008cdb41757cfca7462e440f13f1369a9dc14aafbc1eca77575b6f4753d34458da828f3028fc075e3d6e7b944c3b6b8327dbbd702a1a3ee0c6f8663301a0df5891f347b7525b4369d81238f399a046fb94194e0785d42ba93446b1e27c19a013d863c4c708ff270b60de30faa63e49f073b5ac0e90cdcddb1b10c0736cc923a004deb4be6955e1a300123be48007597f67e4229f8ce70f4f10388de6fd3fa267b90100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000e8314be5c840822d32083014820845b634071a0322e312e33362f6c696e75782d616d6436342f676f312e31302e330000000000a00000000000000000000000000000000000000000000000000000000000000000880000000000000000f901bef86d80847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf8911efc99828edb2600080819ca08f2579f5372bfadca519dfa96b9a7e16689632b65a469c8f96262034d306e104a051692d3de03154e5506afe4616304e6e70bcd7ba7c367b59226045eacd21a8ddf86d01847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf894add3970f31b5f600080819ca0e95246628a5c110356c23a176a97d9dc8a62069b75ef59b090b9074a69eedb63a07882befc79d5e2a5214e33fcc7f2beb45d85551bc4c63588817af2cb83df00b1f86e09847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf8a0b69448bfc2ed424600080819ca08e0645ff002bb43239c22d64cab5dd799b2bb9c56f92d90eb50be7860d4aadd4a00bf2701bb7e1df17921c598954621c724c6a7c50cd1583e117c10708d3c1a537f86e05847735940082520894bf47332f391f995e0e050bc0778a49c61bdd8fdf8a05b7ac452dab97cb600080819ba00a0cacc5f5c556c2435f3e6e43635eda857f71e2f95d73c450ea01a481018011a02cb5ec486060a6f53fd5afdb4f87c57842b7600f3790d66a7bda90fdf0135e93c0")
	var block types.Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatal("decode error: ", err)
	}

	blockFromDb := testBackend.GetBlockByNumber(block.Header().Number.Int64())

	if block.Header().Number.Int64() != blockFromDb.Number || block.Header().Number.Int64() != 1359452 {
		t.Errorf("Block number was incorrect, got: %d, want: %d.", block.Header().Number.Int64(), 1359452)
	}

	if block.Header().Hash().Hex() != blockFromDb.BlockHash {
		t.Errorf("Block hash was incorrect, got: %s, want: %s.", block.Header().Hash().Hex(), blockFromDb.BlockHash)
	}

	if len(block.Transactions()) != blockFromDb.TxCount {
		t.Errorf("Block transactions were incorrect, got: %d, want: %d.", len(block.Transactions()), blockFromDb.TxCount)
	}

	transactionsFromDb := testBackend.GetBlockTransactionsByNumber(block.Header().Number.Int64(), 0, 100)

	if len(block.Transactions()) != len(transactionsFromDb) {
		t.Errorf("Block transactions were incorrect, got: %d, want: %d.", len(block.Transactions()), len(transactionsFromDb))
	}

	if block.Transactions()[0].Hash().Hex() != transactionsFromDb[0].TxHash {
		t.Errorf("Block transaction was incorrect, got: %s, want: %s.", block.Transactions()[0].Hash().Hex(), transactionsFromDb[0].TxHash)
	}

	transactionFromDB := testBackend.GetTransactionByHash(block.Transactions()[0].Hash().Hex())

	if block.Transactions()[0].Hash().Hex() != transactionFromDB.TxHash {
		t.Errorf("Block transaction was incorrect, got: %s, want: %s.", block.Transactions()[0].Hash().Hex(), transactionFromDB.TxHash)
	}

	blockFromDbByHash := testBackend.GetBlockByHash(block.Header().Hash().Hex())

	if block.Header().Number.Int64() != blockFromDbByHash.Number {
		t.Errorf("Block hash was incorrect, got: %d, want: %d.", block.Header().Number.Int64(), blockFromDbByHash.Number)
	}

}

//TODO: cover all following methods:
// func (self *Backend) GetStats() *models.Stats {
// func (self *Backend) GetRichlist(skip, limit int) []*models.Address {
// func (self *Backend) GetTransactionList(address string, skip, limit int) []*models.Transaction {
// func (self *Backend) GetTokenHoldersList(contractAddress string, skip, limit int) []*models.TokenHolder {
// func (self *Backend) GetInternalTransactionsList(contractAddress string, skip, limit int) []*models.InternalTransaction {
// func (self *Backend) GetLatestsBlocks(skip, limit int) []*models.LightBlock {
// func (self *Backend) UpdateStats() {
// func (self *Backend) GetInternalTransactions(address string) []TransferEvent {
// func (self *Backend) NeedReloadBlock(blockNumber int64) bool {
// func (self *Backend) TransactionsConsistent(blockNumber int64) bool {
// func (self *Backend) GetActiveAdresses(fromDate time.Time, onlyContracts bool) []*models.ActiveAddress {
// func (self *Backend) ImportTokenHolder(contractAddress, tokenHolderAddress string, token *TokenDetails) *models.TokenHolder {
// func (self *Backend) ImportInternalTransaction(contractAddress string, transferEvent TransferEvent) *models.InternalTransaction {

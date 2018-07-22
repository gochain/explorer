package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/gochain-io/gochain/common"
	"github.com/gochain-io/gochain/ethclient"
)

func main() {

	client, err := ethclient.Dial("https://rpc.gochain.io")
	if err != nil {
		log.Fatal(err)
	}
	importer := NewImporter(client)
	go listener(client, importer)
	go backfill(client, importer)
	updateAddresses(client, importer)

}

func listener(client *ethclient.Client, importer *ImportMaster) {
	var prevHeader string
	ticker := time.NewTicker(time.Second * 1).C
	for {
		select {
		case <-ticker:
			header, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(header.Number.String())
			if prevHeader != header.Number.String() {
				fmt.Println("Listener is downloading the block:", header.Number.String())
				block, err := client.BlockByNumber(context.Background(), header.Number)
				importer.importBlock(block)
				if err != nil {
					log.Fatal(err)
				}
				checkParentForBlock(client, importer, block.Number().Int64(), 5)
				prevHeader = header.Number.String()
			}
		}
	}
}

func backfill(client *ethclient.Client, importer *ImportMaster) {
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number.String())
	blockNumber := header.Number
	for {
		blocksFromDB := importer.GetBlocksByNumber(blockNumber.String())
		if blocksFromDB == nil {
			fmt.Println("Backfilling the block:", blockNumber.String())
			block, err := client.BlockByNumber(context.Background(), blockNumber)
			importer.importBlock(block)
			if err != nil {
				log.Fatal(err)
			}
		}
		checkParentForBlock(client, importer, blockNumber.Int64(), 5)
		checkTransactionsConsistency(client, importer, blockNumber.Int64())
		blockNumber = big.NewInt(0).Sub(blockNumber, big.NewInt(1))
	}
}

func checkParentForBlock(client *ethclient.Client, importer *ImportMaster, blockNumber int64, numBlocksToCheck int) {
	numBlocksToCheck--
	fmt.Println("Checking the block for it's parent:", blockNumber)
	if importer.needReloadBlock(blockNumber) {
		blockNumber--
		fmt.Println("Redownloading the block because it's corrupted or missing:", blockNumber)
		block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		importer.importBlock(block)
		if err != nil {
			log.Fatal(err)
		}
		if numBlocksToCheck > 0 {
			checkParentForBlock(client, importer, block.Number().Int64(), numBlocksToCheck)
		}
	}
}

func checkTransactionsConsistency(client *ethclient.Client, importer *ImportMaster, blockNumber int64) {
	fmt.Println("Checking a transaction consistency for the block :", blockNumber)
	if !importer.TransactionsConsistent(blockNumber) {
		fmt.Println("Redownloading the block because number of transactions are wrong", blockNumber)
		block, err := client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
		importer.importBlock(block)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func updateAddresses(client *ethclient.Client, importer *ImportMaster) {
	lastUpdatedAt := time.Unix(0, 0)
	for {
		addresses := importer.GetActiveAdresses(lastUpdatedAt)
		fmt.Println("Addresses in db:", len(*addresses), " for date:", lastUpdatedAt)
		for _, address := range *addresses {
			balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Balance of the address:", address, " - ", balance.String())
			importer.importAddress(address, balance)
		}
		lastUpdatedAt = time.Now()
		time.Sleep(120 * time.Second) //sleep for 2 minutes
	}
}

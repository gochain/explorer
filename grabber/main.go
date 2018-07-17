package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	importer := NewImporter()

	client, err := ethclient.Dial("https://rpc.gochain.io")
	if err != nil {
		log.Fatal(err)
	}
	go backfill(client, importer)
	go listener(client, importer)
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
				importer.importBlockIfNotExists(block)
				if err != nil {
					log.Fatal(err)
				}
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
		if len(*blocksFromDB) < 1 {
			fmt.Println("Backfilling the block:", blockNumber.String())
			block, err := client.BlockByNumber(context.Background(), blockNumber)
			importer.importBlockIfNotExists(block)
			if err != nil {
				log.Fatal(err)
			}
		}
		blockNumber = big.NewInt(0).Sub(blockNumber, big.NewInt(1))
	}
}

func updateAddresses(client *ethclient.Client, importer *ImportMaster) {
	lastUpdatedAt := time.Unix(0, 0)
	for {
		addresses := importer.GetActiveAdresses(lastUpdatedAt)
		fmt.Println("Addresses in db:", len(*addresses))
		for _, address := range *addresses {
			balance, err := client.BalanceAt(context.Background(), common.HexToAddress(address), nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("Balance of the address:", address, " - ", balance.String())
			importer.importAddress(address, balance)
		}
		lastUpdatedAt = time.Now()
	}
}

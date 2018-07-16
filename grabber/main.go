package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	importer := NewImporter()

	client, err := ethclient.Dial("https://rpc.gochain.io")
	if err != nil {
		log.Fatal(err)
	}
	// go listener(client, importer)
	backfill(client, importer)

}

func listener(client *ethclient.Client, importer *ImportMaster) {
	var prevHeader string
	ticker := time.NewTicker(time.Second * 1).C
	// go func() {
	for {
		select {
		case <-ticker:
			header, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(header.Number.String())
			if prevHeader != header.Number.String() {
				fmt.Println("Downloading block:", header.Number.String())
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
		blocksFromDB := importer.GetBlockByNumber(blockNumber.String())
		fmt.Println("Blocks in db:", len(*blocksFromDB))
		if len(*blocksFromDB) < 1 {
			fmt.Println("Downloading block:", blockNumber.String())
			block, err := client.BlockByNumber(context.Background(), blockNumber)
			importer.importBlockIfNotExists(block)
			if err != nil {
				log.Fatal(err)
			}
		}
		blockNumber = big.NewInt(0).Sub(blockNumber, big.NewInt(1))
	}
}

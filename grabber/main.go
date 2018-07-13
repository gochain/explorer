package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	importer := NewImporter()
	var prevHeader string
	client, err := ethclient.Dial("https://rpc.gochain.io")
	if err != nil {
		log.Fatal(err)
	}

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
				importer.importBlock(block)
				if err != nil {
					log.Fatal(err)
				}
				prevHeader = header.Number.String()
			}

		}
	}
}

package main

import (
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	app := cli.NewApp()
	app.Name = "sum_fees"
	app.Usage = "Sum gas_fees across all blocks in MongoDB"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "mongo-url",
			Value:  "127.0.0.1:27017",
			Usage:  "MongoDB URL",
			EnvVar: "MONGO_URL",
		},
		cli.StringFlag{
			Name:   "mongo-db",
			Value:  "blocks",
			Usage:  "MongoDB database name",
			EnvVar: "MONGO_DB",
		},
	}

	app.Action = func(c *cli.Context) error {
		logger, _ := zap.NewProduction()
		defer logger.Sync()

		session, err := mgo.Dial(c.String("mongo-url"))
		if err != nil {
			return fmt.Errorf("failed to connect to MongoDB: %v", err)
		}
		defer session.Close()

		db := session.DB(c.String("mongo-db"))

		startTime := time.Now()

		// Query blocks from Darvaza fork (block 17,900,000) onwards where gas_fees is not empty and not "0"
		iter := db.C("Blocks").
			Find(bson.M{"number": bson.M{"$gte": 17900000}, "gas_fees": bson.M{"$nin": []interface{}{"", "0"}}}).
			Select(bson.M{"number": 1, "gas_fees": 1}).
			Iter()
		defer iter.Close()

		var b struct {
			Number  int64  `bson:"number"`
			GasFees string `bson:"gas_fees"`
		}

		totalFees := new(big.Int)
		var blockCount int

		for iter.Next(&b) {
			if fee, ok := new(big.Int).SetString(b.GasFees, 10); ok {
				totalFees.Add(totalFees, fee)
				blockCount++
			}
		}

		if err := iter.Err(); err != nil {
			return fmt.Errorf("iterator error: %v", err)
		}

		// Convert Wei to GO (wei / 10^18)
		weiPerGO := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
		totalGO := new(big.Float).Quo(new(big.Float).SetInt(totalFees), weiPerGO)

		elapsed := time.Since(startTime)
		fmt.Printf("\n==============================================\n")
		fmt.Printf("TOTAL FEES BURNED CALCULATION RESULT:\n")
		fmt.Printf("Blocks with non-zero gas fees: %d\n", blockCount)
		fmt.Printf("Total Fees Burned (Wei): %s\n", totalFees.String())
		fmt.Printf("Total Fees Burned (GO):  %.8f GO\n", totalGO)
		fmt.Printf("Query Duration:          %s\n", elapsed)
		fmt.Printf("==============================================\n\n")

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

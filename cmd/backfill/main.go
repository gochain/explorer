package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/zap"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	app := cli.NewApp()
	app.Name = "backfill"
	app.Usage = "Backfill total_fees_burned on Blocks collection"
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
		cli.IntFlag{
			Name:  "batch-size",
			Value: 10000,
			Usage: "Bulk update batch size",
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

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			logger.Info("Received shutdown signal, stopping gracefully...")
			cancel()
		}()

		var startBlock struct {
			Number          int64  `bson:"number"`
			TotalFeesBurned string `bson:"total_fees_burned"`
		}
		err = db.C("Blocks").
			Find(bson.M{"total_fees_burned": bson.M{"$gt": ""}}).
			Sort("-number").
			Select(bson.M{"number": 1, "total_fees_burned": 1}).One(&startBlock)

		var startNum int64 = 0
		accum := new(big.Int)
		if err == nil && startBlock.TotalFeesBurned != "" && startBlock.TotalFeesBurned != "0" {
			startNum = startBlock.Number
			accum.SetString(startBlock.TotalFeesBurned, 10)
			logger.Info("Starting fee burn backfill", zap.Int64("from_block", startNum), zap.String("start_burned", startBlock.TotalFeesBurned))
		} else {
			logger.Info("No existing fee burn found, starting from block 0")
		}

		iter := db.C("Blocks").
			Find(bson.M{"number": bson.M{"$gt": startNum}}).
			Select(bson.M{"number": 1, "gas_fees": 1}).
			Sort("number").
			Iter()
		defer iter.Close()

		var b struct {
			Number  int64  `bson:"number"`
			GasFees string `bson:"gas_fees"`
		}

		batchSize := c.Int("batch-size")
		bulk := db.C("Blocks").Bulk()
		bulk.Unordered()
		var count, batchCount int
		startTime := time.Now()

		for iter.Next(&b) {
			if ctx.Err() != nil {
				logger.Info("Stopped backfill early", zap.Int("blocks_processed", count), zap.Int64("last_block", b.Number))
				return ctx.Err()
			}
			if b.GasFees != "" && b.GasFees != "0" {
				if fee, ok := new(big.Int).SetString(b.GasFees, 10); ok {
					accum.Add(accum, fee)
				}
			}
			bulk.Update(
				bson.M{"number": b.Number},
				bson.M{"$set": bson.M{"total_fees_burned": accum.String()}},
			)
			count++
			batchCount++

			if count%100000 == 0 {
				elapsed := time.Since(startTime)
				rate := float64(count) / elapsed.Seconds()
				logger.Info("Backfill progress",
					zap.Int("blocks_processed", count),
					zap.Int64("current_block", b.Number),
					zap.String("total_burned", accum.String()),
					zap.Float64("blocks_per_sec", rate),
				)
			}

			if batchCount >= batchSize {
				if _, err := bulk.Run(); err != nil {
					return fmt.Errorf("failed bulk update: %v", err)
				}
				bulk = db.C("Blocks").Bulk()
				bulk.Unordered()
				batchCount = 0
			}
		}
		if batchCount > 0 {
			if _, err := bulk.Run(); err != nil {
				return fmt.Errorf("failed final bulk update: %v", err)
			}
		}

		if err := iter.Err(); err != nil {
			return fmt.Errorf("iterator error: %v", err)
		}

		db.C("Migrations").Upsert(
			bson.M{"ID": 2},
			bson.M{"ID": 2, "comment": "Backfill total_fees_burned on Blocks collection"},
		)

		elapsed := time.Since(startTime)
		logger.Info("Completed total_fees_burned backfill successfully!",
			zap.Int("total_blocks", count),
			zap.String("final_total_burned", accum.String()),
			zap.Duration("duration", elapsed),
		)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

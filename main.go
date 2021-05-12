package main

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kingokeke/go-example-1/binanceUtil"
	"github.com/kingokeke/go-example-1/constants"
	"github.com/kingokeke/go-example-1/mongodbUtil"
	"github.com/kingokeke/go-example-1/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client := mongodbUtil.ConnectToDatabase(&ctx)
	rawStats := client.Database("cs").Collection("raw-stats")

	_, e := rawStats.InsertOne(ctx, mongodbUtil.RawPriceDataStruct{})
	utils.CheckError(e)

	mod := mongo.IndexModel{
		Keys: bson.M{
			"timestamp": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(constants.ONE_MONTH * 7),
	}

	_, e = rawStats.Indexes().CreateOne(ctx, mod)
	utils.CheckError(e)

	c := gocron.NewScheduler(time.UTC)
	c.Cron(constants.CRON_EVERY_FIVE_MINUTES).Do(func() {
		log.Println("Fetching Price Stats...")
		getPriceInfo(&ctx, client, rawStats)
	})

	c.StartBlocking()
}

func getPriceInfo(ctx *context.Context, client *mongo.Client, collection *mongo.Collection) {
	currentTime := time.Now()
	bodyBytes := binanceUtil.GetPriceInfoAsByteArray()

	var priceStatsFromBinance binanceUtil.PriceStats
	json.Unmarshal(bodyBytes, &priceStatsFromBinance)

	stats := []interface{}{}
	for _, stat := range priceStatsFromBinance {
		if strings.HasSuffix(stat.Symbol, constants.BNB) || strings.HasSuffix(stat.Symbol, constants.BTC) || strings.HasSuffix(stat.Symbol, constants.ETH) || strings.HasSuffix(stat.Symbol, constants.USDT) {
			item := mongodbUtil.RawPriceDataStruct{}
			var e error

			item.Symbol = stat.Symbol
			item.Timestamp = currentTime

			item.Price, e = strconv.ParseFloat(stat.Price, 32)
			utils.CheckError(e)

			item.Volume, e = strconv.ParseFloat(stat.Volume, 32)
			utils.CheckError(e)

			stats = append(stats, item)
		}
	}

	_, e := collection.InsertMany(*ctx, stats)
	utils.CheckError(e)
	log.Println("Binance data persisted successfully...")
}

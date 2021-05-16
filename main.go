package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/kingokeke/go-example-1/binanceUtil"
	"github.com/kingokeke/go-example-1/constants"
	"github.com/kingokeke/go-example-1/dotenvUtil"
	"github.com/kingokeke/go-example-1/mongodbUtil"
	"github.com/kingokeke/go-example-1/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	utils.LogToGeneral("App Started...")

	ctx := context.Background()
	client := mongodbUtil.ConnectToDatabase(&ctx)
	defer client.Disconnect(ctx)

	db := client.Database(("cs"))
	timeInSeconds := [5]int32{constants.FIVE_MINUTES, constants.FIFTEEN_MINUTES, constants.ONE_HOUR, constants.FOUR_HOURS, constants.ONE_DAY}
	timeFrames := [5]string{"m5", "m15", "h1", "h4", "d1"}

	collectionNames, e := db.ListCollectionNames(ctx, bson.D{})
	utils.CheckError(e)

	for i, timeframe := range timeFrames {
		collectionExists := false
		collectionName := "raw-data-" + timeframe
		dbCollection := db.Collection(collectionName)
		collectionModel := getModel(timeInSeconds[i])

		for _, name := range collectionNames {
			if name == collectionName {
				collectionExists = true
			}
		}

		if collectionExists == false {
			createTTLIndexedTablesInDB(&ctx, dbCollection, &collectionModel)
		}
	}

	c := gocron.NewScheduler(time.UTC)
	c.Cron(constants.CRON_EVERY_FIVE_MINUTES).Do(func() {
		getPriceInfo(&ctx, client, db)
	})
	
	c.StartBlocking()
	}

func getPriceInfo(ctx *context.Context, client *mongo.Client, database *mongo.Database) {
	currentTime := time.Now().UTC()
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

	_, e := database.Collection("raw-data-m5").InsertMany(*ctx, stats)
	utils.CheckError(e)
	
	if currentTime.Minute() == 0 || currentTime.Minute() % 15 == 0 {
		_, e := database.Collection("raw-data-m15").InsertMany(*ctx, stats)
		utils.CheckError(e)
	}

	if currentTime.Minute() == 0 {
		_, e := database.Collection("raw-data-h1").InsertMany(*ctx, stats)
		utils.CheckError(e)
	}

	if (currentTime.Hour() == 0 || currentTime.Hour() % 4 == 0) && currentTime.Minute() == 0 {
		_, e := database.Collection("raw-data-h4").InsertMany(*ctx, stats)
		utils.CheckError(e)
	}

	if currentTime.Hour() == 0 && currentTime.Minute() == 0 {
		_, e := database.Collection("raw-data-d1").InsertMany(*ctx, stats)
		utils.CheckError(e)
	}

	_, e = http.Get(dotenvUtil.GetValue("MAIN_SERVICE_URL"))
	if e != nil {}

	utils.LogToGeneral("Successfully persisted raw price data...")
}

func getModel(numberOfSeconds int32) mongo.IndexModel {
	return mongo.IndexModel{
		Keys: bson.M{
			"timestamp": 1,
		},
		Options: options.Index().SetExpireAfterSeconds(numberOfSeconds * 300),
	}
}

func createTTLIndexedTablesInDB(appContext *context.Context, collectionName *mongo.Collection, collectionModel *mongo.IndexModel) {
		_, e := collectionName.InsertOne(*appContext, mongodbUtil.RawPriceDataStruct{})
	utils.CheckError(e)

	_, e = collectionName.Indexes().CreateOne(*appContext, *collectionModel)
	utils.CheckError(e)
}

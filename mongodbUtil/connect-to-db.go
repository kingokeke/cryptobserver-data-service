package mongodbUtil

import (
	"context"

	"github.com/kingokeke/go-example-1/dotenvUtil"
	"github.com/kingokeke/go-example-1/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToDatabase(ctx *context.Context) *mongo.Client {
	mongoDbUrl := dotenvUtil.GetValue("MONGODB_URL")
	clientOptions := options.Client().ApplyURI(mongoDbUrl)

	client, e := mongo.Connect(*ctx, clientOptions)
	utils.CheckError(e)
	utils.LogToGeneral("Successfully connected to database...")

	e = client.Ping(*ctx, nil)
	utils.CheckError(e)
	utils.LogToGeneral("Successfully pinged database...")

	return client
}

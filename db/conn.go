package db

import (
	"analytics-api/configs"
	"context"
	"os"

	"analytics-api/internal/pkg/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient open new client to mongodb
func NewClient() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(os.Getenv("URI")).SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.LogError(context.TODO(), err)
	}
	configs.Database.Client = client.Database(os.Getenv("NAME"))
}

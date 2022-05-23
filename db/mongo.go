package db

import (
	"analytics-api/configs"
	"context"

	"analytics-api/internal/pkg/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongo open new client to mongodb
func NewMongo() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(configs.MongoDB.URI).SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.LogError(context.TODO(), err)
	}
	configs.MongoDB.Client = client.Database(configs.MongoDB.Name)
}

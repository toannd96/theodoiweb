package db

import (
	"context"

	"analytics-api/configs"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// NewMongo open new client to mongodb
func NewMongo() {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(configs.MongoDB.URI).SetServerAPIOptions(serverAPIOptions)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logrus.Error(context.TODO(), err)
	}
	configs.MongoDB.Client = client.Database(configs.MongoDB.Name)
}

// CreateSessionCollection create timeseries session collection if not exists
func CreateSessionCollection() error {
	exists, err := checkCollection(configs.MongoDB.SessionCollection)
	if err != nil {
		return err
	}
	if !exists {
		logrus.Info("not exists, create collection name ", configs.MongoDB.SessionCollection)
		ts := options.
			TimeSeries().
			SetMetaField("meta_data").
			SetTimeField("time_report").
			SetGranularity("seconds")
		opts := options.
			CreateCollection().
			SetTimeSeriesOptions(ts).
			SetExpireAfterSeconds(180 * 86400)
		err := configs.MongoDB.Client.CreateCollection(context.TODO(), configs.MongoDB.SessionCollection, opts)
		if err != nil {
			return err
		}

		// created index
		models := []mongo.IndexModel{
			{
				Keys: bson.M{"meta_data.user_id": 1},
			},
			{
				Keys: bson.M{"meta_data.id": 1},
			},
			{
				Keys: bson.M{"meta_data.website_id": 1},
			},
		}

		collection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
		_, CreateIndexErr := collection.Indexes().CreateMany(context.Background(), models)
		if CreateIndexErr != nil {
			return err
		}
	} else {
		logrus.Debug("collection exists")
	}
	return nil
}

func CreateUserCollection() error {
	exists, err := checkCollection(configs.MongoDB.UserCollection)
	if err != nil {
		return err
	}
	if !exists {
		logrus.Info("not exists, create collection name ", configs.MongoDB.UserCollection)
		models := []mongo.IndexModel{
			{
				Keys: bson.M{"id": 1},
			},
			{
				Keys: bson.M{"email": 1},
			},
		}

		collection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
		_, err := collection.Indexes().CreateMany(context.Background(), models)
		if err != nil {
			return err
		}
	} else {
		logrus.Debug("collection exists")
	}
	return nil
}

func CreateWebsiteCollection() error {
	exists, err := checkCollection(configs.MongoDB.WebsiteCollection)
	if err != nil {
		return err
	}
	if !exists {
		logrus.Info("not exists, create collection name ", configs.MongoDB.WebsiteCollection)
		models := []mongo.IndexModel{
			{
				Keys: bson.M{"user_id": 1},
			},
			{
				Keys: bson.M{"id": 1},
			},
		}

		collection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
		_, err := collection.Indexes().CreateMany(context.Background(), models)
		if err != nil {
			return err
		}
	} else {
		logrus.Debug("collection exists")
	}
	return nil
}

// checkCollection check collection exists or not exists
func checkCollection(name string) (bool, error) {
	var exists bool = false
	filter := bson.M{}

	names, err := configs.MongoDB.Client.ListCollectionNames(context.TODO(), filter, nil)
	if err != nil {
		return false, err
	}

	for _, nm := range names {
		if nm == name {
			exists = true
			break
		}
	}
	return exists, nil
}

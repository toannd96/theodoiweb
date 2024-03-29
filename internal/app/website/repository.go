package website

import (
	"context"

	"analytics-api/configs"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	FindWebsite(userID, hostName string) (int64, error)
	FindWebsiteByID(userID, websiteID string) (int64, error)
	InsertWebsite(userID string, website website) error
	GetWebsite(userID, websiteID string, website *website) error
	GetAllWebsite(userID string) (*websites, error)
	DeleteWebsite(userID, websiteID string) error
	DeleteSession(userID, websiteID string) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

func (instance *repository) FindWebsite(userID, hostName string) (int64, error) {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"host_name": hostName},
	}}
	count, err := websiteCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *repository) FindWebsiteByID(userID, websiteID string) (int64, error) {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	count, err := websiteCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *repository) InsertWebsite(userID string, aWebsite website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	docs := website{
		ID:        aWebsite.ID,
		UserID:    userID,
		Category:  aWebsite.Category,
		HostName:  aWebsite.HostName,
		URL:       aWebsite.URL,
		CreatedAt: aWebsite.CreatedAt,
		UpdatedAt: aWebsite.UpdatedAt,
	}
	_, err := websiteCollection.InsertOne(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetWebsite(userID, websiteID string, aWebsite *website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	err := websiteCollection.FindOne(context.TODO(), filter).Decode(&aWebsite)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetAllWebsite(userID string) (*websites, error) {
	var websites websites
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"user_id": userID}
	cursor, err := websiteCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &websites); err != nil {
		return nil, err
	}
	return &websites, nil
}

func (instance *repository) DeleteWebsite(userID, websiteID string) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	deleteResult, err := websiteCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	logrus.Printf("deleted %v documents in the website collection\n", deleteResult.DeletedCount)
	return nil
}

func (instance *repository) DeleteSession(userID, websiteID string) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.website_id": websiteID},
	}}
	deleteResult, err := sessionCollection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return err
	}
	logrus.Printf("deleted %v documents in the session collection\n", deleteResult.DeletedCount)
	return nil
}

package website

import (
	"context"

	"analytics-api/configs"
	"analytics-api/models"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	FindWebsite(userID, hostName string) (int64, error)
	FindWebsiteByID(userID, websiteID string) (int64, error)
	InsertWebsite(userID string, website models.Website) error
	GetWebsite(userID, websiteID string, website *models.Website) error
	GetAllWebsite(userID string) (*models.Websites, error)
	UpdateNameWebsite(userID, websiteID string, website *models.Website) error
	UpdateURLWebsite(userID, websiteID string, website *models.Website) error
	UpdateTrackedWebsite(userID, websiteID string, website *models.Website) error
	UpdateWebsite(userID, websiteID string, website *models.Website) error
	DeleteWebsite(userID, websiteID string) error
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

func (instance *repository) InsertWebsite(userID string, website models.Website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	docs := models.Website{
		ID:        website.ID,
		UserID:    userID,
		Name:      website.Name,
		HostName:  website.HostName,
		URL:       website.URL,
		Tracked:   website.Tracked,
		CreatedAt: website.CreatedAt,
		UpdatedAt: website.UpdatedAt,
	}
	_, err := websiteCollection.InsertOne(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetWebsite(userID, websiteID string, website *models.Website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	err := websiteCollection.FindOne(context.TODO(), filter).Decode(&website)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetAllWebsite(userID string) (*models.Websites, error) {
	var websites models.Websites
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

func (instance *repository) UpdateNameWebsite(userID, websiteID string, website *models.Website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	update := bson.M{
		"$set": bson.M{
			"name":       website.Name,
			"updated_at": website.UpdatedAt,
		},
	}
	result := websiteCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: %v\n", result.Err())
	}
	return nil
}

func (instance *repository) UpdateURLWebsite(userID, websiteID string, website *models.Website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	update := bson.M{
		"$set": bson.M{
			"url":        website.URL,
			"updated_at": website.UpdatedAt,
		},
	}
	result := websiteCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: %v\n", result.Err())
	}
	return nil
}

func (instance *repository) UpdateTrackedWebsite(userID, websiteID string, website *models.Website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	update := bson.M{
		"$set": bson.M{
			"tracked":    website.Tracked,
			"updated_at": website.UpdatedAt,
		},
	}
	result := websiteCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: %v\n", result.Err())
	}
	return nil
}

func (instance *repository) UpdateWebsite(userID, websiteID string, website *models.Website) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	update := bson.M{
		"$set": bson.M{
			"name":       website.Name,
			"url":        website.URL,
			"updated_at": website.UpdatedAt,
		},
	}
	result := websiteCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: %v\n", result.Err())
	}
	return nil
}

func (instance *repository) DeleteWebsite(userID, websiteID string) error {
	websiteCollection := configs.MongoDB.Client.Collection(configs.MongoDB.WebsiteCollection)
	filter := bson.M{"$and": []bson.M{
		{"user_id": userID},
		{"id": websiteID},
	}}
	_, err := websiteCollection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

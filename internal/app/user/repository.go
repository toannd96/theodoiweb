package user

import (
	"context"

	"analytics-api/configs"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	FindUser(email string) (int64, error)
	InsertUser(user user) error
	GetUserByEmail(email string, user *user) error
	GetUserByID(userID string, user *user) error
	UpdateUser(userID string, user *user) error
	UpdateFullName(userID string, user *user) error
	UpdatePassword(userID string, user *user) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

func (instance *repository) FindUser(email string) (int64, error) {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"email": email}
	count, err := userCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *repository) InsertUser(anUser user) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	_, err := userCollection.InsertOne(context.TODO(), anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetUserByEmail(email string, anUser *user) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"email": email}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetUserByID(userID string, anUser *user) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&anUser)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) UpdateFullName(userID string, anUser *user) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{
			"full_name":  anUser.FullName,
			"updated_at": anUser.UpdatedAt,
		},
	}
	result := userCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: ", result.Err())
	}
	return nil
}

func (instance *repository) UpdatePassword(userID string, anUser *user) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{
			"password":   anUser.Password,
			"updated_at": anUser.UpdatedAt,
		},
	}
	result := userCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: ", result.Err())
	}
	return nil
}

func (instance *repository) UpdateUser(userID string, anUser *user) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{
			"full_name":  anUser.FullName,
			"password":   anUser.Password,
			"updated_at": anUser.UpdatedAt,
		},
	}
	result := userCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: ", result.Err())
	}
	return nil
}

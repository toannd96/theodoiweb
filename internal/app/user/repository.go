package user

import (
	"context"

	"analytics-api/configs"
	"analytics-api/models"

	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	FindUser(email string) (int64, error)
	InsertUser(user models.User) error
	GetUserByEmail(email string, user *models.User) error
	GetUserByID(userID string, user *models.User) error
	UpdateUser(userID string, user *models.User) error
	UpdateFullName(userID string, user *models.User) error
	UpdatePassword(userID string, user *models.User) error
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

func (instance *repository) InsertUser(user models.User) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	_, err := userCollection.InsertOne(context.TODO(), user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetUserByEmail(email string, user *models.User) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"email": email}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) GetUserByID(userID string, user *models.User) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return err
	}
	return nil
}

func (instance *repository) UpdateFullName(userID string, user *models.User) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{
			"full_name":  user.FullName,
			"updated_at": user.UpdatedAt,
		},
	}
	result := userCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: ", result.Err())
	}
	return nil
}

func (instance *repository) UpdatePassword(userID string, user *models.User) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		},
	}
	result := userCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: ", result.Err())
	}
	return nil
}

func (instance *repository) UpdateUser(userID string, user *models.User) error {
	userCollection := configs.MongoDB.Client.Collection(configs.MongoDB.UserCollection)
	filter := bson.M{"id": userID}
	update := bson.M{
		"$set": bson.M{
			"full_name":  user.FullName,
			"password":   user.Password,
			"updated_at": user.UpdatedAt,
		},
	}
	result := userCollection.FindOneAndUpdate(context.Background(), filter, update)
	if result.Err() != nil {
		logrus.Error("update failed: ", result.Err())
	}
	return nil
}

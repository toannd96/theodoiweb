package session

import (
	"context"

	"analytics-api/configs"
	"analytics-api/models"

	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	GetAllSession(sessions models.Sessions) (models.Sessions, error)
	GetSessionByID(sessionID string, session *models.Session) error
	Insert(session models.Session) error
	FindSessionID(id string) (int64, error)
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetSessionByID get session by session id
func (instance *repository) GetSessionByID(sessionID string, session *models.Session) error {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	err := sessionCollection.FindOne(context.TODO(), bson.M{"id": sessionID}).Decode(&session)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *repository) GetAllSession(sessions models.Sessions) (models.Sessions, error) {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	cursor, err := sessionCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(context.TODO(), &sessions); err != nil {
		return nil, err
	}
	return sessions, nil
}

// Insert insert session
func (instance *repository) Insert(session models.Session) error {
	_, err := configs.Database.Client.Collection(configs.Database.SessionCollection).InsertOne(context.TODO(), session)
	if err != nil {
		return err
	}
	return nil
}

// FindSessionID find session id to check exists
func (instance *repository) FindSessionID(id string) (int64, error) {
	count, err := configs.Database.Client.Collection(configs.Database.SessionCollection).CountDocuments(context.TODO(), bson.M{"id": id})
	if err != nil {
		return 0, err
	}
	return count, nil
}

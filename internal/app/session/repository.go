package session

import (
	"context"
	"log"

	"analytics-api/configs"
	"analytics-api/models"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	GetAllSession(sessions models.Sessions) (models.Sessions, error)
	GetSession(id string, session *models.Session) error
	InsertSession(session models.Session) error
	FindSession(id string) (int64, error)

	GetCountEvent(id string, session *models.Session) (int, error)
	GetEvent(id string, session *models.Session, limit, skip int) error
	UpdateEvent(session models.Session) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetSession get session by session id
func (instance *repository) GetSession(id string, session *models.Session) error {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	err := sessionCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&session)
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

// InsertSession insert session
func (instance *repository) InsertSession(session models.Session) error {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	_, err := sessionCollection.InsertOne(context.TODO(), session)
	if err != nil {
		return err
	}
	return nil
}

// FindSession find session id to check exists
func (instance *repository) FindSession(id string) (int64, error) {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	count, err := sessionCollection.CountDocuments(context.TODO(), bson.M{"id": id})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetCountEvent get count event of session by session id
func (instance *repository) GetCountEvent(id string, session *models.Session) (int, error) {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	err := sessionCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&session)
	if err != nil {
		return 0, err
	}
	return len(session.Events), nil
}

// GetEvent get limit event of session by session id
func (instance *repository) GetEvent(id string, session *models.Session, limit, skip int) error {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	filter := bson.M{"id": id}
	opt := options.FindOneOptions{
		Projection: bson.M{"events": bson.M{"$slice": []int{skip, limit}}},
	}
	err := sessionCollection.FindOne(context.TODO(), filter, &opt).Decode(&session)
	if err != nil {
		return err
	}
	return nil
}

// UpdateEvent update event of session by session id
func (instance *repository) UpdateEvent(session models.Session) error {
	sessionCollection := configs.Database.Client.Collection(configs.Database.SessionCollection)
	upsert := true
	for _, event := range session.Events {
		filter := bson.M{"id": session.ID}
		update := bson.M{
			"$push": bson.M{"events": event},
		}
		opt := options.FindOneAndUpdateOptions{
			Upsert: &upsert,
		}
		result := sessionCollection.FindOneAndUpdate(context.Background(), filter, update, &opt)
		if result.Err() != nil {
			log.Printf("update failed: %v\n", result.Err())
		}
	}
	return nil
}

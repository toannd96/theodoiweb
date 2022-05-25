package session

import (
	"context"
	"log"
	"strconv"
	"time"

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
	UpdateSession(id string, session models.Session) error

	GetEvents(id string, session *models.Session) (models.Events, error)
	GetEventByLimitSkip(id string, session *models.Session, limit, skip int) error

	GetSessionTimestamp(id string) (int64, error)
	InsertSessionTimestamp(id string, timeStart int64) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetSession get session by session id
func (instance *repository) GetSession(id string, session *models.Session) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	err := sessionCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&session)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *repository) GetAllSession(sessions models.Sessions) (models.Sessions, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
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
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	_, err := sessionCollection.InsertOne(context.TODO(), session)
	if err != nil {
		return err
	}
	return nil
}

// FindSession find session id to check exists
func (instance *repository) FindSession(id string) (int64, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	count, err := sessionCollection.CountDocuments(context.TODO(), bson.M{"id": id})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetEvents get event of session by session id
func (instance *repository) GetEvents(id string, session *models.Session) (models.Events, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	err := sessionCollection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&session)
	if err != nil {
		return nil, err
	}
	return session.Events, nil
}

// GetEventByLimitSkip get limit event of session by session id
func (instance *repository) GetEventByLimitSkip(id string, session *models.Session, limit, skip int) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
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

// UpdateSession update duration session and event by session id
func (instance *repository) UpdateSession(id string, session models.Session) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	upsert := true
	filter := bson.M{"id": id}
	for _, event := range session.Events {
		update := bson.M{
			"$set":  bson.M{"duration": session.Duration},
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

// InsertSessionTimestamp insert first timestamp by session id
func (instance *repository) InsertSessionTimestamp(id string, timeStart int64) error {
	err := configs.Redis.Client.Set(id, timeStart, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetSessionTimestamp get first timestamp by session id
func (instance *repository) GetSessionTimestamp(id string) (int64, error) {
	timeStartStr, err := configs.Redis.Client.Get(id).Result()
	if err != nil {
		return 0, err
	}
	timeStart, err := strconv.ParseInt(timeStartStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return timeStart, nil
}

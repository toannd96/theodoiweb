package session

import (
	"context"
	"strconv"
	"time"

	"analytics-api/configs"
	"analytics-api/models"
	"analytics-api/pkg"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	GetAllSession(listSessionID []string, session models.Session) ([]models.MetaData, error)
	GetAllSessionID(sessionMetaData models.MetaData) ([]string, error)
	GetSession(sessionID string, session *models.Session) error
	GetCountSession(sessionID string) (int64, error)
	InsertSession(session models.Session, event models.Event) error

	GetEventByLimitSkip(sessionID string, limit, skip int) ([]*models.Event, error)

	GetSessionTimestamp(sessionID string) (int64, error)
	InsertSessionTimestamp(sessionID string, timeStart int64) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetSession get session by session id
func (instance *repository) GetSession(sessionID string, session *models.Session) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	err := sessionCollection.FindOne(context.TODO(), bson.M{"id": sessionID}).Decode(&session.MetaData)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *repository) GetAllSession(listSessionID []string, session models.Session) ([]models.MetaData, error) {
	var sessionMetaData []models.MetaData
	opt := options.FindOne()
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)

	for _, sessionID := range listSessionID {
		count, err := sessionCollection.CountDocuments(context.TODO(), bson.M{"id": sessionID})
		if err != nil {
			return nil, err
		}

		opt.SetSkip(count - 1)
		err = sessionCollection.FindOne(context.TODO(), bson.M{"id": sessionID}, opt).Decode(&session.MetaData)
		if err != nil {
			return nil, err
		}
		sessionMetaData = append(sessionMetaData, session.MetaData)
	}
	return sessionMetaData, nil
}

// GetAllSessionID get all id of session
func (instance *repository) GetAllSessionID(sessionMetaData models.MetaData) ([]string, error) {
	var listSessionID []string

	fromDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.UTC)
	toDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 24, 0, 0, 0, time.UTC)

	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	cursor, err := sessionCollection.Find(context.TODO(), bson.M{"created_at": bson.M{
		"$gt": fromDate,
		"$lt": toDate,
	}})
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&sessionMetaData)
		if err != nil {
			return nil, err
		}
		listSessionID = append(listSessionID, sessionMetaData.ID)
	}
	listSessionID = pkg.RemoveDuplicateValues(listSessionID)
	return listSessionID, nil
}

// InsertSession insert session
func (instance *repository) InsertSession(session models.Session, event models.Event) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	docs := bson.M{
		"id":         session.MetaData.ID,
		"website_id": session.MetaData.WebsiteID,
		"country":    session.MetaData.Country,
		"city":       session.MetaData.City,
		"device":     session.MetaData.Device,
		"os":         session.MetaData.OS,
		"browser":    session.MetaData.Browser,
		"version":    session.MetaData.Version,
		"duration":   session.MetaData.Duration,
		"created":    session.MetaData.Created,
		"created_at": session.CreatedAt,
		"event":      event,
	}
	_, err := sessionCollection.InsertOne(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}

// GetCountSession get count session of session id
func (instance *repository) GetCountSession(sessionID string) (int64, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	count, err := sessionCollection.CountDocuments(context.TODO(), bson.M{"id": sessionID})
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *repository) GetEventByLimitSkip(sessionID string, limit, skip int) ([]*models.Event, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)

	filter := bson.M{"id": sessionID}
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip)).SetLimit(int64(limit))

	cur, err := sessionCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		return nil, err
	}

	var events []*models.Event
	for cur.Next(context.TODO()) {
		var session models.Session
		err := cur.Decode(&session)
		if err != nil {
			return nil, err
		}
		events = append(events, &session.Event)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}

	defer cur.Close(context.TODO())
	return events, nil
}

// InsertSessionTimestamp insert first timestamp by session id
func (instance *repository) InsertSessionTimestamp(sessionID string, timeStart int64) error {
	err := configs.Redis.Client.Set(sessionID, timeStart, 24*time.Hour).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetSessionTimestamp get first timestamp by session id
func (instance *repository) GetSessionTimestamp(sessionID string) (int64, error) {
	timeStartStr, err := configs.Redis.Client.Get(sessionID).Result()
	if err != nil {
		return 0, err
	}
	timeStart, err := strconv.ParseInt(timeStartStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return timeStart, nil
}

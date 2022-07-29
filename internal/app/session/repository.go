package session

import (
	"context"
	"strconv"
	"time"

	"analytics-api/configs"
	str "analytics-api/internal/pkg/string"
	"analytics-api/models"

	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

// Repository ...
type Repository interface {
	GetAllSession(userID, websiteID string, listSessionID []string, session models.Session) ([]models.Session, error)
	GetAllSessionID(userID, websiteID string, session models.Session) ([]string, error)
	GetSession(userID, sessionID string, session *models.Session) error
	GetCountSession(userID, sessionID string) (int64, error)
	InsertSession(session models.Session, event models.Event) error

	GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*models.Event, error)

	GetSessionTimestamp(sessionID string) (int64, error)
	InsertSessionTimestamp(sessionID string, timeStart int64) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetSession get session by session id
func (instance *repository) GetSession(userID, sessionID string, session *models.Session) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.id": sessionID},
	}}
	err := sessionCollection.FindOne(context.TODO(), filter).Decode(&session)
	if err != nil {
		return err
	}
	return nil
}

// GetAllSession get all session
func (instance *repository) GetAllSession(userID, websiteID string, listSessionID []string, session models.Session) ([]models.Session, error) {
	var listSession []models.Session
	opt := options.FindOne()
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)

	for _, sessionID := range listSessionID {
		count, err := sessionCollection.CountDocuments(context.TODO(), bson.M{"$and": []bson.M{
			{"meta_data.id": sessionID},
			{"meta_data.website_id": websiteID},
			{"meta_data.user_id": userID},
		}})
		if err != nil {
			return nil, err
		}
		opt.SetSkip(count - 1)
		err = sessionCollection.FindOne(context.TODO(), bson.M{"$and": []bson.M{
			{"meta_data.id": sessionID},
			{"meta_data.website_id": websiteID},
			{"meta_data.user_id": userID},
		}}, opt).Decode(&session)
		if err != nil {
			return nil, err
		}
		listSession = append(listSession, session)
	}
	return listSession, nil
}

// GetAllSessionID get all id of session
func (instance *repository) GetAllSessionID(userID, websiteID string, session models.Session) ([]string, error) {
	var listSessionID []string
	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.website_id": websiteID},
	}}
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	cursor, err := sessionCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(context.TODO()) {
		err := cursor.Decode(&session)
		if err != nil {
			return nil, err
		}
		listSessionID = append(listSessionID, session.MetaData.ID)
	}
	listSessionID = str.RemoveDuplicateValues(listSessionID)
	return listSessionID, nil
}

// InsertSession insert session
func (instance *repository) InsertSession(session models.Session, event models.Event) error {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	docs := models.Session{
		MetaData: models.MetaData{
			ID:        session.MetaData.ID,
			UserID:    session.MetaData.UserID,
			WebsiteID: session.MetaData.WebsiteID,
			Country:   session.MetaData.Country,
			City:      session.MetaData.City,
			Device:    session.MetaData.Device,
			OS:        session.MetaData.OS,
			Browser:   session.MetaData.Browser,
			Version:   session.MetaData.Version,
			CreatedAt: session.MetaData.CreatedAt,
		},
		Duration:   session.Duration,
		Event:      event,
		TimeReport: session.TimeReport,
	}
	_, err := sessionCollection.InsertOne(context.TODO(), docs)
	if err != nil {
		return err
	}
	return nil
}

// GetCountSession get count session of session id
func (instance *repository) GetCountSession(userID, sessionID string) (int64, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)
	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.id": sessionID},
	}}
	count, err := sessionCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (instance *repository) GetEventByLimitSkip(userID, sessionID string, limit, skip int) ([]*models.Event, error) {
	sessionCollection := configs.MongoDB.Client.Collection(configs.MongoDB.SessionCollection)

	filter := bson.M{"$and": []bson.M{
		{"meta_data.user_id": userID},
		{"meta_data.id": sessionID},
	}}
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

package event

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
	GetCountEventBySessionID(sessionID string, events *models.Events) (int, error)
	GetEventBySessionID(sessionID string, events *models.Events, limit, skip int) error
	Insert(sessionID string) error
	FindSessionID(id string) (int64, error)
	Update(events models.Events) error
}

type repository struct{}

// NewRepository ...
func NewRepository() Repository {
	return &repository{}
}

// GetCountEventBySessionID get count event of session by session id
func (instance *repository) GetCountEventBySessionID(sessionID string, events *models.Events) (int, error) {
	eventCollection := configs.Database.Client.Collection(configs.Database.EventCollection)
	err := eventCollection.FindOne(context.TODO(), bson.M{"session_id": sessionID}).Decode(&events)
	if err != nil {
		return 0, err
	}
	return len(events.Events), nil
}

// GetEventBySessionID get limit event of session by session id
func (instance *repository) GetEventBySessionID(sessionID string, events *models.Events, limit, skip int) error {
	eventCollection := configs.Database.Client.Collection(configs.Database.EventCollection)
	filter := bson.M{"session_id": sessionID}
	opt := options.FindOneOptions{
		Projection: bson.M{"events": bson.M{"$slice": []int{skip, limit}}},
	}
	err := eventCollection.FindOne(context.TODO(), filter, &opt).Decode(&events)
	if err != nil {
		return err
	}
	return nil
}

// Insert insert session id by event
func (instance *repository) Insert(sessionID string) error {
	eventCollection := configs.Database.Client.Collection(configs.Database.EventCollection)
	_, err := eventCollection.InsertOne(context.TODO(), bson.M{"session_id": sessionID})
	if err != nil {
		return err
	}
	return nil
}

// FindSessionID find session id if not exists count is 0
func (instance *repository) FindSessionID(id string) (int64, error) {
	eventCollection := configs.Database.Client.Collection(configs.Database.EventCollection)
	count, err := eventCollection.CountDocuments(context.TODO(), bson.M{"session_id": id})
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Update update event of session by session id
func (instance *repository) Update(events models.Events) error {
	upsert := true
	eventCollection := configs.Database.Client.Collection(configs.Database.EventCollection)
	for _, event := range events.Events {
		filter := bson.M{"session_id": events.SessionID}
		update := bson.M{
			"$push": bson.M{"events": event},
		}
		opt := options.FindOneAndUpdateOptions{
			Upsert: &upsert,
		}
		result := eventCollection.FindOneAndUpdate(context.Background(), filter, update, &opt)
		if result.Err() != nil {
			log.Printf("update failed: %v\n", result.Err())
		}
	}
	return nil
}

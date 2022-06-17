package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Session ...
type Session struct {
	MetaData   MetaData  `json:"meta_data" bson:"meta_data"`
	Duration   string    `json:"duration" bson:"duration"`
	Event      Event     `json:"event" bson:"event"`
	TimeReport time.Time `json:"time_report" bson:"time_report"`
}

// MetaData ...
type MetaData struct {
	ID        string `json:"id" bson:"id"`
	WebsiteID string `json:"website_id" bson:"website_id"`
	Country   string `json:"country" bson:"country"`
	City      string `json:"city" bson:"city"`
	Device    string `json:"device" bson:"device"`
	OS        string `json:"os" bson:"os"`
	Browser   string `json:"browser" bson:"browser"`
	Version   string `json:"version" bson:"version"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

// Event ...
type Event struct {
	Type      int64  `json:"type" bson:"type"`
	Data      bson.M `json:"data" bson:"data"`
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
}

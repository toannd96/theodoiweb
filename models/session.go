package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Session ...
type Session struct {
	MetaData  MetaData  `json:"meta_data"`
	CreatedAt time.Time `json:"created_at"`
	Event     Event     `json:"event"`
}

// MetaData ...
type MetaData struct {
	ID        string `json:"id"`
	WebsiteID string `json:"websiteID"`
	Country   string `json:"country"`
	City      string `json:"city"`
	Device    string `json:"device"`
	OS        string `json:"os"`
	Browser   string `json:"browser"`
	Version   string `json:"version"`
	Duration  string `json:"duration"`
	Created   string `json:"created"`
}

// Event ...
type Event struct {
	Type      int64  `json:"type"`
	Data      bson.M `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

package models

import "gopkg.in/mgo.v2/bson"

// Event ...
type Event struct {
	Type      int64  `json:"type"`
	Data      bson.M `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

// Session ...
type Session struct {
	ID        string  `json:"id"`
	WebsiteID string  `json:"website_id"`
	UserAgent string  `json:"user_agent"`
	Browser   string  `json:"browser"`
	Version   string  `json:"version"`
	OS        string  `json:"os"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Events    []Event `json:"events"`
}

// Sessions ...
type Sessions []Session

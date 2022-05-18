package models

import "gopkg.in/mgo.v2/bson"

// Events ...
type Events struct {
	SessionID string  `json:"session_id"`
	Events    []Event `json:"events"`
}

// Event ...
type Event struct {
	Type      int64  `json:"type"`
	Data      bson.M `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

package models

import (
	"gopkg.in/mgo.v2/bson"
)

// Sessions ...
type Sessions []Session

// Session ...
type Session struct {
	ID             string `json:"id"`
	WebsiteID      string `json:"website_id"`
	Country        string `json:"country"`
	City           string `json:"city"`
	Device         string `json:"device"`
	OS             string `json:"os"`
	Browser        string `json:"browser"`
	BrowserVersion string `json:"browser_version"`
	Duration       string `json:"duration"`
	CreatedAt      string `json:"created_at"`
	Events         Events `json:"events"`
}

// Events ...
type Events []Event

// Event ...
type Event struct {
	Type      int64  `json:"type"`
	Data      bson.M `json:"data"`
	Timestamp int64  `json:"timestamp"`
}

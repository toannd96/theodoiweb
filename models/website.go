package models

import "time"

// Website ...
type Website struct {
	ID        string    `json:"id" bson:"id"`
	UserID    string    `json:"user_id" bson:"user_id"`
	Name      string    `json:"name" bson:"name"`
	URL       string    `json:"url" bson:"url"`
	Tracked   bool      `json:"tracked" bson:"tracked"`
	Sessions  []Session `json:"sessions" bson:"sessions"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

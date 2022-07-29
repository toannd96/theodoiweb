package models

import (
	"time"
)

// User ...
type User struct {
	ID           string    `json:"id" bson:"id"`
	FullName     string    `json:"full_name" bson:"full_name"`
	Password     string    `json:"password" bson:"password"`
	Email        string    `json:"email" bson:"email"`
	AccessToken  string    `json:"-" bson:"-"`
	RefreshToken string    `json:"-" bson:"-"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

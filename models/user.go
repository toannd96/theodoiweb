package models

import (
	"time"
)

// User ...
type User struct {
	ID           string    `json:"id" bson:"id"`
	FullName     string    `json:"full_name" validate:"required,min=2,max=100" bson:"full_name"`
	Password     string    `json:"password" validate:"required,min=8" bson:"password"`
	Email        string    `json:"email" validate:"email,required" bson:"email"`
	AccessToken  string    `json:"access_token" bson:"access_token"`
	RefreshToken string    `json:"refresh_token" bson:"refresh_token"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" bson:"updated_at"`
}

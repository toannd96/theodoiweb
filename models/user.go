package models

import (
	"time"
)

// User ...
type User struct {
	ID           string    `json:"id"`
	FullName     string    `json:"full_name" validate:"required,min=2,max=100"`
	Password     string    `json:"password" validate:"required,min=8"`
	Email        string    `json:"email" validate:"email,required"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

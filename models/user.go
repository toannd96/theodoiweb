package models

import (
	"time"
)

// User ...
type User struct {
	ID            string    `json:"id"`
	FullName      string    `json:"full_name" validate:"required,min=2,max=100"`
	Password      string    `json:"password" validate:"required,min=8"`
	Email         string    `json:"email" validate:"email,required"`
	Token         string    `json:"token"`
	Refresh_token string    `json:"refresh_token"`
	Created_at    time.Time `json:"created_at"`
	Updated_at    time.Time `json:"updated_at"`
}

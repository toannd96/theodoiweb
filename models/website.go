package models

import "time"

// Website ...
type Website struct {
	ID         string    `json:"website_id"`
	UserID     string    `json:"user_id"`
	Name       string    `json:"name"`
	URL        string    `json:"url"`
	Tracked    bool      `json:"tracked"`
	Sessions   []Session `json:"sessions,omitempty"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

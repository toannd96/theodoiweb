package models

// Session ...
type Session struct {
	ID        string `json:"id"`
	WebsiteID string `json:"website_id"`
	UserAgent string `json:"user_agent"`
	Browser   string `json:"browser"`
	Version   string `json:"version"`
	OS        string `json:"os"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Sessions ...
type Sessions []Session

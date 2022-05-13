package models

// Session ...
type Session struct {
	ID        string  `json:"session_id"`
	WebsiteID string  `json:"website_id"`
	Name      string  `json:"name"`
	UserAgent string  `json:"user_agent"`
	Browser   string  `json:"browser"`
	Version   string  `json:"version"`
	OS        string  `json:"os"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	Events    []Event `json:"events"`
}

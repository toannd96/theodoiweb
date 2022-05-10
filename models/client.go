package models

// Client ...
type Client struct {
	ClientID     string    `json:"client_id,omitempty"`
	WebsiteID    string    `json:"website_id,omitempty"`
	UserAgent    string    `json:"user_agent,omitempty"`
	Country      string    `json:"country,omitempty"`
	Language     string    `json:"language,omitempty"`
	ReferrerURL  string    `json:"referrer_url,omitempty"`
	TotalSession int64     `json:"total_session,omitempty"`
	Sessions     []Session `json:"sessions,omitempty"`
	CreatedAt    string    `json:"created_at,omitempty"`
}

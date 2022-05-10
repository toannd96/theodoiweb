package models

// Session ...
type Session struct {
	SessionID      string  `json:"session_id,omitempty"`
	ClientID       string  `json:"client_id,omitempty"`
	SessionName    string  `json:"session_name,omitempty"`
	Events         []Event `json:"events,omitempty"`
	UserAgent      string  `json:"user_agent,omitempty"`
	Browser        string  `json:"browser,omitempty"`
	Version        string  `json:"version,omitempty"`
	BrowserVersion string  `json:"browser_version,omitempty"`
	OS             string  `json:"os,omitempty"`
	Device         string  `json:"device,omitempty"`
	ScreenSize     string  `json:"screen_size,omitempty"`
	Duration       string  `json:"duration,omitempty"`
	PageLoadTime   uint64  `json:"page_load_time,omitempty"`
	Source         string  `json:"source,omitempty"`
	CreatedAt      string  `json:"created_at,omitempty"`
	UpdatedAt      string  `json:"updated_at,omitempty"`
}

package models

// Events ...
type Events struct {
	Events []Event `json:"events,omitempty"`
}

// Event ...
type Event struct {
	SessionID string      `json:"session_id"`
	Type      int64       `json:"type"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

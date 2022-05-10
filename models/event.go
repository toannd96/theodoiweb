package models

// Events ...
type Events struct {
	Events []Event `json:"events,omitempty"`
}
// Event ...
type Event struct {
	SessionID string      `json:"session_id,omitempty"`
	Type      int64       `json:"type,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

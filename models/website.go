package models

// Website ...
type Website struct {
	WebsiteID              string `json:"website_id,omitempty"`
	UserID                 string `json:"user_id,omitempty"`
	GroupID                string `json:"group_id,omitempty"`
	Name                   string `json:"name,omitempty"`
	URL                    string `json:"url,omitempty"`
	TrackingCode           bool   `json:"tracking_code,omitempty"`
	AccessSession          int64  `json:"access_session,omitempty"`
	Visitor                int64  `json:"visitor,omitempty"`
	ComputerCount          int64  `json:"computer_count,omitempty"`
	TabletCount            int64  `json:"tablet_count,omitempty"`
	MobileCount            int64  `json:"mobile_count,omitempty"`
	AverageSessionDuration int64  `json:"average_session_duration,omitempty"`
	TotalAccessTime        int64  `json:"total_access_time,omitempty"`
	PageLoadSpeed          uint64 `json:"page_load_speed,omitempty"`
	BounceRate             string `json:"bounce_rate,omitempty"`
	AccessRate             string `json:"access_rate,omitempty"`
	AccessSource           string `json:"access_source,omitempty"`
	CreatedAt              string `json:"created_at,omitempty"`
	UpdatedAt              string `json:"updated_at,omitempty"`
}

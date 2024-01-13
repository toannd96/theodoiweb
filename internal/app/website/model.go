package website

// website ...
type website struct {
	ID        string `json:"id" bson:"id"`
	UserID    string `json:"user_id" bson:"user_id"`
	Category  string `json:"category" bson:"category"`
	HostName  string `json:"host_name" bson:"host_name"`
	URL       string `json:"url" bson:"url"`
	CreatedAt string `json:"created_at" bson:"created_at"`
	UpdatedAt string `json:"updated_at" bson:"updated_at"`
}

// websites ...
type websites []website

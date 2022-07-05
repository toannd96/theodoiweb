package configs

import (
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Port      string
	AppURL    string
	PathGeoDB string

	MongoDB struct {
		Client            *mongo.Database
		URI               string
		Name              string
		UserCollection    string
		WebsiteCollection string
		SessionCollection string
	}

	Redis struct {
		Client *redis.Client
		Host   string
		Port   string
		URL    string
	}
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		return
	}

	Port = os.Getenv("PORT")
	AppURL = os.Getenv("APP_URL")
	PathGeoDB = os.Getenv("PATH_GEO_DB")

	MongoDB.URI = os.Getenv("URI")
	MongoDB.Name = os.Getenv("NAME")
	MongoDB.UserCollection = os.Getenv("USER_COLLECTION")
	MongoDB.WebsiteCollection = os.Getenv("WEBSITE_COLLECTION")
	MongoDB.SessionCollection = os.Getenv("SESSION_COLLECTION")

	if IsDev() {
		Redis.Host = os.Getenv("REDIS_HOST")
		Redis.Port = os.Getenv("REDIS_PORT")
	} else {
		Redis.URL = os.Getenv("REDIS_URL")
	}
}

// IsDev ...
func IsDev() bool {
	return os.Getenv("MODE") == "dev"
}

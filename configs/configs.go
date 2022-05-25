package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Mode       string
	LogLevel   string
	ServerPort string
	AppURL     string
	PathGeoDB  string

	Sentry struct {
		Dsn              string
		TracesSampleRate float64
	}

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
		// Host   string
		// Port   string
		URL string
	}
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		return
	}

	Mode = os.Getenv("MODE")
	LogLevel = os.Getenv("LOG_LEVEL")
	ServerPort = os.Getenv("PORT")
	AppURL = os.Getenv("APP_URL")
	PathGeoDB = os.Getenv("PATH_GEO_DB")

	MongoDB.URI = os.Getenv("URI")
	MongoDB.Name = os.Getenv("NAME")
	MongoDB.UserCollection = os.Getenv("USER_COLLECTION")
	MongoDB.WebsiteCollection = os.Getenv("WEBSITE_COLLECTION")
	MongoDB.SessionCollection = os.Getenv("SESSION_COLLECTION")

	// Redis.Host = os.Getenv("REDIS_HOST")
	// Redis.Port = os.Getenv("REDIS_PORT")
	Redis.URL = os.Getenv("REDIS_URL")

	Sentry.Dsn = os.Getenv("SENTRY_DSN")
	sampleRate, err := strconv.ParseFloat(os.Getenv("SENTRY_TRACE_SAMPLE_RATE"), 32)
	if err != nil {
		sampleRate = 0
	}
	Sentry.TracesSampleRate = sampleRate
	err = sentry.Init(sentry.ClientOptions{
		Dsn:              Sentry.Dsn,
		Environment:      Mode,
		TracesSampleRate: Sentry.TracesSampleRate,
		AttachStacktrace: true,
	})
	if err != nil {
		log.Fatal("sentry.Init:", err)
	}
}

// IsProduction ...
func IsProduction() bool {
	return os.Getenv("MODE") == "production"
}

// IsStaging ...
func IsStaging() bool {
	return os.Getenv("MODE") == "staging"
}

// IsDev ...
func IsDev() bool {
	return !IsProduction() && !IsStaging()
}

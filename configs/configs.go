package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Mode       string
	LogLevel   string
	ServerPort string
	QueryLimit string

	Sentry struct {
		Dsn              string
		TracesSampleRate float64
	}

	Database struct {
		Client            *mongo.Database
		URI               string
		Name              string
		UserCollection    string
		WebsiteCollection string
		SessionCollection string
		EventCollection   string
	}
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		return
	}

	Mode = os.Getenv("MODE")
	LogLevel = os.Getenv("LOG_LEVEL")
	ServerPort = os.Getenv("SERVER_PORT")
	QueryLimit = os.Getenv("QUERY_LIMIT")

	Database.URI = os.Getenv("URI")
	Database.Name = os.Getenv("NAME")
	Database.UserCollection = os.Getenv("USER_COLLECTION")
	Database.WebsiteCollection = os.Getenv("WEBSITE_COLLECTION")
	Database.SessionCollection = os.Getenv("SESSION_COLLECTION")
	Database.EventCollection = os.Getenv("EVENT_COLLECTION")

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

package configs

import (
	"log"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
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
		Client         influxdb2.Client
		URL            string
		Token          string
		Bucket         string
		Measurement    string
		Organization   string
		RequestTimeout uint
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

	Database.URL = os.Getenv("INFLUX_URL")
	Database.Token = os.Getenv("INFLUX_TOKEN")
	Database.Bucket = os.Getenv("INFLUX_BUCKET")
	Database.Measurement = os.Getenv("INFLUX_MEASUREMENT")
	Database.Organization = os.Getenv("INFLUX_ORGANIZATION")
	Database.RequestTimeout = 60

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

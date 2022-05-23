package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"

	"analytics-api/configs"
	"analytics-api/db"
	"analytics-api/internal/app/session"
	"analytics-api/internal/pkg/log"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func init() {
	// Config logger
	if !configs.IsDev() {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				filename := path.Base(f.File)
				return fmt.Sprintf("%s()\n", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
			},
		})
	}
	logLevel, err := logrus.ParseLevel(configs.LogLevel)
	if err != nil {
		logLevel = logrus.TraceLevel
	}
	logrus.SetLevel(logLevel)
	logrus.SetReportCaller(true)
	logrus.SetOutput(os.Stdout)

	if !configs.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}
}

func main() {
	defer sentry.Flush(2 * time.Second)

	var err error

	db.NewMongo()
	db.NewRedis()

	r := gin.Default()

	// Custom recovery to enable Sentry capturing
	r.Use(gin.CustomRecovery(func(c *gin.Context, err interface{}) {
		log.LogPanic(c, err)
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	initializeRoutes(r)

	logrus.Info("starting HTTP server...")
	err = r.Run(":" + configs.ServerPort)
	if err != nil {
		logrus.Fatalln(err)
	}
}

func initializeRoutes(r *gin.Engine) {
	// Register health check handler
	r.GET("healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	delivery := session.NewHTTPDelivery()
	delivery.InitRoutes(r)
}

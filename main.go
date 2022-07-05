package main

import (
	"analytics-api/configs"
	"analytics-api/db"
	"analytics-api/internal/app/session"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error

	db.NewMongo()
	db.CreateSessionCollection()

	db.NewRedis()

	r := gin.Default()
	initializeRoutes(r)

	logrus.Info("starting HTTP server...")
	err = r.Run(":" + configs.Port)
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

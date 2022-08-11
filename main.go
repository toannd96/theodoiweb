package main

import (
	"analytics-api/configs"
	"analytics-api/db"
	"analytics-api/internal/app/session"
	"analytics-api/internal/app/user"
	"analytics-api/internal/app/website"
	"analytics-api/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	var err error

	db.NewMongo()

	userErr := db.CreateUserCollection()
	if userErr != nil {
		logrus.Fatalln(userErr)
	}

	webErr := db.CreateWebsiteCollection()
	if webErr != nil {
		logrus.Fatalln(webErr)
	}

	sessionErr := db.CreateSessionCollection()
	if sessionErr != nil {
		logrus.Fatalln(sessionErr)
	}

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
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "home.html", gin.H{})
	})

	r.LoadHTMLGlob("web/templates/**")
	r.StaticFile("/record.js", "./web/static/js/record.js")

	r.Static("/js", "./web/static/js")
	r.Static("/assets", "./web/static/assets")
	r.Static("/css", "./web/static/css")
	r.Use(middleware.CORSMiddleware())

	g := r.Group("/")
	sessionDelivery := session.NewHTTPDelivery()
	userDelivery := user.NewHTTPDelivery()
	websiteDelivery := website.NewHTTPDelivery()

	sessionDelivery.InitRoutes(g)
	userDelivery.InitRoutes(g)
	websiteDelivery.InitRoutes(g)
}

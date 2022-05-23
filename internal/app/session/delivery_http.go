package session

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"analytics-api/internal/pkg/duration"
	"analytics-api/internal/pkg/log"
	"analytics-api/internal/pkg/middleware"
	"analytics-api/models"

	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
	"github.com/sirupsen/logrus"
)

type httpDelivery struct {
	sessionUseCase UseCase
}

// Request website tracking send to server
type Request struct {
	SessionID string         `json:"session_id"`
	Events    []models.Event `json:"events"`
}

// InitRoutes ...
func (instance *httpDelivery) InitRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("web/template/*")
	r.StaticFile("/record.js", "./web/static/record.js")
	r.Use(middleware.CORSMiddleware())

	r.GET("/", instance.Home)

	// Register routes session
	sessionRoutes := r.Group("session")
	sessionRoutes.GET("/records", instance.RenderListSession)
	sessionRoutes.POST("/save", instance.SaveSession)
	sessionRoutes.GET("/:session_id", instance.RenderSessionPlay)
	sessionRoutes.GET("/event/:session_id", instance.GetEventBySessionID)
}

// GetEventBySessionID streaming all event of session by session id
func (instance *httpDelivery) GetEventBySessionID(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(200)

	var session models.Session

	sessionID := c.Param("session_id")
	limit := 10
	skip := 0

	events, err := instance.sessionUseCase.GetEvents(sessionID, &session)
	if err != nil {
		log.LogError(c, err)
		return
	}
	logrus.Debug("count event ", len(events))

	msgChan := make(chan []models.Event)
	breakLineChan := make(chan string)
	breakChan := make(chan bool)
	defer func() {
		close(msgChan)
		msgChan = nil
		logrus.Debug("client connection is closed")
	}()

	go func() {
		if msgChan != nil {
			for skip <= int(len(events)) {
				var session models.Session
				err := instance.sessionUseCase.GetEventByLimitSkip(sessionID, &session, limit, skip)
				if err != nil {
					log.LogError(c, err)
					return
				}
				skip = skip + limit

				msgChan <- session.Events
				breakLineChan <- "--break--"

				if len(session.Events) < limit {
					breakChan <- true
					return
				}
			}
		}
	}()

	for {
		select {
		case message := <-msgChan:
			enc := json.NewEncoder(c.Writer)
			if err := enc.Encode(message); err != nil {
				logrus.Error("encode msg: ", err)
				return
			}
			c.Writer.Flush()
		case message := <-breakLineChan:
			enc := json.NewEncoder(c.Writer)
			enc.Encode(message)
			c.Writer.Flush()
		case <-breakChan:
			return
		case <-c.Request.Context().Done():
			return
		}
	}
}

// Home
func (instance *httpDelivery) Home(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"URL": "http://localhost:3000",
	})
}

// RenderSessionPlay replay session by session id
func (instance *httpDelivery) RenderSessionPlay(c *gin.Context) {
	sessionID := c.Param("session_id")
	var session models.Session
	err := instance.sessionUseCase.GetSession(sessionID, &session)
	if err != nil {
		log.LogError(c, err)
		return
	}
	c.HTML(http.StatusOK, "session_by_id.html", gin.H{
		"SessionID": sessionID,
		"Session":   session,
	})
}

// RenderListSession show list session record
func (instance *httpDelivery) RenderListSession(c *gin.Context) {
	var sessions models.Sessions
	sessions, err := instance.sessionUseCase.GetAllSession(sessions)
	if err != nil {
		log.LogError(c, err)
		return
	}
	c.HTML(http.StatusOK, "session_list.html", gin.H{
		"Sessions": sessions,
		"URL":      "http://localhost:3000",
	})
}

// SaveSession save session
func (instance *httpDelivery) SaveSession(c *gin.Context) {
	var request Request
	var session models.Session

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	countSessionID, err := instance.sessionUseCase.FindSession(request.SessionID)
	if err != nil {
		log.LogError(c, err)
		return
	}

	// clientIP := net.ParseIP(realip.FromRequest(c.Request))

	// geoDB, err := geodb.Open(configs.PathGeoDB)
	// if err != nil {
	// 	log.LogError(c, err)
	// 	return
	// }
	// defer geoDB.Close()

	// geoData, err := geoDB.City(clientIP)
	// if err != nil {
	// 	log.LogError(c, err)
	// 	return
	// }

	// if session id not exists
	if countSessionID == 0 {
		ua := ua.Parse(c.Request.UserAgent())
		// clientIP := net.ParseIP(realip.FromRequest(c.Request))

		// geoDB, err := geodb.Open(configs.PathGeoDB)
		// if err != nil {
		// 	log.LogError(c, err)
		// 	return
		// }
		// defer geoDB.Close()

		// geoData, err := geoDB.City(clientIP)
		// if err != nil {
		// 	log.LogError(c, err)
		// 	return
		// }

		referrerURL, err := url.Parse(c.Request.Referer())
		if err != nil {
			log.LogError(c, err)
			return
		}

		session.ID = request.SessionID
		session.OS = ua.OS
		session.Browser = ua.Name

		if ua.Mobile {
			session.Device = "Mobile"
		}
		if ua.Tablet {
			session.Device = "Tablet"
		}
		if ua.Desktop {
			session.Device = "Desktop"
		}

		session.Referral = referrerURL.Hostname()
		// session.City = geoData.City.Names["en"]
		// session.Country = geoData.Country.Names["en"]

		session.Events = append(session.Events, request.Events...)

		if len(session.Events) != 0 {
			time1 := session.Events[0].Timestamp / 1000
			time2 := session.Events[len(session.Events)-1].Timestamp / 1000
			duration := duration.Duration(time1, time2)
			session.Duration = duration
			session.CreatedAt = time.Unix(time1, 0).Format("02/01/2006, 15:04:05")

			// save time1 of session id to redis
			err = instance.sessionUseCase.InsertSessionTimestamp(request.SessionID, time1)
			if err != nil {
				log.LogError(c, err)
				return
			}
		}

		// save session
		err = instance.sessionUseCase.InsertSession(session)
		if err != nil {
			log.LogError(c, err)
			return
		}
	}

	// get time1 by session id from redis
	time1, err := instance.sessionUseCase.GetSessionTimestamp(request.SessionID)
	if err != nil {
		log.LogError(c, err)
		return
	}

	// update session
	session.Events = append(session.Events, request.Events...)
	if len(session.Events) != 0 {
		time2 := session.Events[len(session.Events)-1].Timestamp / 1000
		duration := duration.Duration(time1, time2)
		session.Duration = duration
	}

	errEvent := instance.sessionUseCase.UpdateSession(request.SessionID, session)
	if errEvent != nil {
		log.LogError(c, errEvent)
		return
	}
	c.JSON(http.StatusOK, session)
}

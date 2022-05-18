package session

import (
	"encoding/json"
	"net/http"
	"time"

	"analytics-api/internal/app/event"
	"analytics-api/internal/pkg/log"
	"analytics-api/internal/pkg/middleware"
	"analytics-api/models"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
)

type httpDelivery struct {
	sessionUseCase UseCase
	eventUseCase   event.UseCase
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

	r.GET("/", instance.RenderListSession)

	// Register routes session
	sessionRoutes := r.Group("session")
	sessionRoutes.POST("/save", instance.SaveSession)
	sessionRoutes.GET("/:session_id", instance.RenderSessionPlay)
	sessionRoutes.GET("/event/:session_id", instance.GetEventBySessionID)
}

// GetEventBySessionID streaming all event of session by session id
func (instance *httpDelivery) GetEventBySessionID(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(200)

	var events models.Events

	sessionID := c.Param("session_id")
	limit := 10
	skip := 0

	countEvent, err := instance.eventUseCase.GetCountEventBySessionID(sessionID, &events)
	if err != nil {
		log.LogError(c, err)
		return
	}
	logrus.Debug("count event ", countEvent)

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
			for skip <= int(countEvent) {
				var events models.Events
				err := instance.eventUseCase.GetEventBySessionID(sessionID, &events, limit, skip)
				if err != nil {
					log.LogError(c, err)
					return
				}
				skip = skip + limit

				msgChan <- events.Events
				breakLineChan <- "--break--"

				if len(events.Events) < limit {
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

// RenderSessionPlay replay session by session id
func (instance *httpDelivery) RenderSessionPlay(c *gin.Context) {
	sessionID := c.Param("session_id")
	var session models.Session
	err := instance.sessionUseCase.GetSessionByID(sessionID, &session)
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

// SaveSession save session to influxdb
func (instance *httpDelivery) SaveSession(c *gin.Context) {
	var request Request
	var session models.Session
	var events models.Events

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	countSessionID, err := instance.sessionUseCase.FindSessionID(request.SessionID)
	if err != nil {
		log.LogError(c, err)
		return
	}
	// if session id not exists
	if countSessionID == 0 {
		ua := user_agent.New(c.Request.UserAgent())
		browserName, browserVersion := ua.Browser()

		session.OS = ua.OS()
		session.Browser = browserName
		session.Version = browserVersion
		session.ID = request.SessionID
		session.UserAgent = c.Request.UserAgent()
		session.CreatedAt = time.Now().Format("02/01/2006, 15:04:05")

		// save session
		err = instance.sessionUseCase.InsertSession(session)
		if err != nil {
			log.LogError(c, err)
			return
		}
	}

	countID, err := instance.eventUseCase.FindSessionID(request.SessionID)
	if err != nil {
		log.LogError(c, err)
		return
	}
	// if session id in event not exists
	if countID == 0 {
		events.SessionID = request.SessionID
		// save session id in event
		err = instance.eventUseCase.InsertEvent(events.SessionID)
		if err != nil {
			log.LogError(c, err)
			return
		}

		events.Events = append(events.Events, request.Events...)
		// update events
		errEvent := instance.eventUseCase.UpdateEvent(events)
		if errEvent != nil {
			log.LogError(c, errEvent)
			return
		}
	} else {
		events.SessionID = request.SessionID
		events.Events = append(events.Events, request.Events...)

		errEvent := instance.eventUseCase.UpdateEvent(events)
		if errEvent != nil {
			log.LogError(c, errEvent)
			return
		}
	}
	c.JSON(http.StatusOK, "save session success")
}

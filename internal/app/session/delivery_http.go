package session

import (
	"encoding/json"
	"net/http"
	"time"

	"analytics-api/configs"
	"analytics-api/internal/pkg/common"
	"analytics-api/internal/pkg/middleware"
	"analytics-api/models"

	"analytics-api/internal/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/sirupsen/logrus"
)

type httpDelivery struct {
	sessionUseCase UseCase
}

// Request ...
type Request struct {
	SessionID   string         `json:"session_id"`
	SessionName string         `json:"session_name"`
	Events      []models.Event `json:"events"`
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
	sessionRoutes.GET("/event/:session_id", instance.GetAllEventLimitByID)
	sessionRoutes.GET("/info/:session_id", instance.GetAllSessionByID)

}

// GetAllSessionByID return info session by session id
func (instance *httpDelivery) GetAllSessionByID(c *gin.Context) {
	sessionID := c.Param("session_id")
	var session models.Session

	err := instance.sessionUseCase.GetSessionByID(sessionID, &session)
	if err != nil {
		log.LogError(c, err)
		return
	}

	c.JSON(http.StatusOK, session)
}

// GetAllEventLimitByID streaming all event of session by session id
func (instance *httpDelivery) GetAllEventLimitByID(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(200)

	sessionID := c.Param("session_id")
	limit, err := common.StringToInt64(configs.QueryLimit)
	if err != nil {
		logrus.Error("convert string query limit to int64: ", err)
		return
	}
	offset := 0

	totalRecord, err := instance.sessionUseCase.GetTotalColumn(sessionID)
	if err != nil {
		log.LogError(c, err)
		return
	}
	logrus.Debug("total number rows ", totalRecord)

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
			for offset <= int(totalRecord) {
				var events models.Events
				err := instance.sessionUseCase.GetEventLimitBySessionID(sessionID, int(limit), offset, &events)
				if err != nil {
					log.LogError(c, err)
					return
				}
				offset = offset + int(limit)

				msgChan <- events.Events
				breakLineChan <- "--break--"

				if len(events.Events) < int(limit) {
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
	var session models.Session

	listSessionID, err := instance.sessionUseCase.GetAllSessionID()
	if err != nil {
		log.LogError(c, err)
		return
	}

	sessions, err := instance.sessionUseCase.GetAllSession(listSessionID, session)
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

	ua := user_agent.New(c.Request.UserAgent())
	browserName, browserVersion := ua.Browser()
	session.OS = ua.OS()
	session.Browser = browserName
	session.Version = browserVersion
	session.SessionID = request.SessionID
	session.SessionID = request.SessionID
	session.SessionName = request.SessionName
	session.SessionID = request.SessionID
	session.UserAgent = c.Request.UserAgent()
	session.UpdatedAt = time.Now().Format("02/01/2006, 15:04:05")

	events.Events = append(events.Events, request.Events...)

	err = instance.sessionUseCase.InsertSession(session, events)
	if err != nil {
		log.LogError(c, err)
		return
	}

	c.JSON(http.StatusOK, "save session success")
}

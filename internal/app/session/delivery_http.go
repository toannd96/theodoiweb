package session

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"analytics-api/configs"
	dur "analytics-api/internal/pkg/duration"
	"analytics-api/internal/pkg/geodb"
	str "analytics-api/internal/pkg/string"
	"analytics-api/models"

	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
	"github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
)

type httpDelivery struct {
	sessionUseCase UseCase
}

// RequestSession website tracking send to server
type RequestSession struct {
	SessionID string         `json:"session_id"`
	Events    []models.Event `json:"events"`
}

// InitRoutes ...
func (instance *httpDelivery) InitRoutes(r *gin.RouterGroup) {

	// Register routes session
	sessionRoutes := r.Group("session")
	{
		sessionRoutes.GET("/tracking", instance.Tracking)
		sessionRoutes.GET("/records", instance.ListSessionRecord)
		sessionRoutes.POST("/receive", instance.ReceiveSession)
		sessionRoutes.GET("/:session_id", instance.SessionReplay)
		sessionRoutes.GET("/event/:session_id", instance.GetEventBySessionID)
	}
}

// GetEventBySessionID streaming all event of session by session id
func (instance *httpDelivery) GetEventBySessionID(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(200)

	sessionID := c.Param("session_id")
	limit := 10
	skip := 0

	countSession, err := instance.sessionUseCase.GetCountSession(sessionID)
	if err != nil {
		logrus.Error(c, err)
		return
	}
	logrus.Info("count event ", countSession)

	msgChan := make(chan []*models.Event)
	breakLineChan := make(chan string)
	breakChan := make(chan bool)
	defer func() {
		close(msgChan)
		msgChan = nil
		logrus.Info("client connection is closed")
	}()

	go func() {
		if msgChan != nil {
			for skip <= int(countSession) {
				events, err := instance.sessionUseCase.GetEventByLimitSkip(sessionID, limit, skip)
				if err != nil {
					logrus.Error(c, err)
					return
				}
				skip = skip + limit
				logrus.Info("len events ", len(events))
				msgChan <- events
				breakLineChan <- "--break--"

				if len(events) < limit {
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

// Tracking guide tracking code to website
func (instance *httpDelivery) Tracking(c *gin.Context) {
	c.HTML(http.StatusOK, "tracking.html", gin.H{
		"URL": configs.AppURL,
	})
}

// SessionReplay replay session by session id
func (instance *httpDelivery) SessionReplay(c *gin.Context) {
	sessionID := c.Param("session_id")
	var session models.Session
	err := instance.sessionUseCase.GetSession(sessionID, &session)
	if err != nil {
		logrus.Error(c, err)
		return
	}
	c.HTML(http.StatusOK, "video.html", gin.H{
		"SessionID": sessionID,
		"Session":   session.MetaData,
	})
}

// ListSessionRecord show list session record
func (instance *httpDelivery) ListSessionRecord(c *gin.Context) {
	var session models.Session
	listSessionID, err := instance.sessionUseCase.GetAllSessionID(session)
	if err != nil {
		logrus.Error(c, err)
		return
	}
	if len(listSessionID) != 0 {
		listSession, err := instance.sessionUseCase.GetAllSession(listSessionID, session)
		if err != nil {
			logrus.Error(c, err)
			return
		}
		c.HTML(http.StatusOK, "tables.html", gin.H{
			"Sessions": listSession,
		})
	}
}

// ReceiveSession receive session from request client
func (instance *httpDelivery) ReceiveSession(c *gin.Context) {
	var request RequestSession
	var session models.Session

	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ua := ua.Parse(c.Request.UserAgent())
	clientIP := net.ParseIP(realip.FromRequest(c.Request))

	geoDB, err := geodb.Open(configs.PathGeoDB)
	if err != nil {
		logrus.Error(c, err)
		return
	}
	defer geoDB.Close()

	geoData, err := geoDB.City(clientIP)
	if err != nil {
		logrus.Error(c, err)
		return
	}

	session.MetaData.ID = request.SessionID
	session.MetaData.OS = ua.OS
	session.MetaData.Browser = ua.Name
	session.MetaData.Version = ua.Version

	if ua.Mobile {
		session.MetaData.Device = "Mobile"
	}
	if ua.Tablet {
		session.MetaData.Device = "Tablet"
	}
	if ua.Desktop {
		session.MetaData.Device = "Desktop"
	}

	session.MetaData.Country = geoData.Country.Names["en"]
	session.MetaData.City = str.RemoveSubstring(geoData.City.Names["en"], "City")

	events := request.Events

	countSession, err := instance.sessionUseCase.GetCountSession(request.SessionID)
	if err != nil {
		logrus.Error(c, err)
		return
	}
	if countSession == 0 {
		if len(events) != 0 {
			time1 := events[0].Timestamp / 1000
			time2 := events[len(events)-1].Timestamp / 1000
			duration := dur.Duration(time1, time2)

			session.Duration = duration

			timeReport, err := dur.ParseTime(time.Unix(time1, 0).Format("2006-01-02, 15:04:05"))
			if err != nil {
				logrus.Error(c, err)
				return
			}
			session.TimeReport = timeReport
			session.MetaData.CreatedAt = time.Unix(time1, 0).Format("2006-01-02, 15:04:05")

			// save time1 of session id to redis
			err = instance.sessionUseCase.InsertSessionTimestamp(request.SessionID, time1)
			if err != nil {
				logrus.Error(c, err)
				return
			}
		} else {
			session.Duration = "00:00:00"

			timeReport, err := dur.ParseTime(time.Now().Format("2006-01-02, 15:04:05"))
			if err != nil {
				logrus.Error(c, err)
				return
			}
			session.TimeReport = timeReport
			session.MetaData.CreatedAt = time.Now().Format("2006-01-02, 15:04:05")
		}
	} else {
		if len(events) != 0 {
			// get time1 by session id from redis
			time1, err := instance.sessionUseCase.GetSessionTimestamp(request.SessionID)
			if err != nil {
				logrus.Error(c, err)
				return
			}
			time2 := events[len(events)-1].Timestamp / 1000
			duration := dur.Duration(time1, time2)
			session.Duration = duration

			timeReport, err := dur.ParseTime(time.Now().Format("2006-01-02, 15:04:05"))
			if err != nil {
				logrus.Error(c, err)
				return
			}
			session.TimeReport = timeReport
			session.MetaData.CreatedAt = time.Unix(time1, 0).Format("2006-01-02, 15:04:05")
		}
	}

	// save session
	err = instance.sessionUseCase.InsertSession(session, events)
	if err != nil {
		logrus.Error(c, err)
		return
	}
	c.JSON(http.StatusOK, session)
}

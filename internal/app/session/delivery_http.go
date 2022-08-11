package session

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"analytics-api/configs"
	"analytics-api/internal/app/auth"
	"analytics-api/internal/app/website"
	dur "analytics-api/internal/pkg/duration"
	"analytics-api/internal/pkg/geodb"
	"analytics-api/internal/pkg/middleware"
	"analytics-api/internal/pkg/security"
	str "analytics-api/internal/pkg/string"
	"analytics-api/models"

	"github.com/gin-gonic/gin"
	ua "github.com/mileusna/useragent"
	"github.com/sirupsen/logrus"
	"github.com/tomasen/realip"
)

type httpDelivery struct {
	sessionUseCase UseCase
	websiteUseCase website.UseCase
	authUsecase    auth.UseCase
}

// RequestSession website tracking send to server
type RequestSession struct {
	UserID    string         `json:"user_id"`
	WebsiteID string         `json:"website_id"`
	SessionID string         `json:"session_id"`
	Events    []models.Event `json:"events"`
}

// InitRoutes ...
func (instance *httpDelivery) InitRoutes(r *gin.RouterGroup) {

	// Register routes session
	sessionRoutes := r.Group("session")
	{
		sessionRoutes.GET("/heatmaps", middleware.JWTMiddleware(), instance.ShowHeatmaps)
		sessionRoutes.GET("/record", middleware.JWTMiddleware(), instance.ListWebsiteOfSessionRecord)
		sessionRoutes.GET("/record/:website_id", middleware.JWTMiddleware(), instance.ListSessionRecord)
		sessionRoutes.POST("/receive", instance.ReceiveSession)
		sessionRoutes.GET("/:session_id", middleware.JWTMiddleware(), instance.SessionReplay)
		sessionRoutes.GET("/event/:session_id", middleware.JWTMiddleware(), instance.GetEventBySessionID)
	}
}

func (instance *httpDelivery) ShowHeatmaps(c *gin.Context) {
	c.HTML(http.StatusOK, "heatmaps.html", gin.H{})
}

// GetEventBySessionID streaming all event of session by session id
func (instance *httpDelivery) GetEventBySessionID(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(200)

	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "extract token metadata failed"})
		return
	}

	userID, err := instance.authUsecase.GetAuth(tokenAuth.AccessUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "get token auth failed"})
		return
	}

	sessionID := c.Param("session_id")
	limit := 10
	skip := 0

	countSession, err := instance.sessionUseCase.GetCountSession(userID, sessionID)
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
				events, err := instance.sessionUseCase.GetEventByLimitSkip(userID, sessionID, limit, skip)
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
			if err := enc.Encode(message); err != nil {
				logrus.Error("encode msg: ", err)
				return
			}
			c.Writer.Flush()
		case <-breakChan:
			return
		case <-c.Request.Context().Done():
			return
		}
	}
}

// SessionReplay replay session by session id
func (instance *httpDelivery) SessionReplay(c *gin.Context) {
	sessionID := c.Param("session_id")
	var session models.Session

	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "extract token metadata failed"})
		return
	}

	userID, err := instance.authUsecase.GetAuth(tokenAuth.AccessUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "get token auth failed"})
		return
	}

	getSessionErr := instance.sessionUseCase.GetSession(userID, sessionID, &session)
	if getSessionErr != nil {
		logrus.Error(c, err)
		return
	}
	c.HTML(http.StatusOK, "video.html", gin.H{
		"SessionID": sessionID,
		"Session":   session.MetaData,
	})
}

// ListSessionRecord show list session record
func (instance *httpDelivery) ListWebsiteOfSessionRecord(c *gin.Context) {
	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "extract token metadata failed"})
		return
	}

	userID, err := instance.authUsecase.GetAuth(tokenAuth.AccessUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "get token auth failed"})
		return
	}

	websites, err := instance.websiteUseCase.GetAllWebsite(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}
	if len(*websites) == 0 {
		c.HTML(http.StatusOK, "not_website.html", gin.H{})
	} else {
		c.HTML(http.StatusOK, "tables.html", gin.H{
			"Websites": websites,
		})
	}
}

func (instance *httpDelivery) ListSessionRecord(c *gin.Context) {
	var session models.Session
	var listSessionID []string

	tokenAuth, err := security.ExtractAccessTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "extract token metadata failed"})
		return
	}

	userID, err := instance.authUsecase.GetAuth(tokenAuth.AccessUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "get token auth failed"})
		return
	}

	websites, err := instance.websiteUseCase.GetAllWebsite(userID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	websiteID := c.Param("website_id")
	query := c.Query("time")

	switch query {
	case "today":
		listSessionID, err = instance.sessionUseCase.GetSessionIDToday(userID, websiteID, session)
		if err != nil {
			logrus.Error(c, err)
			return
		}
	case "all":
		listSessionID, err = instance.sessionUseCase.GetAllSessionID(userID, websiteID, session)
		if err != nil {
			logrus.Error(c, err)
			return
		}
	default:
		listSessionID, err = instance.sessionUseCase.GetAllSessionID(userID, websiteID, session)
		if err != nil {
			logrus.Error(c, err)
			return
		}
	}

	if len(listSessionID) != 0 {
		listSession, err := instance.sessionUseCase.GetAllSession(userID, websiteID, listSessionID, session)
		if err != nil {
			logrus.Error(c, err)
			return
		}

		c.HTML(http.StatusOK, "tables.html", gin.H{
			"WebsiteID": websiteID,
			"Websites":  websites,
			"Sessions":  listSession,
		})
	} else {
		switch query {
		case "today":
			c.HTML(http.StatusOK, "not_record_today.html", gin.H{
				"Websites": websites,
			})
			return
		default:
			c.HTML(http.StatusOK, "not_record.html", gin.H{
				"Websites": websites,
			})
			return
		}
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

	countSites, err := instance.websiteUseCase.FindWebsiteByID(request.UserID, request.WebsiteID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	if countSites > 0 {
		logrus.Info("receive session from website id ", request.WebsiteID)
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

		session.MetaData.UserID = request.UserID
		session.MetaData.ID = request.SessionID
		session.MetaData.WebsiteID = request.WebsiteID
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

		countSession, err := instance.sessionUseCase.GetCountSession(request.UserID, request.SessionID)
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
	} else {
		logrus.Info("this site id not exists ", request.WebsiteID)
		c.JSON(http.StatusConflict, gin.H{"msg": "this website not exists"})
		return
	}
}

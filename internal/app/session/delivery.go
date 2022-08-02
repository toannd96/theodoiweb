package session

import (
	"analytics-api/internal/app/auth"
	"analytics-api/internal/app/website"

	"github.com/gin-gonic/gin"
)

// HTTPDelivery ...
type HTTPDelivery interface {
	// This function is required for every HTTPDelivery
	InitRoutes(r *gin.RouterGroup)

	// Other functions to handle HTTP requests
	GetEventBySessionID(c *gin.Context)
	SessionReplay(c *gin.Context)

	ListWebsiteOfSessionRecord(c *gin.Context)
	ListSessionRecord(c *gin.Context)
	ReceiveSession(c *gin.Context)
}

// NewHTTPDelivery ...
func NewHTTPDelivery() HTTPDelivery {
	return &httpDelivery{
		sessionUseCase: NewUseCase(),
		websiteUseCase: website.NewUseCase(),
		authUsecase:    auth.NewUseCase(),
	}
}

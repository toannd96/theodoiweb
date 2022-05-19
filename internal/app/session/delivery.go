package session

import (
	"github.com/gin-gonic/gin"
)

// HTTPDelivery ...
type HTTPDelivery interface {
	// This function is required for every HTTPDelivery
	InitRoutes(r *gin.Engine)

	// Other functions to handle HTTP requests
	GetEventBySessionID(c *gin.Context)
	RenderSessionPlay(c *gin.Context)
	RenderListSession(c *gin.Context)
	SaveSession(c *gin.Context)
}

// NewHTTPDelivery ...
func NewHTTPDelivery() HTTPDelivery {
	return &httpDelivery{
		sessionUseCase: NewUseCase(),
	}
}

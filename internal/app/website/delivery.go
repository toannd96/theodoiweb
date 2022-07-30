package website

import (
	"analytics-api/internal/app/auth"

	"github.com/gin-gonic/gin"
)

// HTTPDelivery ...
type HTTPDelivery interface {
	// This function is required for every HTTPDelivery
	InitRoutes(r *gin.RouterGroup)

	// Other functions to handle HTTP requests
	GetWebsite(c *gin.Context)
	GetAllWebsite(c *gin.Context)
	Tracking(c *gin.Context)
	AddWebsite(c *gin.Context)
	DeleteWebsite(c *gin.Context)
}

// NewHTTPDelivery ...
func NewHTTPDelivery() HTTPDelivery {
	return &httpDelivery{
		websiteUseCase: NewUseCase(),
		authUsecase:    auth.NewUseCase(),
	}
}

package user

import (
	"github.com/gin-gonic/gin"

	"analytics-api/internal/app/auth"
)

// HTTPDelivery ...
type HTTPDelivery interface {
	// This function is required for every HTTPDelivery
	InitRoutes(r *gin.RouterGroup)

	// Other functions to handle HTTP requests
	SignUp(c *gin.Context)
	SignIn(c *gin.Context)
	SignOut(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

// NewHTTPDelivery ...
func NewHTTPDelivery() HTTPDelivery {
	return &httpDelivery{
		userUseCase: NewUseCase(),
		authUsecase: auth.NewUseCase(),
	}
}

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
	ShowSignupPage(c *gin.Context)
	ShowLoginPage(c *gin.Context)
	SignUp(c *gin.Context)
	Login(c *gin.Context)
	Logout(c *gin.Context)
	GetUser(c *gin.Context)
	ShowDetailsUserPage(c *gin.Context)
	UpdateUser(c *gin.Context)
}

// NewHTTPDelivery ...
func NewHTTPDelivery() HTTPDelivery {
	return &httpDelivery{
		userUseCase: NewUseCase(),
		authUsecase: auth.NewUseCase(),
	}
}

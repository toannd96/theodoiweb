package middleware

import (
	"net/http"

	"analytics-api/internal/pkg/security"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := security.AccessTokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"msg": "Access is not allowed"})
			c.Abort()
			return
		}
		c.Next()
	}
}

package middleware

import (
	"net/http"

	"analytics-api/internal/pkg/security"

	"github.com/gin-gonic/gin"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessTokenValidErr := security.AccessTokenValid(c.Request)
		if accessTokenValidErr != nil {
			// c.JSON(http.StatusUnauthorized, gin.H{"msg": "Access is not allowed"})
			c.HTML(http.StatusUnauthorized, "401.html", gin.H{})
			c.Abort()
			return
		}
		c.Next()
	}
}

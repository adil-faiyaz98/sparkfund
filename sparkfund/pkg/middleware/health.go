package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheck middleware adds health check endpoint
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

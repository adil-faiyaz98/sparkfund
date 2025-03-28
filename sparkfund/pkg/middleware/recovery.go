package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/pkg/logger"
)

// Recovery middleware recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered",
					"error", err,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"client_ip", c.ClientIP(),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal Server Error",
				})
			}
		}()
		c.Next()
	}
} 
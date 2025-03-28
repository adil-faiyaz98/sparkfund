package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/pkg/logger"
	"github.com/sparkfund/pkg/metrics"
)

// Logger middleware logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		logger.Info("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", statusCode,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)

		metrics.TrackHTTPRequest(c.Request.Method, path, statusCode, latency)
	}
}

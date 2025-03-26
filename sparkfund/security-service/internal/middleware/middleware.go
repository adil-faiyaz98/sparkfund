package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"sparkfund/security-service/internal/metrics"
)

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// LoggerMiddleware logs request information
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		requestID, _ := c.Get("RequestID")

		// Process request
		c.Next()

		// Log request after processing
		duration := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		logger, exists := c.Get("logger")
		if !exists {
			// If no logger is set, use default zap logger
			return
		}

		log := logger.(*zap.Logger)
		log.Info("HTTP request",
			zap.String("request_id", requestID.(string)),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.String("client_ip", clientIP),
			zap.Duration("duration", duration),
		)
	}
}

// CORSMiddleware handles CORS requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// MetricsMiddleware collects metrics for requests
func MetricsMiddleware(metrics *metrics.Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		method := c.Request.Method

		// Track in-flight requests
		metrics.RequestsInFlight.Inc()
		defer metrics.RequestsInFlight.Dec()

		// Process request
		c.Next()

		// Record metrics after processing
		status := c.Writer.Status()
		duration := time.Since(start).Seconds()
		responseSize := float64(c.Writer.Size())

		metrics.RequestCounter.WithLabelValues(method, path, http.StatusText(status)).Inc()
		metrics.RequestDuration.WithLabelValues(method, path).Observe(duration)
		metrics.ResponseSize.WithLabelValues(method, path).Observe(responseSize)
	}
}

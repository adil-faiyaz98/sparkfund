package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sparkfund/kyc-service/internal/logger"
	"github.com/sparkfund/kyc-service/internal/metrics"
	"go.uber.org/zap"
)

type Middleware struct {
	logger  *logger.Logger
	metrics *metrics.MetricsService
}

func New(logger *logger.Logger, metrics *metrics.MetricsService) *Middleware {
	return &Middleware{
		logger:  logger,
		metrics: metrics,
	}
}

// RequestID adds a unique request ID to each request
func (m *Middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger logs request and response details
func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		requestID := c.GetString("request_id")

		// Process request
		c.Next()

		// Log request details
		duration := time.Since(start)
		status := c.Writer.Status()

		logFields := []zap.Field{
			zap.String("path", path),
			zap.String("method", c.Request.Method),
			zap.Int("status", status),
			zap.Duration("duration", duration),
			zap.String("client_ip", c.ClientIP()),
		}

		if query != "" {
			logFields = append(logFields, zap.String("query", query))
		}

		if len(c.Errors) > 0 {
			logFields = append(logFields, zap.String("errors", c.Errors.String()))
			m.logger.WithRequestID(requestID).Error("Request failed", nil, logFields...)
		} else {
			m.logger.WithRequestID(requestID).Info("Request completed", logFields...)
		}
	}
}

// Metrics collects request metrics
func (m *Middleware) Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Record metrics
		duration := time.Since(start).Seconds()
		status := fmt.Sprintf("%d", c.Writer.Status())

		m.metrics.RecordRequest(
			c.Request.Method,
			c.FullPath(),
			status,
			duration,
		)
	}
}

// Recovery handles panics and logs them appropriately
func (m *Middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("request_id")
				m.logger.WithRequestID(requestID).Error(
					"Request panic recovered",
					fmt.Errorf("%v", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)

				m.metrics.RecordError("panic", "500")
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}

// Cors handles CORS headers
func (m *Middleware) Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds security-related headers
func (m *Middleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
}

// Apply applies all middleware in the correct order
func (m *Middleware) Apply(r *gin.Engine) {
	r.Use(m.Recovery())
	r.Use(m.RequestID())
	r.Use(m.Logger())
	r.Use(m.Metrics())
	r.Use(m.Cors())
	r.Use(m.SecurityHeaders())
}

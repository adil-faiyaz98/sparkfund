package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/credit-scoring-service/internal/config"
	"go.uber.org/zap"
)

// SecurityMiddleware provides security-related middleware functions
type SecurityMiddleware struct {
	config *config.Config
	logger *zap.Logger
}

// NewSecurityMiddleware creates a new SecurityMiddleware instance
func NewSecurityMiddleware(config *config.Config, logger *zap.Logger) *SecurityMiddleware {
	return &SecurityMiddleware{
		config: config,
		logger: logger,
	}
}

// CORS returns a CORS middleware
func (sm *SecurityMiddleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range sm.config.Security.CORS.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if !allowed {
			sm.logger.Warn("CORS request rejected", zap.String("origin", origin))
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Methods", strings.Join(sm.config.Security.CORS.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", strings.Join(sm.config.Security.CORS.AllowedHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds security-related headers to responses
func (sm *SecurityMiddleware) SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS
		if sm.config.Security.Headers.EnableHSTS {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}

		// XSS Protection
		if sm.config.Security.Headers.EnableXSSFilter {
			c.Header("X-XSS-Protection", "1; mode=block")
		}

		// Prevent MIME type sniffing
		if sm.config.Security.Headers.EnableNoSniff {
			c.Header("X-Content-Type-Options", "nosniff")
		}

		// Frame options
		c.Header("X-Frame-Options", sm.config.Security.Headers.FrameOptions)

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline';")

		c.Next()
	}
}

// RequestSizeLimit limits the size of incoming requests
func (sm *SecurityMiddleware) RequestSizeLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, int64(sm.config.Server.MaxHeaderBytes))
		c.Next()
	}
}

// ValidateContentType ensures the request has the correct content type
func (sm *SecurityMiddleware) ValidateContentType() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				sm.logger.Warn("invalid content type", 
					zap.String("content_type", contentType),
					zap.String("method", c.Request.Method),
				)
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"error": "Content-Type must be application/json",
				})
				return
			}
		}
		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func (sm *SecurityMiddleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// generateRequestID creates a unique request ID
func generateRequestID() string {
	// Implementation using UUID or similar
	return "req-" + time.Now().Format("20060102150405") + "-" + randString(8)
}

// randString generates a random string of the specified length
func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
} 
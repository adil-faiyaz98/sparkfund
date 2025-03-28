package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/pkg/logger"
	"github.com/sparkfund/pkg/metrics"
)

// Logger logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		// Log request details
		logger.Info("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", statusCode,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)

		// Track metrics
		metrics.TrackHTTPRequest(c.Request.Method, path, statusCode, latency)
	}
}

// Recovery recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", "error", err)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}

// CORS handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24 hours

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RateLimit implements rate limiting
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	// TODO: Implement rate limiting using Redis or in-memory store
	return func(c *gin.Context) {
		c.Next()
	}
}

// Auth validates JWT tokens
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		// TODO: Implement JWT validation
		// token := strings.TrimPrefix(authHeader, "Bearer ")
		// claims, err := jwt.ValidateToken(token)
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 		"error": "Invalid token",
		// 	})
		// 	return
		// }

		// c.Set("user", claims)
		c.Next()
	}
}

// RequireRole checks if the user has the required role
func RequireRole(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement role validation
		// user, exists := c.Get("user")
		// if !exists {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		// 		"error": "User not found",
		// 	})
		// 	return
		// }

		// claims := user.(*jwt.Claims)
		// if !claims.HasRole(role) {
		// 	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		// 		"error": "Insufficient permissions",
		// 	})
		// 	return
		// }

		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// TODO: Generate a unique request ID
			// requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

// HealthCheck provides a health check endpoint
func HealthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now(),
		})
	}
}

// Metrics exposes Prometheus metrics
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement Prometheus metrics endpoint
		c.Next()
	}
}

// Timeout adds a timeout to the request context
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// Compression compresses the response
func Compression() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement response compression
		c.Next()
	}
}

// CacheControl adds cache control headers
func CacheControl(maxAge time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(maxAge.Seconds())))
		c.Next()
	}
}

// Security adds security headers
func Security() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Next()
	}
} 
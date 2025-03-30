package middleware

import (
	"net/http"
	"strings"
	"time"

	"errors"
	"investment-service/internal/config"
	"investment-service/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// RateLimiter implements a simple rate limiter
func RateLimiter() gin.HandlerFunc {
	// Create a map to store limiters for each client
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get or create rate limiter for this client
		limiter, exists := limiters[clientIP]
		if !exists {
			// Allow 5 requests per second with burst of 10
			limiter = rate.NewLimiter(5, 10)
			limiters[clientIP] = limiter
		}

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, models.ErrorResponse{
				Error: "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// JWTAuth validates JWT tokens
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for health endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" || c.Request.URL.Path == "/live" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Invalid signing method")
			}
			return []byte(config.Get().JWT.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Error: "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("userID", claims["sub"])
		c.Set("roles", claims["roles"])

		c.Next()
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// SecurityHeaders adds security headers to responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add security headers
		c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
		c.Writer.Header().Set("X-Frame-Options", "DENY")
		c.Writer.Header().Set("Content-Security-Policy", "default-src 'self'")
		c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// RequestLogger logs request details
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userID, _ := c.Get("userID")

		// Log details
		// In production, use a structured logger
		if statusCode >= 400 {
			// Log errors with more detail
			c.Error(nil).SetMeta(gin.H{
				"status":    statusCode,
				"latency":   latency,
				"client_ip": clientIP,
				"method":    method,
				"path":      path,
				"user":      userID,
			})
		} else {
			// Log success requests more briefly
			c.Set("log_latency", latency)
		}
	}
}

package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/time/rate"
)

// RateLimiterConfig holds configuration for the rate limiter
type RateLimiterConfig struct {
	Enabled  bool
	Requests int
	Window   time.Duration
	Burst    int
}

// DefaultRateLimiterConfig returns default rate limiter configuration
func DefaultRateLimiterConfig() RateLimiterConfig {
	return RateLimiterConfig{
		Enabled:  true,
		Requests: 60,
		Window:   time.Minute,
		Burst:    10,
	}
}

// RateLimiter implements a simple rate limiter
func RateLimiter(cfg RateLimiterConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// Create a map to store limiters for each client
	limiters := make(map[string]*rate.Limiter)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get or create rate limiter for this client
		limiter, exists := limiters[clientIP]
		if !exists {
			// Create rate limiter with configured values
			limiter = rate.NewLimiter(rate.Limit(cfg.Requests)/rate.Limit(cfg.Window.Seconds()), cfg.Burst)
			limiters[clientIP] = limiter
		}

		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, ErrorResponse{
				Error: "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// JWTConfig holds configuration for JWT authentication
type JWTConfig struct {
	Secret  string
	Enabled bool
}

// DefaultJWTConfig returns default JWT configuration
func DefaultJWTConfig() JWTConfig {
	return JWTConfig{
		Secret:  "your-secret-key",
		Enabled: true,
	}
}

// JWTAuth validates JWT tokens
func JWTAuth(cfg JWTConfig) gin.HandlerFunc {
	if !cfg.Enabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		// Skip auth for health endpoints
		if c.Request.URL.Path == "/health" || c.Request.URL.Path == "/ready" || c.Request.URL.Path == "/live" || c.Request.URL.Path == "/metrics" {
			c.Next()
			return
		}

		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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
			return []byte(cfg.Secret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, ErrorResponse{
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

// CORSConfig holds configuration for CORS
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// DefaultCORSConfig returns default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "X-CSRF-Token"},
	}
}

// CORS middleware
func CORS(cfg CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		
		// Check if the origin is allowed
		allowedOrigin := "*"
		if len(cfg.AllowedOrigins) > 0 && cfg.AllowedOrigins[0] != "*" {
			for _, o := range cfg.AllowedOrigins {
				if o == origin {
					allowedOrigin = origin
					break
				}
			}
		}
		
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
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

// CSRFProtection adds CSRF protection
func CSRFProtection() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for GET, HEAD, OPTIONS, TRACE
		if c.Request.Method == "GET" || 
		   c.Request.Method == "HEAD" || 
		   c.Request.Method == "OPTIONS" || 
		   c.Request.Method == "TRACE" {
			c.Next()
			return
		}
		
		// Check for CSRF token in header
		token := c.GetHeader("X-CSRF-Token")
		if token == "" {
			c.JSON(http.StatusForbidden, ErrorResponse{
				Error: "CSRF token required",
			})
			c.Abort()
			return
		}
		
		// In a real implementation, validate the token against the session
		// For now, we'll just check that it's not empty
		
		c.Next()
	}
}

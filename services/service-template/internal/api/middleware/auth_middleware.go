package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	logger *logrus.Logger
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(logger *logrus.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		logger: logger,
	}
}

// Authenticate authenticates a request
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header is in the correct format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Get token
		token := parts[1]
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// TODO: Validate token
		// This is a placeholder for token validation
		// In a real application, you would validate the token using a JWT library
		if token == "invalid" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Set user ID in context
		// In a real application, you would extract the user ID from the token
		c.Set("user_id", "123")

		c.Next()
	}
}

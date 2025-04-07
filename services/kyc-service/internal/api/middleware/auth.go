package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"sparkfund/services/kyc-service/internal/api/dto"
)

// AuthConfig contains authentication configuration
type AuthConfig struct {
	JWTSecret     string
	TokenHeader   string
	TokenPrefix   string
	ExcludedPaths []string
}

// Auth returns a gin middleware for authentication
func Auth(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path is excluded
		path := c.Request.URL.Path
		for _, excludedPath := range config.ExcludedPaths {
			if strings.HasPrefix(path, excludedPath) {
				c.Next()
				return
			}
		}

		// Get token from header
		authHeader := c.GetHeader(config.TokenHeader)
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check token prefix
		if !strings.HasPrefix(authHeader, config.TokenPrefix) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, config.TokenPrefix+" ")

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Invalid token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Invalid token claims",
			})
			c.Abort()
			return
		}

		// Set user ID in context
		userID, ok := claims["sub"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: "Invalid user ID in token",
			})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

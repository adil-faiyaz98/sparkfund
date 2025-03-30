package security

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Roles    []string `json:"roles"`
	Username string   `json:"username"`
	jwt.RegisteredClaims
}

// AuthConfig holds authentication configuration
type AuthConfig struct {
	SecretKey     string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

// DefaultAuthConfig returns default authentication configuration
func DefaultAuthConfig() AuthConfig {
	return AuthConfig{
		SecretKey:     "your-secret-key", // Change this in production
		TokenExpiry:   15 * time.Minute,
		RefreshExpiry: 24 * time.Hour,
	}
}

// AuthMiddleware provides JWT authentication
func AuthMiddleware(config AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for public endpoints
		if isPublicEndpoint(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get token from header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Parse and validate token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(config.SecretKey), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid"})
			c.Abort()
			return
		}

		// Add claims to context
		c.Set("user_id", claims.UserID)
		c.Set("roles", claims.Roles)
		c.Set("username", claims.Username)

		c.Next()
	}
}

// RequireRole middleware checks if the user has the required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "User roles not found"})
			c.Abort()
			return
		}

		userRolesList, ok := userRoles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid user roles format"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			for _, userRole := range userRolesList {
				if role == userRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GenerateToken generates a new JWT token
func GenerateToken(userID, username string, roles []string, config AuthConfig) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// GenerateRefreshToken generates a new refresh token
func GenerateRefreshToken(userID string, config AuthConfig) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.RefreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.SecretKey))
}

// isPublicEndpoint checks if the endpoint is public
func isPublicEndpoint(path string) bool {
	publicPaths := []string{
		"/health",
		"/metrics",
		"/api/v1/auth/login",
		"/api/v1/auth/refresh",
		"/api/v1/auth/register",
	}

	for _, publicPath := range publicPaths {
		if path == publicPath {
			return true
		}
	}
	return false
} 
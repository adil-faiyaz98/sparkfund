package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sparkfund/credit-scoring-service/internal/errors"
	"go.uber.org/zap"
)

type JWTConfig struct {
	SecretKey     string
	TokenExpiry   time.Duration
	Issuer        string
	Audience      string
	AllowedScopes []string
}

func NewJWTConfig(secretKey string) *JWTConfig {
	return &JWTConfig{
		SecretKey:     secretKey,
		TokenExpiry:   24 * time.Hour,
		Issuer:        "sparkfund-credit-service",
		Audience:      "sparkfund-api",
		AllowedScopes: []string{"credit:read", "credit:write"},
	}
}

func AuthMiddleware(config *JWTConfig, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Error("missing authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAPIError(
				errors.ErrUnauthorized,
				"missing authorization header",
			))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Error("invalid authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAPIError(
				errors.ErrUnauthorized,
				"invalid authorization header format",
			))
			return
		}

		token, err := validateToken(parts[1], config)
		if err != nil {
			logger.Error("token validation failed", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAPIError(
				errors.ErrUnauthorized,
				"invalid token",
			))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			logger.Error("invalid token claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, errors.NewAPIError(
				errors.ErrUnauthorized,
				"invalid token claims",
			))
			return
		}

		// Validate scopes
		scopes, ok := claims["scope"].(string)
		if !ok {
			logger.Error("missing scope in token")
			c.AbortWithStatusJSON(http.StatusForbidden, errors.NewAPIError(
				errors.ErrForbidden,
				"missing required scopes",
			))
			return
		}

		userScopes := strings.Split(scopes, " ")
		if !hasRequiredScopes(userScopes, config.AllowedScopes) {
			logger.Error("insufficient scopes", zap.Strings("user_scopes", userScopes))
			c.AbortWithStatusJSON(http.StatusForbidden, errors.NewAPIError(
				errors.ErrForbidden,
				"insufficient scopes",
			))
			return
		}

		// Set user context
		c.Set("user_id", claims["sub"])
		c.Set("scopes", userScopes)
		c.Next()
	}
}

func validateToken(tokenString string, config *JWTConfig) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewAPIError(
				errors.ErrUnauthorized,
				"unexpected signing method",
			)
		}
		return []byte(config.SecretKey), nil
	})
}

func hasRequiredScopes(userScopes, requiredScopes []string) bool {
	for _, required := range requiredScopes {
		found := false
		for _, user := range userScopes {
			if user == required {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
} 
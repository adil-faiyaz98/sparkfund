package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sparkfund/pkg/errors"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.StandardClaims
}

// Config represents the JWT configuration
type Config struct {
	SecretKey     string
	ExpirationTime time.Duration
}

// NewConfig creates a new JWT configuration
func NewConfig(secretKey string, expirationTime time.Duration) *Config {
	return &Config{
		SecretKey:     secretKey,
		ExpirationTime: expirationTime,
	}
}

// GenerateToken generates a new JWT token
func (c *Config) GenerateToken(userID, username string, roles []string) (string, error) {
	claims := &Claims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(c.ExpirationTime).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.SecretKey))
}

// ValidateToken validates a JWT token
func (c *Config) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.SecretKey), nil
	})

	if err != nil {
		return nil, errors.ErrUnauthorized(err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrUnauthorized(fmt.Errorf("invalid token"))
}

// HasRole checks if the claims have a specific role
func (c *Claims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the claims have any of the specified roles
func (c *Claims) HasAnyRole(roles []string) bool {
	for _, role := range roles {
		if c.HasRole(role) {
			return true
		}
	}
	return false
}

// HasAllRoles checks if the claims have all of the specified roles
func (c *Claims) HasAllRoles(roles []string) bool {
	for _, role := range roles {
		if !c.HasRole(role) {
			return false
		}
	}
	return true
} 
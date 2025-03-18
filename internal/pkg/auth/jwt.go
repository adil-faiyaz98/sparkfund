package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenMaker is an interface for managing tokens
type TokenMaker interface {
	// CreateToken creates a new token for a specific user id and duration
	CreateToken(userID uint, duration time.Duration) (string, error)

	// VerifyToken checks if the token is valid and returns the user id
	VerifyToken(token string) (uint, error)
}

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// CustomClaims contains the claims data
type CustomClaims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (TokenMaker, error) {
	if len(secretKey) < 32 {
		return nil, errors.New("secret key must be at least 32 characters")
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new token for a specific user ID and duration
func (maker *JWTMaker) CreateToken(userID uint, duration time.Duration) (string, error) {
	claims := &CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "money-pulse",
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid and returns the user id
func (maker *JWTMaker) VerifyToken(tokenString string) (uint, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(maker.secretKey), nil
		},
	)
	if err != nil {
		return 0, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token claims")
	}

	return claims.UserID, nil
}

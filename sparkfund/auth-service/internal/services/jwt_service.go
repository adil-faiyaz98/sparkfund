package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// JWTService handles JWT token validation and generation
type JWTService struct {
	logger        *zap.Logger
	cognitoClient *cognitoidentityprovider.Client
	userPool      string
	jwtSecret     []byte
	useCognito    bool
}

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTService creates a new JWT service instance
func NewJWTService(logger *zap.Logger) (*JWTService, error) {
	useCognito := os.Getenv("USE_COGNITO") == "true"
	service := &JWTService{
		logger:     logger,
		jwtSecret:  []byte(os.Getenv("JWT_SECRET")),
		useCognito: useCognito,
	}

	if useCognito {
		cfg, err := config.LoadDefaultConfig(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("unable to load AWS config: %v", err)
		}

		service.cognitoClient = cognitoidentityprovider.NewFromConfig(cfg)
		service.userPool = os.Getenv("COGNITO_USER_POOL_ID")
	}

	return service, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	if s.useCognito {
		return s.validateCognitoToken(tokenString)
	}
	return s.validateCustomToken(tokenString)
}

// validateCognitoToken validates a token using AWS Cognito
func (s *JWTService) validateCognitoToken(tokenString string) (*JWTClaims, error) {
	input := &cognitoidentityprovider.GetUserInput{
		AccessToken: &tokenString,
	}

	result, err := s.cognitoClient.GetUser(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to validate Cognito token: %v", err)
	}

	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// Extract user attributes from Cognito response
	for _, attr := range result.UserAttributes {
		switch *attr.Name {
		case "sub":
			claims.UserID = *attr.Value
		case "username":
			claims.Username = *attr.Value
		case "custom:role":
			claims.Role = *attr.Value
		}
	}

	return claims, nil
}

// validateCustomToken validates a custom JWT token
func (s *JWTService) validateCustomToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GenerateToken generates a new JWT token
func (s *JWTService) GenerateToken(userID, username, role string) (string, error) {
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// ExtractTokenFromHeader extracts the JWT token from the Authorization header
func (s *JWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return "", fmt.Errorf("invalid authorization header format")
	}
	return authHeader[7:], nil
}

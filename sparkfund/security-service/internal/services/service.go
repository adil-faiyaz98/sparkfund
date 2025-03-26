package services

import (
	"context"
	"sparkfund/security-service/internal/config"
	"sparkfund/security-service/internal/models"
	"sparkfund/security-service/internal/repositories"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// Service defines the interface for security service operations
type Service interface {
	ValidateToken(ctx context.Context, req models.TokenValidationRequest) (models.TokenValidationResponse, error)
	GenerateToken(ctx context.Context, req models.TokenGenerationRequest) (models.TokenGenerationResponse, error)
	RefreshToken(ctx context.Context, req models.TokenRefreshRequest) (models.TokenRefreshResponse, error)
}

// securityService implements the Service interface
type securityService struct {
	logger     *zap.Logger
	config     *config.Config
	repository repositories.Repository
}

// NewService creates a new security service
func NewService(logger *zap.Logger, config *config.Config, repository repositories.Repository) Service {
	return &securityService{
		logger:     logger,
		config:     config,
		repository: repository,
	}
}

// ValidateToken validates a JWT token
func (s *securityService) ValidateToken(ctx context.Context, req models.TokenValidationRequest) (models.TokenValidationResponse, error) {
	s.logger.Info("Validating token")

	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	resp := models.TokenValidationResponse{
		Valid: false,
	}

	if err != nil {
		s.logger.Error("Failed to parse token", zap.Error(err))
		return resp, err
	}

	if !token.Valid {
		return resp, nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		s.logger.Error("Failed to extract claims")
		return resp, nil
	}

	resp.Valid = true
	resp.Claims = make(map[string]interface{})

	if sub, ok := claims["sub"].(string); ok {
		resp.Subject = sub
	}

	if iss, ok := claims["iss"].(string); ok {
		resp.Issuer = iss
	}

	if iat, ok := claims["iat"].(float64); ok {
		resp.IssuedAt = time.Unix(int64(iat), 0)
	}

	if exp, ok := claims["exp"].(float64); ok {
		resp.ExpiresAt = time.Unix(int64(exp), 0)
	}

	// Copy remaining claims
	for k, v := range claims {
		if k != "sub" && k != "iss" && k != "iat" && k != "exp" {
			resp.Claims[k] = v
		}
	}

	return resp, nil
}

// GenerateToken generates a new JWT token
func (s *securityService) GenerateToken(ctx context.Context, req models.TokenGenerationRequest) (models.TokenGenerationResponse, error) {
	s.logger.Info("Generating token", zap.String("subject", req.Subject))

	// Set default expiration if not provided
	expiresIn := req.ExpiresIn
	if expiresIn == 0 {
		expiresIn = int64(s.config.JWT.ExpireTime)
	}

	expiresAt := time.Now().Add(time.Duration(expiresIn) * time.Second)

	// Create claims
	claims := jwt.MapClaims{
		"sub": req.Subject,
		"iss": "sparkfund-security-service",
		"iat": time.Now().Unix(),
		"exp": expiresAt.Unix(),
	}

	// Add custom claims
	for k, v := range req.Claims {
		claims[k] = v
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		s.logger.Error("Failed to sign token", zap.Error(err))
		return models.TokenGenerationResponse{}, err
	}

	// Generate refresh token if needed
	refreshToken := ""
	// Add refresh token implementation if needed

	return models.TokenGenerationResponse{
		Token:        tokenString,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
	}, nil
}

// RefreshToken refreshes a JWT token
func (s *securityService) RefreshToken(ctx context.Context, req models.TokenRefreshRequest) (models.TokenRefreshResponse, error) {
	s.logger.Info("Refreshing token")

	// Add refresh token implementation
	// This is a placeholder implementation

	expiresAt := time.Now().Add(time.Duration(s.config.JWT.ExpireTime) * time.Second)

	return models.TokenRefreshResponse{
		Token:        "new-token",
		RefreshToken: "new-refresh-token",
		ExpiresAt:    expiresAt,
	}, nil
}

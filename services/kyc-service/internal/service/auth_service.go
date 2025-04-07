package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/pquerna/otp/totp"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"sparkfund/services/kyc-service/internal/model"
	"sparkfund/services/kyc-service/internal/repository"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
	logger      *logrus.Logger
	jwtSecret   []byte
	jwtExpiry   time.Duration
	mfaEnabled  bool
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo *repository.UserRepository,
	sessionRepo *repository.SessionRepository,
	logger *logrus.Logger,
	jwtSecret string,
	jwtExpiry time.Duration,
	mfaEnabled bool,
) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		logger:      logger,
		jwtSecret:   []byte(jwtSecret),
		jwtExpiry:   jwtExpiry,
		mfaEnabled:  mfaEnabled,
	}
}

// Login authenticates a user
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest, deviceInfo model.DeviceInfo) (*model.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		s.logger.WithError(err).WithField("email", req.Email).Error("User not found during login")
		return nil, errors.New("invalid email or password")
	}

	// Check if account is locked
	if user.LockedUntil != nil && time.Now().Before(*user.LockedUntil) {
		s.logger.WithField("email", req.Email).Warn("Attempt to login to locked account")
		return nil, errors.New("account is locked, please try again later")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Increment login attempts
		user.LoginAttempts++
		if user.LoginAttempts >= 5 {
			lockUntil := time.Now().Add(15 * time.Minute)
			user.LockedUntil = &lockUntil
			s.logger.WithField("email", req.Email).Warn("Account locked due to too many failed login attempts")
		}
		s.userRepo.Update(ctx, user)
		
		return nil, errors.New("invalid email or password")
	}

	// Reset login attempts on successful password verification
	user.LoginAttempts = 0
	user.LockedUntil = nil

	// Check if MFA is required
	if s.mfaEnabled && user.MFAEnabled {
		// If MFA code is not provided, return MFA required response
		if req.MFACode == "" {
			return &model.LoginResponse{
				MFARequired: true,
				User: model.User{
					ID:    user.ID,
					Email: user.Email,
				},
			}, nil
		}

		// Verify MFA code
		valid := totp.Validate(req.MFACode, user.MFASecret)
		if !valid {
			s.logger.WithField("email", req.Email).Warn("Invalid MFA code provided")
			return nil, errors.New("invalid MFA code")
		}
	}

	// Generate tokens
	token, expiresAt, err := s.generateJWT(user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate JWT")
		return nil, errors.New("authentication failed")
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate refresh token")
		return nil, errors.New("authentication failed")
	}

	// Create session
	session := &model.Session{
		ID:           uuid.New(),
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    deviceInfo.UserAgent,
		IPAddress:    deviceInfo.IPAddress,
		ExpiresAt:    time.Now().Add(30 * 24 * time.Hour), // 30 days
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.sessionRepo.Create(ctx, session)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to create session")
		return nil, errors.New("authentication failed")
	}

	// Update user's last login information
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = deviceInfo.IPAddress
	user.UpdatedAt = now

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to update user's last login info")
		// Continue despite error
	}

	// Return login response
	return &model.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *user,
		MFARequired:  false,
	}, nil
}

// RefreshToken refreshes an authentication token
func (s *AuthService) RefreshToken(ctx context.Context, req model.RefreshTokenRequest) (*model.LoginResponse, error) {
	// Get session by refresh token
	session, err := s.sessionRepo.GetByRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		s.logger.WithError(err).Error("Session not found during token refresh")
		return nil, errors.New("invalid refresh token")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		s.logger.WithField("session_id", session.ID).Warn("Attempt to use expired refresh token")
		return nil, errors.New("refresh token expired")
	}

	// Get user
	user, err := s.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", session.UserID).Error("User not found during token refresh")
		return nil, errors.New("user not found")
	}

	// Generate new JWT
	token, expiresAt, err := s.generateJWT(user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate JWT during token refresh")
		return nil, errors.New("token refresh failed")
	}

	// Generate new refresh token
	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		s.logger.WithError(err).WithField("user_id", user.ID).Error("Failed to generate refresh token during token refresh")
		return nil, errors.New("token refresh failed")
	}

	// Update session
	session.RefreshToken = refreshToken
	session.ExpiresAt = time.Now().Add(30 * 24 * time.Hour) // 30 days
	session.UpdatedAt = time.Now()

	err = s.sessionRepo.Update(ctx, session)
	if err != nil {
		s.logger.WithError(err).WithField("session_id", session.ID).Error("Failed to update session during token refresh")
		return nil, errors.New("token refresh failed")
	}

	// Return login response
	return &model.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		User:         *user,
		MFARequired:  false,
	}, nil
}

// Logout logs out a user
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// Get session by refresh token
	session, err := s.sessionRepo.GetByRefreshToken(ctx, refreshToken)
	if err != nil {
		s.logger.WithError(err).Error("Session not found during logout")
		return nil // Return success even if session not found
	}

	// Delete session
	err = s.sessionRepo.Delete(ctx, session.ID)
	if err != nil {
		s.logger.WithError(err).WithField("session_id", session.ID).Error("Failed to delete session during logout")
		return errors.New("logout failed")
	}

	return nil
}

// SetupMFA sets up MFA for a user
func (s *AuthService) SetupMFA(ctx context.Context, userID uuid.UUID) (*model.MFASetupResponse, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("User not found during MFA setup")
		return nil, errors.New("user not found")
	}

	// Generate MFA secret
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "KYC Service",
		AccountName: user.Email,
	})
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to generate MFA secret")
		return nil, errors.New("MFA setup failed")
	}

	// Generate recovery codes
	recoveryCodes := make([]string, 10)
	for i := 0; i < 10; i++ {
		code, err := generateRandomString(10)
		if err != nil {
			s.logger.WithError(err).WithField("user_id", userID).Error("Failed to generate recovery codes")
			return nil, errors.New("MFA setup failed")
		}
		recoveryCodes[i] = code
	}

	// Update user with MFA secret
	user.MFASecret = key.Secret()
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to update user with MFA secret")
		return nil, errors.New("MFA setup failed")
	}

	// Return MFA setup response
	return &model.MFASetupResponse{
		Secret:        key.Secret(),
		QRCodeURL:     key.URL(),
		RecoveryCodes: recoveryCodes,
	}, nil
}

// VerifyMFA verifies an MFA code
func (s *AuthService) VerifyMFA(ctx context.Context, userID uuid.UUID, req model.MFAVerifyRequest) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("User not found during MFA verification")
		return errors.New("user not found")
	}

	// Verify MFA code
	valid := totp.Validate(req.Code, user.MFASecret)
	if !valid {
		s.logger.WithField("user_id", userID).Warn("Invalid MFA code provided")
		return errors.New("invalid MFA code")
	}

	// Enable MFA for user
	user.MFAEnabled = true
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to enable MFA for user")
		return errors.New("MFA verification failed")
	}

	return nil
}

// DisableMFA disables MFA for a user
func (s *AuthService) DisableMFA(ctx context.Context, userID uuid.UUID) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("User not found during MFA disabling")
		return errors.New("user not found")
	}

	// Disable MFA for user
	user.MFAEnabled = false
	user.MFASecret = ""
	user.UpdatedAt = time.Now()

	err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.logger.WithError(err).WithField("user_id", userID).Error("Failed to disable MFA for user")
		return errors.New("MFA disabling failed")
	}

	return nil
}

// ValidateToken validates a JWT token
func (s *AuthService) ValidateToken(tokenString string) (*model.JWTClaims, error) {
	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Create JWT claims
	jwtClaims := &model.JWTClaims{
		UserID:    claims["user_id"].(string),
		Email:     claims["email"].(string),
		Role:      claims["role"].(string),
		MFAPassed: claims["mfa_passed"].(bool),
	}

	return jwtClaims, nil
}

// generateJWT generates a JWT token
func (s *AuthService) generateJWT(user *model.User) (string, time.Time, error) {
	// Set expiration time
	expiresAt := time.Now().Add(s.jwtExpiry)

	// Create claims
	claims := jwt.MapClaims{
		"user_id":    user.ID.String(),
		"email":      user.Email,
		"role":       user.Role,
		"mfa_passed": user.MFAEnabled,
		"exp":        expiresAt.Unix(),
		"iat":        time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// generateRefreshToken generates a refresh token
func (s *AuthService) generateRefreshToken() (string, error) {
	// Generate random bytes
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode to base64
	return base64.URLEncoding.EncodeToString(b), nil
}

// generateRandomString generates a random string
func generateRandomString(length int) (string, error) {
	// Generate random bytes
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	// Encode to base64 and trim to desired length
	return base64.URLEncoding.EncodeToString(b)[:length], nil
}

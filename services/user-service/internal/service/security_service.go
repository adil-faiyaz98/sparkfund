package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/sparkfund/user-service/internal/model"
	"github.com/sparkfund/user-service/internal/repository"
	"go.uber.org/zap"
)

type SecurityService struct {
	userRepo     repository.UserRepository
	jwtSecret    []byte
	jwtExpiry    time.Duration
	emailService EmailService
	smsService   SMSService
	logger       *zap.Logger
}

func NewSecurityService(
	userRepo repository.UserRepository,
	jwtSecret string,
	jwtExpiry time.Duration,
	emailService EmailService,
	smsService SMSService,
	logger *zap.Logger,
) *SecurityService {
	return &SecurityService{
		userRepo:     userRepo,
		jwtSecret:    []byte(jwtSecret),
		jwtExpiry:    jwtExpiry,
		emailService: emailService,
		smsService:   smsService,
		logger:       logger,
	}
}

// GenerateTOTP generates a new TOTP secret for MFA
func (s *SecurityService) GenerateTOTP(userID string) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SparkFund",
		AccountName: userID,
		SecretSize:  20,
		Secret:      generateRandomSecret(),
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to generate TOTP: %w", err)
	}

	return key, nil
}

// VerifyTOTP verifies a TOTP code
func (s *SecurityService) VerifyTOTP(secret, code string) bool {
	return totp.Validate(code, secret)
}

// GenerateJWT generates a new JWT token
func (s *SecurityService) GenerateJWT(userID string, roles []string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   userID,
		"roles": roles,
		"exp":   time.Now().Add(s.jwtExpiry).Unix(),
		"iat":   time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

// VerifyJWT verifies a JWT token
func (s *SecurityService) VerifyJWT(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// SendVerificationEmail sends a verification email
func (s *SecurityService) SendVerificationEmail(ctx context.Context, user *model.User) error {
	token, err := s.GenerateJWT(user.ID.String(), []string{"email_verification"})
	if err != nil {
		return err
	}

	verificationLink := fmt.Sprintf("https://api.sparkfund.com/verify-email?token=%s", token)

	return s.emailService.SendVerificationEmail(ctx, user.Email, verificationLink)
}

// SendVerificationSMS sends a verification SMS
func (s *SecurityService) SendVerificationSMS(ctx context.Context, user *model.User) error {
	code := generateRandomCode(6)

	// Store the code in Redis with expiration
	err := s.userRepo.StoreVerificationCode(ctx, user.ID.String(), code, 5*time.Minute)
	if err != nil {
		return err
	}

	return s.smsService.SendVerificationCode(ctx, user.PhoneNumber, code)
}

// SendMFACode sends an MFA code
func (s *SecurityService) SendMFACode(ctx context.Context, user *model.User) error {
	switch user.MFAMethod {
	case model.MFAMethodSMS:
		return s.SendVerificationSMS(ctx, user)
	case model.MFAMethodEmail:
		code := generateRandomCode(6)
		err := s.userRepo.StoreVerificationCode(ctx, user.ID.String(), code, 5*time.Minute)
		if err != nil {
			return err
		}
		return s.emailService.SendMFACode(ctx, user.Email, code)
	default:
		return fmt.Errorf("unsupported MFA method: %s", user.MFAMethod)
	}
}

// GenerateRecoveryCodes generates recovery codes for account recovery
func (s *SecurityService) GenerateRecoveryCodes() []string {
	codes := make([]string, 8)
	for i := range codes {
		codes[i] = generateRandomCode(10)
	}
	return codes
}

// Helper functions
func generateRandomSecret() string {
	bytes := make([]byte, 20)
	rand.Read(bytes)
	return base32.StdEncoding.EncodeToString(bytes)
}

func generateRandomCode(length int) string {
	const charset = "0123456789"
	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

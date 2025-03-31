package service

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/services/user-service/internal/errors"
	"github.com/sparkfund/services/user-service/internal/logger"
	"github.com/sparkfund/services/user-service/internal/models"
	"github.com/sparkfund/services/user-service/internal/repository"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// RegisterUser registers a new user
func (s *UserService) RegisterUser(ctx context.Context, user *models.User) error {
	// Hash password before storing
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return s.userRepo.Create(ctx, user)
}

// AuthenticateUser authenticates a user with email and password
func (s *UserService) AuthenticateUser(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !verifyPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.Get(ctx, userID)
}

// UpdateUser updates user details
func (s *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

// GetUserProfile retrieves a user's profile
func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error) {
	return s.userRepo.GetProfile(ctx, userID)
}

// UpdateUserProfile updates a user's profile
func (s *UserService) UpdateUserProfile(ctx context.Context, userID uuid.UUID, profile *models.UserProfile) error {
	profile.UpdatedAt = time.Now()
	return s.userRepo.UpdateProfile(ctx, userID, profile)
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	user, err := s.userRepo.Get(ctx, userID)
	if err != nil {
		return err
	}

	if !verifyPassword(oldPassword, user.Password) {
		return errors.New("invalid old password")
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(ctx, userID, hashedPassword)
}

// ResetPassword initiates a password reset
func (s *UserService) ResetPassword(ctx context.Context, email string) error {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return err
	}

	token := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	return s.userRepo.StoreResetToken(ctx, user.ID, token, expiresAt)
}

// ConfirmResetPassword confirms a password reset
func (s *UserService) ConfirmResetPassword(ctx context.Context, token, newPassword string) error {
	resetToken, err := s.userRepo.GetResetToken(ctx, token)
	if err != nil {
		return err
	}

	if time.Now().After(resetToken.ExpiresAt) {
		return errors.New("reset token expired")
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePassword(ctx, resetToken.UserID, hashedPassword); err != nil {
		return err
	}

	return s.userRepo.MarkResetTokenUsed(ctx, token)
}

// CreateSession creates a new user session
func (s *UserService) CreateSession(ctx context.Context, session *models.Session) error {
	return s.userRepo.CreateSession(ctx, session)
}

// GetSession retrieves a session by token
func (s *UserService) GetSession(ctx context.Context, token string) (*models.Session, error) {
	return s.userRepo.GetSession(ctx, token)
}

// UpdateSession updates a session's last used time
func (s *UserService) UpdateSession(ctx context.Context, session *models.Session) error {
	session.LastUsed = time.Now()
	return s.userRepo.UpdateSession(ctx, session)
}

// DeleteSession deletes a session
func (s *UserService) DeleteSession(ctx context.Context, token string) error {
	return s.userRepo.DeleteSession(ctx, token)
}

// DeleteExpiredSessions deletes all expired sessions
func (s *UserService) DeleteExpiredSessions(ctx context.Context) error {
	return s.userRepo.DeleteExpiredSessions(ctx)
}

// IncrementFailedAttempts increments failed login attempts for a user
func (s *UserService) IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.IncrementFailedAttempts(ctx, userID)
}

// ResetFailedAttempts resets failed login attempts for a user
func (s *UserService) ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error {
	return s.userRepo.ResetFailedAttempts(ctx, userID)
}

// GetFailedAttempts gets the number of failed login attempts for a user
func (s *UserService) GetFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.userRepo.GetFailedAttempts(ctx, userID)
}

// ListUserSessions lists all active sessions for a user
func (s *UserService) ListUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error) {
	sessions, err := s.userRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get user sessions")
	}

	// Filter out expired sessions
	var activeSessions []models.Session
	for _, session := range sessions {
		if session.ExpiresAt.After(time.Now()) {
			activeSessions = append(activeSessions, session)
		}
	}

	return activeSessions, nil
}

// RevokeSession revokes a specific session
func (s *UserService) RevokeSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	// Verify session belongs to user
	session, err := s.userRepo.GetSession(ctx, sessionID)
	if err != nil {
		return errors.Wrap(err, "Failed to get session")
	}

	if session.UserID != userID {
		return errors.ErrUnauthorized
	}

	if err := s.userRepo.DeleteSession(ctx, sessionID); err != nil {
		return errors.Wrap(err, "Failed to revoke session")
	}

	logger.Info("Session revoked", map[string]interface{}{
		"user_id":    userID,
		"session_id": sessionID,
	})

	return nil
}

// RevokeAllSessions revokes all sessions for a user
func (s *UserService) RevokeAllSessions(ctx context.Context, userID uuid.UUID) error {
	sessions, err := s.userRepo.GetUserSessions(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "Failed to get user sessions")
	}

	for _, session := range sessions {
		if err := s.userRepo.DeleteSession(ctx, session.ID); err != nil {
			return errors.Wrap(err, "Failed to revoke session")
		}
	}

	logger.Info("All sessions revoked", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// EnableMFA enables MFA for a user
func (s *UserService) EnableMFA(ctx context.Context, userID uuid.UUID) (string, error) {
	// Generate a random secret
	secret := make([]byte, 20)
	if _, err := rand.Read(secret); err != nil {
		return "", errors.Wrap(err, "Failed to generate MFA secret")
	}

	// Encode secret as base32
	encodedSecret := base32.StdEncoding.EncodeToString(secret)

	// Store MFA secret
	if err := s.userRepo.StoreMFASecret(ctx, userID, encodedSecret); err != nil {
		return "", errors.Wrap(err, "Failed to store MFA secret")
	}

	logger.Info("MFA enabled", map[string]interface{}{
		"user_id": userID,
	})

	return encodedSecret, nil
}

// DisableMFA disables MFA for a user
func (s *UserService) DisableMFA(ctx context.Context, userID uuid.UUID) error {
	if err := s.userRepo.DeleteMFASecret(ctx, userID); err != nil {
		return errors.Wrap(err, "Failed to disable MFA")
	}

	logger.Info("MFA disabled", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// VerifyMFA verifies an MFA code
func (s *UserService) VerifyMFA(ctx context.Context, userID uuid.UUID, code string) error {
	secret, err := s.userRepo.GetMFASecret(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "Failed to get MFA secret")
	}

	// TODO: Implement TOTP verification
	// For now, just check if code matches a simple pattern
	if code != "123456" {
		return errors.ErrInvalidMFACode
	}

	logger.Info("MFA code verified", map[string]interface{}{
		"user_id": userID,
	})

	return nil
}

// GetSecurityAudit gets security audit logs for a user
func (s *UserService) GetSecurityAudit(ctx context.Context, userID uuid.UUID) ([]models.SecurityAuditLog, error) {
	logs, err := s.userRepo.GetSecurityAuditLogs(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get security audit logs")
	}

	return logs, nil
}

// GetSecurityActivity gets recent security activity for a user
func (s *UserService) GetSecurityActivity(ctx context.Context, userID uuid.UUID) ([]models.SecurityActivity, error) {
	activity, err := s.userRepo.GetSecurityActivity(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get security activity")
	}

	return activity, nil
}

// Helper functions for password hashing and verification
func hashPassword(password string) (string, error) {
	// TODO: Implement proper password hashing
	return password, nil
}

func verifyPassword(password, hash string) bool {
	// TODO: Implement proper password verification
	return password == hash
}

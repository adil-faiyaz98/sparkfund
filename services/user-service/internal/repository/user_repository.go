package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sparkfund/services/user-service/internal/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// User operations
	Create(ctx context.Context, user *models.User) error
	Get(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status models.UserStatus) error
	UpdatePassword(ctx context.Context, id uuid.UUID, hashedPassword string) error

	// Profile operations
	GetProfile(ctx context.Context, userID uuid.UUID) (*models.UserProfile, error)
	UpdateProfile(ctx context.Context, profile *models.UserProfile) error
	DeleteProfile(ctx context.Context, userID uuid.UUID) error

	// Password reset operations
	StoreResetToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error
	GetResetToken(ctx context.Context, token string) (*models.PasswordReset, error)
	MarkResetTokenUsed(ctx context.Context, token string) error

	// Session operations
	CreateSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, token string) (*models.Session, error)
	UpdateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, token string) error
	DeleteExpiredSessions(ctx context.Context) error
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]models.Session, error)

	// MFA operations
	StoreMFASecret(ctx context.Context, userID uuid.UUID, secret string) error
	GetMFASecret(ctx context.Context, userID uuid.UUID) (string, error)
	DeleteMFASecret(ctx context.Context, userID uuid.UUID) error

	// Security audit operations
	GetSecurityAuditLogs(ctx context.Context, userID uuid.UUID) ([]models.SecurityAuditLog, error)
	GetSecurityActivity(ctx context.Context, userID uuid.UUID) ([]models.SecurityActivity, error)

	// Login attempt operations
	IncrementFailedAttempts(ctx context.Context, userID uuid.UUID) error
	ResetFailedAttempts(ctx context.Context, userID uuid.UUID) error
	GetFailedAttempts(ctx context.Context, userID uuid.UUID) (int, error)
}

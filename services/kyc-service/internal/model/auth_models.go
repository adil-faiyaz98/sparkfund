package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID            uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email         string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	PasswordHash  string     `json:"-" gorm:"type:varchar(255);not null"`
	FirstName     string     `json:"first_name" gorm:"type:varchar(100)"`
	LastName      string     `json:"last_name" gorm:"type:varchar(100)"`
	Role          string     `json:"role" gorm:"type:varchar(20);not null;default:'user'"`
	MFAEnabled    bool       `json:"mfa_enabled" gorm:"not null;default:false"`
	MFASecret     string     `json:"-" gorm:"type:varchar(32)"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP   string     `json:"last_login_ip,omitempty" gorm:"type:varchar(45)"`
	LoginAttempts int        `json:"-" gorm:"not null;default:0"`
	LockedUntil   *time.Time `json:"-"`
	CreatedAt     time.Time  `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt     time.Time  `json:"updated_at" gorm:"not null;default:now()"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	MFACode  string `json:"mfa_code,omitempty"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         User      `json:"user"`
	MFARequired  bool      `json:"mfa_required,omitempty"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// MFASetupResponse represents an MFA setup response
type MFASetupResponse struct {
	Secret     string `json:"secret"`
	QRCodeURL  string `json:"qr_code_url"`
	RecoveryCodes []string `json:"recovery_codes"`
}

// MFAVerifyRequest represents an MFA verification request
type MFAVerifyRequest struct {
	Code string `json:"code" binding:"required,len=6"`
}

// JWTClaims represents JWT claims
type JWTClaims struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	MFAPassed bool   `json:"mfa_passed"`
}

// Session represents a user session
type Session struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID       uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	RefreshToken string    `json:"-" gorm:"type:varchar(255);not null"`
	UserAgent    string    `json:"user_agent" gorm:"type:varchar(255)"`
	IPAddress    string    `json:"ip_address" gorm:"type:varchar(45)"`
	ExpiresAt    time.Time `json:"expires_at" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"not null;default:now()"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"not null;default:now()"`
}

package users

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "ADMIN"
	UserRoleUser    UserRole = "USER"
	UserRoleManager UserRole = "MANAGER"
	UserRoleSupport UserRole = "SUPPORT"
)

type UserStatus string

const (
	UserStatusActive    UserStatus = "ACTIVE"
	UserStatusInactive  UserStatus = "INACTIVE"
	UserStatusSuspended UserStatus = "SUSPENDED"
	UserStatusDeleted   UserStatus = "DELETED"
)

type User struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key"`
	Email         string     `json:"email" gorm:"type:varchar(255);unique;not null"`
	PasswordHash  string     `json:"-" gorm:"type:varchar(255);not null"`
	FirstName     string     `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName      string     `json:"last_name" gorm:"type:varchar(100);not null"`
	PhoneNumber   string     `json:"phone_number" gorm:"type:varchar(20)"`
	Role          UserRole   `json:"role" gorm:"type:varchar(20);not null;default:'USER'"`
	Status        UserStatus `json:"status" gorm:"type:varchar(20);not null;default:'ACTIVE'"`
	EmailVerified bool       `json:"email_verified" gorm:"default:false"`
	PhoneVerified bool       `json:"phone_verified" gorm:"default:false"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	UpdatePassword(id uuid.UUID, passwordHash string) error
	UpdateStatus(id uuid.UUID, status UserStatus) error
	UpdateLastLogin(id uuid.UUID) error
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	UpdatePassword(ctx context.Context, id uuid.UUID, password string) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status UserStatus) error
	VerifyEmail(ctx context.Context, id uuid.UUID) error
	VerifyPhone(ctx context.Context, id uuid.UUID) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

package service

import (
	"context"
	"fmt"

	"github.com/adil-faiyaz98/structgen/internal/users"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo users.UserRepository
}

func NewUserService(repo users.UserRepository) users.UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, user *users.User) error {
	// Validate email format
	if !isValidEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Check if email already exists
	existing, err := s.repo.GetByEmail(user.Email)
	if err == nil && existing != nil {
		return fmt.Errorf("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.PasswordHash = string(hashedPassword)

	// Set default role if not provided
	if user.Role == "" {
		user.Role = users.UserRoleUser
	}

	// Set default status
	user.Status = users.UserStatusActive

	return s.repo.Create(user)
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*users.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	return s.repo.GetByEmail(email)
}

func (s *userService) UpdateUser(ctx context.Context, user *users.User) error {
	// Validate email format if changed
	if !isValidEmail(user.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Check if email is already taken by another user
	existing, err := s.repo.GetByEmail(user.Email)
	if err == nil && existing != nil && existing.ID != user.ID {
		return fmt.Errorf("email already registered")
	}

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *userService) UpdatePassword(ctx context.Context, id uuid.UUID, password string) error {
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	return s.repo.UpdatePassword(id, string(hashedPassword))
}

func (s *userService) UpdateStatus(ctx context.Context, id uuid.UUID, status users.UserStatus) error {
	if !isValidUserStatus(status) {
		return fmt.Errorf("invalid user status: %s", status)
	}

	return s.repo.UpdateStatus(id, status)
}

func (s *userService) VerifyEmail(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.EmailVerified = true
	return s.repo.Update(user)
}

func (s *userService) VerifyPhone(ctx context.Context, id uuid.UUID) error {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.PhoneVerified = true
	return s.repo.Update(user)
}

func (s *userService) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	return s.repo.UpdateLastLogin(id)
}

func isValidEmail(email string) bool {
	// In a real application, this would use a more robust email validation
	return len(email) > 0 && len(email) <= 255
}

func isValidUserStatus(status users.UserStatus) bool {
	switch status {
	case users.UserStatusActive,
		users.UserStatusInactive,
		users.UserStatusSuspended,
		users.UserStatusDeleted:
		return true
	default:
		return false
	}
}

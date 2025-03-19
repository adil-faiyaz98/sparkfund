package service

import (
	"context"
	"fmt"

	"github.com/adil-faiyaz98/money-pulse/pkg/auth"
	"github.com/adil-faiyaz98/money-pulse/services/users-service/internal/domain"
	"github.com/adil-faiyaz98/money-pulse/services/users-service/internal/repository"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, firstName, lastName string) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error)
	UpdateUserStatus(ctx context.Context, id uuid.UUID, status domain.UserStatus) error
	Authenticate(ctx context.Context, email, password string) (string, error)
}

type userService struct {
	repo repository.UserRepository
	auth auth.TokenManager
}

func NewUserService(repo repository.UserRepository, auth auth.TokenManager) UserService {
	return &userService{
		repo: repo,
		auth: auth,
	}
}

func (s *userService) CreateUser(ctx context.Context, email, password, firstName, lastName string) (*domain.User, error) {
	// Check if user with email already exists
	existing, err := s.repo.GetByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Create new user
	user, err := domain.NewUser(email, password, firstName, lastName)
	if err != nil {
		return nil, err
	}

	// Hash password before storing
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	// Save user
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *userService) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.repo.GetByEmail(ctx, email)
}

func (s *userService) UpdateUser(ctx context.Context, id uuid.UUID, firstName, lastName string) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := user.Update(firstName, lastName); err != nil {
		return err
	}

	return s.repo.Update(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *userService) ListUsers(ctx context.Context, limit, offset int) ([]*domain.User, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *userService) UpdateUserStatus(ctx context.Context, id uuid.UUID, status domain.UserStatus) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if err := user.UpdateStatus(status); err != nil {
		return err
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("invalid credentials")
	}

	if !auth.VerifyPassword(password, user.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	if user.Status != domain.UserStatusActive {
		return "", fmt.Errorf("user account is not active")
	}

	// Generate JWT token
	token, err := s.auth.GenerateToken(user.ID.String())
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

package service

import (
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/sparkfund/auth-service/internal/model"
	"github.com/sparkfund/auth-service/internal/repository"
)

type AuthService interface {
	CreateUser(user *model.User) error
	GetUserByID(id uuid.UUID) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUser(user *model.User) error
	GetAccessTokenSecret() string
	GetRefreshTokenSecret() string
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (s *authService) CreateUser(user *model.User) error {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	return s.userRepo.Create(user)
}

func (s *authService) GetUserByID(id uuid.UUID) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *authService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *authService) UpdateUser(user *model.User) error {
	return s.userRepo.Update(user)
}

func (s *authService) GetAccessTokenSecret() string {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "your-access-secret-key" // Fallback for development
	}
	return secret
}

func (s *authService) GetRefreshTokenSecret() string {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "your-refresh-secret-key" // Fallback for development
	}
	return secret
}

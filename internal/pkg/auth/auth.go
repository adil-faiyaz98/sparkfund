package auth

import (
	"context"
	"errors"
	"time"

	"github.com/your-username/money-pulse/internal/pkg/models"
)

// Service provides authentication operations
type Service struct {
	userRepo   UserRepository
	tokenMaker TokenMaker
}

// UserRepository defines methods to interact with user storage
type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) error
}

// NewService creates a new auth service
func NewService(userRepo UserRepository, tokenMaker TokenMaker) *Service {
	return &Service{
		userRepo:   userRepo,
		tokenMaker: tokenMaker,
	}
}

// LoginResponse contains the login result
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        *models.User `json:"user"`
}

// Login authenticates a user and returns a token
func (s *Service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("invalid credentials")
	}

	if err := ComparePasswords(user.PasswordHash, password); err != nil {
		return nil, errors.New("invalid credentials")
	}

	accessToken, err := s.tokenMaker.CreateToken(user.ID, time.Hour*24)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: accessToken,
		User:        user,
	}, nil
}

// RegisterRequest contains the registration details
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

// Register creates a new user
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*models.User, error) {
	existingUser, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

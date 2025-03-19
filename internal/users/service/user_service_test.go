package service

import (
	"context"
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/adil-faiyaz98/money-pulse/internal/users"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	testutil.MockRepository
}

func (m *MockUserRepository) Create(user *users.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*users.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*users.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*users.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *users.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	tests := []struct {
		name        string
		user        *users.User
		setupMock   func()
		wantErr     bool
		errContains string
	}{
		{
			name: "successful user creation",
			user: &users.User{
				Email:     "test@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password123",
			},
			setupMock: func() {
				mockRepo.On("Create", mock.AnythingOfType("*users.User")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			user: &users.User{
				Email:     "invalid-email",
				FirstName: "John",
				LastName:  "Doe",
				Password:  "password123",
			},
			setupMock:   func() {},
			wantErr:     true,
			errContains: "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.CreateUser(context.Background(), tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, tt.user.ID)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	expectedUser := &users.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func()
		wantUser  *users.User
		wantErr   bool
	}{
		{
			name: "successful user retrieval",
			id:   userID,
			setupMock: func() {
				mockRepo.On("GetByID", userID).Return(expectedUser, nil)
			},
			wantUser: expectedUser,
			wantErr:  false,
		},
		{
			name: "user not found",
			id:   uuid.New(),
			setupMock: func() {
				mockRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			user, err := service.GetUser(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantUser, user)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)

	userID := uuid.New()
	existingUser := &users.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	updatedUser := &users.User{
		ID:        userID,
		Email:     "test@example.com",
		FirstName: "Jane",
		LastName:  "Doe",
	}

	tests := []struct {
		name        string
		user        *users.User
		setupMock   func()
		wantErr     bool
		errContains string
	}{
		{
			name: "successful user update",
			user: updatedUser,
			setupMock: func() {
				mockRepo.On("GetByID", userID).Return(existingUser, nil)
				mockRepo.On("Update", updatedUser).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "user not found",
			user: &users.User{
				ID:    uuid.New(),
				Email: "test@example.com",
			},
			setupMock: func() {
				mockRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)
			},
			wantErr:     true,
			errContains: "user not found",
		},
		{
			name: "invalid email",
			user: &users.User{
				ID:        userID,
				Email:     "invalid-email",
				FirstName: "Jane",
				LastName:  "Doe",
			},
			setupMock: func() {
				mockRepo.On("GetByID", userID).Return(existingUser, nil)
			},
			wantErr:     true,
			errContains: "invalid email",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.UpdateUser(context.Background(), tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			assert.NoError(t, err)
			mockRepo.AssertExpectations(t)
		})
	}
}

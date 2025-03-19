package service

import (
	"context"
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/accounts"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAccountRepository is a mock implementation of AccountRepository
type MockAccountRepository struct {
	mock.Mock
}

func (m *MockAccountRepository) Create(account *accounts.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) GetByID(id uuid.UUID) (*accounts.Account, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accounts.Account), args.Error(1)
}

func (m *MockAccountRepository) GetByUserID(userID uuid.UUID) ([]*accounts.Account, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*accounts.Account), args.Error(1)
}

func (m *MockAccountRepository) Update(account *accounts.Account) error {
	args := m.Called(account)
	return args.Error(0)
}

func (m *MockAccountRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockAccountRepository) GetByAccountNumber(accountNumber string) (*accounts.Account, error) {
	args := m.Called(accountNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*accounts.Account), args.Error(1)
}

func TestAccountService_CreateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	service := NewAccountService(mockRepo)

	tests := []struct {
		name        string
		account     *accounts.Account
		setupMock   func()
		wantErr     bool
		errContains string
	}{
		{
			name: "successful account creation",
			account: &accounts.Account{
				UserID:  uuid.New(),
				Name:    "Test Account",
				Type:    accounts.AccountTypeSavings,
				Balance: 1000.0,
			},
			setupMock: func() {
				mockRepo.On("Create", mock.AnythingOfType("*accounts.Account")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid account type",
			account: &accounts.Account{
				UserID:  uuid.New(),
				Name:    "Test Account",
				Type:    "invalid_type",
				Balance: 1000.0,
			},
			setupMock:   func() {},
			wantErr:     true,
			errContains: "invalid account type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.CreateAccount(context.Background(), tt.account)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, tt.account.AccountNumber)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	service := NewAccountService(mockRepo)

	accountID := uuid.New()
	expectedAccount := &accounts.Account{
		ID:            accountID,
		UserID:        uuid.New(),
		Name:          "Test Account",
		Type:          accounts.AccountTypeSavings,
		Balance:       1000.0,
		AccountNumber: "ACC123456",
	}

	tests := []struct {
		name        string
		id          uuid.UUID
		setupMock   func()
		wantAccount *accounts.Account
		wantErr     bool
	}{
		{
			name: "successful account retrieval",
			id:   accountID,
			setupMock: func() {
				mockRepo.On("GetByID", accountID).Return(expectedAccount, nil)
			},
			wantAccount: expectedAccount,
			wantErr:     false,
		},
		{
			name: "account not found",
			id:   uuid.New(),
			setupMock: func() {
				mockRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)
			},
			wantAccount: nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			account, err := service.GetAccount(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, account)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantAccount, account)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAccountService_UpdateAccount(t *testing.T) {
	mockRepo := new(MockAccountRepository)
	service := NewAccountService(mockRepo)

	accountID := uuid.New()
	existingAccount := &accounts.Account{
		ID:     accountID,
		UserID: uuid.New(),
		Name:   "Test Account",
		Type:   accounts.AccountTypeSavings,
	}

	updatedAccount := &accounts.Account{
		ID:     accountID,
		UserID: existingAccount.UserID,
		Name:   "Updated Account",
		Type:   accounts.AccountTypeChecking,
	}

	tests := []struct {
		name        string
		account     *accounts.Account
		setupMock   func()
		wantErr     bool
		errContains string
	}{
		{
			name:    "successful account update",
			account: updatedAccount,
			setupMock: func() {
				mockRepo.On("GetByID", accountID).Return(existingAccount, nil)
				mockRepo.On("Update", updatedAccount).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "account not found",
			account: &accounts.Account{
				ID:   uuid.New(),
				Type: accounts.AccountTypeSavings,
			},
			setupMock: func() {
				mockRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)
			},
			wantErr:     true,
			errContains: "account not found",
		},
		{
			name: "invalid account type",
			account: &accounts.Account{
				ID:     accountID,
				UserID: existingAccount.UserID,
				Name:   "Invalid Account",
				Type:   "invalid_type",
			},
			setupMock: func() {
				mockRepo.On("GetByID", accountID).Return(existingAccount, nil)
			},
			wantErr:     true,
			errContains: "invalid account type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.UpdateAccount(context.Background(), tt.account)
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

package service

import (
	"context"
	"testing"
	"time"

	"github.com/adil-faiyaz98/money-pulse/internal/loans"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLoanRepository is a mock implementation of LoanRepository
type MockLoanRepository struct {
	testutil.MockRepository
}

func (m *MockLoanRepository) Create(loan *loans.Loan) error {
	args := m.Called(loan)
	return args.Error(0)
}

func (m *MockLoanRepository) GetByID(id uuid.UUID) (*loans.Loan, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*loans.Loan), args.Error(1)
}

func (m *MockLoanRepository) GetByUserID(userID uuid.UUID) ([]*loans.Loan, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*loans.Loan), args.Error(1)
}

func (m *MockLoanRepository) Update(loan *loans.Loan) error {
	args := m.Called(loan)
	return args.Error(0)
}

func (m *MockLoanRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestLoanService_CreateLoan(t *testing.T) {
	mockRepo := new(MockLoanRepository)
	service := NewLoanService(mockRepo)

	tests := []struct {
		name        string
		loan        *loans.Loan
		setupMock   func()
		wantErr     bool
		errContains string
	}{
		{
			name: "successful loan creation",
			loan: &loans.Loan{
				UserID:       uuid.New(),
				Type:         loans.LoanTypePersonal,
				Amount:       10000.0,
				Term:         12,
				InterestRate: 5.5,
				Purpose:      "Home renovation",
			},
			setupMock: func() {
				mockRepo.On("Create", mock.AnythingOfType("*loans.Loan")).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "invalid loan type",
			loan: &loans.Loan{
				UserID:       uuid.New(),
				Type:         "invalid_type",
				Amount:       10000.0,
				Term:         12,
				InterestRate: 5.5,
				Purpose:      "Home renovation",
			},
			setupMock:   func() {},
			wantErr:     true,
			errContains: "invalid loan type",
		},
		{
			name: "negative amount",
			loan: &loans.Loan{
				UserID:       uuid.New(),
				Type:         loans.LoanTypePersonal,
				Amount:       -1000.0,
				Term:         12,
				InterestRate: 5.5,
				Purpose:      "Home renovation",
			},
			setupMock:   func() {},
			wantErr:     true,
			errContains: "amount must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.CreateLoan(context.Background(), tt.loan)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, tt.loan.ID)
			assert.Equal(t, loans.LoanStatusPending, tt.loan.Status)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLoanService_GetLoan(t *testing.T) {
	mockRepo := new(MockLoanRepository)
	service := NewLoanService(mockRepo)

	loanID := uuid.New()
	expectedLoan := &loans.Loan{
		ID:           loanID,
		UserID:       uuid.New(),
		Type:         loans.LoanTypePersonal,
		Amount:       10000.0,
		Term:         12,
		InterestRate: 5.5,
		Purpose:      "Home renovation",
		Status:       loans.LoanStatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name      string
		id        uuid.UUID
		setupMock func()
		wantLoan  *loans.Loan
		wantErr   bool
	}{
		{
			name: "successful loan retrieval",
			id:   loanID,
			setupMock: func() {
				mockRepo.On("GetByID", loanID).Return(expectedLoan, nil)
			},
			wantLoan: expectedLoan,
			wantErr:  false,
		},
		{
			name: "loan not found",
			id:   uuid.New(),
			setupMock: func() {
				mockRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)
			},
			wantLoan: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			loan, err := service.GetLoan(context.Background(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, loan)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantLoan, loan)
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLoanService_UpdateLoanStatus(t *testing.T) {
	mockRepo := new(MockLoanRepository)
	service := NewLoanService(mockRepo)

	loanID := uuid.New()
	existingLoan := &loans.Loan{
		ID:           loanID,
		UserID:       uuid.New(),
		Type:         loans.LoanTypePersonal,
		Amount:       10000.0,
		Term:         12,
		InterestRate: 5.5,
		Purpose:      "Home renovation",
		Status:       loans.LoanStatusPending,
	}

	updatedLoan := &loans.Loan{
		ID:           loanID,
		UserID:       existingLoan.UserID,
		Type:         existingLoan.Type,
		Amount:       existingLoan.Amount,
		Term:         existingLoan.Term,
		InterestRate: existingLoan.InterestRate,
		Purpose:      existingLoan.Purpose,
		Status:       loans.LoanStatusActive,
	}

	tests := []struct {
		name        string
		loan        *loans.Loan
		setupMock   func()
		wantErr     bool
		errContains string
	}{
		{
			name: "successful status update",
			loan: updatedLoan,
			setupMock: func() {
				mockRepo.On("GetByID", loanID).Return(existingLoan, nil)
				mockRepo.On("Update", updatedLoan).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "loan not found",
			loan: &loans.Loan{
				ID:     uuid.New(),
				Status: loans.LoanStatusActive,
			},
			setupMock: func() {
				mockRepo.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)
			},
			wantErr:     true,
			errContains: "loan not found",
		},
		{
			name: "invalid status transition",
			loan: &loans.Loan{
				ID:           loanID,
				UserID:       existingLoan.UserID,
				Type:         existingLoan.Type,
				Amount:       existingLoan.Amount,
				Term:         existingLoan.Term,
				InterestRate: existingLoan.InterestRate,
				Purpose:      existingLoan.Purpose,
				Status:       loans.LoanStatusRejected,
			},
			setupMock: func() {
				mockRepo.On("GetByID", loanID).Return(existingLoan, nil)
			},
			wantErr:     true,
			errContains: "invalid status transition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.UpdateLoanStatus(context.Background(), tt.loan)
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

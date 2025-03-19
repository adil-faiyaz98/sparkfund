package service

import (
	"context"
	"testing"
	"time"

	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/adil-faiyaz98/money-pulse/internal/transactions"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockTransactionRepository struct {
	testutil.MockRepository
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *transactions.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*transactions.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*transactions.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*transactions.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*transactions.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]*transactions.Transaction, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*transactions.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *transactions.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTransactionService_CreateTransaction(t *testing.T) {
	// Setup
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)
	ctx := context.Background()

	// Test data
	userID := uuid.New()
	accountID := uuid.New()
	testTransaction := &transactions.Transaction{
		UserID:      userID,
		AccountID:   accountID,
		Type:        transactions.TransactionTypeDeposit,
		Amount:      1000.0,
		Description: "Initial deposit",
		Date:        time.Now(),
	}

	t.Run("Create Valid Transaction", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("Create", ctx, mock.AnythingOfType("*transactions.Transaction")).Return(nil)

		// Execute
		err := service.CreateTransaction(ctx, testTransaction)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, testTransaction.ID)
		assert.Equal(t, transactions.TransactionStatusCompleted, testTransaction.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Transaction with Invalid Type", func(t *testing.T) {
		// Setup test data
		invalidTransaction := *testTransaction
		invalidTransaction.Type = "invalid_type"

		// Execute
		err := service.CreateTransaction(ctx, &invalidTransaction)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, invalidTransaction.ID)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Create Transaction with Negative Amount", func(t *testing.T) {
		// Setup test data
		invalidTransaction := *testTransaction
		invalidTransaction.Amount = -100.0

		// Execute
		err := service.CreateTransaction(ctx, &invalidTransaction)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, invalidTransaction.ID)
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestTransactionService_GetTransaction(t *testing.T) {
	// Setup
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)
	ctx := context.Background()

	// Test data
	transactionID := uuid.New()
	testTransaction := &transactions.Transaction{
		ID:          transactionID,
		UserID:      uuid.New(),
		AccountID:   uuid.New(),
		Type:        transactions.TransactionTypeDeposit,
		Amount:      1000.0,
		Description: "Initial deposit",
		Date:        time.Now(),
		Status:      transactions.TransactionStatusCompleted,
	}

	t.Run("Get Existing Transaction", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, transactionID).Return(testTransaction, nil)

		// Execute
		transaction, err := service.GetTransaction(ctx, transactionID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, testTransaction.ID, transaction.ID)
		assert.Equal(t, testTransaction.UserID, transaction.UserID)
		assert.Equal(t, testTransaction.AccountID, transaction.AccountID)
		assert.Equal(t, testTransaction.Type, transaction.Type)
		assert.Equal(t, testTransaction.Amount, transaction.Amount)
		assert.Equal(t, testTransaction.Description, transaction.Description)
		assert.Equal(t, testTransaction.Date, transaction.Date)
		assert.Equal(t, testTransaction.Status, transaction.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get Non-Existing Transaction", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, transactionID).Return(nil, transactions.ErrTransactionNotFound)

		// Execute
		transaction, err := service.GetTransaction(ctx, transactionID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, transaction)
		assert.Equal(t, transactions.ErrTransactionNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestTransactionService_UpdateTransactionStatus(t *testing.T) {
	// Setup
	mockRepo := new(MockTransactionRepository)
	service := NewTransactionService(mockRepo)
	ctx := context.Background()

	// Test data
	transactionID := uuid.New()
	testTransaction := &transactions.Transaction{
		ID:          transactionID,
		UserID:      uuid.New(),
		AccountID:   uuid.New(),
		Type:        transactions.TransactionTypeDeposit,
		Amount:      1000.0,
		Description: "Initial deposit",
		Date:        time.Now(),
		Status:      transactions.TransactionStatusPending,
	}

	t.Run("Update Transaction Status Successfully", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, transactionID).Return(testTransaction, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*transactions.Transaction")).Return(nil)

		// Execute
		err := service.UpdateTransactionStatus(ctx, transactionID, transactions.TransactionStatusCompleted, "Transaction completed successfully")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, transactions.TransactionStatusCompleted, testTransaction.Status)
		assert.Equal(t, "Transaction completed successfully", testTransaction.StatusReason)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update Non-Existing Transaction Status", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, transactionID).Return(nil, transactions.ErrTransactionNotFound)

		// Execute
		err := service.UpdateTransactionStatus(ctx, transactionID, transactions.TransactionStatusCompleted, "Transaction completed successfully")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, transactions.ErrTransactionNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update Transaction Status with Invalid Status", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, transactionID).Return(testTransaction, nil)

		// Execute
		err := service.UpdateTransactionStatus(ctx, transactionID, "invalid_status", "Invalid status")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, transactions.TransactionStatusPending, testTransaction.Status)
		mockRepo.AssertNotCalled(t, "Update")
	})
}

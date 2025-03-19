package service

import (
	"testing"
	"time"

	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/adil-faiyaz98/money-pulse/internal/transactions"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup test database
	testDB := testutil.NewTestDB(t)
	defer testDB.Close(t)

	// Create repository
	repo := NewPostgresTransactionRepository(testDB.DB)
	service := NewTransactionService(repo)

	// Create test context
	ctx := testutil.CreateTestContext(t)

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

	t.Run("Create and Retrieve Transaction", func(t *testing.T) {
		// Create transaction
		err := service.CreateTransaction(ctx, testTransaction)
		require.NoError(t, err)
		require.NotEmpty(t, testTransaction.ID)
		require.Equal(t, transactions.TransactionStatusCompleted, testTransaction.Status)

		// Retrieve transaction
		transaction, err := service.GetTransaction(ctx, testTransaction.ID)
		require.NoError(t, err)
		assert.Equal(t, testTransaction.ID, transaction.ID)
		assert.Equal(t, testTransaction.UserID, transaction.UserID)
		assert.Equal(t, testTransaction.AccountID, transaction.AccountID)
		assert.Equal(t, testTransaction.Type, transaction.Type)
		assert.Equal(t, testTransaction.Amount, transaction.Amount)
		assert.Equal(t, testTransaction.Description, transaction.Description)
		assert.Equal(t, testTransaction.Date, transaction.Date)
		assert.Equal(t, testTransaction.Status, transaction.Status)
	})

	t.Run("Get User Transactions", func(t *testing.T) {
		// Create another transaction for the same user
		anotherTransaction := &transactions.Transaction{
			UserID:      userID,
			AccountID:   accountID,
			Type:        transactions.TransactionTypeWithdrawal,
			Amount:      500.0,
			Description: "ATM withdrawal",
			Date:        time.Now(),
		}
		err := service.CreateTransaction(ctx, anotherTransaction)
		require.NoError(t, err)

		// Get all transactions for the user
		transactions, err := service.GetUserTransactions(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, transactions, 2)

		// Verify transaction details
		transactionMap := make(map[uuid.UUID]*transactions.Transaction)
		for _, transaction := range transactions {
			transactionMap[transaction.ID] = transaction
		}

		assert.Contains(t, transactionMap, testTransaction.ID)
		assert.Contains(t, transactionMap, anotherTransaction.ID)
	})

	t.Run("Get Account Transactions", func(t *testing.T) {
		// Get all transactions for the account
		transactions, err := service.GetAccountTransactions(ctx, accountID)
		require.NoError(t, err)
		assert.Len(t, transactions, 2)

		// Verify transaction details
		transactionMap := make(map[uuid.UUID]*transactions.Transaction)
		for _, transaction := range transactions {
			transactionMap[transaction.ID] = transaction
		}

		assert.Contains(t, transactionMap, testTransaction.ID)
	})

	t.Run("Update Transaction Status", func(t *testing.T) {
		// Update transaction status
		err := service.UpdateTransactionStatus(ctx, testTransaction.ID, transactions.TransactionStatusFailed, "Insufficient funds")
		require.NoError(t, err)

		// Verify update
		updated, err := service.GetTransaction(ctx, testTransaction.ID)
		require.NoError(t, err)
		assert.Equal(t, transactions.TransactionStatusFailed, updated.Status)
		assert.Equal(t, "Insufficient funds", updated.StatusReason)
	})

	t.Run("Delete Transaction", func(t *testing.T) {
		// Delete transaction
		err := service.DeleteTransaction(ctx, testTransaction.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetTransaction(ctx, testTransaction.ID)
		assert.Error(t, err)
	})
}

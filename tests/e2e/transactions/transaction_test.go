package transactions

import (
	"testing"
	"time"

	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/adil-faiyaz98/money-pulse/internal/transactions"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransactionAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

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

	t.Run("Create Transaction", func(t *testing.T) {
		// Send create transaction request
		var response transactions.Transaction
		err := server.SendRequest("POST", "/api/v1/transactions", testTransaction, &response)
		require.NoError(t, err)
		require.Equal(t, 201, response.StatusCode)

		// Verify response
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testTransaction.UserID, response.UserID)
		assert.Equal(t, testTransaction.AccountID, response.AccountID)
		assert.Equal(t, testTransaction.Type, response.Type)
		assert.Equal(t, testTransaction.Amount, response.Amount)
		assert.Equal(t, testTransaction.Description, response.Description)
		assert.Equal(t, testTransaction.Date, response.Date)
		assert.Equal(t, transactions.TransactionStatusCompleted, response.Status)
	})

	t.Run("Get Transaction", func(t *testing.T) {
		// First create a transaction
		var createResponse transactions.Transaction
		err := server.SendRequest("POST", "/api/v1/transactions", testTransaction, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Get the transaction
		var getResponse transactions.Transaction
		err = server.SendRequest("GET", "/api/v1/transactions/"+createResponse.ID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 200, getResponse.StatusCode)

		// Verify response
		assert.Equal(t, createResponse.ID, getResponse.ID)
		assert.Equal(t, createResponse.UserID, getResponse.UserID)
		assert.Equal(t, createResponse.AccountID, getResponse.AccountID)
		assert.Equal(t, createResponse.Type, getResponse.Type)
		assert.Equal(t, createResponse.Amount, getResponse.Amount)
		assert.Equal(t, createResponse.Description, getResponse.Description)
		assert.Equal(t, createResponse.Date, getResponse.Date)
		assert.Equal(t, createResponse.Status, getResponse.Status)
	})

	t.Run("Get User Transactions", func(t *testing.T) {
		// First create a transaction
		var createResponse transactions.Transaction
		err := server.SendRequest("POST", "/api/v1/transactions", testTransaction, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Get all transactions for the user
		var getResponse []transactions.Transaction
		err = server.SendRequest("GET", "/api/v1/transactions/user/"+userID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 200, getResponse[0].StatusCode)

		// Verify response
		assert.NotEmpty(t, getResponse)
		transactionMap := make(map[uuid.UUID]*transactions.Transaction)
		for _, transaction := range getResponse {
			transactionMap[transaction.ID] = &transaction
		}
		assert.Contains(t, transactionMap, createResponse.ID)
	})

	t.Run("Get Account Transactions", func(t *testing.T) {
		// First create a transaction
		var createResponse transactions.Transaction
		err := server.SendRequest("POST", "/api/v1/transactions", testTransaction, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Get all transactions for the account
		var getResponse []transactions.Transaction
		err = server.SendRequest("GET", "/api/v1/transactions/account/"+accountID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 200, getResponse[0].StatusCode)

		// Verify response
		assert.NotEmpty(t, getResponse)
		transactionMap := make(map[uuid.UUID]*transactions.Transaction)
		for _, transaction := range getResponse {
			transactionMap[transaction.ID] = &transaction
		}
		assert.Contains(t, transactionMap, createResponse.ID)
	})

	t.Run("Update Transaction Status", func(t *testing.T) {
		// First create a transaction
		var createResponse transactions.Transaction
		err := server.SendRequest("POST", "/api/v1/transactions", testTransaction, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Update transaction status
		updateRequest := struct {
			Status       string `json:"status"`
			StatusReason string `json:"status_reason"`
		}{
			Status:       string(transactions.TransactionStatusFailed),
			StatusReason: "Insufficient funds",
		}

		var updateResponse transactions.Transaction
		err = server.SendRequest("PATCH", "/api/v1/transactions/"+createResponse.ID.String()+"/status", updateRequest, &updateResponse)
		require.NoError(t, err)
		require.Equal(t, 200, updateResponse.StatusCode)

		// Verify response
		assert.Equal(t, transactions.TransactionStatusFailed, updateResponse.Status)
		assert.Equal(t, "Insufficient funds", updateResponse.StatusReason)
	})

	t.Run("Delete Transaction", func(t *testing.T) {
		// First create a transaction
		var createResponse transactions.Transaction
		err := server.SendRequest("POST", "/api/v1/transactions", testTransaction, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Delete the transaction
		var deleteResponse struct{}
		err = server.SendRequest("DELETE", "/api/v1/transactions/"+createResponse.ID.String(), nil, &deleteResponse)
		require.NoError(t, err)
		require.Equal(t, 204, deleteResponse.StatusCode)

		// Try to get the deleted transaction
		var getResponse transactions.Transaction
		err = server.SendRequest("GET", "/api/v1/transactions/"+createResponse.ID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 404, getResponse.StatusCode)
	})
}

package loans

import (
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/loans"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoanAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

	// Test data
	userID := uuid.New()
	testLoan := &loans.Loan{
		UserID:       userID,
		Type:         loans.LoanTypePersonal,
		Amount:       10000.0,
		Term:         12,
		InterestRate: 5.5,
		Purpose:      "Home renovation",
	}

	t.Run("Create Loan", func(t *testing.T) {
		// Send create loan request
		var response loans.Loan
		err := server.SendRequest("POST", "/api/v1/loans", testLoan, &response)
		require.NoError(t, err)
		require.Equal(t, 201, response.StatusCode)

		// Verify response
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testLoan.UserID, response.UserID)
		assert.Equal(t, testLoan.Type, response.Type)
		assert.Equal(t, testLoan.Amount, response.Amount)
		assert.Equal(t, testLoan.Term, response.Term)
		assert.Equal(t, testLoan.InterestRate, response.InterestRate)
		assert.Equal(t, testLoan.Purpose, response.Purpose)
		assert.Equal(t, loans.LoanStatusPending, response.Status)
	})

	t.Run("Get Loan", func(t *testing.T) {
		// First create a loan
		var createResponse loans.Loan
		err := server.SendRequest("POST", "/api/v1/loans", testLoan, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Get the loan
		var getResponse loans.Loan
		err = server.SendRequest("GET", "/api/v1/loans/"+createResponse.ID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 200, getResponse.StatusCode)

		// Verify response
		assert.Equal(t, createResponse.ID, getResponse.ID)
		assert.Equal(t, createResponse.UserID, getResponse.UserID)
		assert.Equal(t, createResponse.Type, getResponse.Type)
		assert.Equal(t, createResponse.Amount, getResponse.Amount)
		assert.Equal(t, createResponse.Term, getResponse.Term)
		assert.Equal(t, createResponse.InterestRate, getResponse.InterestRate)
		assert.Equal(t, createResponse.Purpose, getResponse.Purpose)
		assert.Equal(t, createResponse.Status, getResponse.Status)
	})

	t.Run("Get User Loans", func(t *testing.T) {
		// First create a loan
		var createResponse loans.Loan
		err := server.SendRequest("POST", "/api/v1/loans", testLoan, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Get all loans for the user
		var getResponse []loans.Loan
		err = server.SendRequest("GET", "/api/v1/loans/user/"+userID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 200, getResponse[0].StatusCode)

		// Verify response
		assert.NotEmpty(t, getResponse)
		loanMap := make(map[uuid.UUID]*loans.Loan)
		for _, loan := range getResponse {
			loanMap[loan.ID] = &loan
		}
		assert.Contains(t, loanMap, createResponse.ID)
	})

	t.Run("Update Loan Status", func(t *testing.T) {
		// First create a loan
		var createResponse loans.Loan
		err := server.SendRequest("POST", "/api/v1/loans", testLoan, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Update loan status
		updateRequest := struct {
			Status string `json:"status"`
		}{
			Status: string(loans.LoanStatusActive),
		}

		var updateResponse loans.Loan
		err = server.SendRequest("PATCH", "/api/v1/loans/"+createResponse.ID.String()+"/status", updateRequest, &updateResponse)
		require.NoError(t, err)
		require.Equal(t, 200, updateResponse.StatusCode)

		// Verify response
		assert.Equal(t, loans.LoanStatusActive, updateResponse.Status)
	})

	t.Run("Delete Loan", func(t *testing.T) {
		// First create a loan
		var createResponse loans.Loan
		err := server.SendRequest("POST", "/api/v1/loans", testLoan, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Delete the loan
		var deleteResponse struct{}
		err = server.SendRequest("DELETE", "/api/v1/loans/"+createResponse.ID.String(), nil, &deleteResponse)
		require.NoError(t, err)
		require.Equal(t, 204, deleteResponse.StatusCode)

		// Try to get the deleted loan
		var getResponse loans.Loan
		err = server.SendRequest("GET", "/api/v1/loans/"+createResponse.ID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 404, getResponse.StatusCode)
	})
}

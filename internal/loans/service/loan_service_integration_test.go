package service

import (
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/loans"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoanService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup test database
	testDB := testutil.NewTestDB(t)
	defer testDB.Close(t)

	// Create repository
	repo := NewPostgresLoanRepository(testDB.DB)
	service := NewLoanService(repo)

	// Create test context
	ctx := testutil.CreateTestContext(t)

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

	t.Run("Create and Retrieve Loan", func(t *testing.T) {
		// Create loan
		err := service.CreateLoan(ctx, testLoan)
		require.NoError(t, err)
		require.NotEmpty(t, testLoan.ID)
		require.Equal(t, loans.LoanStatusPending, testLoan.Status)

		// Retrieve loan
		loan, err := service.GetLoan(ctx, testLoan.ID)
		require.NoError(t, err)
		assert.Equal(t, testLoan.ID, loan.ID)
		assert.Equal(t, testLoan.UserID, loan.UserID)
		assert.Equal(t, testLoan.Type, loan.Type)
		assert.Equal(t, testLoan.Amount, loan.Amount)
		assert.Equal(t, testLoan.Term, loan.Term)
		assert.Equal(t, testLoan.InterestRate, loan.InterestRate)
		assert.Equal(t, testLoan.Purpose, loan.Purpose)
		assert.Equal(t, testLoan.Status, loan.Status)
	})

	t.Run("Get User Loans", func(t *testing.T) {
		// Create another loan for the same user
		anotherLoan := &loans.Loan{
			UserID:       userID,
			Type:         loans.LoanTypeBusiness,
			Amount:       20000.0,
			Term:         24,
			InterestRate: 6.0,
			Purpose:      "Business expansion",
		}
		err := service.CreateLoan(ctx, anotherLoan)
		require.NoError(t, err)

		// Get all loans for the user
		loans, err := service.GetUserLoans(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, loans, 2)

		// Verify loan details
		loanMap := make(map[uuid.UUID]*loans.Loan)
		for _, loan := range loans {
			loanMap[loan.ID] = loan
		}

		assert.Contains(t, loanMap, testLoan.ID)
		assert.Contains(t, loanMap, anotherLoan.ID)
	})

	t.Run("Update Loan Status", func(t *testing.T) {
		// Update loan status
		testLoan.Status = loans.LoanStatusActive
		err := service.UpdateLoanStatus(ctx, testLoan)
		require.NoError(t, err)

		// Verify update
		updated, err := service.GetLoan(ctx, testLoan.ID)
		require.NoError(t, err)
		assert.Equal(t, loans.LoanStatusActive, updated.Status)
	})

	t.Run("Delete Loan", func(t *testing.T) {
		// Delete loan
		err := service.DeleteLoan(ctx, testLoan.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetLoan(ctx, testLoan.ID)
		assert.Error(t, err)
	})
}

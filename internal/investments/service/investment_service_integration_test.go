package service

import (
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/investments"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvestmentService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup test database
	testDB := testutil.NewTestDB(t)
	defer testDB.Close(t)

	// Create repository
	repo := NewPostgresInvestmentRepository(testDB.DB)
	service := NewInvestmentService(repo)

	// Create test context
	ctx := testutil.CreateTestContext(t)

	// Test data
	userID := uuid.New()
	testInvestment := &investments.Investment{
		UserID:      userID,
		Type:        investments.InvestmentTypeStocks,
		Amount:      5000.0,
		RiskLevel:   investments.RiskLevelMedium,
		Duration:    36,
		Description: "Stock portfolio investment",
	}

	t.Run("Create and Retrieve Investment", func(t *testing.T) {
		// Create investment
		err := service.CreateInvestment(ctx, testInvestment)
		require.NoError(t, err)
		require.NotEmpty(t, testInvestment.ID)
		require.Equal(t, investments.InvestmentStatusActive, testInvestment.Status)

		// Retrieve investment
		investment, err := service.GetInvestment(ctx, testInvestment.ID)
		require.NoError(t, err)
		assert.Equal(t, testInvestment.ID, investment.ID)
		assert.Equal(t, testInvestment.UserID, investment.UserID)
		assert.Equal(t, testInvestment.Type, investment.Type)
		assert.Equal(t, testInvestment.Amount, investment.Amount)
		assert.Equal(t, testInvestment.RiskLevel, investment.RiskLevel)
		assert.Equal(t, testInvestment.Duration, investment.Duration)
		assert.Equal(t, testInvestment.Description, investment.Description)
		assert.Equal(t, testInvestment.Status, investment.Status)
	})

	t.Run("Get User Investments", func(t *testing.T) {
		// Create another investment for the same user
		anotherInvestment := &investments.Investment{
			UserID:      userID,
			Type:        investments.InvestmentTypeBonds,
			Amount:      3000.0,
			RiskLevel:   investments.RiskLevelLow,
			Duration:    24,
			Description: "Government bonds investment",
		}
		err := service.CreateInvestment(ctx, anotherInvestment)
		require.NoError(t, err)

		// Get all investments for the user
		investments, err := service.GetUserInvestments(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, investments, 2)

		// Verify investment details
		investmentMap := make(map[uuid.UUID]*investments.Investment)
		for _, investment := range investments {
			investmentMap[investment.ID] = investment
		}

		assert.Contains(t, investmentMap, testInvestment.ID)
		assert.Contains(t, investmentMap, anotherInvestment.ID)
	})

	t.Run("Update Investment Status", func(t *testing.T) {
		// Update investment status
		err := service.UpdateInvestmentStatus(ctx, testInvestment.ID, investments.InvestmentStatusCompleted, "Investment completed successfully")
		require.NoError(t, err)

		// Verify update
		updated, err := service.GetInvestment(ctx, testInvestment.ID)
		require.NoError(t, err)
		assert.Equal(t, investments.InvestmentStatusCompleted, updated.Status)
		assert.Equal(t, "Investment completed successfully", updated.StatusReason)
	})

	t.Run("Delete Investment", func(t *testing.T) {
		// Delete investment
		err := service.DeleteInvestment(ctx, testInvestment.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetInvestment(ctx, testInvestment.ID)
		assert.Error(t, err)
	})
}

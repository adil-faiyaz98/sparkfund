package investments

import (
	"testing"

	"github.com/adil-faiyaz98/money-pulse/internal/investments"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInvestmentAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

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

	t.Run("Create Investment", func(t *testing.T) {
		// Send create investment request
		var response investments.Investment
		err := server.SendRequest("POST", "/api/v1/investments", testInvestment, &response)
		require.NoError(t, err)
		require.Equal(t, 201, response.StatusCode)

		// Verify response
		assert.NotEmpty(t, response.ID)
		assert.Equal(t, testInvestment.UserID, response.UserID)
		assert.Equal(t, testInvestment.Type, response.Type)
		assert.Equal(t, testInvestment.Amount, response.Amount)
		assert.Equal(t, testInvestment.RiskLevel, response.RiskLevel)
		assert.Equal(t, testInvestment.Duration, response.Duration)
		assert.Equal(t, testInvestment.Description, response.Description)
		assert.Equal(t, investments.InvestmentStatusActive, response.Status)
	})

	t.Run("Get Investment", func(t *testing.T) {
		// First create an investment
		var createResponse investments.Investment
		err := server.SendRequest("POST", "/api/v1/investments", testInvestment, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Get the investment
		var getResponse investments.Investment
		err = server.SendRequest("GET", "/api/v1/investments/"+createResponse.ID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 200, getResponse.StatusCode)

		// Verify response
		assert.Equal(t, createResponse.ID, getResponse.ID)
		assert.Equal(t, createResponse.UserID, getResponse.UserID)
		assert.Equal(t, createResponse.Type, getResponse.Type)
		assert.Equal(t, createResponse.Amount, getResponse.Amount)
		assert.Equal(t, createResponse.RiskLevel, getResponse.RiskLevel)
		assert.Equal(t, createResponse.Duration, getResponse.Duration)
		assert.Equal(t, createResponse.Description, getResponse.Description)
		assert.Equal(t, createResponse.Status, getResponse.Status)
	})

	t.Run("Get User Investments", func(t *testing.T) {
		// First create an investment
		var createResponse investments.Investment
		err := server.SendRequest("POST", "/api/v1/investments", testInvestment, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Get all investments for the user
		var getResponse []investments.Investment
		err = server.SendRequest("GET", "/api/v1/investments/user/"+userID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 200, getResponse[0].StatusCode)

		// Verify response
		assert.NotEmpty(t, getResponse)
		investmentMap := make(map[uuid.UUID]*investments.Investment)
		for _, investment := range getResponse {
			investmentMap[investment.ID] = &investment
		}
		assert.Contains(t, investmentMap, createResponse.ID)
	})

	t.Run("Update Investment Status", func(t *testing.T) {
		// First create an investment
		var createResponse investments.Investment
		err := server.SendRequest("POST", "/api/v1/investments", testInvestment, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Update investment status
		updateRequest := struct {
			Status       string `json:"status"`
			StatusReason string `json:"status_reason"`
		}{
			Status:       string(investments.InvestmentStatusCompleted),
			StatusReason: "Investment completed successfully",
		}

		var updateResponse investments.Investment
		err = server.SendRequest("PATCH", "/api/v1/investments/"+createResponse.ID.String()+"/status", updateRequest, &updateResponse)
		require.NoError(t, err)
		require.Equal(t, 200, updateResponse.StatusCode)

		// Verify response
		assert.Equal(t, investments.InvestmentStatusCompleted, updateResponse.Status)
		assert.Equal(t, "Investment completed successfully", updateResponse.StatusReason)
	})

	t.Run("Delete Investment", func(t *testing.T) {
		// First create an investment
		var createResponse investments.Investment
		err := server.SendRequest("POST", "/api/v1/investments", testInvestment, &createResponse)
		require.NoError(t, err)
		require.Equal(t, 201, createResponse.StatusCode)

		// Delete the investment
		var deleteResponse struct{}
		err = server.SendRequest("DELETE", "/api/v1/investments/"+createResponse.ID.String(), nil, &deleteResponse)
		require.NoError(t, err)
		require.Equal(t, 204, deleteResponse.StatusCode)

		// Try to get the deleted investment
		var getResponse investments.Investment
		err = server.SendRequest("GET", "/api/v1/investments/"+createResponse.ID.String(), nil, &getResponse)
		require.NoError(t, err)
		require.Equal(t, 404, getResponse.StatusCode)
	})
}

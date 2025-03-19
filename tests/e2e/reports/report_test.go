package reports

import (
	"testing"
	"time"

	"github.com/adil-faiyaz98/money-pulse/internal/reports"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestReportAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test")
	}

	// Setup test server
	server := testutil.NewTestServer(t)
	defer server.Close()

	// Create test user
	userID := uuid.New()
	user := map[string]interface{}{
		"id":        userID,
		"email":     "test@example.com",
		"firstName": "Test",
		"lastName":  "User",
	}
	server.SendRequest(t, "POST", "/api/v1/users", user)

	// Test data
	testReport := map[string]interface{}{
		"userId":      userID,
		"type":        reports.ReportTypeMonthly,
		"startDate":   time.Now().AddDate(0, -1, 0),
		"endDate":     time.Now(),
		"description": "Monthly financial report",
	}

	t.Run("Create Report", func(t *testing.T) {
		var response reports.Report
		status := server.SendRequest(t, "POST", "/api/v1/reports", testReport)
		assert.Equal(t, 201, status)
		server.DecodeResponse(t, &response)

		assert.NotEmpty(t, response.ID)
		assert.Equal(t, userID, response.UserID)
		assert.Equal(t, reports.ReportTypeMonthly, response.Type)
		assert.Equal(t, testReport["description"], response.Description)
		assert.Equal(t, reports.ReportStatusGenerated, response.Status)
	})

	t.Run("Get Report", func(t *testing.T) {
		// First create a report
		var createResponse reports.Report
		status := server.SendRequest(t, "POST", "/api/v1/reports", testReport)
		assert.Equal(t, 201, status)
		server.DecodeResponse(t, &createResponse)

		// Then get it
		var getResponse reports.Report
		status = server.SendRequest(t, "GET", "/api/v1/reports/"+createResponse.ID.String(), nil)
		assert.Equal(t, 200, status)
		server.DecodeResponse(t, &getResponse)

		assert.Equal(t, createResponse.ID, getResponse.ID)
		assert.Equal(t, createResponse.UserID, getResponse.UserID)
		assert.Equal(t, createResponse.Type, getResponse.Type)
		assert.Equal(t, createResponse.Description, getResponse.Description)
		assert.Equal(t, createResponse.Status, getResponse.Status)
	})

	t.Run("Get User Reports", func(t *testing.T) {
		// Create another report
		anotherReport := map[string]interface{}{
			"userId":      userID,
			"type":        reports.ReportTypeAnnual,
			"startDate":   time.Now().AddDate(-1, 0, 0),
			"endDate":     time.Now(),
			"description": "Annual financial report",
		}
		status := server.SendRequest(t, "POST", "/api/v1/reports", anotherReport)
		assert.Equal(t, 201, status)

		// Get all reports for the user
		var response []reports.Report
		status = server.SendRequest(t, "GET", "/api/v1/reports/user/"+userID.String(), nil)
		assert.Equal(t, 200, status)
		server.DecodeResponse(t, &response)

		assert.Len(t, response, 2)
	})

	t.Run("Update Report Status", func(t *testing.T) {
		// First create a report
		var createResponse reports.Report
		status := server.SendRequest(t, "POST", "/api/v1/reports", testReport)
		assert.Equal(t, 201, status)
		server.DecodeResponse(t, &createResponse)

		// Update its status
		updateRequest := map[string]interface{}{
			"status":       reports.ReportStatusFailed,
			"statusReason": "Failed to generate report",
		}
		status = server.SendRequest(t, "PATCH", "/api/v1/reports/"+createResponse.ID.String()+"/status", updateRequest)
		assert.Equal(t, 200, status)

		// Verify update
		var getResponse reports.Report
		status = server.SendRequest(t, "GET", "/api/v1/reports/"+createResponse.ID.String(), nil)
		assert.Equal(t, 200, status)
		server.DecodeResponse(t, &getResponse)

		assert.Equal(t, reports.ReportStatusFailed, getResponse.Status)
		assert.Equal(t, "Failed to generate report", getResponse.StatusReason)
	})

	t.Run("Delete Report", func(t *testing.T) {
		// First create a report
		var createResponse reports.Report
		status := server.SendRequest(t, "POST", "/api/v1/reports", testReport)
		assert.Equal(t, 201, status)
		server.DecodeResponse(t, &createResponse)

		// Delete it
		status = server.SendRequest(t, "DELETE", "/api/v1/reports/"+createResponse.ID.String(), nil)
		assert.Equal(t, 204, status)

		// Verify deletion
		status = server.SendRequest(t, "GET", "/api/v1/reports/"+createResponse.ID.String(), nil)
		assert.Equal(t, 404, status)
	})
}

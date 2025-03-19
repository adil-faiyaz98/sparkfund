package service

import (
	"testing"
	"time"

	"github.com/adil-faiyaz98/money-pulse/internal/reports"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReportService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// Setup test database
	testDB := testutil.NewTestDB(t)
	defer testDB.Close(t)

	// Create repository
	repo := NewPostgresReportRepository(testDB.DB)
	service := NewReportService(repo)

	// Create test context
	ctx := testutil.CreateTestContext(t)

	// Test data
	userID := uuid.New()
	testReport := &reports.Report{
		UserID:      userID,
		Type:        reports.ReportTypeMonthly,
		StartDate:   time.Now().AddDate(0, -1, 0),
		EndDate:     time.Now(),
		Description: "Monthly financial report",
	}

	t.Run("Create and Retrieve Report", func(t *testing.T) {
		// Create report
		err := service.CreateReport(ctx, testReport)
		require.NoError(t, err)
		require.NotEmpty(t, testReport.ID)
		require.Equal(t, reports.ReportStatusGenerated, testReport.Status)

		// Retrieve report
		report, err := service.GetReport(ctx, testReport.ID)
		require.NoError(t, err)
		assert.Equal(t, testReport.ID, report.ID)
		assert.Equal(t, testReport.UserID, report.UserID)
		assert.Equal(t, testReport.Type, report.Type)
		assert.Equal(t, testReport.StartDate, report.StartDate)
		assert.Equal(t, testReport.EndDate, report.EndDate)
		assert.Equal(t, testReport.Description, report.Description)
		assert.Equal(t, testReport.Status, report.Status)
	})

	t.Run("Get User Reports", func(t *testing.T) {
		// Create another report for the same user
		anotherReport := &reports.Report{
			UserID:      userID,
			Type:        reports.ReportTypeAnnual,
			StartDate:   time.Now().AddDate(-1, 0, 0),
			EndDate:     time.Now(),
			Description: "Annual financial report",
		}
		err := service.CreateReport(ctx, anotherReport)
		require.NoError(t, err)

		// Get all reports for the user
		reports, err := service.GetUserReports(ctx, userID)
		require.NoError(t, err)
		assert.Len(t, reports, 2)

		// Verify report details
		reportMap := make(map[uuid.UUID]*reports.Report)
		for _, report := range reports {
			reportMap[report.ID] = report
		}

		assert.Contains(t, reportMap, testReport.ID)
		assert.Contains(t, reportMap, anotherReport.ID)
	})

	t.Run("Update Report Status", func(t *testing.T) {
		// Update report status
		err := service.UpdateReportStatus(ctx, testReport.ID, reports.ReportStatusFailed, "Failed to generate report")
		require.NoError(t, err)

		// Verify update
		updated, err := service.GetReport(ctx, testReport.ID)
		require.NoError(t, err)
		assert.Equal(t, reports.ReportStatusFailed, updated.Status)
		assert.Equal(t, "Failed to generate report", updated.StatusReason)
	})

	t.Run("Delete Report", func(t *testing.T) {
		// Delete report
		err := service.DeleteReport(ctx, testReport.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = service.GetReport(ctx, testReport.ID)
		assert.Error(t, err)
	})
}

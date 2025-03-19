package service

import (
	"context"
	"testing"
	"time"

	"github.com/adil-faiyaz98/money-pulse/internal/reports"
	"github.com/adil-faiyaz98/money-pulse/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockReportRepository struct {
	testutil.MockRepository
}

func (m *MockReportRepository) Create(ctx context.Context, report *reports.Report) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *MockReportRepository) GetByID(ctx context.Context, id uuid.UUID) (*reports.Report, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*reports.Report), args.Error(1)
}

func (m *MockReportRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*reports.Report, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*reports.Report), args.Error(1)
}

func (m *MockReportRepository) Update(ctx context.Context, report *reports.Report) error {
	args := m.Called(ctx, report)
	return args.Error(0)
}

func (m *MockReportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestReportService_CreateReport(t *testing.T) {
	// Setup
	mockRepo := new(MockReportRepository)
	service := NewReportService(mockRepo)
	ctx := context.Background()

	// Test data
	userID := uuid.New()
	testReport := &reports.Report{
		UserID:      userID,
		Type:        reports.ReportTypeMonthly,
		StartDate:   time.Now().AddDate(0, -1, 0),
		EndDate:     time.Now(),
		Description: "Monthly financial report",
	}

	t.Run("Create Valid Report", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("Create", ctx, mock.AnythingOfType("*reports.Report")).Return(nil)

		// Execute
		err := service.CreateReport(ctx, testReport)

		// Assert
		require.NoError(t, err)
		assert.NotEmpty(t, testReport.ID)
		assert.Equal(t, reports.ReportStatusGenerated, testReport.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create Report with Invalid Type", func(t *testing.T) {
		// Setup test data
		invalidReport := *testReport
		invalidReport.Type = "invalid_type"

		// Execute
		err := service.CreateReport(ctx, &invalidReport)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, invalidReport.ID)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Create Report with Invalid Date Range", func(t *testing.T) {
		// Setup test data
		invalidReport := *testReport
		invalidReport.StartDate = time.Now()
		invalidReport.EndDate = time.Now().AddDate(0, -1, 0)

		// Execute
		err := service.CreateReport(ctx, &invalidReport)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, invalidReport.ID)
		mockRepo.AssertNotCalled(t, "Create")
	})
}

func TestReportService_GetReport(t *testing.T) {
	// Setup
	mockRepo := new(MockReportRepository)
	service := NewReportService(mockRepo)
	ctx := context.Background()

	// Test data
	reportID := uuid.New()
	testReport := &reports.Report{
		ID:          reportID,
		UserID:      uuid.New(),
		Type:        reports.ReportTypeMonthly,
		StartDate:   time.Now().AddDate(0, -1, 0),
		EndDate:     time.Now(),
		Description: "Monthly financial report",
		Status:      reports.ReportStatusGenerated,
	}

	t.Run("Get Existing Report", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, reportID).Return(testReport, nil)

		// Execute
		report, err := service.GetReport(ctx, reportID)

		// Assert
		require.NoError(t, err)
		assert.Equal(t, testReport.ID, report.ID)
		assert.Equal(t, testReport.UserID, report.UserID)
		assert.Equal(t, testReport.Type, report.Type)
		assert.Equal(t, testReport.StartDate, report.StartDate)
		assert.Equal(t, testReport.EndDate, report.EndDate)
		assert.Equal(t, testReport.Description, report.Description)
		assert.Equal(t, testReport.Status, report.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Get Non-Existing Report", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, reportID).Return(nil, reports.ErrReportNotFound)

		// Execute
		report, err := service.GetReport(ctx, reportID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, report)
		assert.Equal(t, reports.ErrReportNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestReportService_UpdateReportStatus(t *testing.T) {
	// Setup
	mockRepo := new(MockReportRepository)
	service := NewReportService(mockRepo)
	ctx := context.Background()

	// Test data
	reportID := uuid.New()
	testReport := &reports.Report{
		ID:          reportID,
		UserID:      uuid.New(),
		Type:        reports.ReportTypeMonthly,
		StartDate:   time.Now().AddDate(0, -1, 0),
		EndDate:     time.Now(),
		Description: "Monthly financial report",
		Status:      reports.ReportStatusGenerating,
	}

	t.Run("Update Report Status Successfully", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, reportID).Return(testReport, nil)
		mockRepo.On("Update", ctx, mock.AnythingOfType("*reports.Report")).Return(nil)

		// Execute
		err := service.UpdateReportStatus(ctx, reportID, reports.ReportStatusGenerated, "Report generated successfully")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, reports.ReportStatusGenerated, testReport.Status)
		assert.Equal(t, "Report generated successfully", testReport.StatusReason)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update Non-Existing Report Status", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, reportID).Return(nil, reports.ErrReportNotFound)

		// Execute
		err := service.UpdateReportStatus(ctx, reportID, reports.ReportStatusGenerated, "Report generated successfully")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, reports.ErrReportNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update Report Status with Invalid Status", func(t *testing.T) {
		// Setup expectations
		mockRepo.On("GetByID", ctx, reportID).Return(testReport, nil)

		// Execute
		err := service.UpdateReportStatus(ctx, reportID, "invalid_status", "Invalid status")

		// Assert
		assert.Error(t, err)
		assert.Equal(t, reports.ReportStatusGenerating, testReport.Status)
		mockRepo.AssertNotCalled(t, "Update")
	})
}

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/adil-faiyaz98/structgen/internal/reports"
	"github.com/google/uuid"
)

type reportService struct {
	repo reports.ReportRepository
}

func NewReportService(repo reports.ReportRepository) reports.ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) CreateReport(ctx context.Context, report *reports.Report) error {
	// Validate report type
	if !isValidReportType(report.Type) {
		return fmt.Errorf("invalid report type: %s", report.Type)
	}

	// Validate report format
	if !isValidReportFormat(report.Format) {
		return fmt.Errorf("invalid report format: %s", report.Format)
	}

	// Set initial status
	report.Status = reports.ReportStatusPending

	return s.repo.Create(report)
}

func (s *reportService) GetReport(ctx context.Context, id uuid.UUID) (*reports.Report, error) {
	return s.repo.GetByID(id)
}

func (s *reportService) GetUserReports(ctx context.Context, userID uuid.UUID) ([]*reports.Report, error) {
	return s.repo.GetByUserID(userID)
}

func (s *reportService) UpdateReportStatus(ctx context.Context, id uuid.UUID, status reports.ReportStatus, fileURL string, err error) error {
	report, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("report not found: %w", err)
	}

	report.Status = status
	report.FileURL = fileURL

	now := time.Now()
	switch status {
	case reports.ReportStatusCompleted:
		report.GeneratedAt = &now
	case reports.ReportStatusFailed:
		report.FailedAt = &now
		if err != nil {
			report.Error = err.Error()
		}
	}

	return s.repo.Update(report)
}

func (s *reportService) DeleteReport(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(id)
}

func isValidReportType(reportType reports.ReportType) bool {
	switch reportType {
	case reports.ReportTypeBalance,
		reports.ReportTypeTransaction,
		reports.ReportTypeInvestment,
		reports.ReportTypeLoan,
		reports.ReportTypeTax:
		return true
	default:
		return false
	}
}

func isValidReportFormat(format reports.ReportFormat) bool {
	switch format {
	case reports.ReportFormatPDF,
		reports.ReportFormatCSV,
		reports.ReportFormatJSON,
		reports.ReportFormatXLSX:
		return true
	default:
		return false
	}
}

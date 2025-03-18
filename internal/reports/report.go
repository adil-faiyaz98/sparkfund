package reports

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type ReportType string

const (
	ReportTypeBalance     ReportType = "BALANCE"
	ReportTypeTransaction ReportType = "TRANSACTION"
	ReportTypeInvestment  ReportType = "INVESTMENT"
	ReportTypeLoan        ReportType = "LOAN"
	ReportTypeTax         ReportType = "TAX"
)

type ReportFormat string

const (
	ReportFormatPDF  ReportFormat = "PDF"
	ReportFormatCSV  ReportFormat = "CSV"
	ReportFormatJSON ReportFormat = "JSON"
	ReportFormatXLSX ReportFormat = "XLSX"
)

type ReportStatus string

const (
	ReportStatusPending    ReportStatus = "PENDING"
	ReportStatusGenerating ReportStatus = "GENERATING"
	ReportStatusCompleted  ReportStatus = "COMPLETED"
	ReportStatusFailed     ReportStatus = "FAILED"
)

type Report struct {
	ID          uuid.UUID    `json:"id" gorm:"type:uuid;primary_key"`
	UserID      uuid.UUID    `json:"user_id" gorm:"type:uuid;not null"`
	Type        ReportType   `json:"type" gorm:"type:varchar(20);not null"`
	Format      ReportFormat `json:"format" gorm:"type:varchar(10);not null"`
	Status      ReportStatus `json:"status" gorm:"type:varchar(20);not null;default:'PENDING'"`
	Parameters  string       `json:"parameters" gorm:"type:jsonb"`
	FileURL     string       `json:"file_url" gorm:"type:text"`
	Error       string       `json:"error" gorm:"type:text"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	GeneratedAt *time.Time   `json:"generated_at"`
	FailedAt    *time.Time   `json:"failed_at"`
}

type ReportRepository interface {
	Create(report *Report) error
	GetByID(id uuid.UUID) (*Report, error)
	GetByUserID(userID uuid.UUID) ([]*Report, error)
	Update(report *Report) error
	Delete(id uuid.UUID) error
}

type ReportService interface {
	CreateReport(ctx context.Context, report *Report) error
	GetReport(ctx context.Context, id uuid.UUID) (*Report, error)
	GetUserReports(ctx context.Context, userID uuid.UUID) ([]*Report, error)
	UpdateReportStatus(ctx context.Context, id uuid.UUID, status ReportStatus, fileURL string, err error) error
	DeleteReport(ctx context.Context, id uuid.UUID) error
}

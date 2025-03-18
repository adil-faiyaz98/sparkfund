package postgres

import (
	"fmt"

	"github.com/adil-faiyaz98/structgen/internal/reports"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) reports.ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(report *reports.Report) error {
	return r.db.Create(report).Error
}

func (r *reportRepository) GetByID(id uuid.UUID) (*reports.Report, error) {
	var report reports.Report
	if err := r.db.First(&report, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("report not found: %w", err)
	}
	return &report, nil
}

func (r *reportRepository) GetByUserID(userID uuid.UUID) ([]*reports.Report, error) {
	var reportList []*reports.Report
	if err := r.db.Where("user_id = ?", userID).Find(&reportList).Error; err != nil {
		return nil, fmt.Errorf("failed to get user reports: %w", err)
	}
	return reportList, nil
}

func (r *reportRepository) Update(report *reports.Report) error {
	return r.db.Save(report).Error
}

func (r *reportRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&reports.Report{}, "id = ?", id).Error
}

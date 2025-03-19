package repositories

import (
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
)

// Repository defines the interface for data access
type Repository interface {
	// EmailLog operations
	CreateEmailLog(log *models.EmailLog) error
	GetEmailLog(id string) (*models.EmailLog, error)
	UpdateEmailLog(log *models.EmailLog) error
	GetEmailLogs() ([]models.EmailLog, error)
	GetEmailStats() (*models.EmailStats, error)

	// Template operations
	CreateTemplate(template *models.Template) error
	GetTemplate(id string) (*models.Template, error)
	UpdateTemplate(template *models.Template) error
	DeleteTemplate(id string) error
	GetTemplates() ([]models.Template, error)
}

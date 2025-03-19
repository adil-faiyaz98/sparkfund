package repositories

import (
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
)

// EmailRepository defines the interface for email persistence
type EmailRepository interface {
	// Email Log Operations
	CreateEmailLog(log models.EmailLog) error
	UpdateEmailLog(log models.EmailLog) error
	GetEmailLog(id string) (*models.EmailLog, error)
	GetEmailLogs() ([]models.EmailLog, error)

	// Template Operations
	CreateTemplate(template models.Template) error
	GetTemplate(id string) (*models.Template, error)
	UpdateTemplate(template models.Template) error
	DeleteTemplate(id string) error
	GetTemplates() ([]models.Template, error)
}

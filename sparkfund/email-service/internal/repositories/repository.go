package repositories

import (
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
)

// Repository defines the interface for email service data operations
type Repository interface {
	// Email logs
	CreateEmailLog(log models.EmailLog) error
	GetEmailLog(id string) (*models.EmailLog, error)
	UpdateEmailLog(log models.EmailLog) error
	GetEmailLogs() ([]models.EmailLog, error)

	// Templates
	CreateTemplate(template models.Template) error
	GetTemplate(id string) (*models.Template, error)
	UpdateTemplate(template models.Template) error
	DeleteTemplate(id string) error
	GetTemplates() ([]models.Template, error)
}

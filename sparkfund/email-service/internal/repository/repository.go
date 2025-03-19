package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
	"github.com/jmoiron/sqlx"
)

// Repository defines the interface for database operations
type Repository interface {
	CreateEmailLog(ctx context.Context, log models.EmailLog) error
	GetEmailLogs(ctx context.Context) ([]models.EmailLog, error)
	CreateTemplate(ctx context.Context, template *models.Template) error
	GetTemplate(ctx context.Context, id string) (*models.Template, error)
	UpdateTemplate(ctx context.Context, template *models.Template) error
	DeleteTemplate(ctx context.Context, id string) error
}

// postgresRepository implements the Repository interface
type postgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL repository instance
func NewPostgresRepository(db *sqlx.DB) Repository {
	return &postgresRepository{db: db}
}

// CreateEmailLog creates a new email log entry
func (r *postgresRepository) CreateEmailLog(ctx context.Context, log models.EmailLog) error {
	query := `
		INSERT INTO email_logs (id, to_addresses, cc_addresses, bcc_addresses, subject, body, status, error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.To,
		log.Cc,
		log.Bcc,
		log.Subject,
		log.Body,
		log.Status,
		log.Error,
		log.CreatedAt,
		log.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create email log: %w", err)
	}

	return nil
}

// GetEmailLogs retrieves all email logs
func (r *postgresRepository) GetEmailLogs(ctx context.Context) ([]models.EmailLog, error) {
	query := `
		SELECT id, to_addresses, cc_addresses, bcc_addresses, subject, body, status, error, created_at, updated_at
		FROM email_logs
		ORDER BY created_at DESC
	`

	var logs []models.EmailLog
	if err := r.db.SelectContext(ctx, &logs, query); err != nil {
		return nil, fmt.Errorf("failed to get email logs: %w", err)
	}

	return logs, nil
}

// CreateTemplate creates a new email template
func (r *postgresRepository) CreateTemplate(ctx context.Context, template *models.Template) error {
	query := `
		INSERT INTO templates (id, name, subject, body, variables, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		template.ID,
		template.Name,
		template.Subject,
		template.Body,
		template.Variables,
		template.Description,
		template.CreatedAt,
		template.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}

// GetTemplate retrieves a template by ID
func (r *postgresRepository) GetTemplate(ctx context.Context, id string) (*models.Template, error) {
	query := `
		SELECT id, name, subject, body, variables, description, created_at, updated_at
		FROM templates
		WHERE id = $1
	`

	var template models.Template
	if err := r.db.GetContext(ctx, &template, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("template not found")
		}
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &template, nil
}

// UpdateTemplate updates an existing template
func (r *postgresRepository) UpdateTemplate(ctx context.Context, template *models.Template) error {
	query := `
		UPDATE templates
		SET name = $1, subject = $2, body = $3, variables = $4, description = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(ctx, query,
		template.Name,
		template.Subject,
		template.Body,
		template.Variables,
		template.Description,
		template.UpdatedAt,
		template.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

// DeleteTemplate deletes a template by ID
func (r *postgresRepository) DeleteTemplate(ctx context.Context, id string) error {
	query := `DELETE FROM templates WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete template: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("template not found")
	}

	return nil
}

package repositories

import (
	"context"
	"database/sql"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/errors"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sqlx.DB) Repository {
	return &PostgresRepository{
		db:     db,
		logger: zap.L().Named("postgres_repository"),
	}
}

// CreateEmailLog creates a new email log entry
func (r *PostgresRepository) CreateEmailLog(log models.EmailLog) error {
	query := `
		INSERT INTO email_logs (id, to_addresses, from_address, subject, body, content_type, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(context.Background(), query,
		log.ID,
		log.To,
		log.From,
		log.Subject,
		log.Body,
		log.ContentType,
		log.Status,
		log.CreatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to create email log",
			zap.Error(err),
			zap.String("id", log.ID),
			zap.Strings("to", log.To),
		)
		return errors.NewDatabaseError(err)
	}

	r.logger.Info("Email log created successfully",
		zap.String("id", log.ID),
		zap.Strings("to", log.To),
	)
	return nil
}

// GetEmailLog retrieves an email log by ID
func (r *PostgresRepository) GetEmailLog(id string) (*models.EmailLog, error) {
	var log models.EmailLog
	query := `SELECT * FROM email_logs WHERE id = $1`

	err := r.db.GetContext(context.Background(), &log, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		r.logger.Error("Failed to get email log",
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, errors.NewDatabaseError(err)
	}

	return &log, nil
}

// GetEmailLogs retrieves all email logs
func (r *PostgresRepository) GetEmailLogs() ([]models.EmailLog, error) {
	var logs []models.EmailLog
	query := `SELECT * FROM email_logs ORDER BY created_at DESC`

	err := r.db.SelectContext(context.Background(), &logs, query)
	if err != nil {
		r.logger.Error("Failed to get email logs",
			zap.Error(err),
		)
		return nil, errors.NewDatabaseError(err)
	}

	return logs, nil
}

// UpdateEmailLog updates an existing email log
func (r *PostgresRepository) UpdateEmailLog(log models.EmailLog) error {
	query := `
		UPDATE email_logs
		SET status = $1, error = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(context.Background(), query,
		log.Status,
		log.Error,
		log.UpdatedAt,
		log.ID,
	)
	if err != nil {
		r.logger.Error("Failed to update email log",
			zap.Error(err),
			zap.String("id", log.ID),
		)
		return errors.NewDatabaseError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	if rows == 0 {
		return errors.ErrNotFound
	}

	r.logger.Info("Email log updated successfully",
		zap.String("id", log.ID),
		zap.String("status", log.Status),
	)
	return nil
}

// CreateTemplate creates a new email template
func (r *PostgresRepository) CreateTemplate(template models.Template) error {
	query := `
		INSERT INTO templates (id, name, subject, body, variables, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(context.Background(), query,
		template.ID,
		template.Name,
		template.Subject,
		template.Body,
		template.Variables,
		template.Description,
		template.CreatedAt,
	)
	if err != nil {
		r.logger.Error("Failed to create template",
			zap.Error(err),
			zap.String("id", template.ID),
			zap.String("name", template.Name),
		)
		return errors.NewDatabaseError(err)
	}

	r.logger.Info("Template created successfully",
		zap.String("id", template.ID),
		zap.String("name", template.Name),
	)
	return nil
}

// GetTemplate retrieves a template by ID
func (r *PostgresRepository) GetTemplate(id string) (*models.Template, error) {
	var template models.Template
	query := `SELECT * FROM templates WHERE id = $1`

	err := r.db.GetContext(context.Background(), &template, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.ErrNotFound
		}
		r.logger.Error("Failed to get template",
			zap.Error(err),
			zap.String("id", id),
		)
		return nil, errors.NewDatabaseError(err)
	}

	return &template, nil
}

// UpdateTemplate updates an existing template
func (r *PostgresRepository) UpdateTemplate(template models.Template) error {
	query := `
		UPDATE templates
		SET name = $1, subject = $2, body = $3, variables = $4, description = $5, updated_at = $6
		WHERE id = $7
	`

	result, err := r.db.ExecContext(context.Background(), query,
		template.Name,
		template.Subject,
		template.Body,
		template.Variables,
		template.Description,
		template.UpdatedAt,
		template.ID,
	)
	if err != nil {
		r.logger.Error("Failed to update template",
			zap.Error(err),
			zap.String("id", template.ID),
		)
		return errors.NewDatabaseError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	if rows == 0 {
		return errors.ErrNotFound
	}

	r.logger.Info("Template updated successfully",
		zap.String("id", template.ID),
		zap.String("name", template.Name),
	)
	return nil
}

// DeleteTemplate deletes a template by ID
func (r *PostgresRepository) DeleteTemplate(id string) error {
	query := `DELETE FROM templates WHERE id = $1`

	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		r.logger.Error("Failed to delete template",
			zap.Error(err),
			zap.String("id", id),
		)
		return errors.NewDatabaseError(err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	if rows == 0 {
		return errors.ErrNotFound
	}

	r.logger.Info("Template deleted successfully",
		zap.String("id", id),
	)
	return nil
}

// GetTemplates retrieves all templates
func (r *PostgresRepository) GetTemplates() ([]models.Template, error) {
	var templates []models.Template
	query := `SELECT * FROM templates ORDER BY created_at DESC`

	err := r.db.SelectContext(context.Background(), &templates, query)
	if err != nil {
		r.logger.Error("Failed to get templates",
			zap.Error(err),
		)
		return nil, errors.NewDatabaseError(err)
	}

	return templates, nil
}

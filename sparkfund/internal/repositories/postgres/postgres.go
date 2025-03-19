package postgres

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/errors"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
	"github.com/adil-faiyaz98/sparkfund/email-service/internal/repositories"
	"github.com/jmoiron/sqlx"
)

// PostgresRepository implements the Repository interface using PostgreSQL
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *sqlx.DB) repositories.Repository {
	return &PostgresRepository{db: db}
}

// CreateEmailLog creates a new email log entry
func (r *PostgresRepository) CreateEmailLog(log *models.EmailLog) error {
	query := `
		INSERT INTO email_logs (
			id, recipients, "from", subject, body, content_type,
			status, error, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`

	recipientsJSON, err := json.Marshal(log.Recipients)
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	_, err = r.db.Exec(query,
		log.ID,
		recipientsJSON,
		log.From,
		log.Subject,
		log.Body,
		log.ContentType,
		log.Status,
		log.Error,
		log.CreatedAt,
		log.UpdatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err)
	}

	return nil
}

// GetEmailLog retrieves an email log entry by ID
func (r *PostgresRepository) GetEmailLog(id string) (*models.EmailLog, error) {
	var log models.EmailLog
	var recipientsJSON []byte

	query := `
		SELECT id, recipients, "from", subject, body, content_type,
			status, error, created_at, updated_at
		FROM email_logs
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&log.ID,
		&recipientsJSON,
		&log.From,
		&log.Subject,
		&log.Body,
		&log.ContentType,
		&log.Status,
		&log.Error,
		&log.CreatedAt,
		&log.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("email log not found")
	}

	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	if err := json.Unmarshal(recipientsJSON, &log.Recipients); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return &log, nil
}

// UpdateEmailLog updates an existing email log entry
func (r *PostgresRepository) UpdateEmailLog(log *models.EmailLog) error {
	query := `
		UPDATE email_logs
		SET recipients = $1, "from" = $2, subject = $3,
			body = $4, content_type = $5, status = $6,
			error = $7, updated_at = $8
		WHERE id = $9
	`

	recipientsJSON, err := json.Marshal(log.Recipients)
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	result, err := r.db.Exec(query,
		recipientsJSON,
		log.From,
		log.Subject,
		log.Body,
		log.ContentType,
		log.Status,
		log.Error,
		time.Now(),
		log.ID,
	)

	if err != nil {
		return errors.NewDatabaseError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("email log not found")
	}

	return nil
}

// GetEmailLogs retrieves all email log entries
func (r *PostgresRepository) GetEmailLogs() ([]models.EmailLog, error) {
	var logs []models.EmailLog

	query := `
		SELECT id, recipients, "from", subject, body, content_type,
			status, error, created_at, updated_at
		FROM email_logs
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var log models.EmailLog
		var recipientsJSON []byte

		err := rows.Scan(
			&log.ID,
			&recipientsJSON,
			&log.From,
			&log.Subject,
			&log.Body,
			&log.ContentType,
			&log.Status,
			&log.Error,
			&log.CreatedAt,
			&log.UpdatedAt,
		)

		if err != nil {
			return nil, errors.NewDatabaseError(err)
		}

		if err := json.Unmarshal(recipientsJSON, &log.Recipients); err != nil {
			return nil, errors.NewDatabaseError(err)
		}

		logs = append(logs, log)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return logs, nil
}

// GetEmailStats retrieves email statistics
func (r *PostgresRepository) GetEmailStats() (*models.EmailStats, error) {
	var stats models.EmailStats

	query := `
		SELECT
			COUNT(*) FILTER (WHERE status = 'sent') as total_sent,
			COUNT(*) FILTER (WHERE status = 'failed') as total_failed,
			COUNT(*) FILTER (WHERE status = 'pending') as total_pending,
			EXTRACT(EPOCH FROM AVG(updated_at - created_at)) as average_latency
		FROM email_logs
	`

	err := r.db.QueryRow(query).Scan(
		&stats.TotalSent,
		&stats.TotalFailed,
		&stats.TotalPending,
		&stats.AverageLatency,
	)

	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return &stats, nil
}

// CreateTemplate creates a new email template
func (r *PostgresRepository) CreateTemplate(template *models.Template) error {
	query := `
		INSERT INTO templates (
			id, name, subject, content, content_type,
			data, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`

	dataJSON, err := json.Marshal(template.Data)
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	_, err = r.db.Exec(query,
		template.ID,
		template.Name,
		template.Subject,
		template.Content,
		template.ContentType,
		dataJSON,
		template.CreatedAt,
		template.UpdatedAt,
	)

	if err != nil {
		return errors.NewDatabaseError(err)
	}

	return nil
}

// GetTemplate retrieves an email template by ID
func (r *PostgresRepository) GetTemplate(id string) (*models.Template, error) {
	var template models.Template
	var dataJSON []byte

	query := `
		SELECT id, name, subject, content, content_type,
			data, created_at, updated_at
		FROM templates
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&template.ID,
		&template.Name,
		&template.Subject,
		&template.Content,
		&template.ContentType,
		&dataJSON,
		&template.CreatedAt,
		&template.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.NewNotFoundError("template not found")
	}

	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	if err := json.Unmarshal(dataJSON, &template.Data); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return &template, nil
}

// UpdateTemplate updates an existing email template
func (r *PostgresRepository) UpdateTemplate(template *models.Template) error {
	query := `
		UPDATE templates
		SET name = $1, subject = $2, content = $3,
			content_type = $4, data = $5, updated_at = $6
		WHERE id = $7
	`

	dataJSON, err := json.Marshal(template.Data)
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	result, err := r.db.Exec(query,
		template.Name,
		template.Subject,
		template.Content,
		template.ContentType,
		dataJSON,
		time.Now(),
		template.ID,
	)

	if err != nil {
		return errors.NewDatabaseError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("template not found")
	}

	return nil
}

// DeleteTemplate deletes an email template
func (r *PostgresRepository) DeleteTemplate(id string) error {
	query := `DELETE FROM templates WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewDatabaseError(err)
	}

	if rowsAffected == 0 {
		return errors.NewNotFoundError("template not found")
	}

	return nil
}

// GetTemplates retrieves all email templates
func (r *PostgresRepository) GetTemplates() ([]models.Template, error) {
	var templates []models.Template

	query := `
		SELECT id, name, subject, content, content_type,
			data, created_at, updated_at
		FROM templates
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var template models.Template
		var dataJSON []byte

		err := rows.Scan(
			&template.ID,
			&template.Name,
			&template.Subject,
			&template.Content,
			&template.ContentType,
			&dataJSON,
			&template.CreatedAt,
			&template.UpdatedAt,
		)

		if err != nil {
			return nil, errors.NewDatabaseError(err)
		}

		if err := json.Unmarshal(dataJSON, &template.Data); err != nil {
			return nil, errors.NewDatabaseError(err)
		}

		templates = append(templates, template)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	return templates, nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/sparkfund/email-service/internal/models"
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

// postgresRepository implements the Repository interface using PostgreSQL
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
		INSERT INTO email_logs (id, recipients, cc, bcc, from_address, subject, body, content_type, status, error, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.Recipients,
		log.Cc,
		log.Bcc,
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
		return fmt.Errorf("failed to create email log: %w", err)
	}

	return nil
}

// GetEmailLogs retrieves all email logs
func (r *postgresRepository) GetEmailLogs(ctx context.Context) ([]models.EmailLog, error) {
	query := `
		SELECT id, recipients, cc, bcc, from_address as from, subject, body, content_type, status, error, created_at, updated_at
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

// memoryRepository implements the Repository interface using in-memory storage
type memoryRepository struct {
	mu        sync.RWMutex
	logs      map[string]models.EmailLog
	templates map[string]*models.Template
}

// NewMemoryRepository creates a new in-memory repository instance
func NewMemoryRepository() Repository {
	return &memoryRepository{
		logs:      make(map[string]models.EmailLog),
		templates: make(map[string]*models.Template),
	}
}

// CreateEmailLog creates a new email log entry
func (r *memoryRepository) CreateEmailLog(ctx context.Context, log models.EmailLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.logs[log.ID]; exists {
		return fmt.Errorf("log with ID %s already exists", log.ID)
	}

	r.logs[log.ID] = log
	return nil
}

// GetEmailLogs retrieves all email logs
func (r *memoryRepository) GetEmailLogs(ctx context.Context) ([]models.EmailLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	logs := make([]models.EmailLog, 0, len(r.logs))
	for _, log := range r.logs {
		logs = append(logs, log)
	}

	return logs, nil
}

// CreateTemplate creates a new email template
func (r *memoryRepository) CreateTemplate(ctx context.Context, template *models.Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[template.ID]; exists {
		return fmt.Errorf("template with ID %s already exists", template.ID)
	}

	r.templates[template.ID] = template
	return nil
}

// GetTemplate retrieves a template by ID
func (r *memoryRepository) GetTemplate(ctx context.Context, id string) (*models.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.templates[id]
	if !exists {
		return nil, fmt.Errorf("template with ID %s does not exist", id)
	}

	return template, nil
}

// UpdateTemplate updates an existing template
func (r *memoryRepository) UpdateTemplate(ctx context.Context, template *models.Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[template.ID]; !exists {
		return fmt.Errorf("template with ID %s does not exist", template.ID)
	}

	r.templates[template.ID] = template
	return nil
}

// DeleteTemplate deletes a template by ID
func (r *memoryRepository) DeleteTemplate(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[id]; !exists {
		return fmt.Errorf("template with ID %s does not exist", id)
	}

	delete(r.templates, id)
	return nil
}

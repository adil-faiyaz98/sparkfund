package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sparkfund/email-service/internal/errors"
	"github.com/sparkfund/email-service/internal/models"
	"go.uber.org/zap" // Import zap
)

// Repository defines the interface for database operations
type Repository interface {
	CreateEmailLog(ctx context.Context, log models.EmailLog) error
	GetEmailLogs(ctx context.Context) ([]models.EmailLog, error)
	CreateTemplate(ctx context.Context, template *models.Template) error
	GetTemplate(ctx context.Context, id uuid.UUID) (*models.Template, error)
	UpdateTemplate(ctx context.Context, template *models.Template) error
	DeleteTemplate(ctx context.Context, id uuid.UUID) error
}

// postgresRepository implements the Repository interface using PostgreSQL
type postgresRepository struct {
	db     *sqlx.DB
	logger *zap.Logger // Add logger
}

// NewPostgresRepository creates a new PostgreSQL repository instance
func NewPostgresRepository(db *sqlx.DB, logger *zap.Logger) Repository {
	return &postgresRepository{db: db, logger: logger}
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
		r.logger.Error("Failed to create email log", zap.Error(err))
		return errors.NewDatabaseError(err, errors.ComponentRepository)
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
		r.logger.Error("Failed to get email logs", zap.Error(err))
		return nil, errors.NewDatabaseError(err, errors.ComponentRepository)
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
		r.logger.Error("Failed to create template", zap.Error(err))
		return errors.NewDatabaseError(err, errors.ComponentRepository)
	}

	return nil
}

// GetTemplate retrieves a template by ID
func (r *postgresRepository) GetTemplate(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	query := `
		SELECT id, name, subject, body, variables, description, created_at, updated_at
		FROM templates
		WHERE id = $1
	`

	var template models.Template
	err := r.db.GetContext(ctx, &template, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Template not found", zap.String("template_id", id.String()))
			return nil, errors.NewNotFoundError("template not found")
		}
		r.logger.Error("Failed to get template", zap.Error(err), zap.String("template_id", id.String()))
		return nil, errors.NewDatabaseError(err, errors.ComponentRepository)
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
		r.logger.Error("Failed to update template", zap.Error(err), zap.String("template_id", template.ID.String()))
		return errors.NewDatabaseError(err, errors.ComponentRepository)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err), zap.String("template_id", template.ID.String()))
		return errors.NewDatabaseError(err, errors.ComponentRepository)
	}

	if rows == 0 {
		r.logger.Warn("Template not found", zap.String("template_id", template.ID.String()))
		return errors.NewNotFoundError("template not found")
	}

	return nil
}

// DeleteTemplate deletes a template by ID
func (r *postgresRepository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM templates WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete template", zap.Error(err), zap.String("template_id", id.String()))
		return errors.NewDatabaseError(err, errors.ComponentRepository)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err), zap.String("template_id", id.String()))
		return errors.NewDatabaseError(err, errors.ComponentRepository)
	}

	if rows == 0 {
		r.logger.Warn("Template not found", zap.String("template_id", id.String()))
		return errors.NewNotFoundError("template not found")
	}

	return nil
}

// memoryRepository implements the Repository interface using in-memory storage
type memoryRepository struct {
	mu        sync.RWMutex
	logs      map[uuid.UUID]models.EmailLog
	templates map[uuid.UUID]*models.Template
	logger    *zap.Logger // Add logger
}

// NewMemoryRepository creates a new in-memory repository instance
func NewMemoryRepository(logger *zap.Logger) Repository {
	return &memoryRepository{
		logs:      make(map[uuid.UUID]models.EmailLog),
		templates: make(map[uuid.UUID]*models.Template),
		logger:    logger,
	}
}

// CreateEmailLog creates a new email log entry
func (r *memoryRepository) CreateEmailLog(ctx context.Context, log models.EmailLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.logs[log.ID]; exists {
		err := fmt.Errorf("log with ID %s already exists", log.ID)
		r.logger.Error("Failed to create email log", zap.Error(err), zap.String("email_log_id", log.ID.String()))
		return err
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
		err := fmt.Errorf("template with ID %s already exists", template.ID)
		r.logger.Error("Failed to create template", zap.Error(err), zap.String("template_id", template.ID.String()))
		return err
	}

	r.templates[template.ID] = template
	return nil
}

// GetTemplate retrieves a template by ID
func (r *memoryRepository) GetTemplate(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.templates[id]
	if !exists {
		err := fmt.Errorf("template with ID %s does not exist", id)
		r.logger.Warn("Template not found", zap.Error(err), zap.String("template_id", id.String()))
		return nil, errors.NewNotFoundError("template not found")
	}

	return template, nil
}

// UpdateTemplate updates an existing template
func (r *memoryRepository) UpdateTemplate(ctx context.Context, template *models.Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[template.ID]; !exists {
		err := fmt.Errorf("template with ID %s does not exist", template.ID)
		r.logger.Error("Failed to update template", zap.Error(err), zap.String("template_id", template.ID.String()))
		return err
	}

	r.templates[template.ID] = template
	return nil
}

// DeleteTemplate deletes a template by ID
func (r *memoryRepository) DeleteTemplate(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[id]; !exists {
		err := fmt.Errorf("template with ID %s does not exist", id)
		r.logger.Error("Failed to delete template", zap.Error(err), zap.String("template_id", id.String()))
		return err
	}

	delete(r.templates, id)
	return nil
}

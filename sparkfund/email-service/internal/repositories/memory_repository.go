package repositories

import (
	"fmt"
	"sync"
	"time"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
)

// MemoryRepository implements EmailRepository using in-memory storage
type MemoryRepository struct {
	mu        sync.RWMutex
	logs      map[string]*models.EmailLog
	templates map[string]*models.Template
}

// NewMemoryRepository creates a new in-memory repository instance
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		logs:      make(map[string]*models.EmailLog),
		templates: make(map[string]*models.Template),
	}
}

// CreateLog creates a new email log entry
func (r *MemoryRepository) CreateLog(log *models.EmailLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.logs[log.ID]; exists {
		return fmt.Errorf("log with ID %s already exists", log.ID)
	}

	r.logs[log.ID] = log
	return nil
}

// UpdateLog updates an existing email log entry
func (r *MemoryRepository) UpdateLog(log *models.EmailLog) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.logs[log.ID]; !exists {
		return fmt.Errorf("log with ID %s does not exist", log.ID)
	}

	r.logs[log.ID] = log
	return nil
}

// GetLogByID retrieves an email log by ID
func (r *MemoryRepository) GetLogByID(id string) (*models.EmailLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	log, exists := r.logs[id]
	if !exists {
		return nil, fmt.Errorf("log with ID %s does not exist", id)
	}

	return log, nil
}

// GetLogsByStatus retrieves all email logs with the specified status
func (r *MemoryRepository) GetLogsByStatus(status models.EmailStatus) ([]*models.EmailLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var logs []*models.EmailLog
	for _, log := range r.logs {
		if log.Status == status {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

// GetLogsByDateRange retrieves all email logs within the specified date range
func (r *MemoryRepository) GetLogsByDateRange(start, end time.Time) ([]*models.EmailLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var logs []*models.EmailLog
	for _, log := range r.logs {
		if log.CreatedAt.After(start) && log.CreatedAt.Before(end) {
			logs = append(logs, log)
		}
	}

	return logs, nil
}

// GetEmailLogs retrieves all email logs
func (r *MemoryRepository) GetEmailLogs() ([]models.EmailLog, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	logs := make([]models.EmailLog, 0, len(r.logs))
	for _, log := range r.logs {
		logs = append(logs, *log)
	}

	return logs, nil
}

// SaveTemplate saves a new email template
func (r *MemoryRepository) SaveTemplate(template *models.Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[template.ID]; exists {
		return fmt.Errorf("template with ID %s already exists", template.ID)
	}

	r.templates[template.ID] = template
	return nil
}

// GetTemplate retrieves an email template by ID
func (r *MemoryRepository) GetTemplate(id string) (*models.Template, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.templates[id]
	if !exists {
		return nil, fmt.Errorf("template with ID %s does not exist", id)
	}

	return template, nil
}

// UpdateTemplate updates an existing email template
func (r *MemoryRepository) UpdateTemplate(template *models.Template) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[template.ID]; !exists {
		return fmt.Errorf("template with ID %s does not exist", template.ID)
	}

	r.templates[template.ID] = template
	return nil
}

// DeleteTemplate deletes an email template
func (r *MemoryRepository) DeleteTemplate(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.templates[id]; !exists {
		return fmt.Errorf("template with ID %s does not exist", id)
	}

	delete(r.templates, id)
	return nil
}

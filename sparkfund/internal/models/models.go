package models

import "time"

// SendEmailRequest represents a request to send an email
type SendEmailRequest struct {
	Recipients  []string          `json:"recipients"`
	From        string            `json:"from"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	ContentType string            `json:"content_type"`
	Attachments []EmailAttachment `json:"attachments,omitempty"`
	TemplateID  string            `json:"template_id,omitempty"`
	Data        map[string]string `json:"data,omitempty"`
}

// EmailAttachment represents an email attachment
type EmailAttachment struct {
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Content     []byte `json:"content"`
}

// Template represents an email template
type Template struct {
	ID          string            `json:"id" db:"id"`
	Name        string            `json:"name" db:"name"`
	Subject     string            `json:"subject" db:"subject"`
	Content     string            `json:"content" db:"content"`
	ContentType string            `json:"content_type" db:"content_type"`
	Data        map[string]string `json:"data,omitempty" db:"data"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}

// EmailLog represents a log entry for an email
type EmailLog struct {
	ID          string    `json:"id" db:"id"`
	Recipients  []string  `json:"recipients" db:"recipients"`
	From        string    `json:"from" db:"from"`
	Subject     string    `json:"subject" db:"subject"`
	Body        string    `json:"body" db:"body"`
	ContentType string    `json:"content_type" db:"content_type"`
	Status      string    `json:"status" db:"status"`
	Error       string    `json:"error,omitempty" db:"error"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// EmailStats represents email statistics
type EmailStats struct {
	TotalSent      int64 `json:"total_sent" db:"total_sent"`
	TotalFailed    int64 `json:"total_failed" db:"total_failed"`
	TotalPending   int64 `json:"total_pending" db:"total_pending"`
	AverageLatency int64 `json:"average_latency" db:"average_latency"`
}

package validation

import (
	"fmt"
	"net/mail"
	"strings"

	"github.com/adil-faiyaz98/sparkfund/email-service/internal/models"
)

// ValidationError represents a validation error with field and message
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidateEmailRequest validates an email request
func ValidateEmailRequest(req *models.SendEmailRequest) error {
	if req == nil {
		return &ValidationError{
			Field:   "request",
			Message: "request cannot be nil",
		}
	}

	if len(req.Recipients) == 0 {
		return &ValidationError{
			Field:   "recipients",
			Message: "at least one recipient is required",
		}
	}

	for i, recipient := range req.Recipients {
		if err := validateEmail(recipient); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: err.Error(),
			}
		}
	}

	if err := validateEmail(req.From); err != nil {
		return &ValidationError{
			Field:   "from",
			Message: err.Error(),
		}
	}

	if strings.TrimSpace(req.Subject) == "" {
		return &ValidationError{
			Field:   "subject",
			Message: "subject cannot be empty",
		}
	}

	if strings.TrimSpace(req.Body) == "" {
		return &ValidationError{
			Field:   "body",
			Message: "body cannot be empty",
		}
	}

	if req.ContentType != "text/plain" && req.ContentType != "text/html" {
		return &ValidationError{
			Field:   "content_type",
			Message: "content type must be either text/plain or text/html",
		}
	}

	for i, attachment := range req.Attachments {
		if err := validateAttachment(attachment); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("attachments[%d]", i),
				Message: err.Error(),
			}
		}
	}

	return nil
}

// ValidateTemplate validates an email template
func ValidateTemplate(template *models.Template) error {
	if template == nil {
		return &ValidationError{
			Field:   "template",
			Message: "template cannot be nil",
		}
	}

	if strings.TrimSpace(template.Name) == "" {
		return &ValidationError{
			Field:   "name",
			Message: "template name cannot be empty",
		}
	}

	if strings.TrimSpace(template.Subject) == "" {
		return &ValidationError{
			Field:   "subject",
			Message: "template subject cannot be empty",
		}
	}

	if strings.TrimSpace(template.Content) == "" {
		return &ValidationError{
			Field:   "content",
			Message: "template content cannot be empty",
		}
	}

	if template.ContentType != "text/plain" && template.ContentType != "text/html" {
		return &ValidationError{
			Field:   "content_type",
			Message: "content type must be either text/plain or text/html",
		}
	}

	return nil
}

// ValidateEmailLog validates an email log entry
func ValidateEmailLog(log *models.EmailLog) error {
	if log == nil {
		return &ValidationError{
			Field:   "log",
			Message: "log cannot be nil",
		}
	}

	if strings.TrimSpace(log.ID) == "" {
		return &ValidationError{
			Field:   "id",
			Message: "ID cannot be empty",
		}
	}

	if len(log.Recipients) == 0 {
		return &ValidationError{
			Field:   "recipients",
			Message: "at least one recipient is required",
		}
	}

	for i, recipient := range log.Recipients {
		if err := validateEmail(recipient); err != nil {
			return &ValidationError{
				Field:   fmt.Sprintf("recipients[%d]", i),
				Message: err.Error(),
			}
		}
	}

	if err := validateEmail(log.From); err != nil {
		return &ValidationError{
			Field:   "from",
			Message: err.Error(),
		}
	}

	if strings.TrimSpace(log.Subject) == "" {
		return &ValidationError{
			Field:   "subject",
			Message: "subject cannot be empty",
		}
	}

	if strings.TrimSpace(log.Body) == "" {
		return &ValidationError{
			Field:   "body",
			Message: "body cannot be empty",
		}
	}

	if log.ContentType != "text/plain" && log.ContentType != "text/html" {
		return &ValidationError{
			Field:   "content_type",
			Message: "content type must be either text/plain or text/html",
		}
	}

	if log.Status != "sent" && log.Status != "failed" && log.Status != "pending" {
		return &ValidationError{
			Field:   "status",
			Message: "status must be one of: sent, failed, pending",
		}
	}

	return nil
}

// validateEmail validates an email address
func validateEmail(email string) error {
	if strings.TrimSpace(email) == "" {
		return &ValidationError{
			Field:   "email",
			Message: "email cannot be empty",
		}
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return &ValidationError{
			Field:   "email",
			Message: "invalid email format",
		}
	}

	if addr.Address != email {
		return &ValidationError{
			Field:   "email",
			Message: "email address contains display name",
		}
	}

	return nil
}

// validateAttachment validates an email attachment
func validateAttachment(attachment []byte) error {
	if len(attachment) == 0 {
		return &ValidationError{
			Field:   "attachment",
			Message: "attachment content cannot be empty",
		}
	}

	return nil
}

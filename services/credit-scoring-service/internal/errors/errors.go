package errors

import (
	"fmt"
	"net/http"
)

// Error codes
const (
	ErrInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrUnauthorized       = "UNAUTHORIZED"
	ErrForbidden         = "FORBIDDEN"
	ErrNotFound          = "NOT_FOUND"
	ErrValidation        = "VALIDATION_ERROR"
	ErrRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
	ErrDatabase          = "DATABASE_ERROR"
	ErrExternalService   = "EXTERNAL_SERVICE_ERROR"
)

// APIError represents a structured API error response
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// NewAPIError creates a new APIError with the given code and message
func NewAPIError(code string, message string) *APIError {
	status := getStatusFromCode(code)
	return &APIError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// getStatusFromCode maps error codes to HTTP status codes
func getStatusFromCode(code string) int {
	switch code {
	case ErrInternalServer:
		return http.StatusInternalServerError
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrNotFound:
		return http.StatusNotFound
	case ErrValidation:
		return http.StatusBadRequest
	case ErrRateLimitExceeded:
		return http.StatusTooManyRequests
	case ErrDatabase:
		return http.StatusInternalServerError
	case ErrExternalService:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// ValidationError represents a field-level validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

func (e *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed: %d errors", len(e.Errors))
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *ValidationErrors {
	return &ValidationErrors{
		Errors: []ValidationError{
			{
				Field:   field,
				Message: message,
			},
		},
	}
}

// AddError adds a new validation error to the collection
func (e *ValidationErrors) AddError(field, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Field:   field,
		Message: message,
	})
}

// IsValidationError checks if an error is a ValidationErrors
func IsValidationError(err error) bool {
	_, ok := err.(*ValidationErrors)
	return ok
}

// DatabaseError represents a database operation error
type DatabaseError struct {
	Operation string
	Err       error
}

func (e *DatabaseError) Error() string {
	return fmt.Sprintf("database error during %s: %v", e.Operation, e.Errors)
}

// NewDatabaseError creates a new database error
func NewDatabaseError(operation string, err error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Err:       err,
	}
}

// ExternalServiceError represents an error from an external service
type ExternalServiceError struct {
	Service string
	Err     error
}

func (e *ExternalServiceError) Error() string {
	return fmt.Sprintf("external service error from %s: %v", e.Service, e.Errors)
}

// NewExternalServiceError creates a new external service error
func NewExternalServiceError(service string, err error) *ExternalServiceError {
	return &ExternalServiceError{
		Service: service,
		Err:     err,
	}
} 
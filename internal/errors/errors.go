package errors

import (
	"fmt"
	"net/http"
)

// Error types
const (
	ErrValidation     = "VALIDATION_ERROR"
	ErrNotFound       = "NOT_FOUND"
	ErrUnauthorized   = "UNAUTHORIZED"
	ErrForbidden      = "FORBIDDEN"
	ErrInternal       = "INTERNAL_ERROR"
	ErrBadRequest     = "BAD_REQUEST"
	ErrConflict       = "CONFLICT"
	ErrTooManyRequests = "TOO_MANY_REQUESTS"
)

// AppError represents a structured application error
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewAppError creates a new AppError
func NewAppError(code string, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap wraps an existing error with AppError
func Wrap(err error, code string, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// Common error constructors
func NewValidationError(message string) *AppError {
	return NewAppError(ErrValidation, message, http.StatusBadRequest)
}

func NewNotFoundError(message string) *AppError {
	return NewAppError(ErrNotFound, message, http.StatusNotFound)
}

func NewUnauthorizedError(message string) *AppError {
	return NewAppError(ErrUnauthorized, message, http.StatusUnauthorized)
}

func NewForbiddenError(message string) *AppError {
	return NewAppError(ErrForbidden, message, http.StatusForbidden)
}

func NewInternalError(message string) *AppError {
	return NewAppError(ErrInternal, message, http.StatusInternalServerError)
}

func NewBadRequestError(message string) *AppError {
	return NewAppError(ErrBadRequest, message, http.StatusBadRequest)
}

func NewConflictError(message string) *AppError {
	return NewAppError(ErrConflict, message, http.StatusConflict)
}

func NewTooManyRequestsError(message string) *AppError {
	return NewAppError(ErrTooManyRequests, message, http.StatusTooManyRequests)
}

// ErrorResponse represents the structure of error responses
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// HandleError handles application errors and returns appropriate HTTP responses
func HandleError(err error) (int, ErrorResponse) {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Status, ErrorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		}
	}

	// Default to internal server error for unknown errors
	return http.StatusInternalServerError, ErrorResponse{
		Code:    ErrInternal,
		Message: "An unexpected error occurred",
	}
} 
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Component identifies the component where the error occurred.
type Component string

const (
	ComponentAPI       Component = "API"
	ComponentService   Component = "Service"
	ComponentRepository Component = "Repository"
	ComponentKafka      Component = "Kafka"
	ComponentSMTP       Component = "SMTP"
)

// Error represents a custom error with additional context
type Error struct {
	Code      string
	Message   string
	Status    int
	Err       error
	Component Component // Add component field
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %s: %v", e.Component, e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s: %s",  e.Component, e.Code, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// Common error codes
const (
	ErrCodeValidation      = "VALIDATION_ERROR"
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeDuplicate       = "DUPLICATE"
	ErrCodeDatabase        = "DATABASE_ERROR"
	ErrCodeKafka           = "KAFKA_ERROR"
	ErrCodeSMTP            = "SMTP_ERROR"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeBadRequest      = "BAD_REQUEST"
	ErrCodeTooManyRequests = "TOO_MANY_REQUESTS"
	ErrCodeTransactionFailed = "TRANSACTION_FAILED" //Repository Specific
)

// Common errors (using factory functions)

func NewValidationError(err error) *Error {
	return &Error{
		Code:      ErrCodeValidation,
		Message:   "validation failed",
		Status:    http.StatusBadRequest,
		Err:       err,
		Component: ComponentAPI, //Or wherever validation happens
	}
}

func NewDatabaseError(err error, component Component) *Error {
	return &Error{
		Code:      ErrCodeDatabase,
		Message:   "database operation failed",
		Status:    http.StatusInternalServerError,
		Err:       err,
		Component: component,
	}
}

func NewKafkaError(err error) *Error {
	return &Error{
		Code:      ErrCodeKafka,
		Message:   "Kafka operation failed",
		Status:    http.StatusInternalServerError,
		Err:       err,
		Component: ComponentKafka,
	}
}

func NewSMTPError(err error) *Error {
	Code:      ErrCodeSMTP,
		Message:   "SMTP operation failed",
		Status:    http.StatusInternalServerError,
		Err:       err,
		Component: ComponentSMTP,
	}
}

func NewInternalError(err error) *Error {
	return &Error{
		Code:      ErrCodeInternal,
		Message:   "internal server error",
		Status:    http.StatusInternalServerError,
		Err:       err,
		Component: ComponentService, // Or wherever it originates
	}
}

func NewNotFoundError(message string) *Error {
    return &Error{
        Code:      ErrCodeNotFound,
        Message:   message,
        Status:    http.StatusNotFound,
        Component: ComponentRepository, // Or relevant component
    }
}

func NewDuplicateError(message string) *Error {
    return &Error{
        Code:      ErrCodeDuplicate,
        Message:   message,
        Status:    http.StatusConflict,
        Component: ComponentRepository, // Or relevant component
    }
}

func NewTransactionError(err error) *Error {
	return &Error{
		Code:      ErrCodeTransactionFailed,
		Message:   "transaction failed",
		Status:    http.StatusInternalServerError,
		Err:       err,
		Component: ComponentRepository,
	}
}

// Define common error instances using the factory functions
var (
	ErrInvalidInput = NewValidationError(errors.New("invalid input")) //Example
	ErrTemplateNotFound = NewNotFoundError("template not found")
	ErrTemplateExists = NewDuplicateError("template already exists")
)

// Use errors.Is and errors.As from the standard library
// Example usage:
// if errors.Is(err, ErrTemplateNotFound) { ... }
// var myError *Error
// if errors.As(err, &myError) { ... }
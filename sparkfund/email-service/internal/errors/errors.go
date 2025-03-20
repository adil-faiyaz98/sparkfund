package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// Error represents a custom error with additional context
type Error struct {
	Code    string
	Message string
	Status  int
	Err     error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
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
)

// Common errors
var (
	ErrInvalidInput = &Error{
		Code:    ErrCodeValidation,
		Message: "invalid input",
		Status:  http.StatusBadRequest,
	}

	ErrNotFound = &Error{
		Code:    ErrCodeNotFound,
		Message: "resource not found",
		Status:  http.StatusNotFound,
	}

	ErrDuplicate = &Error{
		Code:    ErrCodeDuplicate,
		Message: "resource already exists",
		Status:  http.StatusConflict,
	}

	ErrDatabase = &Error{
		Code:    ErrCodeDatabase,
		Message: "database operation failed",
		Status:  http.StatusInternalServerError,
	}

	ErrKafka = &Error{
		Code:    ErrCodeKafka,
		Message: "Kafka operation failed",
		Status:  http.StatusInternalServerError,
	}

	ErrSMTP = &Error{
		Code:    ErrCodeSMTP,
		Message: "SMTP operation failed",
		Status:  http.StatusInternalServerError,
	}

	ErrInternal = &Error{
		Code:    ErrCodeInternal,
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
	}

	ErrUnauthorized = &Error{
		Code:    ErrCodeUnauthorized,
		Message: "unauthorized",
		Status:  http.StatusUnauthorized,
	}

	ErrForbidden = &Error{
		Code:    ErrCodeForbidden,
		Message: "forbidden",
		Status:  http.StatusForbidden,
	}

	ErrTooManyRequests = &Error{
		Code:    ErrCodeTooManyRequests,
		Message: "too many requests",
		Status:  http.StatusTooManyRequests,
	}
)

// NewValidationError creates a new validation error
func NewValidationError(err error) *Error {
	return &Error{
		Code:    ErrCodeValidation,
		Message: "validation failed",
		Status:  http.StatusBadRequest,
		Err:     err,
	}
}

// NewDatabaseError creates a new database error
func NewDatabaseError(err error) *Error {
	return &Error{
		Code:    ErrCodeDatabase,
		Message: "database operation failed",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// NewKafkaError creates a new Kafka error
func NewKafkaError(err error) *Error {
	return &Error{
		Code:    ErrCodeKafka,
		Message: "Kafka operation failed",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// NewSMTPError creates a new SMTP error
func NewSMTPError(err error) *Error {
	return &Error{
		Code:    ErrCodeSMTP,
		Message: "SMTP operation failed",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// NewInternalError creates a new internal error
func NewInternalError(err error) *Error {
	return &Error{
		Code:    ErrCodeInternal,
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// IsValidationError checks if the error is a validation error
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeValidation
	}
	return false
}

// IsDatabaseError checks if the error is a database error
func IsDatabaseError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeDatabase
	}
	return false
}

// IsKafkaError checks if the error is a Kafka error
func IsKafkaError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeKafka
	}
	return false
}

// IsSMTPError checks if the error is an SMTP error
func IsSMTPError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeSMTP
	}
	return false
}

// IsInternalError checks if the error is an internal error
func IsInternalError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeInternal
	}
	return false
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeNotFound
	}
	return false
}

// IsDuplicate checks if the error is a duplicate error
func IsDuplicate(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrCodeDuplicate
	}
	return false
}

// As is a helper function to check if an error can be converted to a target type
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// RepositoryError represents a base error type for repository operations
type RepositoryError struct {
	Code    string
	Message string
	Err     error
}

func (e *RepositoryError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("repository error [%s]: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("repository error [%s]: %s", e.Code, e.Message)
}

// Common repository error codes
const (
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeDuplicate         = "DUPLICATE"
	ErrCodeInvalidInput      = "INVALID_INPUT"
	ErrCodeDatabaseError     = "DATABASE_ERROR"
	ErrCodeTransactionFailed = "TRANSACTION_FAILED"
)

// Common repository errors
var (
	ErrTemplateNotFound = &RepositoryError{
		Code:    ErrCodeNotFound,
		Message: "template not found",
	}

	ErrEmailLogNotFound = &RepositoryError{
		Code:    ErrCodeNotFound,
		Message: "email log not found",
	}

	ErrTemplateExists = &RepositoryError{
		Code:    ErrCodeDuplicate,
		Message: "template already exists",
	}

	ErrEmailLogExists = &RepositoryError{
		Code:    ErrCodeDuplicate,
		Message: "email log already exists",
	}

	ErrInvalidID = &RepositoryError{
		Code:    ErrCodeInvalidInput,
		Message: "invalid ID format",
	}

	ErrEmptyContent = &RepositoryError{
		Code:    ErrCodeInvalidInput,
		Message: "content cannot be empty",
	}

	ErrInvalidStatus = &RepositoryError{
		Code:    ErrCodeInvalidInput,
		Message: "invalid email status",
	}
)

// NewDatabaseError creates a new database error
func NewDatabaseError(err error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrCodeDatabaseError,
		Message: "database operation failed",
		Err:     err,
	}
}

// NewTransactionError creates a new transaction error
func NewTransactionError(err error) *RepositoryError {
	return &RepositoryError{
		Code:    ErrCodeTransactionFailed,
		Message: "transaction failed",
		Err:     err,
	}
}

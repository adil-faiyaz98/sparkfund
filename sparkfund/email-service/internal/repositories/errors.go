package repositories

import "fmt"

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

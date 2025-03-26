package errors

import "fmt"

// Error represents a custom application error
type Error struct {
	Status  int    // HTTP status code
	Message string // Error message
}

// Error implements the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("status: %d, message: %s", e.Status, e.Message)
}

// NewError creates a new Error instance
func NewError(status int, message string) *Error {
	return &Error{
		Status:  status,
		Message: message,
	}
}
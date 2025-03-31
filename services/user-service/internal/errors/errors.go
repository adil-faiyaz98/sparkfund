package errors

import (
	"fmt"
	"net/http"
)

// Error represents a custom application error
type Error struct {
	Code    int
	Message string
	Err     error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap implements the unwrap interface
func (e *Error) Unwrap() error {
	return e.Err
}

// Common error types
var (
	// Authentication errors
	ErrInvalidCredentials = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Invalid credentials",
	}
	ErrAccountLocked = &Error{
		Code:    http.StatusForbidden,
		Message: "Account is locked",
	}
	ErrSessionExpired = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Session has expired",
	}
	ErrInvalidToken = &Error{
		Code:    http.StatusUnauthorized,
		Message: "Invalid authentication token",
	}

	// Authorization errors
	ErrInsufficientPermissions = &Error{
		Code:    http.StatusForbidden,
		Message: "Insufficient permissions",
	}
	ErrAccessDenied = &Error{
		Code:    http.StatusForbidden,
		Message: "Access denied",
	}

	// Validation errors
	ErrInvalidInput = &Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid input",
	}
	ErrPasswordTooWeak = &Error{
		Code:    http.StatusBadRequest,
		Message: "Password does not meet security requirements",
	}
	ErrInvalidEmail = &Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid email format",
	}
	ErrInvalidPhoneNumber = &Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid phone number format",
	}

	// Resource errors
	ErrUserNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "User not found",
	}
	ErrProfileNotFound = &Error{
		Code:    http.StatusNotFound,
		Message: "User profile not found",
	}
	ErrUserAlreadyExists = &Error{
		Code:    http.StatusConflict,
		Message: "User already exists",
	}
	ErrEmailAlreadyInUse = &Error{
		Code:    http.StatusConflict,
		Message: "Email is already in use",
	}

	// Security errors
	ErrTooManyLoginAttempts = &Error{
		Code:    http.StatusTooManyRequests,
		Message: "Too many login attempts",
	}
	ErrPasswordResetExpired = &Error{
		Code:    http.StatusBadRequest,
		Message: "Password reset token has expired",
	}
	ErrInvalidResetToken = &Error{
		Code:    http.StatusBadRequest,
		Message: "Invalid password reset token",
	}
	ErrResetTokenUsed = &Error{
		Code:    http.StatusBadRequest,
		Message: "Password reset token has already been used",
	}

	// System errors
	ErrInternalServer = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Internal server error",
	}
	ErrDatabase = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Database error",
	}
	ErrEmailService = &Error{
		Code:    http.StatusInternalServerError,
		Message: "Email service error",
	}
)

// Wrap wraps an error with additional context
func Wrap(err error, message string) *Error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*Error); ok {
		return &Error{
			Code:    e.Code,
			Message: message,
			Err:     e,
		}
	}

	return &Error{
		Code:    http.StatusInternalServerError,
		Message: message,
		Err:     err,
	}
}

// IsNotFound checks if an error is a not found error
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusNotFound
	}
	return false
}

// IsValidationError checks if an error is a validation error
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusBadRequest
	}
	return false
}

// IsAuthenticationError checks if an error is an authentication error
func IsAuthenticationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusUnauthorized
	}
	return false
}

// IsAuthorizationError checks if an error is an authorization error
func IsAuthorizationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == http.StatusForbidden
	}
	return false
}

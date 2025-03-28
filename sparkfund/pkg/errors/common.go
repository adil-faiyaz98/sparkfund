package errors

import "net/http"

// Common error types
var (
	ErrBadRequest = func(err error) *AppError {
		return NewAppError(http.StatusBadRequest, "Bad request", err)
	}

	ErrUnauthorized = func(err error) *AppError {
		return NewAppError(http.StatusUnauthorized, "Unauthorized", err)
	}

	ErrForbidden = func(err error) *AppError {
		return NewAppError(http.StatusForbidden, "Forbidden", err)
	}

	ErrNotFound = func(err error) *AppError {
		return NewAppError(http.StatusNotFound, "Not found", err)
	}

	ErrConflict = func(err error) *AppError {
		return NewAppError(http.StatusConflict, "Conflict", err)
	}

	ErrInternalServer = func(err error) *AppError {
		return NewAppError(http.StatusInternalServerError, "Internal server error", err)
	}

	ErrServiceUnavailable = func(err error) *AppError {
		return NewAppError(http.StatusServiceUnavailable, "Service unavailable", err)
	}
)

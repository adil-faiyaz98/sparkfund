package errors

import "net/http"

// IsAppError checks if an error is an AppError
func IsAppError(err error) bool {
	_, ok := err.(*AppError)
	return ok
}

// GetStatusCode gets the status code from an error
func GetStatusCode(err error) int {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Code
	}
	return http.StatusInternalServerError
}

// GetMessage gets the message from an error
func GetMessage(err error) string {
	if appErr, ok := err.(*AppError); ok {
		return appErr.Message
	}
	return "Internal server error"
}

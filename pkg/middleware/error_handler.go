package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Common error types
var (
	ErrNotFound          = errors.New("resource not found")
	ErrBadRequest        = errors.New("bad request")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrInternalServer    = errors.New("internal server error")
	ErrServiceUnavailable = errors.New("service unavailable")
)

// ErrorHandlerConfig holds configuration for error handling
type ErrorHandlerConfig struct {
	Logger *logrus.Logger
}

// ErrorHandler middleware handles errors consistently
func ErrorHandler(cfg ErrorHandlerConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Get error and log appropriately
			statusCode, response := handleError(err, cfg.Logger)

			// Return error response
			c.JSON(statusCode, response)
		}
	}
}

// handleError maps errors to status codes and response messages
func handleError(err error, logger *logrus.Logger) (int, ErrorResponse) {
	var statusCode int
	var message string

	// Check for known errors
	switch {
	case errors.Is(err, ErrNotFound):
		statusCode = http.StatusNotFound
		message = "Resource not found"

	case errors.Is(err, ErrBadRequest):
		statusCode = http.StatusBadRequest
		message = "Bad request"

	case errors.Is(err, ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		message = "Unauthorized"

	case errors.Is(err, ErrForbidden):
		statusCode = http.StatusForbidden
		message = "Forbidden"

	case errors.Is(err, ErrServiceUnavailable):
		statusCode = http.StatusServiceUnavailable
		message = "Service unavailable"

	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"

		// Log unexpected errors
		if logger != nil {
			logger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Unexpected error")
		}
	}

	return statusCode, ErrorResponse{
		Error: message,
	}
}

// RequestLogger logs request details
func RequestLogger(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log request
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		path := c.Request.URL.Path
		userID, _ := c.Get("userID")

		// Log details
		if statusCode >= 400 {
			// Log errors with more detail
			logger.WithFields(logrus.Fields{
				"status":    statusCode,
				"latency":   latency,
				"client_ip": clientIP,
				"method":    method,
				"path":      path,
				"user":      userID,
			}).Error("Request error")
		} else {
			// Log success requests more briefly
			logger.WithFields(logrus.Fields{
				"status":    statusCode,
				"latency":   latency,
				"method":    method,
				"path":      path,
			}).Info("Request processed")
		}
	}
}

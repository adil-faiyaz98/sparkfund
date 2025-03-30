package middleware

import (
	"errors"
	"net/http"

	"investment-service/internal/models"
	"investment-service/internal/services"
	"investment-service/internal/validation"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandler middleware handles errors consistently
func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			// Get error and log appropriately
			statusCode, response := handleError(err, logger)

			// Return error response
			c.JSON(statusCode, response)
		}
	}
}

// handleError maps errors to status codes and response messages
func handleError(err error, logger *logrus.Logger) (int, models.ErrorResponse) {
	var statusCode int
	var message string

	// Check for known errors
	switch {
	case errors.Is(err, services.ErrInvestmentNotFound):
		statusCode = http.StatusNotFound
		message = "Investment not found"

	case errors.Is(err, services.ErrInsufficientQuantity):
		statusCode = http.StatusBadRequest
		message = "Insufficient quantity for transaction"

	case errors.Is(err, validation.ErrMissingUserID),
		errors.Is(err, validation.ErrInvalidType),
		errors.Is(err, validation.ErrInvalidStatus),
		errors.Is(err, validation.ErrNonPositiveQuantity),
		errors.Is(err, validation.ErrNonPositivePrice),
		errors.Is(err, validation.ErrMissingSymbol):
		statusCode = http.StatusBadRequest
		message = err.Error()

	default:
		statusCode = http.StatusInternalServerError
		message = "Internal server error"

		// Log unexpected errors
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unexpected error")
	}

	return statusCode, models.ErrorResponse{
		Error: message,
	}
}

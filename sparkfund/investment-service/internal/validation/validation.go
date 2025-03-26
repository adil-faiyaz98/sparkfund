package validation

import (
	"fmt"
	"strconv"

	"github.com/sparkfund/investment-service/internal/models"
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: %s - %s", e.Field, e.Message)
}

func ValidatePortfolio(portfolio *models.Portfolio) error {
	if portfolio == nil {
		return &ValidationError{
			Field:   "portfolio",
			Message: "portfolio cannot be nil",
		}
	}

	if portfolio.ClientId <= 0 {
		return &ValidationError{
			Field:   "clientId",
			Message: "clientId must be a positive integer",
		}
	}

	if portfolio.Name == "" {
		return &ValidationError{
			Field:   "name",
			Message: "name cannot be empty",
		}
	}

	return nil
}

func ValidateClientId(clientId string) error {
	if clientId == "" {
		return &ValidationError{
			Field:   "clientId",
			Message: "clientId cannot be empty",
		}
	}

	// Basic check to ensure clientId is a number
	_, err := strconv.Atoi(clientId)
	if err != nil {
		return &ValidationError{
			Field:   "clientId",
			Message: "clientId must be a number",
		}
	}

	return nil
}

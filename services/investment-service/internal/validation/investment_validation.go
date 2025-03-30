package validation

import (
	"errors"
	"strings"

	"investment-service/internal/models"
)

// Validation errors
var (
	ErrMissingUserID       = errors.New("user ID is required")
	ErrInvalidType         = errors.New("invalid investment type")
	ErrInvalidStatus       = errors.New("invalid investment status")
	ErrNonPositiveQuantity = errors.New("quantity must be positive")
	ErrNonPositivePrice    = errors.New("price must be positive")
	ErrMissingSymbol       = errors.New("symbol is required")
)

// Valid investment types
var validTypes = map[string]bool{
	"STOCK":       true,
	"BOND":        true,
	"ETF":         true,
	"MUTUAL_FUND": true,
	"CRYPTO":      true,
	"REAL_ESTATE": true,
	"OTHER":       true,
}

// Valid investment statuses
var validStatuses = map[string]bool{
	"ACTIVE":  true,
	"SOLD":    true,
	"PENDING": true,
}

// ValidateInvestment validates investment data
func ValidateInvestment(investment *models.Investment) error {
	// Check required fields
	if investment.UserID == 0 {
		return ErrMissingUserID
	}

	if investment.Symbol == "" {
		return ErrMissingSymbol
	}

	// Normalize type and status to uppercase
	investment.Type = strings.ToUpper(investment.Type)
	investment.Status = strings.ToUpper(investment.Status)

	// Validate type
	if !validTypes[investment.Type] {
		return ErrInvalidType
	}

	// Validate status
	if investment.Status != "" && !validStatuses[investment.Status] {
		return ErrInvalidStatus
	}

	// Validate numeric fields
	if investment.Quantity <= 0 {
		return ErrNonPositiveQuantity
	}

	if investment.PurchasePrice <= 0 {
		return ErrNonPositivePrice
	}

	return nil
}

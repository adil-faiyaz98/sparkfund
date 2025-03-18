package docs

import (
	"time"

	"github.com/google/uuid"
)

// @title           Money Pulse Investments Service
// @version         1.0
// @description     A microservice for managing investment portfolios in the Money Pulse application.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// Investment represents a financial investment
// @Description Investment details including type, status, and pricing information
type Investment struct {
	// @Description Unique identifier for the investment
	// @example 123e4567-e89b-12d3-a456-426614174000
	ID uuid.UUID `json:"id"`

	// @Description ID of the user who owns this investment
	// @example 123e4567-e89b-12d3-a456-426614174001
	UserID uuid.UUID `json:"user_id"`

	// @Description ID of the account this investment belongs to
	// @example 123e4567-e89b-12d3-a456-426614174002
	AccountID uuid.UUID `json:"account_id"`

	// @Description Type of investment (STOCK, BOND, MUTUAL_FUND, ETF, CRYPTO)
	// @example STOCK
	Type string `json:"type"`

	// @Description Current status of the investment
	// @example ACTIVE
	Status string `json:"status"`

	// @Description Stock symbol or investment identifier
	// @example AAPL
	Symbol string `json:"symbol"`

	// @Description Number of units held
	// @example 10.5
	Quantity float64 `json:"quantity"`

	// @Description Price per unit at purchase
	// @example 150.25
	PurchasePrice float64 `json:"purchase_price"`

	// @Description Current market price per unit
	// @example 155.75
	CurrentPrice float64 `json:"current_price"`

	// @Description Currency code (e.g., USD, EUR)
	// @example USD
	Currency string `json:"currency"`

	// @Description Date when the investment was purchased
	// @example 2024-01-15T10:30:00Z
	PurchaseDate time.Time `json:"purchase_date"`

	// @Description Last time the investment data was updated
	// @example 2024-01-16T15:45:00Z
	LastUpdated time.Time `json:"last_updated"`

	// @Description Creation timestamp
	// @example 2024-01-15T10:30:00Z
	CreatedAt time.Time `json:"created_at"`

	// @Description Last update timestamp
	// @example 2024-01-16T15:45:00Z
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateInvestmentRequest represents the request body for creating a new investment
// @Description Request body for creating a new investment
type CreateInvestmentRequest struct {
	// @Description ID of the user who owns this investment
	// @example 123e4567-e89b-12d3-a456-426614174001
	UserID uuid.UUID `json:"user_id" binding:"required"`

	// @Description ID of the account this investment belongs to
	// @example 123e4567-e89b-12d3-a456-426614174002
	AccountID uuid.UUID `json:"account_id" binding:"required"`

	// @Description Type of investment (STOCK, BOND, MUTUAL_FUND, ETF, CRYPTO)
	// @example STOCK
	Type string `json:"type" binding:"required"`

	// @Description Stock symbol or investment identifier
	// @example AAPL
	Symbol string `json:"symbol" binding:"required"`

	// @Description Number of units to purchase
	// @example 10.5
	Quantity float64 `json:"quantity" binding:"required"`

	// @Description Price per unit
	// @example 150.25
	PurchasePrice float64 `json:"purchase_price" binding:"required"`

	// @Description Currency code (e.g., USD, EUR)
	// @example USD
	Currency string `json:"currency" binding:"required"`
}

// UpdateInvestmentRequest represents the request body for updating an investment
// @Description Request body for updating an existing investment
type UpdateInvestmentRequest struct {
	// @Description Type of investment (STOCK, BOND, MUTUAL_FUND, ETF, CRYPTO)
	// @example STOCK
	Type string `json:"type"`

	// @Description Current status of the investment
	// @example ACTIVE
	Status string `json:"status"`

	// @Description Stock symbol or investment identifier
	// @example AAPL
	Symbol string `json:"symbol"`

	// @Description Number of units held
	// @example 10.5
	Quantity float64 `json:"quantity"`

	// @Description Current market price per unit
	// @example 155.75
	CurrentPrice float64 `json:"current_price"`
}

// ErrorResponse represents a standard error response
// @Description Standard error response structure
type ErrorResponse struct {
	// @Description Error message
	// @example Invalid investment type
	Error string `json:"error"`
}

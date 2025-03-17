package model

import "time"

// Account represents a financial account in the system
type Account struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	AccountNumber string    `json:"account_number"`
	AccountType   string    `json:"account_type"`
	Balance       float64   `json:"balance"`
	Currency      string    `json:"currency"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AccountTypes defines the supported account types
var AccountTypes = struct {
	Checking   string
	Savings    string
	Credit     string
	Investment string
}{
	Checking:   "checking",
	Savings:    "savings",
	Credit:     "credit",
	Investment: "investment",
}

// Currencies defines the supported currencies
var Currencies = struct {
	USD string
	EUR string
	GBP string
}{
	USD: "USD",
	EUR: "EUR",
	GBP: "GBP",
}

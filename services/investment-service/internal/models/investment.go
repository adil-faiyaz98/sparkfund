package models

import (
	"time"
)

// Investment represents an investment made by a user
type Investment struct {
	ID            uint       `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UserID        uint       `gorm:"not null" json:"user_id"`
	PortfolioID   uint       `json:"portfolio_id"` // Add this field to match with the foreignKey in Portfolio
	Amount        float64    `gorm:"not null" json:"amount"`
	Type          string     `gorm:"not null" json:"type" example:"STOCK"`    // e.g., "STOCK", "CRYPTO", "REAL_ESTATE"
	Status        string     `gorm:"not null" json:"status" example:"ACTIVE"` // e.g., "ACTIVE", "SOLD", "PENDING"
	PurchaseDate  time.Time  `gorm:"not null" json:"purchase_date"`
	SellDate      *time.Time `json:"sell_date,omitempty"`
	PurchasePrice float64    `gorm:"not null" json:"purchase_price"`
	SellPrice     *float64   `json:"sell_price,omitempty"`
	Symbol        string     `gorm:"not null" json:"symbol" example:"AAPL"` // e.g., "AAPL", "BTC", "ETH"
	Quantity      float64    `gorm:"not null" json:"quantity"`
	Notes         string     `json:"notes,omitempty"`
}

// Transaction represents a transaction related to an investment
type Transaction struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	UserID        uint      `gorm:"not null" json:"user_id"`
	InvestmentID  uint      `gorm:"not null" json:"investment_id"`
	Type          string    `gorm:"not null" json:"type" example:"BUY"` // e.g., "BUY", "SELL"
	Amount        float64   `gorm:"not null" json:"amount"`
	Price         float64   `gorm:"not null" json:"price"`
	Quantity      float64   `gorm:"not null" json:"quantity"`
	Timestamp     time.Time `gorm:"not null" json:"timestamp"`
	Status        string    `gorm:"not null" json:"status" example:"COMPLETED"` // e.g., "COMPLETED", "PENDING", "FAILED"
	TransactionID string    `gorm:"unique;not null" json:"transaction_id"`
}

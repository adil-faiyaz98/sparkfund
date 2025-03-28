package models

import (
	"time"

	"gorm.io/gorm"
)

type Investment struct {
	gorm.Model
	UserID        uint      `gorm:"not null"`
	PortfolioID   uint      // Add this field to match with the foreignKey in Portfolio
	Amount        float64   `gorm:"not null"`
	Type          string    `gorm:"not null"` // e.g., "STOCK", "CRYPTO", "REAL_ESTATE"
	Status        string    `gorm:"not null"` // e.g., "ACTIVE", "SOLD", "PENDING"
	PurchaseDate  time.Time `gorm:"not null"`
	SellDate      *time.Time
	PurchasePrice float64 `gorm:"not null"`
	SellPrice     *float64
	Symbol        string  `gorm:"not null"` // e.g., "AAPL", "BTC", "ETH"
	Quantity      float64 `gorm:"not null"`
	Notes         string
}

type Transaction struct {
	gorm.Model
	UserID        uint      `gorm:"not null"`
	InvestmentID  uint      `gorm:"not null"`
	Type          string    `gorm:"not null"` // e.g., "BUY", "SELL"
	Amount        float64   `gorm:"not null"`
	Price         float64   `gorm:"not null"`
	Quantity      float64   `gorm:"not null"`
	Timestamp     time.Time `gorm:"not null"`
	Status        string    `gorm:"not null"` // e.g., "COMPLETED", "PENDING", "FAILED"
	TransactionID string    `gorm:"unique;not null"`
}

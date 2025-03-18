package models

import (
	"time"

	"gorm.io/gorm"
)

// Account represents a financial account
type Account struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID         uint    `json:"user_id" gorm:"not null;index"`
	Name           string  `json:"name" gorm:"not null"`
	Type           string  `json:"type" gorm:"not null"` // checking, savings, credit, investment, etc.
	InitialBalance float64 `json:"initial_balance"`
	CurrentBalance float64 `json:"current_balance"`
	Currency       string  `json:"currency" gorm:"not null"`

	// Relations
	Transactions []Transaction `json:"-" gorm:"foreignKey:AccountID"`
}

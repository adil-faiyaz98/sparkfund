package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"index;not null" json:"user_id"`
	CategoryID  uint           `gorm:"index" json:"category_id"`
	AccountID   uint           `gorm:"index;not null" json:"account_id"`
	Amount      float64        `gorm:"not null" json:"amount"`
	Description string         `json:"description"`
	Date        time.Time      `gorm:"index;not null" json:"date"`
	Type        string         `gorm:"not null" json:"type"` // "income" or "expense"
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID" json:"-"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Account  Account  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
}

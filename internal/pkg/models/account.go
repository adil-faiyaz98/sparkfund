package models

import (
	"time"

	"gorm.io/gorm"
)

// Account represents a financial account
type Account struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	Type      string         `gorm:"not null" json:"type"` // checking, savings, credit, cash
	Balance   float64        `gorm:"not null" json:"balance"`
	Currency  string         `gorm:"not null" json:"currency"`
	Color     string         `json:"color"`
	Icon      string         `json:"icon"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User         User          `gorm:"foreignKey:UserID" json:"-"`
	Transactions []Transaction `json:"-"`
}

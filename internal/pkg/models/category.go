package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a transaction category
type Category struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	Name      string         `gorm:"not null" json:"name"`
	Color     string         `json:"color"`
	Icon      string         `json:"icon"`
	Type      string         `gorm:"not null" json:"type"` // "income" or "expense"
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User         User          `gorm:"foreignKey:UserID" json:"-"`
	Transactions []Transaction `json:"-"`
}

package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a transaction category
type Category struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	UserID   uint   `json:"user_id"` // Zero for system categories
	Name     string `json:"name" gorm:"not null"`
	Color    string `json:"color"`
	Icon     string `json:"icon"`
	IsSystem bool   `json:"is_system" gorm:"default:false"` // True for default system categories

	// Relations
	Transactions []Transaction `json:"-" gorm:"foreignKey:CategoryID"`
}

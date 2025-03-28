package models

import (
	"time"
)

// Portfolio represents a user's investment portfolio
type Portfolio struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description,omitempty"`
	TotalValue  float64   `gorm:"not null;default:0" json:"total_value"`
	LastUpdated time.Time `gorm:"not null" json:"last_updated"`
}

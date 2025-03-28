package models

import (
	"time"

	"gorm.io/gorm"
)

type Portfolio struct {
	gorm.Model
	UserID      uint   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Description string
	Investments []Investment `gorm:"foreignKey:PortfolioID"`
	TotalValue  float64      `gorm:"not null;default:0"`
	LastUpdated time.Time    `gorm:"not null"`
}

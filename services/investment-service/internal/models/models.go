package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Investment struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	Type            string           `gorm:"not null" json:"type"`
	Amount          float64          `gorm:"not null" json:"amount"`
	Currency        string           `gorm:"not null" json:"currency"`
	RiskLevel       string           `gorm:"not null" json:"risk_level"`
	Status          string           `gorm:"not null" json:"status"`
	PortfolioID     uuid.UUID        `gorm:"type:uuid" json:"portfolio_id"`
	AssetID         string           `gorm:"not null" json:"asset_id"`
	AssetName       string           `gorm:"not null" json:"asset_name"`
	AssetType       string           `gorm:"not null" json:"asset_type"`
	AssetSymbol     string           `gorm:"not null" json:"asset_symbol"`
	AssetPrice      float64          `gorm:"not null" json:"asset_price"`
	AssetQuantity   float64          `gorm:"not null" json:"asset_quantity"`
	TransactionFee  float64          `gorm:"not null" json:"transaction_fee"`
	TransactionType string           `gorm:"not null" json:"transaction_type"`
	CreatedAt       time.Time        `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
}

type Portfolio struct {
	ID            uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	Name          string           `gorm:"not null" json:"name"`
	Description   string           `gorm:"type:text" json:"description"`
	RiskLevel     string           `gorm:"not null" json:"risk_level"`
	TotalValue    float64          `gorm:"not null" json:"total_value"`
	Currency      string           `gorm:"not null" json:"currency"`
	CreatedAt     time.Time        `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time        `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt   `gorm:"index" json:"-"`
	Investments   []Investment     `gorm:"foreignKey:PortfolioID" json:"investments"`
}

type Asset struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Symbol      string         `gorm:"not null;unique" json:"symbol"`
	Type        string         `gorm:"not null" json:"type"`
	Description string         `gorm:"type:text" json:"description"`
	RiskLevel   string         `gorm:"not null" json:"risk_level"`
	Currency    string         `gorm:"not null" json:"currency"`
	Price       float64        `gorm:"not null" json:"price"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Transaction struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	InvestmentID    uuid.UUID      `gorm:"type:uuid;not null" json:"investment_id"`
	Type            string         `gorm:"not null" json:"type"`
	Amount          float64        `gorm:"not null" json:"amount"`
	Currency        string         `gorm:"not null" json:"currency"`
	Status          string         `gorm:"not null" json:"status"`
	TransactionFee  float64        `gorm:"not null" json:"transaction_fee"`
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserPreference struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;unique" json:"user_id"`
	RiskTolerance string         `gorm:"not null" json:"risk_tolerance"`
	Currency      string         `gorm:"not null" json:"currency"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
} 
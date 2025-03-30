package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	ID              uuid.UUID        `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID        `gorm:"type:uuid;not null" json:"user_id"`
	Amount          float64          `gorm:"not null" json:"amount"`
	Currency        string           `gorm:"not null" json:"currency"`
	Type            string           `gorm:"not null" json:"type"`
	Status          string           `gorm:"not null" json:"status"`
	RiskScore       float64          `gorm:"not null" json:"risk_score"`
	RiskLevel       string           `gorm:"not null" json:"risk_level"`
	RiskFactors     []RiskFactor     `gorm:"foreignKey:TransactionID" json:"risk_factors"`
	Alerts          []Alert          `gorm:"foreignKey:TransactionID" json:"alerts"`
	ScreeningResult *ScreeningResult `gorm:"foreignKey:TransactionID" json:"screening_result"`
	CreatedAt       time.Time        `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt   `gorm:"index" json:"-"`
}

type RiskFactor struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID      `gorm:"type:uuid;not null" json:"transaction_id"`
	Type          string         `gorm:"not null" json:"type"`
	Description   string         `gorm:"not null" json:"description"`
	Weight        float64        `gorm:"not null" json:"weight"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type Alert struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID      `gorm:"type:uuid;not null" json:"transaction_id"`
	Type          string         `gorm:"not null" json:"type"`
	Severity      string         `gorm:"not null" json:"severity"`
	Description   string         `gorm:"not null" json:"description"`
	Status        string         `gorm:"not null" json:"status"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type ScreeningResult struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID      `gorm:"type:uuid;not null;unique" json:"transaction_id"`
	SanctionList  bool           `gorm:"not null" json:"sanction_list"`
	PEPList       bool           `gorm:"not null" json:"pep_list"`
	WatchList     bool           `gorm:"not null" json:"watch_list"`
	Details       string         `gorm:"type:text" json:"details"`
	CreatedAt     time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type RiskProfile struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;unique" json:"user_id"`
	RiskScore   float64        `gorm:"not null" json:"risk_score"`
	RiskLevel   string         `gorm:"not null" json:"risk_level"`
	LastUpdated time.Time      `gorm:"not null" json:"last_updated"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

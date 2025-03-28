package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
	TransactionTypeInvestment TransactionType = "investment"
)

type TransactionStatus string

const (
	TransactionStatusPending   TransactionStatus = "pending"
	TransactionStatusApproved  TransactionStatus = "approved"
	TransactionStatusRejected  TransactionStatus = "rejected"
	TransactionStatusFlagged   TransactionStatus = "flagged"
	TransactionStatusCompleted TransactionStatus = "completed"
)

type RiskLevel string

const (
	RiskLevelLow     RiskLevel = "low"
	RiskLevelMedium  RiskLevel = "medium"
	RiskLevelHigh    RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

type Transaction struct {
	ID              uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID         `gorm:"type:uuid;not null"`
	Type            TransactionType   `gorm:"type:varchar(20);not null"`
	Amount          float64           `gorm:"not null"`
	Currency        string            `gorm:"type:varchar(3);not null"`
	Status          TransactionStatus `gorm:"type:varchar(20);not null;default:'pending'"`
	RiskLevel       RiskLevel         `gorm:"type:varchar(20);not null;default:'low'"`
	SourceAccount   string            `gorm:"not null"`
	DestinationAccount string         `gorm:"not null"`
	Description     string            `gorm:"type:text"`
	IPAddress       string            `gorm:"type:varchar(45)"`
	DeviceID        string            `gorm:"type:varchar(100)"`
	Location        string            `gorm:"type:varchar(100)"`
	Metadata        string            `gorm:"type:jsonb"`
	FlaggedBy       *string           `gorm:"type:varchar(100)"`
	FlagReason      *string           `gorm:"type:text"`
	ReviewedBy      *string           `gorm:"type:varchar(100)"`
	ReviewNotes     *string           `gorm:"type:text"`
	CreatedAt       time.Time         `gorm:"not null"`
	UpdatedAt       time.Time         `gorm:"not null"`
	DeletedAt       gorm.DeletedAt    `gorm:"index"`
}

type TransactionRequest struct {
	Type              TransactionType `json:"type" binding:"required,oneof=deposit withdrawal transfer investment"`
	Amount            float64         `json:"amount" binding:"required,gt=0"`
	Currency          string          `json:"currency" binding:"required,len=3"`
	SourceAccount     string          `json:"sourceAccount" binding:"required"`
	DestinationAccount string         `json:"destinationAccount" binding:"required"`
	Description       string          `json:"description"`
	IPAddress         string          `json:"ipAddress"`
	DeviceID          string          `json:"deviceId"`
	Location          string          `json:"location"`
	Metadata          map[string]interface{} `json:"metadata"`
}

type TransactionResponse struct {
	ID                uuid.UUID         `json:"id"`
	UserID            uuid.UUID         `json:"userId"`
	Type              TransactionType   `json:"type"`
	Amount            float64           `json:"amount"`
	Currency          string            `json:"currency"`
	Status            TransactionStatus `json:"status"`
	RiskLevel         RiskLevel         `json:"riskLevel"`
	SourceAccount     string            `json:"sourceAccount"`
	DestinationAccount string           `json:"destinationAccount"`
	Description       string            `json:"description"`
	IPAddress         string            `json:"ipAddress"`
	DeviceID          string            `json:"deviceId"`
	Location          string            `json:"location"`
	Metadata          map[string]interface{} `json:"metadata"`
	FlaggedBy         *string           `json:"flaggedBy,omitempty"`
	FlagReason        *string           `json:"flagReason,omitempty"`
	ReviewedBy        *string           `json:"reviewedBy,omitempty"`
	ReviewNotes       *string           `json:"reviewNotes,omitempty"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

type RiskAssessment struct {
	TransactionID uuid.UUID `json:"transactionId"`
	RiskLevel     RiskLevel `json:"riskLevel"`
	RiskScore     float64   `json:"riskScore"`
	Factors       []string  `json:"factors"`
	Recommendation string   `json:"recommendation"`
}

type TransactionFilter struct {
	UserID            *uuid.UUID        `json:"userId,omitempty"`
	Type              *TransactionType  `json:"type,omitempty"`
	Status            *TransactionStatus `json:"status,omitempty"`
	RiskLevel         *RiskLevel        `json:"riskLevel,omitempty"`
	StartDate         *time.Time        `json:"startDate,omitempty"`
	EndDate           *time.Time        `json:"endDate,omitempty"`
	MinAmount         *float64          `json:"minAmount,omitempty"`
	MaxAmount         *float64          `json:"maxAmount,omitempty"`
	Currency          *string           `json:"currency,omitempty"`
	FlaggedOnly       *bool             `json:"flaggedOnly,omitempty"`
} 
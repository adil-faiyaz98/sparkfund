package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KYCStatus string

const (
	KYCStatusPending  KYCStatus = "pending"
	KYCStatusVerified KYCStatus = "verified"
	KYCStatusRejected KYCStatus = "rejected"
	KYCStatusExpired  KYCStatus = "expired"
)

type KYC struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID          uuid.UUID      `gorm:"type:uuid;not null"`
	FirstName       string         `gorm:"not null"`
	LastName        string         `gorm:"not null"`
	DateOfBirth     string         `gorm:"not null"` // Format: YYYY-MM-DD
	Address         string         `gorm:"not null"`
	City            string         `gorm:"not null"`
	Country         string         `gorm:"not null"`
	PostalCode      string         `gorm:"not null"`
	DocumentType    string         `gorm:"not null"`
	DocumentNumber  string         `gorm:"not null"`
	DocumentFront   string         `gorm:"not null"`
	DocumentBack    string         `gorm:"not null"`
	SelfieImage     string         `gorm:"not null"`
	Status          KYCStatus      `gorm:"type:varchar(20);not null;default:'pending'"`
	RejectionReason string         `gorm:"type:text;default:null"`
	VerifiedBy      *uuid.UUID     `gorm:"type:uuid;default:null"`
	VerifiedAt      *time.Time     `gorm:"default:null"`
	CreatedAt       time.Time      `gorm:"not null"`
	UpdatedAt       time.Time      `gorm:"not null"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type KYCRequest struct {
	FirstName      string `json:"firstName" binding:"required"`
	LastName       string `json:"lastName" binding:"required"`
	DateOfBirth    string `json:"dateOfBirth" binding:"required"` // Format: YYYY-MM-DD
	Address        string `json:"address" binding:"required"`
	City           string `json:"city" binding:"required"`
	Country        string `json:"country" binding:"required"`
	PostalCode     string `json:"postalCode" binding:"required"`
	DocumentType   string `json:"documentType" binding:"required"`
	DocumentNumber string `json:"documentNumber" binding:"required"`
	DocumentFront  string `json:"documentFront" binding:"required"`
	DocumentBack   string `json:"documentBack" binding:"required"`
	SelfieImage    string `json:"selfieImage" binding:"required"`
}

type KYCResponse struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"userId"`
	Status          KYCStatus  `json:"status"`
	RejectionReason string     `json:"rejectionReason,omitempty"`
	VerifiedAt      *time.Time `json:"verifiedAt,omitempty"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

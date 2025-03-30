package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type KYCVerification struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Status         string         `gorm:"not null" json:"status"`        // PENDING, APPROVED, REJECTED
	DocumentType   string         `gorm:"not null" json:"document_type"` // PASSPORT, DRIVERS_LICENSE, NATIONAL_ID
	DocumentNumber string         `gorm:"not null" json:"document_number"`
	DocumentURL    string         `gorm:"not null" json:"document_url"`
	FirstName      string         `gorm:"not null" json:"first_name"`
	LastName       string         `gorm:"not null" json:"last_name"`
	DateOfBirth    time.Time      `gorm:"not null" json:"date_of_birth"`
	Address        string         `gorm:"not null" json:"address"`
	City           string         `gorm:"not null" json:"city"`
	Country        string         `gorm:"not null" json:"country"`
	PostalCode     string         `gorm:"not null" json:"postal_code"`
	PhoneNumber    string         `gorm:"not null" json:"phone_number"`
	Email          string         `gorm:"not null" json:"email"`
	RiskLevel      string         `gorm:"not null" json:"risk_level"` // LOW, MEDIUM, HIGH
	RiskScore      int            `gorm:"not null" json:"risk_score"`
	Notes          string         `gorm:"type:text" json:"notes"`
	ReviewedBy     uuid.UUID      `gorm:"type:uuid" json:"reviewed_by"`
	ReviewedAt     time.Time      `json:"reviewed_at"`
	CreatedAt      time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type KYCReview struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VerificationID  uuid.UUID      `gorm:"type:uuid;not null" json:"verification_id"`
	ReviewerID      uuid.UUID      `gorm:"type:uuid;not null" json:"reviewer_id"`
	Status          string         `gorm:"not null" json:"status"` // APPROVED, REJECTED
	Reason          string         `gorm:"type:text" json:"reason"`
	RiskAssessment  string         `gorm:"type:text" json:"risk_assessment"`
	AdditionalNotes string         `gorm:"type:text" json:"additional_notes"`
	CreatedAt       time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type KYCDocument struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VerificationID     uuid.UUID      `gorm:"type:uuid;not null" json:"verification_id"`
	DocumentType       string         `gorm:"not null" json:"document_type"`
	DocumentNumber     string         `gorm:"not null" json:"document_number"`
	DocumentURL        string         `gorm:"not null" json:"document_url"`
	VerificationStatus string         `gorm:"not null" json:"verification_status"` // PENDING, VERIFIED, FAILED
	VerificationNotes  string         `gorm:"type:text" json:"verification_notes"`
	CreatedAt          time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

type KYCWatchlist struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Reason     string         `gorm:"type:text;not null" json:"reason"`
	RiskLevel  string         `gorm:"not null" json:"risk_level"` // LOW, MEDIUM, HIGH
	ExpiryDate time.Time      `json:"expiry_date"`
	CreatedAt  time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

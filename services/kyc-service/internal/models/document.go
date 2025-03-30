package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Document represents a KYC document
type Document struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Type       DocumentType   `gorm:"type:varchar(50);not null" json:"type"`
	Number     string         `gorm:"not null" json:"number"`
	URL        string         `gorm:"not null" json:"url"`
	IssueDate  time.Time      `gorm:"not null" json:"issue_date"`
	ExpiryDate time.Time      `gorm:"not null" json:"expiry_date"`
	Status     DocumentStatus `gorm:"type:varchar(20);not null" json:"status"`
	CreatedAt  time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

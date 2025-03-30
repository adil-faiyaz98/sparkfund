package models

import (
	"time"

	"gorm.io/gorm"
)

// KYC represents a Know Your Customer record
type KYC struct {
	gorm.Model
	UserID          string    `gorm:"type:uuid;not null;uniqueIndex"`
	Status          string    `gorm:"type:varchar(20);not null;default:'pending'"`
	DocumentType    string    `gorm:"type:varchar(50);not null"`
	DocumentNumber  string    `gorm:"type:varchar(100);not null"`
	DocumentURL     string    `gorm:"type:text;not null"`
	FirstName       string    `gorm:"type:varchar(100);not null"`
	LastName        string    `gorm:"type:varchar(100);not null"`
	DateOfBirth     time.Time `gorm:"type:date;not null"`
	Address         string    `gorm:"type:text;not null"`
	City            string    `gorm:"type:varchar(100);not null"`
	Country         string    `gorm:"type:varchar(100);not null"`
	PostalCode      string    `gorm:"type:varchar(20);not null"`
	PhoneNumber     string    `gorm:"type:varchar(20);not null"`
	Email           string    `gorm:"type:varchar(255);not null"`
	RejectionReason string    `gorm:"type:text"`
	ReviewedBy      string    `gorm:"type:uuid"`
	ReviewedAt      time.Time
	Notes           string `gorm:"type:text"`
}

// TableName specifies the table name for the KYC model
func (KYC) TableName() string {
	return "kyc_records"
}

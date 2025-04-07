package repository

import (
	"gorm.io/gorm"
)

// Repositories holds all repository instances
type Repositories struct {
	Document     *DocumentRepository
	KYC          *KYCRepository
	Verification *VerificationRepository
}

// NewRepositories creates a new Repositories instance
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Document:     NewDocumentRepository(db),
		KYC:          NewKYCRepository(db),
		Verification: NewVerificationRepository(db),
	}
}

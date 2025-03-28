package database

import (
	"github.com/sparkfund/kyc-service/internal/model"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}

	// Run migrations
	return db.AutoMigrate(
		&model.KYC{},
	)
}

// RollbackMigrations rolls back all database migrations
func RollbackMigrations(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&model.KYC{},
	)
}

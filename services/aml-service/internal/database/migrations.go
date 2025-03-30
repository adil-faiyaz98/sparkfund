package database

import (
	"aml-service/internal/model"

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
		&model.Transaction{},
	)
}

// RollbackMigrations rolls back all database migrations
func RollbackMigrations(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&model.Transaction{},
	)
}

// CreateIndexes creates necessary indexes for better query performance
func CreateIndexes(db *gorm.DB) error {
	// Create indexes for transactions
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
		CREATE INDEX IF NOT EXISTS idx_transactions_status ON transactions(status);
		CREATE INDEX IF NOT EXISTS idx_transactions_risk_level ON transactions(risk_level);
		CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions(created_at);
		CREATE INDEX IF NOT EXISTS idx_transactions_amount ON transactions(amount);
		CREATE INDEX IF NOT EXISTS idx_transactions_currency ON transactions(currency);
		CREATE INDEX IF NOT EXISTS idx_transactions_type ON transactions(type);
	`).Error; err != nil {
		return err
	}

	return nil
}

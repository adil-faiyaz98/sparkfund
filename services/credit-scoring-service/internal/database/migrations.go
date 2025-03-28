package database

import (
	"github.com/sparkfund/credit-scoring-service/internal/model"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}

	// Run migrations
	if err := db.AutoMigrate(
		&model.CreditScore{},
		&model.CreditHistory{},
	); err != nil {
		return err
	}

	// Create indexes
	return CreateIndexes(db)
}

func RollbackMigrations(db *gorm.DB) error {
	// Drop tables
	if err := db.Migrator().DropTable(
		&model.CreditScore{},
		&model.CreditHistory{},
	); err != nil {
		return err
	}

	return nil
}

func CreateIndexes(db *gorm.DB) error {
	// Credit Score indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_credit_scores_user_id ON credit_scores(user_id);
		CREATE INDEX IF NOT EXISTS idx_credit_scores_score_range ON credit_scores(score_range);
		CREATE INDEX IF NOT EXISTS idx_credit_scores_last_updated ON credit_scores(last_updated);
	`).Error; err != nil {
		return err
	}

	// Credit History indexes
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_credit_histories_user_id ON credit_histories(user_id);
		CREATE INDEX IF NOT EXISTS idx_credit_histories_account_type ON credit_histories(account_type);
		CREATE INDEX IF NOT EXISTS idx_credit_histories_status ON credit_histories(status);
		CREATE INDEX IF NOT EXISTS idx_credit_histories_institution ON credit_histories(institution);
		CREATE INDEX IF NOT EXISTS idx_credit_histories_open_date ON credit_histories(open_date);
		CREATE INDEX IF NOT EXISTS idx_credit_histories_close_date ON credit_histories(close_date);
		CREATE INDEX IF NOT EXISTS idx_credit_histories_last_payment_date ON credit_histories(last_payment_date);
	`).Error; err != nil {
		return err
	}

	return nil
}

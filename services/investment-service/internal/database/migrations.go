package database

import (
	"fmt"
	"time"

	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/sparkfund/services/investment-service/internal/config"
	"github.com/sparkfund/services/investment-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// RunMigrations runs all database migrations in order
func RunMigrations(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, []*gormigrate.Migration{
		{
			ID: "202503281200",
			Migrate: func(tx *gorm.DB) error {
				// Create portfolio table
				err := tx.AutoMigrate(&models.Portfolio{})
				if err != nil {
					return err
				}

				// Seed some default portfolios for testing
				if tx.Migrator().HasTable(&models.Portfolio{}) {
					// Only seed in development environment
					if config.Get().Environment != "production" {
						defaultPortfolios := []models.Portfolio{
							{
								UserID:      1,
								Name:        "Retirement",
								Description: "Long-term retirement investments",
							},
							{
								UserID:      1,
								Name:        "Growth",
								Description: "High-growth technology investments",
							},
						}

						for _, portfolio := range defaultPortfolios {
							if err := tx.Create(&portfolio).Error; err != nil {
								return err
							}
						}
					}
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("portfolios")
			},
		},
		{
			ID: "202503281201",
			Migrate: func(tx *gorm.DB) error {
				// Create investment table
				return tx.AutoMigrate(&models.Investment{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("investments")
			},
		},
		{
			ID: "202503281202",
			Migrate: func(tx *gorm.DB) error {
				// Create transaction table
				return tx.AutoMigrate(&models.Transaction{})
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable("transactions")
			},
		},
		{
			ID: "202503281203",
			Migrate: func(tx *gorm.DB) error {
				// Add indexes for common queries
				if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_investments_user_id ON investments(user_id)").Error; err != nil {
					return err
				}
				if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_investments_portfolio_id ON investments(portfolio_id)").Error; err != nil {
					return err
				}
				if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_investments_symbol ON investments(symbol)").Error; err != nil {
					return err
				}
				if err := tx.Exec("CREATE INDEX IF NOT EXISTS idx_transactions_investment_id ON transactions(investment_id)").Error; err != nil {
					return err
				}
				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				if err := tx.Exec("DROP INDEX IF EXISTS idx_investments_user_id").Error; err != nil {
					return err
				}
				if err := tx.Exec("DROP INDEX IF EXISTS idx_investments_portfolio_id").Error; err != nil {
					return err
				}
				if err := tx.Exec("DROP INDEX IF EXISTS idx_investments_symbol").Error; err != nil {
					return err
				}
				if err := tx.Exec("DROP INDEX IF EXISTS idx_transactions_investment_id").Error; err != nil {
					return err
				}
				return nil
			},
		},
	})

	return m.Migrate()
}

// InitDBMigrations initializes the database and runs migrations
func InitDBMigrations() error {
	if err := InitDB(); err != nil {
		return err
	}

	return RunMigrations(DB)
}

// InitDB initializes the database connection
func InitDB() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Get().Database.Host,
		config.Get().Database.Port,
		config.Get().Database.User,
		config.Get().Database.Password,
		config.Get().Database.Name,
		config.Get().Database.SSLMode,
	)

	// Setup DB connection with retry
	var db *gorm.DB
	var err error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			break
		}
		time.Sleep(time.Second * 2)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after %d retries: %w", maxRetries, err)
	}

	DB = db
	return nil
}

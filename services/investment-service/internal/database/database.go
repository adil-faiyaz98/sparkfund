package database

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"investment-service/internal/config"
	"investment-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database connection
var DB *gorm.DB

// MaxRetries is the maximum number of database connection attempts
const MaxRetries = 5

// RetryTimeout is the timeout for each retry attempt
const RetryTimeout = 5 * time.Second

// InitDB initializes the database connection
func InitDB() error {
	cfg := config.Get()

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	// Configure GORM logger
	gormLogger := logger.New(
		log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level (change to Warn in production)
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error
			Colorful:                  false,       // Disable color in production
		},
	)

	// Configure GORM
	config := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC() // Use UTC for all timestamps
		},
		PrepareStmt:                              true, // Cache prepared statements for better performance
		DisableForeignKeyConstraintWhenMigrating: false,
	}

	// Connect to database with retry mechanism
	var db *gorm.DB
	var err error

	for i := 0; i < MaxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), config)
		if err == nil {
			break
		}

		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, MaxRetries, err)

		// Add jitter to prevent thundering herd problem
		jitter := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(RetryTimeout + jitter)
	}

	if err != nil {
		return fmt.Errorf("failed to connect to database after %d attempts: %w", MaxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)                  // Maximum idle connections
	sqlDB.SetMaxOpenConns(100)                 // Maximum open connections
	sqlDB.SetConnMaxLifetime(30 * time.Minute) // Connection reuse timeout
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // How long connections can remain idle

	// Set global DB variable
	DB = db
	log.Println("Database connected successfully")

	return nil
}

// InitDBMigrations initializes the database and runs migrations
func InitDBMigrations() error {
	// First initialize DB connection
	if err := InitDB(); err != nil {
		return err
	}

	// Then run migrations
	if err := RunMigrations(DB); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Migrate performs database migrations using AutoMigrate
func Migrate(db *gorm.DB) error {
	// First migrate Portfolio as it's referenced by Investment
	if err := db.AutoMigrate(&models.Portfolio{}); err != nil {
		return fmt.Errorf("failed to migrate Portfolio model: %w", err)
	}

	// Then migrate Investment as it's referenced by Transaction
	if err := db.AutoMigrate(&models.Investment{}); err != nil {
		return fmt.Errorf("failed to migrate Investment model: %w", err)
	}

	// Finally migrate Transaction
	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		return fmt.Errorf("failed to migrate Transaction model: %w", err)
	}

	return nil
}

// WithTransaction executes operations within a database transaction
func WithTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

package database

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config holds database configuration
type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns default database configuration
func DefaultConfig() Config {
	return Config{
		Host:            "localhost",
		Port:            "5432",
		User:            "postgres",
		Password:        "postgres",
		Name:            "kyc_service",
		SSLMode:         "disable",
		MaxIdleConns:    10,
		MaxOpenConns:    100,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

// MaxRetries is the maximum number of database connection attempts
const MaxRetries = 5

// RetryTimeout is the timeout for each retry attempt
const RetryTimeout = 5 * time.Second

// Connect initializes the database connection
func Connect(cfg Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SSLMode,
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
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", MaxRetries, err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	log.Println("Database connected successfully")

	return db, nil
}

// WithTransaction executes operations within a database transaction
func WithTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.WithContext(ctx).Begin()
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

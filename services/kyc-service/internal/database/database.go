package database

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"kyc-service/internal/config"

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

package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sparkfund/aml-service/internal/database"
	"github.com/sparkfund/aml-service/internal/handler"
	"github.com/sparkfund/aml-service/internal/middleware"
	"github.com/sparkfund/aml-service/internal/repository"
	"github.com/sparkfund/aml-service/internal/service"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	// DefaultDatabaseURL is a default connection string for local development. This is insecure and should not be used in production.
	DefaultDatabaseURL = "host=localhost user=postgres password=postgres dbname=sparkfund port=5432 sslmode=disable"
)
type Config struct {
	Port        string
	DatabaseURL string
}

func loadConfig() Config {
	config := Config{
		Port:        os.Getenv("PORT"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}
	if config.DatabaseURL == "" {
		config.DatabaseURL = DefaultDatabaseURL
	}

	return config
}

func initLogger() (*zap.Logger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	return logger, nil
}

func main() {
	config := loadConfig()

	logger, err := initLogger()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info("Starting aml-service")
	sugar.Infof("Loaded config: %+v", config)
	
	// Initialize database
	db, err := initDB(config.DatabaseURL)
	if err != nil{
		sugar.Fatalf("Failed to initialize database: %v", err)
	}
	sugar.Info("Database initialized successfully")

	// Run migration
	if err := database.RunMigrations(db); err != nil {
		sugar.Fatalf("Failed to run migrations: %v", err)
	}

	// Create indexes
	if err := database.CreateIndexes(db); err != nil {
		log.Fatalf("Failed to create indexes: %v", err)
	}	

	// Initialize dependencies
	txRepo := repository.NewTransactionRepository(db)
	amlService := service.NewAMLService(txRepo)
	amlHandler := handler.NewAMLHandler(amlService)

	// Initialize Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Recovery())
	router.Use(middleware.RateLimit(100, time.Minute, 10)) // 100 requests per minute, burst of 10
	router.Use(middleware.AuthMiddleware())                // JWT authentication

	// Register routes
	amlHandler.RegisterRoutes(router)

	// Start server
	if err := router.Run(":" + config.Port); err != nil {
		sugar.Fatalf("Failed to start server: %v", err)
	}
}

func initDB(databaseURL string) (db *gorm.DB, err error) {
	maxRetries := 3
	retryDelay := 2 * time.Second
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
		if err == nil {
			break // Connection successful
		}

		if errors.Is(err, fmt.Errorf("failed to connect to `host=localhost`")) {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		fmt.Printf("Failed to connect to database, retrying in %v... (attempt %d/%d)\n", retryDelay, i+1, maxRetries)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}	
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

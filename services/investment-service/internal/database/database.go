package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"investment-service/internal/models"
)

var DB *gorm.DB

func InitDB() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := Migrate(db); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func Migrate(db *gorm.DB) error {
	// First migrate Portfolio as it's referenced by Investment
	if err := db.AutoMigrate(&models.Portfolio{}); err != nil {
		return err
	}

	// Then migrate Investment as it's referenced by Transaction
	if err := db.AutoMigrate(&models.Investment{}); err != nil {
		return err
	}

	// Finally migrate Transaction
	return db.AutoMigrate(&models.Transaction{})
}

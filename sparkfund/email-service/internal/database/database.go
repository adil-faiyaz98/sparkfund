package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"github.com/sparkfund/email-service/internal/config"
)

// NewDB creates a new database connection.
func NewDB(cfg config.DatabaseConfig) *sqlx.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Connected to database")
	return db
}

// CloseDB closes the database connection
func CloseDB(db *sqlx.DB) error {
	if db != nil {
		return db.Close()
	}
	return nil
}

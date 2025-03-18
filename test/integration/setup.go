package integration

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
	// Import your database driver
	// _ "github.com/lib/pq"
)

// TestDB provides a database connection for integration tests
var TestDB *sql.DB

// SetupTestDatabase sets up a test database connection
func SetupTestDatabase() (*sql.DB, error) {
	// Get database connection details from environment variables or use defaults
	dbHost := getEnv("TEST_DB_HOST", "localhost")
	dbPort := getEnv("TEST_DB_PORT", "5432")
	dbUser := getEnv("TEST_DB_USER", "postgres")
	dbPass := getEnv("TEST_DB_PASS", "postgres")
	dbName := getEnv("TEST_DB_NAME", "money_pulse_test")

	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping test database: %w", err)
	}

	// Set connection pool parameters
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// CleanupTestDatabase performs cleanup operations after tests
func CleanupTestDatabase(db *sql.DB) {
	if db != nil {
		// Here you could truncate tables if needed
		// Example: db.Exec("TRUNCATE accounts, transactions, users CASCADE")

		db.Close()
	}
}

// TestMain sets up the test environment
func TestMain(m *testing.M) {
	// Skip setup if running in short mode (unit tests only)
	if testing.Short() {
		log.Println("Running in short mode, skipping integration test setup")
		os.Exit(m.Run())
	}

	// Setup test database
	var err error
	TestDB, err = SetupTestDatabase()
	if err != nil {
		log.Printf("Error setting up test database: %v", err)
		os.Exit(1)
	}

	// Run tests
	code := m.Run()

	// Cleanup
	CleanupTestDatabase(TestDB)

	os.Exit(code)
}

// Helper function to get environment variable with fallback
func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

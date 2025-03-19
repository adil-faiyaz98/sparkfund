package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestDB represents a test database connection
type TestDB struct {
	DB   *gorm.DB
	SQL  *sql.DB
	Name string
}

// NewTestDB creates a new test database connection
func NewTestDB(t *testing.T) *TestDB {
	t.Helper()

	// Get database connection details from environment variables
	host := getEnvOrDefault("TEST_DB_HOST", "localhost")
	port := getEnvOrDefault("TEST_DB_PORT", "5432")
	user := getEnvOrDefault("TEST_DB_USER", "postgres")
	password := getEnvOrDefault("TEST_DB_PASSWORD", "postgres")
	dbName := fmt.Sprintf("test_%s", uuid.New().String())

	// Create the database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		host, port, user, password)
	sqlDB, err := sql.Open("postgres", dsn)
	require.NoError(t, err)

	_, err = sqlDB.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	require.NoError(t, err)

	// Connect to the test database
	dsn = fmt.Sprintf("%s dbname=%s", dsn, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	return &TestDB{
		DB:   db,
		SQL:  sqlDB,
		Name: dbName,
	}
}

// Close closes the test database connection and drops the database
func (tdb *TestDB) Close(t *testing.T) {
	t.Helper()

	// Close the database connection
	sqlDB, err := tdb.DB.DB()
	require.NoError(t, err)
	require.NoError(t, sqlDB.Close())

	// Drop the test database
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		getEnvOrDefault("TEST_DB_HOST", "localhost"),
		getEnvOrDefault("TEST_DB_PORT", "5432"),
		getEnvOrDefault("TEST_DB_USER", "postgres"),
		getEnvOrDefault("TEST_DB_PASSWORD", "postgres"))

	sqlDB, err = sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer sqlDB.Close()

	_, err = sqlDB.Exec(fmt.Sprintf("DROP DATABASE %s", tdb.Name))
	require.NoError(t, err)
}

// CleanupTables cleans up all tables in the test database
func (tdb *TestDB) CleanupTables(t *testing.T) {
	t.Helper()

	// Get all table names
	var tables []string
	err := tdb.DB.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&tables).Error
	require.NoError(t, err)

	// Truncate each table
	for _, table := range tables {
		err := tdb.DB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error
		require.NoError(t, err)
	}
}

// GetEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// CreateTestContext creates a test context with timeout
func CreateTestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)
	return ctx
}

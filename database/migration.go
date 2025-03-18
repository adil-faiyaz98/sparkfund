package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jinzhu/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID        string
	UpSQL     string
	DownSQL   string
	Service   string
	Timestamp int64
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *gorm.DB) *MigrationManager {
	return &MigrationManager{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// AddMigration adds a migration to the manager
func (m *MigrationManager) AddMigration(service string, upSQL, downSQL string) error {
	// Extract timestamp from filename
	timestamp := extractTimestamp(filepath.Base(upSQL))
	if timestamp == 0 {
		return fmt.Errorf("invalid migration filename: %s", upSQL)
	}

	// Read migration files
	upContent, err := ioutil.ReadFile(upSQL)
	if err != nil {
		return fmt.Errorf("failed to read up migration: %w", err)
	}

	downContent, err := ioutil.ReadFile(downSQL)
	if err != nil {
		return fmt.Errorf("failed to read down migration: %w", err)
	}

	migration := Migration{
		ID:        fmt.Sprintf("%s_%d", service, timestamp),
		UpSQL:     string(upContent),
		DownSQL:   string(downContent),
		Service:   service,
		Timestamp: timestamp,
	}

	m.migrations = append(m.migrations, migration)
	return nil
}

// Migrate runs all pending migrations
func (m *MigrationManager) Migrate() error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Sort migrations by timestamp
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Timestamp < m.migrations[j].Timestamp
	})

	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if !applied[migration.ID] {
			log.Printf("Running migration: %s", migration.ID)
			if err := m.runMigration(migration); err != nil {
				return fmt.Errorf("failed to run migration %s: %w", migration.ID, err)
			}
		}
	}

	return nil
}

// Rollback rolls back the last migration
func (m *MigrationManager) Rollback() error {
	// Get applied migrations
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Find the last applied migration
	var lastMigration Migration
	for _, migration := range m.migrations {
		if applied[migration.ID] {
			lastMigration = migration
		}
	}

	if lastMigration.ID == "" {
		return fmt.Errorf("no migrations to rollback")
	}

	// Run down migration
	log.Printf("Rolling back migration: %s", lastMigration.ID)
	if err := m.runDownMigration(lastMigration); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastMigration.ID, err)
	}

	return nil
}

// createMigrationsTable creates the migrations table if it doesn't exist
func (m *MigrationManager) createMigrationsTable() error {
	return m.db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id VARCHAR(255) PRIMARY KEY,
			service VARCHAR(255) NOT NULL,
			timestamp BIGINT NOT NULL,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`).Error
}

// getAppliedMigrations returns a map of applied migration IDs
func (m *MigrationManager) getAppliedMigrations() (map[string]bool, error) {
	var applied []struct {
		ID string
	}
	if err := m.db.Table("migrations").Select("id").Find(&applied).Error; err != nil {
		return nil, err
	}

	result := make(map[string]bool)
	for _, migration := range applied {
		result[migration.ID] = true
	}
	return result, nil
}

// runMigration runs a migration
func (m *MigrationManager) runMigration(migration Migration) error {
	// Start transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Run up migration
	if err := tx.Exec(migration.UpSQL).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Record migration
	if err := tx.Exec(
		"INSERT INTO migrations (id, service, timestamp) VALUES (?, ?, ?)",
		migration.ID, migration.Service, migration.Timestamp,
	).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// runDownMigration runs a down migration
func (m *MigrationManager) runDownMigration(migration Migration) error {
	// Start transaction
	tx := m.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Run down migration
	if err := tx.Exec(migration.DownSQL).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Remove migration record
	if err := tx.Exec("DELETE FROM migrations WHERE id = ?", migration.ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// extractTimestamp extracts timestamp from migration filename
func extractTimestamp(filename string) int64 {
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return 0
	}

	var timestamp int64
	fmt.Sscanf(parts[0], "%d", &timestamp)
	return timestamp
}

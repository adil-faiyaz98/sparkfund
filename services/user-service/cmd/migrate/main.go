package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

var (
	up   = flag.Bool("up", false, "Run migrations up")
	down = flag.Bool("down", false, "Run migrations down")
)

func main() {
	flag.Parse()

	if !*up && !*down {
		log.Fatal("Please specify either -up or -down flag")
	}

	if *up && *down {
		log.Fatal("Cannot specify both -up and -down flags")
	}

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Get migrations directory
	migrationsDir := "migrations"
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("Migrations directory not found: %s", migrationsDir)
	}

	// Read migration files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Filter and sort migration files
	var migrationFiles []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	if *up {
		// Run migrations up
		for _, file := range migrationFiles {
			if strings.HasSuffix(file, ".up.sql") {
				if err := runMigration(db, filepath.Join(migrationsDir, file)); err != nil {
					log.Fatalf("Failed to run migration %s: %v", file, err)
				}
				log.Printf("Successfully ran migration: %s", file)
			}
		}
	} else {
		// Run migrations down in reverse order
		for i := len(migrationFiles) - 1; i >= 0; i-- {
			file := migrationFiles[i]
			if strings.HasSuffix(file, ".down.sql") {
				if err := runMigration(db, filepath.Join(migrationsDir, file)); err != nil {
					log.Fatalf("Failed to run migration %s: %v", file, err)
				}
				log.Printf("Successfully ran migration: %s", file)
			}
		}
	}
}

func runMigration(db *sql.DB, filepath string) error {
	// Read migration file
	content, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	return nil
}

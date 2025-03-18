package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/adil-faiyaz98/structgen/database"
)

func main() {
	// Parse command line flags
	action := flag.String("action", "up", "Migration action (up/down)")
	service := flag.String("service", "", "Service name (users/accounts/transactions/investments/loans/reports)")
	flag.Parse()

	if *service == "" {
		log.Fatal("Service name is required")
	}

	// Get database URL from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := database.ConnectFromURL(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Create migration manager
	mgr := database.NewMigrationManager(db)

	// Find migration files
	migrationsDir := filepath.Join("internal", *service, "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		log.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Add migrations to manager
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".sql") {
			continue
		}

		if strings.HasSuffix(filename, ".up.sql") {
			downFile := strings.Replace(filename, ".up.sql", ".down.sql", 1)
			upPath := filepath.Join(migrationsDir, filename)
			downPath := filepath.Join(migrationsDir, downFile)

			if err := mgr.AddMigration(*service, upPath, downPath); err != nil {
				log.Printf("Warning: Failed to add migration %s: %v", filename, err)
			}
		}
	}

	// Run migration
	switch *action {
	case "up":
		if err := mgr.Migrate(); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations completed successfully")
	case "down":
		if err := mgr.Rollback(); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Println("Rollback completed successfully")
	default:
		log.Fatalf("Invalid action: %s", *action)
	}
}

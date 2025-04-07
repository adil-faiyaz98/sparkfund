package database

import (
	"log"

	"gorm.io/gorm"
)

// RunMigrations runs all database migrations
func RunMigrations(db *gorm.DB) error {
	// Enable UUID extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		return err
	}

	// Enable pgcrypto extension
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return err
	}

	// Run migrations for existing tables
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS kycs (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`).Error; err != nil {
		return err
	}

	// Create additional tables for enhanced KYC service
	log.Println("Creating enhanced KYC tables...")

	// Create users table
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) NOT NULL UNIQUE,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		role VARCHAR(50) NOT NULL DEFAULT 'user',
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`).Error; err != nil {
		return err
	}

	// Create selfies table
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS selfies (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		file_url VARCHAR(255) NOT NULL,
		file_name VARCHAR(255) NOT NULL,
		file_size BIGINT NOT NULL,
		content_type VARCHAR(100) NOT NULL,
		status VARCHAR(50) NOT NULL DEFAULT 'PENDING',
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`).Error; err != nil {
		return err
	}

	// Create document_analyses table
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS document_analyses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
		document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
		document_type VARCHAR(50) NOT NULL,
		is_authentic BOOLEAN NOT NULL,
		confidence FLOAT NOT NULL,
		extracted_data JSONB NOT NULL,
		issues TEXT[] DEFAULT '{}',
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`).Error; err != nil {
		return err
	}

	// Create face_matches table
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS face_matches (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
		document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
		selfie_id UUID NOT NULL REFERENCES selfies(id) ON DELETE CASCADE,
		is_match BOOLEAN NOT NULL,
		confidence FLOAT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`).Error; err != nil {
		return err
	}

	// Create risk_analyses table
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS risk_analyses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		risk_score FLOAT NOT NULL,
		risk_level VARCHAR(50) NOT NULL,
		risk_factors TEXT[] DEFAULT '{}',
		device_info JSONB NOT NULL,
		ip_address VARCHAR(50) NOT NULL,
		location VARCHAR(255),
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`).Error; err != nil {
		return err
	}

	// Create anomaly_detections table
	if err := db.Exec(`CREATE TABLE IF NOT EXISTS anomaly_detections (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
		user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		is_anomaly BOOLEAN NOT NULL,
		anomaly_score FLOAT NOT NULL,
		anomaly_type VARCHAR(100),
		reasons TEXT[] DEFAULT '{}',
		device_info JSONB NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);`).Error; err != nil {
		return err
	}

	// Create indexes
	log.Println("Creating indexes...")

	// Create index on status in verifications
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_verifications_status ON verifications (status);`).Error; err != nil {
		return err
	}

	// Create index on created_at in verifications
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_verifications_created_at ON verifications (created_at);`).Error; err != nil {
		return err
	}

	// Create composite index on verification_id and document_id in document_analyses
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_document_analyses_verification_document ON document_analyses (verification_id, document_id);`).Error; err != nil {
		return err
	}

	// Create composite index on verification_id and selfie_id in face_matches
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_face_matches_verification_selfie ON face_matches (verification_id, selfie_id);`).Error; err != nil {
		return err
	}

	// Create index on risk_level in risk_analyses
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_risk_analyses_risk_level ON risk_analyses (risk_level);`).Error; err != nil {
		return err
	}

	// Create index on is_anomaly in anomaly_detections
	if err := db.Exec(`CREATE INDEX IF NOT EXISTS idx_anomaly_detections_is_anomaly ON anomaly_detections (is_anomaly);`).Error; err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// RollbackMigrations rolls back all database migrations
func RollbackMigrations(db *gorm.DB) error {
	// Drop tables in reverse order to avoid foreign key constraints
	if err := db.Exec(`DROP TABLE IF EXISTS
		anomaly_detections,
		risk_analyses,
		face_matches,
		document_analyses,
		selfies,
		audit_logs,
		verification_results,
		verifications,
		documents,
		users,
		kycs CASCADE;`).Error; err != nil {
		return err
	}

	return nil
}

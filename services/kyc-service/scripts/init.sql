-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create schema
CREATE SCHEMA IF NOT EXISTS kyc;

-- Set search path
SET search_path TO kyc, public;

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    mfa_secret VARCHAR(32),
    last_login_at TIMESTAMP,
    last_login_ip VARCHAR(45),
    login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    refresh_token VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255),
    ip_address VARCHAR(45),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create documents table
CREATE TABLE IF NOT EXISTS documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    name VARCHAR(255) NOT NULL,
    path VARCHAR(255) NOT NULL,
    size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    metadata JSONB,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create verifications table
CREATE TABLE IF NOT EXISTS verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kyc_id UUID NOT NULL,
    document_id UUID REFERENCES documents(id) ON DELETE SET NULL,
    selfie_id UUID REFERENCES documents(id) ON DELETE SET NULL,
    verifier_id UUID REFERENCES users(id) ON DELETE SET NULL,
    method VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    notes TEXT,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create document_analysis_results table
CREATE TABLE IF NOT EXISTS document_analysis_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    document_type VARCHAR(50) NOT NULL,
    is_authentic BOOLEAN NOT NULL,
    confidence FLOAT NOT NULL,
    extracted_data JSONB,
    issues JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create face_match_results table
CREATE TABLE IF NOT EXISTS face_match_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
    document_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    selfie_id UUID NOT NULL REFERENCES documents(id) ON DELETE CASCADE,
    is_match BOOLEAN NOT NULL,
    confidence FLOAT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create risk_analysis_results table
CREATE TABLE IF NOT EXISTS risk_analysis_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    risk_score FLOAT NOT NULL,
    risk_level VARCHAR(20) NOT NULL,
    risk_factors JSONB,
    device_info JSONB,
    ip_address VARCHAR(45),
    location VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create anomaly_detection_results table
CREATE TABLE IF NOT EXISTS anomaly_detection_results (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    verification_id UUID NOT NULL REFERENCES verifications(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_anomaly BOOLEAN NOT NULL,
    anomaly_score FLOAT NOT NULL,
    anomaly_type VARCHAR(50),
    reasons JSONB,
    device_info JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create ai_model_info table
CREATE TABLE IF NOT EXISTS ai_model_info (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    version VARCHAR(20) NOT NULL,
    type VARCHAR(50) NOT NULL,
    accuracy FLOAT NOT NULL,
    last_trained_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX IF NOT EXISTS idx_documents_user_id ON documents(user_id);
CREATE INDEX IF NOT EXISTS idx_documents_type ON documents(type);
CREATE INDEX IF NOT EXISTS idx_verifications_user_id ON verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_verifications_status ON verifications(status);
CREATE INDEX IF NOT EXISTS idx_document_analysis_verification_id ON document_analysis_results(verification_id);
CREATE INDEX IF NOT EXISTS idx_face_match_verification_id ON face_match_results(verification_id);
CREATE INDEX IF NOT EXISTS idx_risk_analysis_verification_id ON risk_analysis_results(verification_id);
CREATE INDEX IF NOT EXISTS idx_anomaly_detection_verification_id ON anomaly_detection_results(verification_id);
CREATE INDEX IF NOT EXISTS idx_ai_model_type ON ai_model_info(type);

-- Create test user (password: password123)
INSERT INTO users (email, password_hash, first_name, last_name, role)
VALUES ('admin@example.com', '$2a$10$1qAz2wSx3eDc4rFv5tGb5edDmJwVYVvxwc7VoOH.FNn4A0ftqCTDm', 'Admin', 'User', 'admin');

-- Create test AI models
INSERT INTO ai_model_info (name, version, type, accuracy, last_trained_at)
VALUES 
    ('Document Verification Model', '1.0.0', 'DOCUMENT', 0.98, NOW() - INTERVAL '1 day'),
    ('Face Recognition Model', '1.0.0', 'FACE', 0.95, NOW() - INTERVAL '2 days'),
    ('Risk Analysis Model', '1.0.0', 'RISK', 0.92, NOW() - INTERVAL '3 days'),
    ('Anomaly Detection Model', '1.0.0', 'ANOMALY', 0.90, NOW() - INTERVAL '4 days');

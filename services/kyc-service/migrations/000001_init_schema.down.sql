-- Drop triggers
DROP TRIGGER IF EXISTS update_documents_updated_at ON documents;
DROP TRIGGER IF EXISTS update_verification_details_updated_at ON verification_details;
DROP TRIGGER IF EXISTS update_kyc_profiles_updated_at ON kyc_profiles;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_documents_user_id;
DROP INDEX IF EXISTS idx_documents_status;
DROP INDEX IF EXISTS idx_documents_type;
DROP INDEX IF EXISTS idx_verification_details_document_id;
DROP INDEX IF EXISTS idx_kyc_profiles_status;
DROP INDEX IF EXISTS idx_kyc_profiles_risk_level;

-- Drop tables
DROP TABLE IF EXISTS verification_details;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS kyc_profiles; 
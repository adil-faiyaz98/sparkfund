-- Drop indexes
DROP INDEX IF EXISTS idx_security_activities_created_at;
DROP INDEX IF EXISTS idx_security_activities_user_id;
DROP INDEX IF EXISTS idx_security_audit_logs_created_at;
DROP INDEX IF EXISTS idx_security_audit_logs_user_id;
DROP INDEX IF EXISTS idx_password_resets_expires_at;
DROP INDEX IF EXISTS idx_password_resets_user_id;
DROP INDEX IF EXISTS idx_sessions_expires_at;
DROP INDEX IF EXISTS idx_sessions_user_id;
DROP INDEX IF EXISTS idx_sessions_token;
DROP INDEX IF EXISTS idx_users_email;

-- Drop tables
DROP TABLE IF EXISTS security_activities;
DROP TABLE IF EXISTS security_audit_logs;
DROP TABLE IF EXISTS mfa_configs;
DROP TABLE IF EXISTS password_resets;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS user_profiles;
DROP TABLE IF EXISTS users; 
-- Drop indexes
DROP INDEX IF EXISTS idx_email_logs_status;
DROP INDEX IF EXISTS idx_email_logs_created_at;
DROP INDEX IF EXISTS idx_templates_name;
DROP INDEX IF EXISTS idx_templates_created_at;

-- Drop tables
DROP TABLE IF EXISTS email_logs;
DROP TABLE IF EXISTS templates; 
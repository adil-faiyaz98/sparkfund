-- Drop triggers
DROP TRIGGER IF EXISTS update_email_logs_updated_at ON email_logs;
DROP TRIGGER IF EXISTS update_templates_updated_at ON templates;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_email_logs_status;
DROP INDEX IF EXISTS idx_email_logs_created_at;
DROP INDEX IF EXISTS idx_templates_name;

-- Drop tables
DROP TABLE IF EXISTS email_logs;
DROP TABLE IF EXISTS templates; 
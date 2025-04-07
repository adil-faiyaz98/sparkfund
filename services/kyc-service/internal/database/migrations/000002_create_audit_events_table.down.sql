-- Drop the view
DROP VIEW IF EXISTS audit_events_summary;

-- Remove the scheduled job
SELECT cron.unschedule('SELECT cleanup_old_audit_events()');

-- Drop the cleanup function
DROP FUNCTION IF EXISTS cleanup_old_audit_events();

-- Drop the table
DROP TABLE IF EXISTS audit_events;

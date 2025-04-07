CREATE TABLE IF NOT EXISTS audit_events (
    id UUID PRIMARY KEY,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    user_id VARCHAR(255),
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(50) NOT NULL,
    resource_id VARCHAR(255),
    status VARCHAR(50) NOT NULL,
    client_ip VARCHAR(50),
    user_agent TEXT,
    request_id VARCHAR(255),
    request_method VARCHAR(10),
    request_path TEXT,
    request_params JSONB,
    response_code INTEGER,
    error_message TEXT,
    changes JSONB,
    metadata JSONB
);

-- Create indexes for efficient querying
CREATE INDEX IF NOT EXISTS idx_audit_events_timestamp ON audit_events (timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_events_user_id ON audit_events (user_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_resource_id ON audit_events (resource_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_resource ON audit_events (resource);
CREATE INDEX IF NOT EXISTS idx_audit_events_action ON audit_events (action);
CREATE INDEX IF NOT EXISTS idx_audit_events_status ON audit_events (status);
CREATE INDEX IF NOT EXISTS idx_audit_events_request_id ON audit_events (request_id);

-- Create a function to automatically clean up old audit events
CREATE OR REPLACE FUNCTION cleanup_old_audit_events()
RETURNS void AS $$
BEGIN
    -- Delete audit events older than 1 year
    DELETE FROM audit_events
    WHERE timestamp < NOW() - INTERVAL '1 year';
END;
$$ LANGUAGE plpgsql;

-- Create a scheduled job to run the cleanup function
CREATE EXTENSION IF NOT EXISTS pg_cron;

-- Schedule the cleanup job to run once a day at 2 AM
SELECT cron.schedule('0 2 * * *', 'SELECT cleanup_old_audit_events()');

-- Create a view for common audit queries
CREATE OR REPLACE VIEW audit_events_summary AS
SELECT
    date_trunc('day', timestamp) AS day,
    action,
    resource,
    status,
    COUNT(*) AS count
FROM audit_events
GROUP BY 1, 2, 3, 4
ORDER BY 1 DESC, 5 DESC;

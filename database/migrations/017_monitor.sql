-- Migration: 017_monitor
-- Description: Create monitoring and alerting tables

-- System metrics table
CREATE TABLE system_metrics (
    id VARCHAR(26) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    cpu_usage DECIMAL(5,2),
    memory_usage DECIMAL(5,2),
    memory_total BIGINT,
    memory_used BIGINT,
    disk_usage DECIMAL(5,2),
    disk_total BIGINT,
    disk_used BIGINT,
    network_in BIGINT,
    network_out BIGINT,
    db_connections INTEGER,
    api_requests BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- API metrics table
CREATE TABLE api_metrics (
    id VARCHAR(26) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    endpoint VARCHAR(255),
    method VARCHAR(10),
    duration BIGINT,
    status_code INTEGER,
    user_id VARCHAR(26),
    ip_address VARCHAR(45),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Log entries table
CREATE TABLE log_entries (
    id VARCHAR(26) PRIMARY KEY,
    timestamp TIMESTAMP NOT NULL,
    level VARCHAR(20) CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR')),
    message TEXT,
    source VARCHAR(100),
    module VARCHAR(100),
    user_id VARCHAR(26),
    request_id VARCHAR(100),
    metadata TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Alert rules table
CREATE TABLE alert_rules (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metric VARCHAR(100),
    condition VARCHAR(10) CHECK (condition IN ('>', '<', '==', '!=', '>=', '<=')),
    threshold DECIMAL(10,2),
    duration INTEGER,
    severity VARCHAR(20) CHECK (severity IN ('warning', 'critical')),
    is_active BOOLEAN DEFAULT TRUE,
    notify_channels TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Alert history table
CREATE TABLE alert_history (
    id VARCHAR(26) PRIMARY KEY,
    rule_id VARCHAR(26) REFERENCES alert_rules(id),
    rule_name VARCHAR(255),
    severity VARCHAR(20),
    message TEXT,
    value DECIMAL(10,2),
    threshold DECIMAL(10,2),
    status VARCHAR(20) CHECK (status IN ('firing', 'resolved')),
    resolved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_system_metrics_timestamp ON system_metrics(timestamp);
CREATE INDEX idx_api_metrics_timestamp ON api_metrics(timestamp);
CREATE INDEX idx_api_metrics_endpoint ON api_metrics(endpoint);
CREATE INDEX idx_log_entries_timestamp ON log_entries(timestamp);
CREATE INDEX idx_log_entries_level ON log_entries(level);
CREATE INDEX idx_log_entries_source ON log_entries(source);
CREATE INDEX idx_alert_history_rule ON alert_history(rule_id);
CREATE INDEX idx_alert_history_status ON alert_history(status);
CREATE INDEX idx_alert_history_created ON alert_history(created_at);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_alert_rules_updated_at BEFORE UPDATE ON alert_rules
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Partitioning for log entries (optional, for large deployments)
-- CREATE TABLE log_entries_partitioned (LIKE log_entries) PARTITION BY RANGE (timestamp);

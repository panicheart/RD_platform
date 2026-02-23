-- Migration: 018_zotero_connection
-- Description: Create table for Zotero API connections

-- Zotero connection settings for users
CREATE TABLE zotero_connections (
    id VARCHAR(26) PRIMARY KEY,
    user_id VARCHAR(26) NOT NULL UNIQUE,
    api_key VARCHAR(255) NOT NULL,
    zotero_user_id VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_zotero_connections_user ON zotero_connections(user_id);
CREATE INDEX idx_zotero_connections_active ON zotero_connections(is_active);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_zotero_connections_updated_at BEFORE UPDATE ON zotero_connections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

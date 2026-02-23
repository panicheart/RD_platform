-- Migration: 010_technologies
-- Description: Create technologies table for technical shelf

-- Technologies table for technical shelf (tech tree)
CREATE TABLE technologies (
    id CHAR(26) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    trl_level INTEGER CHECK (trl_level >= 1 AND trl_level <= 9),
    category VARCHAR(100),
    parent_id CHAR(26) REFERENCES technologies(id),
    owner_id CHAR(26) REFERENCES users(id),
    is_published BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(26) REFERENCES users(id)
);

CREATE INDEX idx_technologies_trl_level ON technologies(trl_level);
CREATE INDEX idx_technologies_category ON technologies(category);
CREATE INDEX idx_technologies_is_published ON technologies(is_published);
CREATE INDEX idx_technologies_parent ON technologies(parent_id);

CREATE TRIGGER update_technologies_updated_at BEFORE UPDATE ON technologies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE technologies IS 'Technical shelf - technology tree';

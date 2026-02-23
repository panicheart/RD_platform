-- Migration: 014_knowledge
-- Description: Create knowledge base tables

-- Categories table for 3-level tree structure
CREATE TABLE categories (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    parent_id VARCHAR(26) REFERENCES categories(id),
    level INTEGER DEFAULT 1 CHECK (level IN (1, 2, 3)),
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Knowledge entries table
CREATE TABLE knowledge (
    id VARCHAR(26) PRIMARY KEY,
    title VARCHAR(500) NOT NULL,
    content TEXT,
    category_id VARCHAR(26) REFERENCES categories(id),
    author_id VARCHAR(26) NOT NULL,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    version INTEGER DEFAULT 1,
    parent_id VARCHAR(26) REFERENCES knowledge(id),
    view_count INTEGER DEFAULT 0,
    source VARCHAR(50),
    source_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at TIMESTAMP
);

-- Tags table
CREATE TABLE tags (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    color VARCHAR(7) DEFAULT '#1890ff',
    count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Many-to-many relationship between knowledge and tags
CREATE TABLE knowledge_tags (
    knowledge_id VARCHAR(26) REFERENCES knowledge(id) ON DELETE CASCADE,
    tag_id VARCHAR(26) REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (knowledge_id, tag_id)
);

-- Knowledge reviews for approval workflow
CREATE TABLE knowledge_reviews (
    id VARCHAR(26) PRIMARY KEY,
    knowledge_id VARCHAR(26) REFERENCES knowledge(id) ON DELETE CASCADE,
    reviewer_id VARCHAR(26) NOT NULL,
    status VARCHAR(20) CHECK (status IN ('pending', 'approved', 'rejected')),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Zotero items integration
CREATE TABLE zotero_items (
    id VARCHAR(26) PRIMARY KEY,
    zotero_key VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(500),
    item_type VARCHAR(50),
    authors TEXT,
    abstract TEXT,
    publication VARCHAR(255),
    volume VARCHAR(50),
    issue VARCHAR(50),
    pages VARCHAR(50),
    date VARCHAR(50),
    doi VARCHAR(100),
    url TEXT,
    pdf_path TEXT,
    tags TEXT,
    synced_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Obsidian vault mappings
CREATE TABLE obsidian_mappings (
    id VARCHAR(26) PRIMARY KEY,
    vault_path TEXT NOT NULL,
    local_path TEXT NOT NULL,
    category_id VARCHAR(26) REFERENCES categories(id),
    auto_sync BOOLEAN DEFAULT FALSE,
    last_sync_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_knowledge_category ON knowledge(category_id);
CREATE INDEX idx_knowledge_author ON knowledge(author_id);
CREATE INDEX idx_knowledge_status ON knowledge(status);
CREATE INDEX idx_knowledge_source ON knowledge(source);
CREATE INDEX idx_categories_parent ON categories(parent_id);
CREATE INDEX idx_knowledge_reviews_knowledge ON knowledge_reviews(knowledge_id);
CREATE INDEX idx_knowledge_reviews_status ON knowledge_reviews(status);
CREATE INDEX idx_obsidian_mappings_category ON obsidian_mappings(category_id);

-- Trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_knowledge_updated_at BEFORE UPDATE ON knowledge
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_zotero_items_updated_at BEFORE UPDATE ON zotero_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_obsidian_mappings_updated_at BEFORE UPDATE ON obsidian_mappings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

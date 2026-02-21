-- Product Shelf Tables Migration
-- Migration: 009_products.sql

-- TRL level enum
CREATE TYPE trl_level AS ENUM ('1', '2', '3', '4', '5', '6', '7', '8', '9');

-- Products table
CREATE TABLE products (
    id CHAR(26) PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    trl_level INTEGER CHECK (trl_level >= 1 AND trl_level <= 9),
    category VARCHAR(100),
    version VARCHAR(50),
    source_project_id CHAR(26) REFERENCES projects(id),
    owner_id CHAR(26) REFERENCES users(id),
    is_published BOOLEAN DEFAULT false,
    published_at TIMESTAMP WITH TIME ZONE,
    download_count INTEGER DEFAULT 0,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(26) REFERENCES users(id)
);

CREATE INDEX idx_products_trl_level ON products(trl_level);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_is_published ON products(is_published);
CREATE INDEX idx_products_source_project ON products(source_project_id);

-- Product versions table
CREATE TABLE product_versions (
    id CHAR(26) PRIMARY KEY,
    product_id CHAR(26) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    version VARCHAR(50) NOT NULL,
    parent_version_id CHAR(26) REFERENCES product_versions(id),
    changelog TEXT,
    files JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by CHAR(26) REFERENCES users(id),
    UNIQUE(product_id, version)
);

CREATE INDEX idx_product_versions_product_id ON product_versions(product_id);

-- Cart table
CREATE TABLE cart_items (
    id CHAR(26) PRIMARY KEY,
    user_id CHAR(26) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id CHAR(26) NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    project_id CHAR(26) REFERENCES projects(id),
    quantity INTEGER DEFAULT 1,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_id)
);

CREATE INDEX idx_cart_items_user_id ON cart_items(user_id);

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

COMMENT ON TABLE products IS 'Product catalog from projects';
COMMENT ON TABLE product_versions IS 'Product version history';
COMMENT ON TABLE cart_items IS 'User shopping cart for products';

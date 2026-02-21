-- Migration: 000_init_extensions
-- Description: Initialize required PostgreSQL extensions
-- Created: 2026-02-22

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pgcrypto for additional crypto functions
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Enable trigram similarity for text search
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

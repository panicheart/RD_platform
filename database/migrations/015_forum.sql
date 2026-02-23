-- Migration: 015_forum
-- Description: Create forum tables

-- Forum boards table
CREATE TABLE forum_boards (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(50),
    icon VARCHAR(100),
    sort_order INTEGER DEFAULT 0,
    topic_count INTEGER DEFAULT 0,
    post_count INTEGER DEFAULT 0,
    last_post_at TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Forum posts table
CREATE TABLE forum_posts (
    id VARCHAR(26) PRIMARY KEY,
    board_id VARCHAR(26) NOT NULL REFERENCES forum_boards(id) ON DELETE CASCADE,
    title VARCHAR(500) NOT NULL,
    content TEXT,
    author_id VARCHAR(26) NOT NULL,
    author_name VARCHAR(100),
    view_count INTEGER DEFAULT 0,
    reply_count INTEGER DEFAULT 0,
    is_pinned BOOLEAN DEFAULT FALSE,
    is_locked BOOLEAN DEFAULT FALSE,
    is_best_answer BOOLEAN DEFAULT FALSE,
    tags TEXT,
    knowledge_id VARCHAR(26) REFERENCES knowledge(id),
    last_reply_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Forum replies table
CREATE TABLE forum_replies (
    id VARCHAR(26) PRIMARY KEY,
    post_id VARCHAR(26) NOT NULL REFERENCES forum_posts(id) ON DELETE CASCADE,
    parent_id VARCHAR(26) REFERENCES forum_replies(id),
    content TEXT,
    author_id VARCHAR(26) NOT NULL,
    author_name VARCHAR(100),
    is_best_answer BOOLEAN DEFAULT FALSE,
    mentions TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Forum tags table
CREATE TABLE forum_tags (
    id VARCHAR(26) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    color VARCHAR(7) DEFAULT '#1890ff',
    count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX idx_forum_posts_board ON forum_posts(board_id);
CREATE INDEX idx_forum_posts_author ON forum_posts(author_id);
CREATE INDEX idx_forum_posts_pinned ON forum_posts(is_pinned);
CREATE INDEX idx_forum_posts_created ON forum_posts(created_at);
CREATE INDEX idx_forum_replies_post ON forum_replies(post_id);
CREATE INDEX idx_forum_replies_parent ON forum_replies(parent_id);

-- Trigger to update updated_at timestamp
CREATE TRIGGER update_forum_boards_updated_at BEFORE UPDATE ON forum_boards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_forum_posts_updated_at BEFORE UPDATE ON forum_posts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_forum_replies_updated_at BEFORE UPDATE ON forum_replies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Trigger to update topic/reply counts
CREATE OR REPLACE FUNCTION update_forum_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE forum_boards 
        SET post_count = post_count + 1,
            last_post_at = NEW.created_at
        WHERE id = (SELECT board_id FROM forum_posts WHERE id = NEW.post_id);
        
        UPDATE forum_posts 
        SET reply_count = reply_count + 1,
            last_reply_at = NEW.created_at
        WHERE id = NEW.post_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE forum_boards 
        SET post_count = post_count - 1
        WHERE id = (SELECT board_id FROM forum_posts WHERE id = OLD.post_id);
        
        UPDATE forum_posts 
        SET reply_count = reply_count - 1
        WHERE id = OLD.post_id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER trigger_update_forum_counts
    AFTER INSERT OR DELETE ON forum_replies
    FOR EACH ROW EXECUTE FUNCTION update_forum_counts();

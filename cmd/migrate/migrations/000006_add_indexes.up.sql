-- Enable pg_trgm extension for text similarity indexing
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Create indexes for faster search
CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);

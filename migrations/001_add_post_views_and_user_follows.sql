-- Add view_count column to posts table
ALTER TABLE posts ADD COLUMN view_count BIGINT DEFAULT 0;

-- Add follower/following counts to users table
ALTER TABLE users ADD COLUMN followers_count BIGINT DEFAULT 0;
ALTER TABLE users ADD COLUMN following_count BIGINT DEFAULT 0;

-- Create post_views table
CREATE TABLE IF NOT EXISTS post_views (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    post_id UUID NOT NULL,
    user_id UUID, -- nullable for anonymous views
    ip_address VARCHAR(45),
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for post_views
CREATE INDEX IF NOT EXISTS idx_post_views_post_id ON post_views(post_id);
CREATE INDEX IF NOT EXISTS idx_post_views_user_id ON post_views(user_id);
CREATE INDEX IF NOT EXISTS idx_post_views_deleted_at ON post_views(deleted_at);
CREATE INDEX IF NOT EXISTS idx_post_views_created_at ON post_views(created_at);

-- Create unique index to prevent duplicate views from same user on same post
CREATE UNIQUE INDEX IF NOT EXISTS idx_post_views_unique_user_post 
ON post_views(post_id, user_id) 
WHERE user_id IS NOT NULL AND deleted_at IS NULL;

-- Create user_follows table
CREATE TABLE IF NOT EXISTS user_follows (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    follower_id UUID NOT NULL,
    following_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes for user_follows
CREATE INDEX IF NOT EXISTS idx_user_follows_follower_id ON user_follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_user_follows_following_id ON user_follows(following_id);
CREATE INDEX IF NOT EXISTS idx_user_follows_deleted_at ON user_follows(deleted_at);
CREATE INDEX IF NOT EXISTS idx_user_follows_created_at ON user_follows(created_at);

-- Create unique index to prevent duplicate follows
CREATE UNIQUE INDEX IF NOT EXISTS idx_user_follows_unique 
ON user_follows(follower_id, following_id) 
WHERE deleted_at IS NULL;

-- Add foreign key constraints
ALTER TABLE post_views 
ADD CONSTRAINT fk_post_views_post_id 
FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE;

ALTER TABLE post_views 
ADD CONSTRAINT fk_post_views_user_id 
FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_follows 
ADD CONSTRAINT fk_user_follows_follower_id 
FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_follows 
ADD CONSTRAINT fk_user_follows_following_id 
FOREIGN KEY (following_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add constraint to prevent self-following
ALTER TABLE user_follows 
ADD CONSTRAINT chk_user_follows_no_self_follow 
CHECK (follower_id != following_id);

-- Create function to update follow counts
CREATE OR REPLACE FUNCTION update_user_follow_counts()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        -- Increment follower's following count
        UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;
        -- Increment following user's followers count
        UPDATE users SET followers_count = followers_count + 1 WHERE id = NEW.following_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        -- Decrement follower's following count
        UPDATE users SET following_count = following_count - 1 WHERE id = OLD.follower_id;
        -- Decrement following user's followers count
        UPDATE users SET followers_count = followers_count - 1 WHERE id = OLD.following_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for follow count updates
DROP TRIGGER IF EXISTS trigger_update_follow_counts_insert ON user_follows;
CREATE TRIGGER trigger_update_follow_counts_insert
    AFTER INSERT ON user_follows
    FOR EACH ROW
    EXECUTE FUNCTION update_user_follow_counts();

DROP TRIGGER IF EXISTS trigger_update_follow_counts_delete ON user_follows;
CREATE TRIGGER trigger_update_follow_counts_delete
    AFTER DELETE ON user_follows
    FOR EACH ROW
    EXECUTE FUNCTION update_user_follow_counts();

-- Create function to update post view counts
CREATE OR REPLACE FUNCTION update_post_view_count()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE posts SET view_count = view_count + 1 WHERE id = NEW.post_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE posts SET view_count = view_count - 1 WHERE id = OLD.post_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for view count updates
DROP TRIGGER IF EXISTS trigger_update_view_count_insert ON post_views;
CREATE TRIGGER trigger_update_view_count_insert
    AFTER INSERT ON post_views
    FOR EACH ROW
    EXECUTE FUNCTION update_post_view_count();

DROP TRIGGER IF EXISTS trigger_update_view_count_delete ON post_views;
CREATE TRIGGER trigger_update_view_count_delete
    AFTER DELETE ON post_views
    FOR EACH ROW
    EXECUTE FUNCTION update_post_view_count();
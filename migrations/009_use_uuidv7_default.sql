-- +goose Up
-- Switch UUID primary key defaults from uuid_generate_v4() to native uuidv7() (PostgreSQL 18+).
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE posts ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE post_comments ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE files ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE post_views ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE user_follows ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE post_likes ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE bookmark_folders ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE post_bookmarks ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE chat_conversations ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE chat_messages ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE notifications ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE password_reset_tokens ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE auth_activity_logs ALTER COLUMN id SET DEFAULT uuidv7();

-- +goose Down
ALTER TABLE auth_activity_logs ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE password_reset_tokens ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE notifications ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE chat_messages ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE chat_conversations ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE post_bookmarks ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE bookmark_folders ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE post_likes ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE user_follows ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE post_views ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE files ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE post_comments ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE posts ALTER COLUMN id SET DEFAULT uuidv4();
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuidv4();

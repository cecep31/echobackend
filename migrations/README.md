# Database Migrations

Managed via [goose](https://github.com/pressly/goose). Env vars configured in `.env`.

## Running Migrations

```bash
# Apply all pending migrations
goose up

# Rollback one migration
goose down

# Check current status
goose status

# Create a new migration
goose create nama_migration sql
```

## Migration Order

| # | File | Description |
|---|------|-------------|
| 001 | `001_init_schema.sql` | Core tables: users, posts, tags, posts_to_tags, profiles, sessions, post_comments, files |
| 002 | `002_add_post_views_and_user_follows.sql` | post_views table, user_follows table, view/follow count triggers |
| 003 | `003_add_post_likes.sql` | post_likes table, like_count trigger |
| 004 | `004_add_bookmark_folders_and_post_bookmarks.sql` | bookmark_folders, post_bookmarks, bookmark_count trigger |
| 005 | `005_add_chat_conversations_and_messages.sql` | chat_conversations, chat_messages |

## Notes

- All `CREATE TABLE` and `ADD COLUMN` statements use `IF NOT EXISTS` / `IF NOT EXISTS` for idempotency
- Foreign key constraints with `ON DELETE CASCADE` ensure data integrity
- Database triggers automatically maintain count fields (view_count, like_count, bookmark_count, followers_count, following_count)
- Soft deletes are supported via `deleted_at` timestamps
- UUID v4 (`uuid_generate_v4()`) is used for primary keys
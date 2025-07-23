# Database Migrations

This directory contains database migration scripts for the Echo Backend application.

## Available Migrations

### 001_add_post_views_and_user_follows.sql

This migration adds support for:
- **Post View Tracking**: Records when users view posts, including anonymous views
- **User Following System**: Allows users to follow/unfollow each other

#### Features Added:

**Post Views:**
- `post_views` table to track individual post views
- `view_count` column added to `posts` table
- Support for both authenticated and anonymous views
- IP address and user agent tracking
- Automatic view count updates via database triggers

**User Following:**
- `user_follows` table to track follower/following relationships
- `followers_count` and `following_count` columns added to `users` table
- Prevention of self-following via database constraints
- Automatic follower/following count updates via database triggers

#### Database Objects Created:

**Tables:**
- `post_views` - Records individual post views
- `user_follows` - Records user following relationships

**Indexes:**
- Optimized indexes for query performance
- Unique constraints to prevent duplicate views/follows

**Functions & Triggers:**
- `update_user_follow_counts()` - Maintains follower/following counts
- `update_post_view_count()` - Maintains post view counts
- Automatic triggers for count updates

## Running Migrations

To apply the migration, run the SQL script against your PostgreSQL database:

```bash
# Using psql
psql -d your_database_name -f migrations/001_add_post_views_and_user_follows.sql

# Or using a migration tool like migrate
migrate -path migrations -database "postgres://user:password@localhost/dbname?sslmode=disable" up
```

## API Endpoints Added

### Post View Endpoints
- `POST /api/posts/:id/view` - Record a post view
- `GET /api/posts/:id/views` - Get post views (paginated)
- `GET /api/posts/:id/view-stats` - Get post view statistics
- `GET /api/posts/:id/viewed` - Check if current user viewed the post

### User Follow Endpoints
- `POST /api/users/:id/follow` - Follow a user
- `DELETE /api/users/:id/follow` - Unfollow a user
- `GET /api/users/:id/followers` - Get user's followers (paginated)
- `GET /api/users/:id/following` - Get users that user is following (paginated)
- `GET /api/users/:id/follow-status` - Check if current user follows the user
- `GET /api/users/:id/mutual-follows` - Get mutual follows (paginated)
- `GET /api/users/:id/follow-stats` - Get follow statistics

## Notes

- All migrations use `IF NOT EXISTS` clauses to be safely re-runnable
- Foreign key constraints ensure data integrity
- Database triggers automatically maintain count fields
- Soft deletes are supported via `deleted_at` timestamps
- UUID v7 is used for primary keys for better performance
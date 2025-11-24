# Backend API Documentation

## Overview

This is a comprehensive REST API for a blog/content management system built with Go and Echo framework. The API provides endpoints for user management, authentication, posts, comments, tags, workspaces, pages, user interactions, and chat conversations.

## Base URL

```
http://localhost:8080/v1
```

Note: Replace with your actual deployment URL when deployed.

## Authentication

The API uses JWT (JSON Web Token) for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Standard Response Format

All API responses follow a consistent format:

```json
{
  "success": true|false,
  "message": "Response message",
  "data": {}, // Response data (optional)
  "error": "Error message", // Only present on errors
  "meta": {} // Pagination metadata (optional)
}
```

### Pagination Response

For paginated endpoints, the response includes metadata:

```json
{
  "success": true,
  "message": "Successfully retrieved data",
  "data": [],
  "meta": {
    "total_items": 100,
    "offset": 0,
    "limit": 10,
    "total_pages": 10
  }
}
```

## HTTP Status Codes

- `200` - OK: Successful request
- `201` - Created: Resource successfully created
- `400` - Bad Request: Invalid request format or parameters
- `401` - Unauthorized: Authentication required or invalid
- `403` - Forbidden: Access denied
- `404` - Not Found: Resource not found
- `422` - Unprocessable Entity: Validation errors
- `500` - Internal Server Error: Server error

---

## Authentication Endpoints

### Register User

**POST** `/v1/auth/register`

Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123"
}
```

**Validation Rules:**
- `email`: Required, valid email format
- `username`: Required, 3-30 characters
- `password`: Required, minimum 6 characters

**Response (201):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "username"
  }
}
```

**Error Response (409):**
```json
{
  "success": false,
  "message": "Registration failed",
  "error": "Email or username already exists"
}
```

### Login

**POST** `/v1/auth/login`

Authenticate user and receive JWT token.

**Rate Limit:** 5 requests per 5 minutes

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "jwt-token-here",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "username"
    }
  }
}
```

### Check Username Availability

**POST** `/v1/auth/check-username`

Check if a username is available for registration.

**Request Body:**
```json
{
  "username": "desired-username"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Username availability checked",
  "data": {
    "username": "desired-username",
    "available": true
  }
}
```

### Forgot Password

**POST** `/v1/auth/forgot-password`

Request a password reset email.

**Rate Limit:** 5 requests per 5 minutes

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Password reset email sent if user exists"
}
```

### Reset Password

**POST** `/v1/auth/reset-password`

Reset password using reset token.

**Request Body:**
```json
{
  "token": "reset-token-from-email",
  "password": "new-password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

### Refresh Token

**POST** `/v1/auth/refresh`

Refresh JWT token using refresh token.

**Request Body:**
```json
{
  "refresh_token": "refresh-token"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "new-jwt-token"
  }
}
```

### Change Password

**PUT** `/v1/auth/change-password`

🔒 **Authentication Required**

Change user password.

**Request Body:**
```json
{
  "current_password": "current-password",
  "new_password": "new-password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Password changed successfully"
}
```

---

## User Endpoints

### Get User by ID

**GET** `/v1/users/{id}`

Retrieve user information by ID. Includes follow status if authenticated.

**Parameters:**
- `id` (path): User ID

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved user",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "name": "First Last",
    "username": "username",
    "image": "profile-image-url",
    "first_name": "First",
    "last_name": "Last",
    "followers_count": 10,
    "following_count": 5,
    "is_following": true,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### Get Current User

**GET** `/v1/users/me`

🔒 **Authentication Required**

Retrieve current authenticated user's information.

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved current user",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "first_name": "First",
    "last_name": "Last",
    "image": "profile-image-url",
    "is_super_admin": false,
    "followers_count": 10,
    "following_count": 5
  }
}
```

### Get All Users (Admin Only)

**GET** `/v1/users`

🔒 **Authentication Required** | 👑 **Admin Only**

Retrieve paginated list of all users.

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved users",
  "data": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "username": "username",
      "name": "First Last"
    }
  ],
  "meta": {
    "total_items": 100,
    "offset": 0,
    "limit": 10,
    "total_pages": 10
  }
}
```

### Delete User (Admin Only)

**DELETE** `/v1/users/{id}`

🔒 **Authentication Required** | 👑 **Admin Only**

Delete a user account.

**Parameters:**
- `id` (path): User ID to delete

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully deleted user",
  "data": null
}
```

---

## User Follow Endpoints

### Follow User

**POST** `/v1/users/follow`

🔒 **Authentication Required**

Follow another user.

**Request Body:**
```json
{
  "user_id": "uuid-of-user-to-follow"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully followed user"
}
```

### Unfollow User

**DELETE** `/v1/users/{id}/follow`

🔒 **Authentication Required**

Unfollow a user.

**Parameters:**
- `id` (path): User ID to unfollow

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully unfollowed user"
}
```

### Check Follow Status

**GET** `/v1/users/{id}/follow-status`

🔒 **Authentication Required**

Check if current user follows the specified user.

**Parameters:**
- `id` (path): User ID to check

**Response (200):**
```json
{
  "success": true,
  "message": "Follow status retrieved",
  "data": {
    "is_following": true
  }
}
```

### Get Mutual Follows

**GET** `/v1/users/{id}/mutual-follows`

🔒 **Authentication Required**

Get users that both current user and specified user follow.

**Parameters:**
- `id` (path): User ID

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Mutual follows retrieved",
  "data": [
    {
      "id": "uuid",
      "username": "mutual-friend",
      "name": "Mutual Friend"
    }
  ],
  "meta": {
    "total_items": 5,
    "offset": 0,
    "limit": 10,
    "total_pages": 1
  }
}
```

### Get User Followers

**GET** `/v1/users/{id}/followers`

Get list of users following the specified user.

**Parameters:**
- `id` (path): User ID

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Followers retrieved",
  "data": [
    {
      "id": "uuid",
      "username": "follower",
      "name": "Follower Name"
    }
  ],
  "meta": {
    "total_items": 25,
    "offset": 0,
    "limit": 10,
    "total_pages": 3
  }
}
```

### Get User Following

**GET** `/v1/users/{id}/following`

Get list of users that the specified user follows.

**Parameters:**
- `id` (path): User ID

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Following list retrieved",
  "data": [
    {
      "id": "uuid",
      "username": "followed-user",
      "name": "Followed User"
    }
  ],
  "meta": {
    "total_items": 15,
    "offset": 0,
    "limit": 10,
    "total_pages": 2
  }
}
```

### Get Follow Statistics

**GET** `/v1/users/{id}/follow-stats`

Get follower and following counts for a user.

**Parameters:**
- `id` (path): User ID

**Response (200):**
```json
{
  "success": true,
  "message": "Follow statistics retrieved",
  "data": {
    "user_id": "uuid",
    "followers_count": 100,
    "following_count": 50
  }
}
```

---

## Post Endpoints

### Create Post

**POST** `/v1/posts`

🔒 **Authentication Required**

Create a new blog post.

**Request Body:**
```json
{
  "title": "Post Title",
  "photo_url": "https://example.com/image.jpg",
  "slug": "post-slug",
  "body": "Post content here...",
  "tags": ["tag1", "tag2"]
}
```

**Validation Rules:**
- `title`: Required, minimum 7 characters
- `slug`: Required, minimum 7 characters, must be unique
- `body`: Required, minimum 10 characters
- `photo_url`: Optional
- `tags`: Optional array of tag names

**Response (201):**
```json
{
  "success": true,
  "message": "Successfully created post",
  "data": {
    "id": "uuid"
  }
}
```

### Get All Posts

**GET** `/v1/posts`

Retrieve paginated list of all posts with advanced filtering and search capabilities.

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10, max: 100)
- `search` (optional): Search term to filter posts by title only (case-insensitive)
- `sort_by` (optional): Sort field - `id`, `title`, `created_at`, `updated_at`, `view_count`, `like_count` (default: `created_at`)
- `sort_order` (optional): Sort order - `asc` or `desc` (default: `desc`)
- `start_date` (optional): Filter posts created after this date (YYYY-MM-DD format)
- `end_date` (optional): Filter posts created before this date (YYYY-MM-DD format)
- `created_by` (optional): Filter posts by author ID
- `tags` (optional): Filter posts by tag names (comma-separated, e.g., "golang,webdev,api")
- `published` (optional): Filter by publication status - `true` or `false`

**Examples:**
```
GET /v1/posts?search=golang&limit=20
GET /v1/posts?sort_by=view_count&sort_order=desc&limit=10
GET /v1/posts?tags=golang,webdev&start_date=2023-01-01&end_date=2023-12-31
GET /v1/posts?search=api&sort_by=created_at&sort_order=desc&limit=20&offset=10
```

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved posts",
  "data": [
    {
      "id": "uuid",
      "title": "Post Title",
      "photo_url": "image-url",
      "body": "Truncated content...",
      "slug": "post-slug",
      "view_count": 100,
      "creator": {
        "id": "uuid",
        "username": "author",
        "name": "Author Name"
      },
      "tags": [
        {
          "id": "uuid",
          "name": "tag-name"
        }
      ],
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total_items": 100,
    "offset": 0,
    "limit": 10,
    "total_pages": 10
  }
}
```

### Get Post by ID

**GET** `/v1/posts/{id}`

Retrieve a specific post by ID. Automatically records a view.

**Parameters:**
- `id` (path): Post ID

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved post",
  "data": {
    "id": "uuid",
    "title": "Post Title",
    "photo_url": "image-url",
    "body": "Full post content...",
    "slug": "post-slug",
    "view_count": 100,
    "creator": {
      "id": "uuid",
      "username": "author",
      "name": "Author Name"
    },
    "tags": [],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### Get Post by Slug and Username

**GET** `/v1/posts/u/{username}/{slug}`

Retrieve a post by username and slug. Automatically records a view.

**Parameters:**
- `username` (path): Author's username
- `slug` (path): Post slug

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved post",
  "data": {
    "id": "uuid",
    "title": "Post Title",
    "photo_url": "image-url",
    "body": "Full post content...",
    "slug": "post-slug",
    "view_count": 100,
    "creator": {
      "id": "uuid",
      "username": "author",
      "name": "Author Name"
    },
    "tags": [],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### Update Post

**PUT** `/v1/posts/{id}`

🔒 **Authentication Required**

Update an existing post. Only the post author can update.

**Parameters:**
- `id` (path): Post ID

**Request Body:**
```json
{
  "title": "Updated Title",
  "photo_url": "new-image-url",
  "slug": "updated-slug",
  "body": "Updated content...",
  "tags": ["new-tag"]
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Post updated successfully",
  "data": {
    "id": "uuid",
    "title": "Updated Title",
    "photo_url": "new-image-url",
    "slug": "updated-slug",
    "body": "Updated content...",
    "view_count": 100,
    "creator": {
      "id": "uuid",
      "username": "author",
      "name": "Author Name"
    },
    "tags": [
      {
        "id": "uuid",
        "name": "new-tag"
      }
    ],
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z"
  }
}
```

### Delete Post

**DELETE** `/v1/posts/{id}`

🔒 **Authentication Required**

Delete a post. Only the post author can delete.

**Parameters:**
- `id` (path): Post ID

**Response (200):**
```json
{
  "success": true,
  "message": "Post deleted successfully",
  "data": null
}
```

### Get Random Posts

**GET** `/v1/posts/random`

Retrieve random posts for discovery.

**Query Parameters:**
- `limit` (optional): Number of posts to return (default: 9, max: 50)

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved random posts",
  "data": [
    {
      "id": "uuid",
      "title": "Random Post Title",
      "photo_url": "image-url",
      "slug": "random-post-slug",
      "view_count": 50,
      "creator": {
        "id": "uuid",
        "username": "author",
        "name": "Author Name"
      }
    }
  ]
}
```

### Get My Posts

**GET** `/v1/posts/mine`

🔒 **Authentication Required**

Retrieve current user's posts.

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved user's posts",
  "data": [
    {
      "id": "uuid",
      "title": "User's Post Title",
      "photo_url": "image-url",
      "slug": "user-post-slug",
      "body": "Post content...",
      "view_count": 25,
      "like_count": 3,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total_items": 5,
    "offset": 0,
    "limit": 10,
    "total_pages": 1
  }
}
```

### Get Posts by Username

**GET** `/v1/posts/username/{username}`

Retrieve posts by a specific user.

**Parameters:**
- `username` (path): Author's username

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved posts by username",
  "data": [
    {
      "id": "uuid",
      "title": "Post Title",
      "photo_url": "image-url",
      "slug": "post-slug",
      "view_count": 75,
      "creator": {
        "id": "uuid",
        "username": "author",
        "name": "Author Name"
      },
      "tags": [],
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total_items": 15,
    "offset": 0,
    "limit": 10,
    "total_pages": 2
  }
}
```

### Get Posts by Tag

**GET** `/v1/posts/tag/{tag}`

Retrieve posts with a specific tag.

**Parameters:**
- `tag` (path): Tag name

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved posts by tag",
  "data": [
    {
      "id": "uuid",
      "title": "Post Title",
      "photo_url": "image-url",
      "slug": "post-slug",
      "view_count": 45,
      "creator": {
        "id": "uuid",
        "username": "author",
        "name": "Author Name"
      },
      "tags": [
        {
          "id": "uuid",
          "name": "tag-name"
        }
      ],
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total_items": 8,
    "offset": 0,
    "limit": 10,
    "total_pages": 1
  }
}
```

### Upload Post Image

**POST** `/v1/posts/image`

🔒 **Authentication Required**

Upload an image for use in posts.

**Request:** Multipart form data
- `image` (file): Image file to upload

**Response (200):**
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "image_url": "https://your-domain.com/uploads/image-uuid.jpg",
    "filename": "original-filename.jpg",
    "size": 1024000,
    "mime_type": "image/jpeg"
  }
}
```

---

## Comment Endpoints

### Get Post Comments

**GET** `/v1/posts/{id}/comments`

Retrieve all comments for a specific post.

**Parameters:**
- `id` (path): Post ID

**Response (200):**
```json
{
  "success": true,
  "message": "Comments fetched successfully",
  "data": [
    {
      "id": "uuid",
      "text": "Comment text",
      "author": {
        "id": "uuid",
        "username": "commenter",
        "name": "Commenter Name"
      },
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### Create Comment

**POST** `/v1/posts/{id}/comments`

🔒 **Authentication Required**

Create a new comment on a post.

**Parameters:**
- `id` (path): Post ID

**Request Body:**
```json
{
  "text": "Comment text here..."
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Comment created successfully",
  "data": {
    "id": "uuid",
    "text": "Comment text",
    "created_at": "2023-01-01T00:00:00Z"
  }
}
```

### Update Comment

**PUT** `/v1/posts/{id}/comments/{comment_id}`

🔒 **Authentication Required** | 👤 **Author Only**

Update an existing comment. Only the comment author can update.

**Parameters:**
- `id` (path): Post ID
- `comment_id` (path): Comment ID

**Request Body:**
```json
{
  "text": "Updated comment text..."
}
```

### Delete Comment

**DELETE** `/v1/posts/{id}/comments/{comment_id}`

🔒 **Authentication Required** | 👤 **Author Only**

Delete a comment. Only the comment author can delete.

**Parameters:**
- `id` (path): Post ID
- `comment_id` (path): Comment ID

---

## Post View Endpoints

### Record Post View

**POST** `/v1/posts/{id}/view`

Record a view for a post. Can be called by anonymous users.

**Parameters:**
- `id` (path): Post ID

### Get Post Views

**GET** `/v1/posts/{id}/views`

🔒 **Authentication Required**

Get detailed view information for a post.

### Get Post View Statistics

**GET** `/v1/posts/{id}/view-stats`

Get view statistics for a post.

### Check if User Viewed Post

**GET** `/v1/posts/{id}/viewed`

🔒 **Authentication Required**

Check if current user has viewed a specific post.

---

## Post Like Endpoints

### Like Post

**POST** `/v1/posts/{id}/like`

🔒 **Authentication Required**

Like a post. Users can only like a post once.

**Parameters:**
- `id` (path): Post ID

**Response (200):**
```json
{
  "success": true,
  "message": "Post liked successfully",
  "data": null
}
```

**Error Responses:**
- `400`: User has already liked this post
- `401`: Authentication required
- `404`: Post not found

### Unlike Post

**DELETE** `/v1/posts/{id}/like`

🔒 **Authentication Required**

Remove a like from a post.

**Parameters:**
- `id` (path): Post ID

**Response (200):**
```json
{
  "success": true,
  "message": "Post unliked successfully",
  "data": null
}
```

**Error Responses:**
- `400`: User has not liked this post
- `401`: Authentication required
- `404`: Post not found

### Get Post Likes

**GET** `/v1/posts/{id}/likes`

Get users who liked a specific post with pagination.

**Parameters:**
- `id` (path): Post ID
- `limit` (query, optional): Number of results per page (default: 10)
- `offset` (query, optional): Number of results to skip (default: 0)

**Response (200):**
```json
{
  "success": true,
  "message": "Post likes retrieved successfully",
  "data": {
    "likes": [
      {
        "id": "uuid",
        "post_id": "uuid",
        "user_id": "uuid",
        "user": {
          "id": "uuid",
          "username": "username",
          "email": "user@example.com"
        },
        "created_at": "2023-01-01T00:00:00Z"
      }
    ],
    "total": 25,
    "limit": 10,
    "offset": 0
  }
}
```

### Get Post Like Statistics

**GET** `/v1/posts/{id}/like-stats`

Get like statistics for a post.

**Parameters:**
- `id` (path): Post ID

**Response (200):**
```json
{
  "success": true,
  "message": "Like stats retrieved successfully",
  "data": {
    "post_id": "uuid",
    "total_likes": 42
  }
}
```

### Check if User Liked Post

**GET** `/v1/posts/{id}/liked`

🔒 **Authentication Required**

Check if current user has liked a specific post.

**Parameters:**
- `id` (path): Post ID

**Response (200):**
```json
{
  "success": true,
  "message": "Like status retrieved successfully",
  "data": {
    "has_liked": true,
    "post_id": "uuid",
    "user_id": "uuid"
  }
}
```

---

## Tag Endpoints

### Create Tag

**POST** `/v1/tags`

🔒 **Authentication Required**

Create a new tag.

**Request Body:**
```json
{
  "name": "tag-name",
  "description": "Tag description"
}
```

### Get All Tags

**GET** `/v1/tags`

Retrieve all available tags.

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved tags",
  "data": [
    {
      "id": "uuid",
      "name": "tag-name",
      "description": "Tag description",
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

### Get Tag by ID

**GET** `/v1/tags/{id}`

Retrieve a specific tag by ID.

### Update Tag

**PUT** `/v1/tags/{id}`

🔒 **Authentication Required** | 👑 **Admin Only**

Update a tag. Admin access required.

### Delete Tag

**DELETE** `/v1/tags/{id}`

🔒 **Authentication Required** | 👑 **Admin Only**

Delete a tag. Admin access required.

---

## Workspace Endpoints

### Create Workspace

**POST** `/v1/workspaces`

🔒 **Authentication Required**

Create a new workspace.

### Get All Workspaces

**GET** `/v1/workspaces`

🔒 **Authentication Required**

Retrieve all workspaces.

### Get User Workspaces

**GET** `/v1/workspaces/me`

🔒 **Authentication Required**

Retrieve current user's workspaces.

### Get Workspace by ID

**GET** `/v1/workspaces/{id}`

🔒 **Authentication Required**

Retrieve a specific workspace.

### Update Workspace

**PUT** `/v1/workspaces/{id}`

🔒 **Authentication Required**

Update a workspace.

### Delete Workspace

**DELETE** `/v1/workspaces/{id}`

🔒 **Authentication Required**

Delete a workspace.

### Add Workspace Member

**POST** `/v1/workspaces/{id}/members`

🔒 **Authentication Required**

Add a member to a workspace.

### Get Workspace Members

**GET** `/v1/workspaces/{id}/members`

🔒 **Authentication Required**

Retrieve workspace members.

### Update Member Role

**PUT** `/v1/workspaces/{id}/members/{user_id}`

🔒 **Authentication Required**

Update a member's role in the workspace.

### Remove Workspace Member

**DELETE** `/v1/workspaces/{id}/members/{user_id}`

🔒 **Authentication Required**

Remove a member from the workspace.

---

## Page Endpoints

### Create Page

**POST** `/v1/pages`

🔒 **Authentication Required**

Create a new page within a workspace.

### Get Page

**GET** `/v1/pages/{id}`

🔒 **Authentication Required**

Retrieve a specific page.

### Update Page

**PUT** `/v1/pages/{id}`

🔒 **Authentication Required**

Update a page.

### Delete Page

**DELETE** `/v1/pages/{id}`

🔒 **Authentication Required**

Delete a page.

### Get Workspace Pages

**GET** `/v1/pages/workspace/{workspace_id}`

🔒 **Authentication Required**

Retrieve all pages in a workspace.

### Get Child Pages

**GET** `/v1/pages/children/{parent_id}`

🔒 **Authentication Required**

Retrieve child pages of a parent page.

---

## Chat Conversation Endpoints

### Get My Conversations

**GET** `/v1/chat/conversations`

🔒 **Authentication Required**

Retrieve all chat conversations for the authenticated user with pagination support.

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved conversations",
  "data": [
    {
      "id": "uuid",
      "title": "Conversation Title",
      "user_id": "user-uuid",
      "message_count": 5,
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total_items": 15,
    "offset": 0,
    "limit": 10,
    "total_pages": 2
  }
}
```

### Get Conversation by ID

**GET** `/v1/chat/conversations/{id}`

🔒 **Authentication Required**

Retrieve a specific chat conversation by ID. Users can only access their own conversations.

**Parameters:**
- `id` (path): Conversation ID

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved conversation",
  "data": {
    "id": "uuid",
    "title": "Conversation Title",
    "user_id": "user-uuid",
    "message_count": 5,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### Create Conversation

**POST** `/v1/chat/conversations`

🔒 **Authentication Required**

Create a new chat conversation.

**Request Body:**
```json
{
  "title": "New Conversation"
}
```

**Validation Rules:**
- `title`: Required, maximum 255 characters

**Response (201):**
```json
{
  "success": true,
  "message": "Successfully created conversation",
  "data": {
    "id": "uuid",
    "title": "New Conversation",
    "user_id": "user-uuid",
    "message_count": 0,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### Update Conversation

**PUT** `/v1/chat/conversations/{id}`

🔒 **Authentication Required**

Update an existing chat conversation. Users can only update their own conversations.

**Parameters:**
- `id` (path): Conversation ID

**Request Body:**
```json
{
  "title": "Updated Conversation Title"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Conversation updated successfully",
  "data": {
    "id": "uuid",
    "title": "Updated Conversation Title",
    "user_id": "user-uuid",
    "message_count": 5,
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-02T00:00:00Z"
  }
}
```

### Delete Conversation

**DELETE** `/v1/chat/conversations/{id}`

🔒 **Authentication Required**

Delete a chat conversation. Users can only delete their own conversations.

**Parameters:**
- `id` (path): Conversation ID

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully deleted conversation",
  "data": null
}
```

---

## Chat Message Endpoints

### Create Message

**POST** `/v1/chat/conversations/{id}/messages`

🔒 **Authentication Required**

Add a new message to a chat conversation.

**Parameters:**
- `id` (path): Conversation ID

**Request Body:**
```json
{
  "content": "Hello, how are you?",
  "role": "user"
}
```

**Validation Rules:**
- `content`: Required
- `role`: Required, must be "user" or "assistant"

**Response (201):**
```json
{
  "success": true,
  "message": "Message created successfully",
  "data": {
    "id": "uuid",
    "conversation_id": "conversation-uuid",
    "user_id": "user-uuid",
    "role": "user",
    "content": "Hello, how are you?",
    "created_at": "2023-01-01T00:00:00Z",
    "updated_at": "2023-01-01T00:00:00Z"
  }
}
```

### Get Messages

**GET** `/v1/chat/conversations/{id}/messages`

🔒 **Authentication Required**

Retrieve all messages from a specific conversation.

**Parameters:**
- `id` (path): Conversation ID

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 50)

**Response (200):**
```json
{
  "success": true,
  "message": "Successfully retrieved messages",
  "data": [
    {
      "id": "uuid",
      "conversation_id": "conversation-uuid",
      "user_id": "user-uuid",
      "role": "user",
      "content": "Hello, how are you?",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ],
  "meta": {
    "total_items": 5,
    "offset": 0,
    "limit": 50,
    "total_pages": 1
  }
}
```

---

## Debug Endpoints (Development Only)

These endpoints are only available when `DEBUG=true` in the configuration.

### Performance Profiling

**GET** `/v1/debug/pprof/*`

Access Go pprof profiling endpoints for performance analysis.

Available endpoints:
- `/v1/debug/pprof/` - Index
- `/v1/debug/pprof/cmdline` - Command line
- `/v1/debug/pprof/profile` - CPU profile
- `/v1/debug/pprof/symbol` - Symbol lookup
- `/v1/debug/pprof/trace` - Execution trace
- `/v1/debug/pprof/heap` - Heap profile
- `/v1/debug/pprof/goroutine` - Goroutine profile
- `/v1/debug/pprof/allocs` - Memory allocations
- `/v1/debug/pprof/block` - Block profile
- `/v1/debug/pprof/mutex` - Mutex profile

---

## Error Handling

### Common Error Responses

**Validation Error (422):**
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Validation error details",
  "data": {
    "field_name": ["error message"]
  }
}
```

**Authentication Error (401):**
```json
{
  "success": false,
  "message": "Authentication required",
  "error": "Unauthorized access"
}
```

**Authorization Error (403):**
```json
{
  "success": false,
  "message": "Access denied",
  "error": "Access forbidden"
}
```

**Not Found Error (404):**
```json
{
  "success": false,
  "message": "Resource not found",
  "error": "Resource not found"
}
```

**Server Error (500):**
```json
{
  "success": false,
  "message": "Internal server error",
  "error": "Error details"
}
```

---

## Rate Limiting

Certain endpoints have rate limiting applied to prevent abuse:

- **Login endpoint**: 5 requests per 5 minutes per IP address
- **Forgot password endpoint**: 5 requests per 5 minutes per IP address
- **Register endpoint**: No specific rate limiting (handled by general API limits)
- **General API**: Default rate limits apply based on server configuration

Rate limits are subject to change and may be adjusted based on server load and usage patterns.

---

## Data Models

### User Model
```json
{
  "id": "uuid",
  "email": "user@example.com",
  "name": "First Last",
  "username": "username",
  "image": "profile-image-url",
  "first_name": "First",
  "last_name": "Last",
  "is_super_admin": false,
  "followers_count": 0,
  "following_count": 0,
  "is_following": null,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Post Model
```json
{
  "id": "uuid",
  "title": "Post Title",
  "photo_url": "image-url",
  "body": "Post content",
  "slug": "post-slug",
  "view_count": 0,
  "like_count": 0,
  "creator": {
    "id": "uuid",
    "username": "author",
    "name": "Author Name"
  },
  "tags": [
    {
      "id": "uuid",
      "name": "tag-name"
    }
  ],
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Tag Model
```json
{
  "id": "uuid",
  "name": "tag-name",
  "description": "Tag description",
  "created_at": "2023-01-01T00:00:00Z"
}
```

### PostLike Model
```json
{
  "id": "uuid",
  "post_id": "uuid",
  "user_id": "uuid",
  "user": {
    "id": "uuid",
    "username": "username",
    "email": "user@example.com",
    "name": "User Name"
  },
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "deleted_at": null
}
```

### ChatConversation Model
```json
{
  "id": "uuid",
  "title": "Conversation Title",
  "user_id": "user-uuid",
  "message_count": 5,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### ChatMessage Model
```json
{
  "id": "uuid",
  "conversation_id": "conversation-uuid",
  "user_id": "user-uuid",
  "role": "user",
  "content": "Message content",
  "model": "gpt-3.5-turbo",
  "prompt_tokens": 10,
  "completion_tokens": 20,
  "total_tokens": 30,
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

---

## Best Practices

### Request Headers
Always include appropriate headers:
```
Content-Type: application/json
Authorization: Bearer <jwt-token>
```

### Pagination
Use `offset` and `limit` parameters for pagination:
- Default `limit`: 10
- Default `offset`: 0
- Maximum `limit`: Varies by endpoint

### Error Handling
Always check the `success` field in responses and handle errors appropriately.

### Authentication
Store JWT tokens securely and include them in the Authorization header for protected endpoints.

### Rate Limiting
Respect rate limits and implement appropriate retry logic with exponential backoff.

---

## Changelog

### Current Version
- **Updated API Documentation**: Comprehensive review and update of all API endpoints
- **Added Missing Endpoints**: Added forgot password, reset password, refresh token, and change password endpoints
- **Enhanced Response Examples**: Added detailed response examples for all endpoints
- **Improved Authentication Requirements**: Updated authentication requirements for all endpoints
- **Removed Non-existent Features**: Removed WebSocket documentation (not implemented)
- **Updated Rate Limiting**: Added comprehensive rate limiting information

### Version 1.0.0
- Initial API release
- User authentication and management
- Post creation and management
- Comment system
- Tag system
- Workspace and page management
- User follow system
- Post view tracking
- Chat conversations and messages
- Post likes and views tracking
- Enhanced security with rate limiting

---

*This documentation is generated based on the Echo Backend API codebase analysis.*

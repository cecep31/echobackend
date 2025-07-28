# Backend API Documentation

## Overview

This is a comprehensive REST API for a blog/content management system. The API provides endpoints for user management, authentication, posts, comments, tags, workspaces, and more.

## Base URL

```
http://echo.pilput.me/v1
http://nest.pilput.me/v1
```

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

### Unfollow User

**DELETE** `/v1/users/{id}/follow`

🔒 **Authentication Required**

Unfollow a user.

**Parameters:**
- `id` (path): User ID to unfollow

### Check Follow Status

**GET** `/v1/users/{id}/follow-status`

🔒 **Authentication Required**

Check if current user follows the specified user.

### Get Mutual Follows

**GET** `/v1/users/{id}/mutual-follows`

🔒 **Authentication Required**

Get users that both current user and specified user follow.

### Get User Followers

**GET** `/v1/users/{id}/followers`

Get list of users following the specified user.

### Get User Following

**GET** `/v1/users/{id}/following`

Get list of users that the specified user follows.

### Get Follow Statistics

**GET** `/v1/users/{id}/follow-stats`

Get follower and following counts for a user.

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

Retrieve paginated list of all posts.

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

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

### Update Post

**PUT** `/v1/posts/{id}`

🔒 **Authentication Required** | 👤 **Author Only**

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

### Delete Post

**DELETE** `/v1/posts/{id}`

🔒 **Authentication Required** | 👤 **Author Only**

Delete a post. Only the post author can delete.

**Parameters:**
- `id` (path): Post ID

### Get Random Posts

**GET** `/v1/posts/random`

Retrieve random posts for discovery.

**Query Parameters:**
- `limit` (optional): Number of posts to return (default: 9, max: 50)

### Get My Posts

**GET** `/v1/posts/mine`

🔒 **Authentication Required**

Retrieve current user's posts.

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

### Get Posts by Username

**GET** `/v1/posts/username/{username}`

Retrieve posts by a specific user.

**Parameters:**
- `username` (path): Author's username

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

### Get Posts by Tag

**GET** `/v1/posts/tag/{tag}`

Retrieve posts with a specific tag.

**Parameters:**
- `tag` (path): Tag name

**Query Parameters:**
- `offset` (optional): Number of records to skip (default: 0)
- `limit` (optional): Number of records to return (default: 10)

### Upload Post Image

**POST** `/v1/posts/image`

🔒 **Authentication Required**

Upload an image for use in posts.

**Request:** Multipart form data
- `image` (file): Image file to upload

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
      "content": "Comment content",
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
  "content": "Comment content here..."
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Comment created successfully",
  "data": {
    "id": "uuid",
    "content": "Comment content",
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
  "content": "Updated comment content..."
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

## WebSocket Gateway

### Connection

**WebSocket** `wss://nest.pilput.me/ws/posts`

Connect to the WebSocket gateway for real-time post interactions. Requires authentication via JWT token and a post ID.

**Connection Parameters:**
- `token` (query): JWT authentication token
- `post_id` (query): ID of the post to connect to

**Example Connection URL:**
```
wss://nest.pilput.me/ws/posts?token=your-jwt-token&post_id=post-uuid
```

### Authentication

WebSocket connections require a valid JWT token passed as a query parameter. The token is verified upon connection.

### Events

#### Client to Server Events

**`sendComment`**
Send a new comment to a post.

**Payload:**
```json
{
  "content": "Comment content here..."
}
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": "comment-uuid",
      "content": "Comment content",
      "created_at": "2023-01-01T00:00:00Z",
      "creator": {
        "id": "user-uuid",
        "username": "commenter",
        "email": "commenter@example.com",
        "first_name": "Commenter",
        "last_name": "User"
      }
    }
  ]
}
```

**`typing`**
Send typing indicator status.

**Payload:**
```json
{
  "postId": "post-uuid",
  "isTyping": true
}
```

**Response:**
```json
{
  "status": "success"
}
```

**`markAsRead`**
Mark a comment as read.

**Payload:**
```json
{
  "commentId": "comment-uuid"
}
```

**Response:**
```json
{
  "status": "success"
}
```

**`getAllComments`**
Request all comments for the post.

**Response:**
```json
{
  "status": "success"
}
```

#### Server to Client Events

**`newComment`**
Broadcast when new comments are available or updated.

**Payload:**
```json
[
  {
    "id": "comment-uuid",
    "content": "Comment content",
    "created_at": "2023-01-01T00:00:00Z",
    "creator": {
      "id": "user-uuid",
      "username": "commenter",
      "email": "commenter@example.com",
      "first_name": "Commenter",
      "last_name": "User"
    }
  }
]
```

**`userTyping`**
Broadcast when a user starts or stops typing.

**Payload:**
```json
{
  "userId": "user-uuid",
  "isTyping": true
}
```

**`error`**
Sent when an error occurs during connection or processing.

**Payload:**
```json
{
  "status": "error",
  "message": "Error description"
}
```

### WebSocket Connection Lifecycle

1. **Connection**: Client connects with token and post_id parameters
2. **Authentication**: Server verifies JWT token
3. **Initialization**: Server loads initial comments and sends them via `newComment` event
4. **Real-time Interaction**: Client and server exchange events for comments and typing indicators
5. **Disconnection**: Client or server terminates connection

### Error Handling

WebSocket errors are sent via the `error` event. Common errors include:
- Invalid or missing authentication token
- Invalid post ID
- Unauthorized access
- Connection timeout

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

Certain endpoints have rate limiting applied:

- **Login endpoint**: 5 requests per 5 minutes per IP address

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

### Version 1.0.0
- Initial API release
- User authentication and management
- Post creation and management
- Comment system
- Tag system
- Workspace and page management
- User follow system
- Post view tracking

---

*This documentation is generated based on the Echo Backend API codebase analysis.*

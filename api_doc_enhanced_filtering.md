# Enhanced Post Query Filtering API Documentation

## Overview

The `/v1/posts` endpoint supports advanced query filtering to provide more flexible and powerful post retrieval capabilities.

## Base Endpoint

```
GET /v1/posts
```

**Authentication:** Not required for published posts, but required for unpublished posts when using `published=false`

## Query Parameters

| Parameter | Type | Description | Default | Example |
|-----------|------|-------------|---------|---------|
| `limit` | `integer` | Number of posts to return (max: 100) | `10` | `?limit=20` |
| `offset` | `integer` | Number of posts to skip for pagination | `0` | `?offset=10` |
| `search` | `string` | Search term to match in title or body | `""` | `?search=golang` |
| `sort_by` | `string` | Field to sort by | `created_at` | `?sort_by=title` |
| `sort_order` | `string` | Sort direction | `desc` | `?sort_order=asc` |
| `start_date` | `string` | Filter posts created after this date (YYYY-MM-DD) | `""` | `?start_date=2023-01-01` |
| `end_date` | `string` | Filter posts created before this date (YYYY-MM-DD) | `""` | `?end_date=2023-12-31` |
| `created_by` | `string` | Filter posts by author user ID | `""` | `?created_by=550e8400-e29b-41d4-a716-446655440000` |
| `published` | `boolean` | Filter by publication status | `true` (when no search) | `?published=false` |
| `tags` | `string` | Comma-separated list of tags | `""` | `?tags=golang,webdev` |

**Note:** When `search` parameter is provided, only published posts are returned regardless of the `published` parameter value.

## Sort Options

### Sortable Fields (`sort_by`)

- `id` - Sort by post ID
- `title` - Sort by post title
- `created_at` - Sort by creation date (default)
- `updated_at` - Sort by last update date
- `view_count` - Sort by view count
- `like_count` - Sort by like count

### Sort Orders (`sort_order`)

- `asc` - Ascending order
- `desc` - Descending order (default)

## Usage Examples

### 1. Basic Pagination
```bash
GET /api/v1/posts?limit=10&offset=0
```

### 2. Search Posts
```bash
GET /api/v1/posts?search=api%20development&limit=20
```

### 3. Sort by Title (Ascending)
```bash
GET /api/v1/posts?sort_by=title&sort_order=asc&limit=15
```

### 4. Filter by Date Range
```bash
GET /api/v1/posts?start_date=2023-01-01&end_date=2023-12-31&limit=25
```

### 5. Filter by Author
```bash
GET /api/v1/posts?created_by=550e8400-e29b-41d4-a716-446655440000&limit=10
```

### 6. Filter by Tags
```bash
GET /api/v1/posts?tags=javascript,react,frontend&limit=20
```

### 7. Filter by Publication Status
```bash
GET /api/v1/posts?published=false&limit=5
```

### 8. Complex Filtering
```bash
GET /api/v1/posts?search=web&sort_by=view_count&sort_order=desc&start_date=2023-01-01&tags=golang,api&limit=15&offset=5
```

## Response Format

```json
{
  "success": true,
  "message": "Successfully retrieved posts",
  "data": {
    "posts": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "title": "Getting Started with Golang",
        "photo_url": "https://example.com/image.jpg",
        "body": "Go is a statically typed, compiled programming language...",
        "slug": "getting-started-with-golang",
        "view_count": 1540,
        "like_count": 45,
        "published": true,
        "creator": {
          "id": "550e8400-e29b-41d4-a716-446655440001",
          "username": "johndoe",
          "name": "John Doe"
        },
        "tags": [
          {
            "id": "550e8400-e29b-41d4-a716-446655440002",
            "name": "golang",
            "slug": "golang"
          }
        ],
        "created_at": "2023-06-15T10:30:00Z",
        "updated_at": "2023-06-15T10:30:00Z",
        "deleted_at": null
      }
    ],
    "meta": {
      "total_items": 150,
      "offset": 0,
      "limit": 10,
      "total_pages": 15
    }
  }
}
```

## Query Logic

### Search Functionality
- Searches in both `title` and `body` fields using case-insensitive matching
- Uses `ILIKE` operator with `%term%` pattern
- When search is provided, only published posts are returned regardless of `published` parameter

### Date Filtering
- `start_date`: Filters posts created on or after the specified date
- `end_date`: Filters posts created on or before the specified date
- Date format: `YYYY-MM-DD` (e.g., `2023-12-31`)
- **Note:** Date filtering applies to `created_at` field

### Tag Filtering
- Accepts comma-separated list of tag names
- Posts must contain ANY of the specified tags (OR logic, not AND as previously documented)
- Example: `tags=golang,webdev` returns posts that have either "golang" OR "webdev" tags
- Whitespace around tag names is automatically trimmed

### Publication Status
- When no search term is provided and `published` is not specified: defaults to `published=true`
- When search term is provided: only searches published posts regardless of `published` parameter
- Can be explicitly set to `true` or `false` to override defaults (requires authentication for `false`)

### Sorting Options
**Valid sort fields (`sort_by`):**
- `id` - Sort by post ID
- `title` - Sort by post title
- `created_at` - Sort by creation date (default)
- `updated_at` - Sort by last update date
- `view_count` - Sort by view count
- `like_count` - Sort by like count

**Valid sort orders (`sort_order`):**
- `asc` - Ascending order
- `desc` - Descending order (default)

### Sorting Priority
1. Primary sort field (specified by `sort_by`)
2. Secondary sort by `created_at` (for consistent pagination)

## Performance Considerations

- All query parameters are optional
- Use `limit` to control response size (max: 100, default: 10)
- Use `offset` for pagination (default: 0)
- Complex queries with multiple filters may impact performance
- Date range queries are optimized for time-based indexing
- Search queries are case-insensitive and may be slower on large datasets

## Error Handling

The endpoint returns standard HTTP status codes:
- `200 OK` - Successful request
- `500 Internal Server Error` - Server error during query processing

Error responses follow the standard API error format:
```json
{
  "success": false,
  "message": "Failed to get posts",
  "error": "Database query failed",
  "data": null
}
```

## Integration Examples

### JavaScript (Fetch API)
```javascript
const response = await fetch('/v1/posts?search=api&sort_by=created_at&limit=20');
const data = await response.json();
console.log(data.data.posts);
```

### Python (Requests)
```python
import requests

# For published posts (no auth required)
params = {
    'search': 'web development',
    'sort_by': 'view_count',
    'sort_order': 'desc',
    'limit': 10
}

response = requests.get('https://echo.pilput.me/v1/posts', params=params)
posts = response.json()['data']['posts']

# For unpublished posts (auth required)
headers = {'Authorization': 'Bearer your-jwt-token'}
response = requests.get('https://echo.pilput.me/v1/posts?published=false', headers=headers)
```

### cURL
```bash
# Get published posts with filtering
curl "https://echo.pilput.me/v1/posts?search=golang&tags=web,api&sort_by=like_count&sort_order=desc&limit=15"

# Get unpublished posts (requires authentication)
curl -H "Authorization: Bearer your-jwt-token" "https://echo.pilput.me/v1/posts?published=false"
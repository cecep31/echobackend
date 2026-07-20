# Tags Module - `/api/tags`

Tag management for posts. Create requires login; update/delete requires **super admin** access.

| Method | Path | Auth |
|--------|------|------|
| POST | `` | Bearer |
| GET | `` | No |
| GET | `/trending` | No |
| GET | `/sitemap` | No |
| GET | `/:id` | No |
| PUT | `/:id` | Bearer + super admin |
| DELETE | `/:id` | Bearer + super admin |

## Data Types

### `TagResponse` (all tag endpoints: list, create, get by id, update, nested in posts)

```json
{ "id": 1, "name": "golang" }
```

Create, get-by-id, and update all return this shape. The tag model's `created_at` timestamp is not exposed by these endpoints.

### `TrendingTagResponse`

```json
{
  "id": 1,
  "name": "golang",
  "total_views": 1500,
  "total_likes": 80,
  "trending_score": 1710
}
```

`trending_score` is calculated from published post aggregation: `like_count * 2 + bookmark_count * 2 + view_count`.

### `SitemapTag`

```json
{ "name": "golang", "created_at": "2026-01-01T00:00:00Z" }
```

---

## POST `/api/tags`

**Body (`CreateTagRequest`)**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 1-30 characters |

**Success - 201** - `data`: `TagResponse` (`{ "id", "name" }`).

**Errors**

| HTTP | Condition |
|------|-----------|
| 400 | Invalid JSON body |
| 422 | Validation failed (`name` missing or not 1-30 characters) |
| 500 | Server error |

---

## GET `/api/tags`

**Success - 200** - `data`: `TagResponse[]`.

---

## GET `/api/tags/trending`

Returns the 5 most trending tags from published posts. This endpoint does not accept a `limit` query parameter.

**Cache:** the response is cached in Valkey/Redis for 30 minutes with key `tags:trending` (using the configured cache prefix, if any).

**Success - 200** - `data`: `TrendingTagResponse[]`.

Example response:

```json
{
  "success": true,
  "message": "Successfully retrieved trending tags",
  "data": [
    {
      "id": 1,
      "name": "golang",
      "total_views": 1500,
      "total_likes": 80,
      "trending_score": 1710
    }
  ]
}
```

---

## GET `/api/tags/sitemap`

**Success - 200** - `data`: `SitemapTag[]`.

---

## GET `/api/tags/:id`

**Path:** numeric `id`.

| HTTP | Condition |
|------|-----------|
| 400 | Invalid ID |
| 404 | Tag not found |
| 200 | `data`: tag |

---

## PUT `/api/tags/:id`

**Body (`UpdateTagRequest`)**

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | Yes | 1-30 characters |

Super admin only.

**Success - 200** - `data`: `TagResponse`.

| HTTP | Condition |
|------|-----------|
| 400 | Invalid ID |
| 404 | Tag not found |
| 422 | Validation failed |

---

## DELETE `/api/tags/:id`

**Success - 200** - `data`: `null`.

| HTTP | Condition |
|------|-----------|
| 400 | Invalid ID |
| 404 | Tag not found |

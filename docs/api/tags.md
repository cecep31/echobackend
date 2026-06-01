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

### `TagResponse` (list / nested in posts)

```json
{ "id": 1, "name": "golang" }
```

### Tag model (create / get by id / update)

| Field | Type |
|-------|------|
| `id` | number |
| `name` | string |
| `created_at` | string (optional in response) |

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

**Body:** tag JSON - primary field is `name` (string). The handler does not have struct-tag validation; business errors return 500.

**Success - 201** - `data`: tag object.

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

**Body:** tag object (for example, a new `name`). Super admin only.

---

## DELETE `/api/tags/:id`

**Success - 200** - `data`: `null`.

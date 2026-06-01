# Bookmarks Module - `/api/bookmarks`

Post bookmarking and folder organization. **All routes require a Bearer token.**

## Route Summary

| Method | Path | Description |
|--------|------|-------------|
| POST | `/:post_id` | Toggle bookmark (create/delete) |
| GET | `` | List user bookmarks |
| PATCH | `/:bookmark_id` | Update name/notes |
| PATCH | `/:bookmark_id/move` | Move to folder |
| POST | `/folders` | Create folder |
| GET | `/folders` | List folders |
| PATCH | `/folders/:folder_id` | Update folder |
| DELETE | `/folders/:folder_id` | Delete folder |

---

## Data Types

### `BookmarkResponse`

| Field | Type |
|-------|------|
| `id` | string (UUID) |
| `post_id` | string (UUID) |
| `user_id` | string (UUID) |
| `folder_id` | string (UUID) \| null |
| `name` | string \| null |
| `notes` | string \| null |
| `post` | `PostResponse` \| null |
| `folder` | `BookmarkFolderResponse` \| null |
| `created_at` | string \| null |
| `updated_at` | string \| null |

### `BookmarkFolderResponse`

| Field | Type |
|-------|------|
| `id` | string (UUID) |
| `user_id` | string (UUID) |
| `name` | string |
| `description` | string \| null |
| `bookmark_count` | number |
| `created_at` | string \| null |
| `updated_at` | string \| null |

### `ToggleBookmarkResponse`

| Field | Type |
|-------|------|
| `action` | string (`added` / `removed`) |
| `bookmark` | `BookmarkResponse` \| null |

---

## POST `/api/bookmarks/:post_id`

Toggle a bookmark on a post. If it already exists, it is removed; otherwise it is created.

**Body (`ToggleBookmarkRequest`)** - all fields are optional

| Field | Validation |
|-------|------------|
| `folder_id` | UUID |
| `name` | max 255 |
| `notes` | max 2000 |

**Success - 200** - `data`: `ToggleBookmarkResponse`.

---

## GET `/api/bookmarks`

**Query**

| Param | Description |
|-------|-------------|
| `folder_id` | Folder UUID; `null` = bookmarks without a folder |
| `limit`, `offset` | Pagination (default limit 50, max 100) |

**Success - 200** - `data`: `BookmarkResponse[]`, `meta`: pagination.

---

## PATCH `/api/bookmarks/:bookmark_id`

**Body (`UpdateBookmarkRequest`)**

| Field | Validation |
|-------|------------|
| `name` | max 255 |
| `notes` | max 2000 |

**Success - 200** - `data`: `BookmarkResponse`.

---

## PATCH `/api/bookmarks/:bookmark_id/move`

**Body**

| Field | Required | Validation |
|-------|----------|------------|
| `folder_id` | No | UUID; `null` = remove from folder |

**Success - 200** - `data`: `BookmarkResponse`.

---

## POST `/api/bookmarks/folders`

**Body**

| Field | Required | Validation |
|-------|----------|------------|
| `name` | Yes | 1-100 characters |
| `description` | No | max 1000 |

**Success - 201** - `data`: `BookmarkFolderResponse`.

---

## GET `/api/bookmarks/folders`

**Success - 200** - `data`: `BookmarkFolderResponse[]`.

---

## PATCH `/api/bookmarks/folders/:folder_id`

**Body** - optional fields: `name` (1-100), `description` (max 1000).

**Success - 200** - `data`: `BookmarkFolderResponse`.

---

## DELETE `/api/bookmarks/folders/:folder_id`

**Success - 200** - `data`: `null`.

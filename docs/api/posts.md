# Modul Posts — `/api/posts`

CRUD post, upload gambar, feed, komentar, view, dan like. Sub-resource memakai path `:id` = UUID post.

## Tipe data: `PostResponse`

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | string (UUID) | |
| `title` | string \| null | |
| `photo_url` | string \| null | |
| `body` | string \| null | List feed dipotong ~250 rune + `" ..."` |
| `slug` | string \| null | |
| `view_count` | number | |
| `like_count` | number | |
| `bookmark_count` | number | |
| `published` | boolean \| null | |
| `published_at` | string \| null | |
| `user` | `UserResponse` \| null | Penulis |
| `tags` | `TagResponse[]` | `{ id, name }` |
| `created_at` | string \| null | |
| `updated_at` | string \| null | |
| `deleted_at` | string \| null | Soft delete |

### `TagResponse`

`{ "id": number, "name": string }`

---

## CRUD & listing

| Method | Path | Auth |
|--------|------|------|
| POST | `` | Bearer |
| GET | `` | Tidak |
| GET | `/random` | Tidak |
| GET | `/trending` | Tidak |
| GET | `/mine` | Bearer |
| GET | `/for-you` | Bearer |
| POST | `/image` | Bearer |
| GET | `/sitemap` | Tidak |
| GET | `/username/:username` | Tidak |
| GET | `/u/:username/:slug` | Tidak |
| GET | `/tag/:tag` | Tidak |
| GET | `/:id` | Tidak |
| PUT | `/:id` | Bearer |
| DELETE | `/:id` | Bearer |

### POST `/api/posts`

**Body (`CreatePostRequest`)**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `title` | string | Ya | min 7 |
| `slug` | string | Ya | min 7 |
| `body` | string | Ya | min 10 |
| `photo_url` | string | Tidak | |
| `published` | boolean | Tidak | default false |
| `tags` | string[] | Tidak | Nama tag |

**Sukses — 201** — `data`: `{ "id": "uuid" }`.

### GET `/api/posts`

**Query**

| Param | Keterangan |
|-------|------------|
| `search` | Teks pencarian |
| `sort_by` | `id`, `title`, `created_at`, `updated_at`, `view_count`, `like_count` |
| `sort_order` | `asc` / `desc` (default `desc`) |
| `start_date`, `end_date` | Filter tanggal |
| `created_by` | UUID penulis |
| `published` | `true` / `false` |
| `tags` | Comma-separated |
| `limit`, `offset` | Paginasi (default limit 10) |

**Sukses — 200** — `data`: `PostResponse[]`, `meta`: paginasi.

### GET `/api/posts/random`

**Query:** `limit` (default 9, max 20).

### GET `/api/posts/trending`

**Query:** `limit` (default 10).

### GET `/api/posts/mine` dan `/for-you`

**Query:** `limit`, `offset`. **Auth wajib.**

### POST `/api/posts/image`

**Content-Type:** `multipart/form-data`

| Field | Wajib |
|-------|-------|
| `image` | Ya (file) |

**Sukses — 200** — URL/file tersimpan (lihat implementasi service; `data` dapat `null` dengan message sukses).

**Error:** 400 jika file kosong atau storage tidak tersedia.

### GET `/api/posts/sitemap`

**Sukses — 200** — `data`: `SitemapPost[]`:

```json
{ "username": "...", "slug": "...", "created_at": "...", "updated_at": "..." }
```

### GET `/api/posts/u/:username/:slug`

Detail penuh satu post (body tidak dipotong).

### PUT `/api/posts/:id`

**Body (`UpdatePostRequest`)** — semua field opsional; `published` boolean pointer.

**Error umum**

| HTTP | Kondisi |
|------|---------|
| 403 | Bukan penulis |
| 404 | Post tidak ada |

### DELETE `/api/posts/:id`

**Sukses — 200** — `data`: `null`.

---

## Komentar — `/api/posts/:id/comments`

| Method | Path | Auth |
|--------|------|------|
| GET | `/:id/comments` | Tidak |
| POST | `/:id/comments` | Bearer |
| PUT | `/:id/comments/:comment_id` | Bearer |
| DELETE | `/:id/comments/:comment_id` | Bearer |

### `CommentResponse`

| Field | Tipe |
|-------|------|
| `id` | string (UUID) |
| `post_id` | string |
| `parent_comment_id` | string \| null |
| `text` | string |
| `user` | `UserResponse` \| null |
| `created_at` | string \| null |
| `updated_at` | string \| null |

### POST — body

| Field | Wajib | Validasi |
|-------|-------|----------|
| `text` | Ya | 1–1000 karakter |

**Sukses — 201** — `data`: `CommentResponse`.

---

## View — `/api/posts/:id/view*`

| Method | Path | Auth |
|--------|------|------|
| POST | `/:id/view` | Bearer |
| GET | `/:id/views` | Bearer |
| GET | `/:id/view-stats` | Tidak |
| GET | `/:id/viewed` | Bearer |

### POST `/api/posts/:id/view`

Catat view untuk user login.

### GET `/api/posts/:id/view-stats`

**Sukses — 200** — `data`:

```json
{
  "post_id": "uuid",
  "total_views": 0,
  "unique_views": 0,
  "anonymous_views": 0,
  "authenticated_views": 0
}
```

### GET `/api/posts/:id/viewed`

**Sukses — 200** — `data`: `{ "has_viewed": true }`.

### GET `/api/posts/:id/views`

List view (model internal: `id`, `post_id`, `user_id`, `ip_address`, `user_agent`, timestamps) + `meta` paginasi.

---

## Like — `/api/posts/:id/like*`

| Method | Path | Auth |
|--------|------|------|
| POST | `/:id/like` | Bearer |
| DELETE | `/:id/like` | Bearer |
| GET | `/:id/likes` | Tidak |
| GET | `/:id/like-stats` | Tidak |
| GET | `/:id/liked` | Bearer |

### POST / DELETE like

**Error 400** jika sudah like / belum like.

### GET `/api/posts/:id/likes`

**Sukses — 200** — `data`:

```json
{
  "likes": [ /* PostLikeResponse */ ],
  "total": 0,
  "limit": 10,
  "offset": 0
}
```

`PostLikeResponse`: `id`, `post_id`, `user_id`, `user`, `created_at`.

### GET `/api/posts/:id/like-stats`

`data`: `{ "post_id", "total_likes" }`.

### GET `/api/posts/:id/liked`

`data`: `{ "has_liked", "post_id", "user_id" }`.

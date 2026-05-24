# Modul Bookmarks — `/api/bookmarks`

Bookmark post dan organisasi folder. **Semua route membutuhkan Bearer token.**

## Ringkasan route

| Method | Path | Keterangan |
|--------|------|------------|
| POST | `/:post_id` | Toggle bookmark (buat/hapus) |
| GET | `` | List bookmark user |
| PATCH | `/:bookmark_id` | Update nama/catatan |
| PATCH | `/:bookmark_id/move` | Pindah ke folder |
| POST | `/folders` | Buat folder |
| GET | `/folders` | List folder |
| PATCH | `/folders/:folder_id` | Update folder |
| DELETE | `/folders/:folder_id` | Hapus folder |

---

## Tipe data

### `BookmarkResponse`

| Field | Tipe |
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

| Field | Tipe |
|-------|------|
| `id` | string (UUID) |
| `user_id` | string (UUID) |
| `name` | string |
| `description` | string \| null |
| `bookmark_count` | number |
| `created_at` | string \| null |
| `updated_at` | string \| null |

### `ToggleBookmarkResponse`

| Field | Tipe |
|-------|------|
| `action` | string (`added` / `removed`) |
| `bookmark` | `BookmarkResponse` \| null |

---

## POST `/api/bookmarks/:post_id`

Toggle bookmark pada post. Jika sudah ada, dihapus; jika belum, dibuat.

**Body (`ToggleBookmarkRequest`)** — semua opsional

| Field | Validasi |
|-------|----------|
| `folder_id` | UUID |
| `name` | max 255 |
| `notes` | max 2000 |

**Sukses — 200** — `data`: `ToggleBookmarkResponse`.

---

## GET `/api/bookmarks`

**Query**

| Param | Keterangan |
|-------|------------|
| `folder_id` | UUID folder; `null` = bookmark tanpa folder |
| `limit`, `offset` | Paginasi (default limit 50, max 100) |

**Sukses — 200** — `data`: `BookmarkResponse[]`, `meta`: paginasi.

---

## PATCH `/api/bookmarks/:bookmark_id`

**Body (`UpdateBookmarkRequest`)**

| Field | Validasi |
|-------|----------|
| `name` | max 255 |
| `notes` | max 2000 |

**Sukses — 200** — `data`: `BookmarkResponse`.

---

## PATCH `/api/bookmarks/:bookmark_id/move`

**Body**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `folder_id` | Tidak | UUID; `null` = keluarkan dari folder |

**Sukses — 200** — `data`: `BookmarkResponse`.

---

## POST `/api/bookmarks/folders`

**Body**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `name` | Ya | 1–100 karakter |
| `description` | Tidak | max 1000 |

**Sukses — 201** — `data`: `BookmarkFolderResponse`.

---

## GET `/api/bookmarks/folders`

**Sukses — 200** — `data`: `BookmarkFolderResponse[]`.

---

## PATCH `/api/bookmarks/folders/:folder_id`

**Body** — field opsional: `name` (1–100), `description` (max 1000).

**Sukses — 200** — `data`: `BookmarkFolderResponse`.

---

## DELETE `/api/bookmarks/folders/:folder_id`

**Sukses — 200** — `data`: `null`.

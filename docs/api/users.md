# Modul Users & Follow — `/api/users`

Profil user, daftar admin, follow/unfollow, dan statistik sosial.

## Tipe data: `UserResponse`

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | string (UUID) | |
| `email` | string | |
| `name` | string | Gabungan first + last name |
| `username` | string \| null | |
| `image` | string \| null | URL avatar |
| `first_name` | string \| null | |
| `last_name` | string \| null | |
| `followers_count` | number | |
| `following_count` | number | |
| `is_following` | boolean \| null | Hanya terisi pada route yang memuat konteks auth (mis. admin `GET /:id`) |
| `is_super_admin` | boolean \| null | Hanya pada route admin (`GET /`, `GET /:id`) |
| `profile` | object \| null | Tidak di-load pada `GET /` (admin list); ada pada route lain |
| `created_at` | string (ISO) \| null | |
| `updated_at` | string (ISO) \| null | |
| `deleted_at` | string (ISO) \| null | Hanya pada route admin; terisi jika user sudah di-soft-delete |

### `Profile`

| Field | Tipe |
|-------|------|
| `id` | number |
| `user_id` | string (UUID) |
| `bio` | string \| null |
| `website` | string \| null |
| `phone` | string \| null |
| `location` | string \| null |
| `created_at` | string \| null |
| `updated_at` | string \| null |

---

## Profil & admin

| Method | Path | Auth | Keterangan |
|--------|------|------|------------|
| GET | `/:id` | Bearer + **super admin** | By UUID |
| GET | `/username/:username` | Tidak | By username |
| GET | `/me` | Bearer | User dari token |
| GET | `` | Bearer + **super admin** | List user (paginated); filter soft-delete via query |
| DELETE | `/:id` | Bearer + **super admin** | Soft-delete user |
| POST | `/:id/restore` | Bearer + **super admin** | Restore user yang sudah di-soft-delete |

### GET `/api/users/me`

**Sukses — 200** — `data`: satu `UserResponse`.

### GET `/api/users` (admin)

**Query**

| Param | Tipe | Default | Keterangan |
|-------|------|---------|------------|
| `limit` | number | 10 | Max 100 |
| `offset` | number | 0 | |
| `deleted` | string | *(kosong)* | Filter soft-delete: `false` atau kosong = hanya user aktif; `true` = hanya user terhapus; `all` = semua user |

Nilai `deleted` selain `true`, `false`, `all`, atau kosong → **400 Bad Request**.

**Sukses — 200** — `data`: array `UserResponse`, `meta`: paginasi.

Contoh:

```http
GET /api/users?deleted=true&limit=10&offset=0
GET /api/users?deleted=all
```

### GET `/api/users/:id` (admin)

**Query**

| Param | Tipe | Default | Keterangan |
|-------|------|---------|------------|
| `deleted` | string | *(kosong)* | `true` = ambil user yang sudah di-soft-delete by UUID |

Tanpa `deleted=true`, hanya user **aktif** yang ditemukan. User terhapus tidak muncul kecuali query di atas dipakai.

**Sukses — 200** — `data`: `UserResponse`. `is_following` terisi jika requester login (route admin sudah memakai auth). `deleted_at` terisi jika user sudah di-soft-delete.

Contoh:

```http
GET /api/users/<uuid>?deleted=true
```

### GET `/api/users/username/:username`

**Sukses — 200** — `data`: `UserResponse`. Route publik — `is_following` tidak terisi.

### DELETE `/api/users/:id` (admin)

Soft-delete user (set `deleted_at`; baris tidak dihapus permanen dari database).

**Sukses — 200** — `data`: `null`.

### POST `/api/users/:id/restore` (admin)

Mengembalikan user yang sudah di-soft-delete (`deleted_at` → `null`).

**Sukses — 200** — `data`: `UserResponse` user yang sudah aktif kembali.

**Error**

| Status | Kondisi |
|--------|---------|
| 404 | User tidak ditemukan atau belum di-soft-delete |
| 409 | Email atau username sudah dipakai user aktif lain |

Contoh:

```http
POST /api/users/<uuid>/restore
```

---

## Follow

| Method | Path | Auth |
|--------|------|------|
| POST | `/follow` | Bearer |
| DELETE | `/:id/follow` | Bearer |
| GET | `/:id/follow-status` | Bearer |
| GET | `/:id/mutual-follows` | Bearer |
| GET | `/:id/followers` | Tidak |
| GET | `/:id/following` | Tidak |
| GET | `/:id/follow-stats` | Tidak |

### POST `/api/users/follow`

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `user_id` | string | Ya | UUID |

**Sukses — 200** — `data`:

```json
{
  "is_following": true,
  "message": "Pesan dari service"
}
```

### DELETE `/api/users/:id/follow`

Unfollow user dengan UUID di path.

**Sukses — 200** — `data`: `FollowResponse` (bentuk sama seperti follow).

### GET `/api/users/:id/follow-status`

**Sukses — 200**

```json
{
  "data": { "is_following": false }
}
```

### GET `/api/users/:id/mutual-follows`

**Sukses — 200** — `data`: array `UserResponse`.

### GET `/api/users/:id/followers` dan `/:id/following`

**Query:** `limit`, `offset`.

**Sukses — 200** — `data`: array `UserResponse`, `meta`: paginasi.

### GET `/api/users/:id/follow-stats`

**Sukses — 200** — `data`:

```json
{
  "user_id": "uuid",
  "followers_count": 10,
  "following_count": 5
}
```

**Catatan:** Beberapa error domain follow (user tidak ditemukan, sudah follow, dll.) saat ini dapat mengembalikan **500** di layer handler, bukan 4xx khusus.

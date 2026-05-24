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
| GET | `` | Bearer + **super admin** | List semua user (paginated) |
| DELETE | `/:id` | Bearer + **super admin** | Hapus user |

### GET `/api/users/me`

**Sukses — 200** — `data`: satu `UserResponse`.

### GET `/api/users` (admin)

**Query:** `limit`, `offset` (default limit 10, max 100).

**Sukses — 200** — `data`: array `UserResponse`, `meta`: paginasi.

### GET `/api/users/:id` (admin)

**Sukses — 200** — `data`: `UserResponse`. `is_following` terisi jika requester login (route admin sudah memakai auth).

### GET `/api/users/username/:username`

**Sukses — 200** — `data`: `UserResponse`. Route publik — `is_following` tidak terisi.

### DELETE `/api/users/:id` (admin)

**Sukses — 200** — `data`: `null`.

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

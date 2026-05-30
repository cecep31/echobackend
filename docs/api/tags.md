# Modul Tags — `/api/tags`

Manajemen tag untuk post. Create butuh login; update/delete butuh **super admin**.

| Method | Path | Auth |
|--------|------|------|
| POST | `` | Bearer |
| GET | `` | Tidak |
| GET | `/trending` | Tidak |
| GET | `/sitemap` | Tidak |
| GET | `/:id` | Tidak |
| PUT | `/:id` | Bearer + super admin |
| DELETE | `/:id` | Bearer + super admin |

## Tipe data

### `TagResponse` (list / nested di post)

```json
{ "id": 1, "name": "golang" }
```

### Model tag (create / get by id / update)

| Field | Tipe |
|-------|------|
| `id` | number |
| `name` | string |
| `created_at` | string (opsional di respons) |

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

`trending_score` dihitung dari agregasi post published: `like_count * 2 + bookmark_count * 2 + view_count`.

### `SitemapTag`

```json
{ "name": "golang", "created_at": "2026-01-01T00:00:00Z" }
```

---

## POST `/api/tags`

**Body:** JSON tag — field utama `name` (string). Tidak ada validasi struct tag di handler; error bisnis → 500.

**Sukses — 201** — `data`: objek tag.

---

## GET `/api/tags`

**Sukses — 200** — `data`: array `TagResponse`.

---

## GET `/api/tags/trending`

Mengambil 5 tag paling trending dari post published. Endpoint ini tidak menerima query parameter limit.

**Cache:** hasil response di-cache di Valkey/Redis selama 30 menit dengan key `tags:trending` (mengikuti prefix cache jika dikonfigurasi).

**Sukses — 200** — `data`: array `TrendingTagResponse`.

Contoh response:

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

**Sukses — 200** — `data`: array `SitemapTag`.

---

## GET `/api/tags/:id`

**Path:** `id` numerik.

| HTTP | Kondisi |
|------|---------|
| 400 | ID tidak valid |
| 404 | Tag tidak ada |
| 200 | `data`: tag |

---

## PUT `/api/tags/:id`

**Body:** objek tag (mis. `name` baru). Hanya super admin.

---

## DELETE `/api/tags/:id`

**Sukses — 200** — `data`: `null`.

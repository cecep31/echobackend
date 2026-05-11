# Modul Tags — `/api/tags`

Manajemen tag untuk post. Create butuh login; update/delete butuh **super admin**.

| Method | Path | Auth |
|--------|------|------|
| POST | `` | Bearer |
| GET | `` | Tidak |
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

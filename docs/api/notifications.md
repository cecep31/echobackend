# Modul Notifications — `/api/notifications`

Notifikasi in-app untuk user login. **Semua route membutuhkan Bearer token.**

Notifikasi dibuat otomatis oleh sistem (mis. saat ada komentar baru atau follow baru). Tidak ada endpoint HTTP publik untuk membuat notifikasi.

## Ringkasan route

| Method | Path | Keterangan |
|--------|------|------------|
| GET | `` | List notifikasi |
| GET | `/unread-count` | Jumlah belum dibaca |
| PATCH | `/:id/read` | Tandai satu sebagai dibaca |
| PATCH | `/read-all` | Tandai semua sebagai dibaca |

---

## Tipe data: `NotificationResponse`

| Field | Tipe |
|-------|------|
| `id` | string (UUID) |
| `user_id` | string (UUID) |
| `type` | string |
| `title` | string |
| `message` | string \| null |
| `read` | boolean |
| `data` | object \| null |
| `created_at` | string \| null |
| `updated_at` | string \| null |

### Tipe notifikasi yang dipakai saat ini

| `type` | Pemicu |
|--------|--------|
| `comment` | Komentar baru pada post milik user |
| `follow` | User lain mengikuti akun |

---

## GET `/api/notifications`

**Query**

| Param | Default | Keterangan |
|-------|---------|------------|
| `unread` | — | `true` = hanya belum dibaca |
| `limit` | 20 | max 100 |
| `offset` | 0 | |

**Sukses — 200** — `data`: `NotificationResponse[]`, `meta`: paginasi.

---

## GET `/api/notifications/unread-count`

**Sukses — 200** — `data`:

```json
{ "unread_count": 3 }
```

---

## PATCH `/api/notifications/:id/read`

**Sukses — 200** — `data`: `NotificationResponse` yang sudah `read: true`.

---

## PATCH `/api/notifications/read-all`

**Sukses — 200** — `data`:

```json
{ "updated_count": 5 }
```

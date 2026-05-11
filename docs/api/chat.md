# Modul Chat — `/api/chat/conversations`

CRUD percakapan chat per user. Semua route **wajib** `Authorization: Bearer <token>`.

| Method | Path |
|--------|------|
| POST | `` |
| GET | `` |
| GET | `/:id` |
| PUT | `/:id` |
| DELETE | `/:id` |

## Tipe data: `ChatConversationResponse`

| Field | Tipe |
|-------|------|
| `id` | string (UUID) |
| `title` | string |
| `user_id` | string (UUID) |
| `is_pinned` | boolean |
| `pinned_at` | string \| null |
| `message_count` | number |
| `created_at` | string (ISO) |
| `updated_at` | string (ISO) |

---

## POST `/api/chat/conversations`

**Body**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `title` | Ya | max 255 karakter |

**Sukses — 201** — `data`: `ChatConversationResponse`.

---

## GET `/api/chat/conversations`

**Query:** `limit`, `offset` (default limit 10, max 100).

**Sukses — 200** — `data`: array percakapan, `meta`: paginasi.

---

## GET `/api/chat/conversations/:id`

**Sukses — 200** — `data`: satu `ChatConversationResponse`.

Hanya percakapan milik user dari token.

---

## PUT `/api/chat/conversations/:id`

**Body**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `title` | Tidak | max 255 |

**Sukses — 200** — `data`: percakapan terbaru.

---

## DELETE `/api/chat/conversations/:id`

**Sukses — 200** — `data`: `null`.

**Catatan:** Error domain (tidak ditemukan / bukan pemilik) saat ini dapat mengembalikan **500** di handler.

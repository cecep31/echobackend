# Modul Chat — `/api/chat/conversations` & `/api/chat/messages`

CRUD percakapan chat per user, pesan, dan streaming respons AI. Semua route **wajib** `Authorization: Bearer <token>`.

## Ringkasan route

### Conversations — `/api/chat/conversations`

| Method | Path | Keterangan |
|--------|------|------------|
| POST | `` | Buat percakapan |
| POST | `/stream` | Buat percakapan + pesan pertama (SSE) |
| GET | `` | List percakapan user |
| GET | `/:id` | Detail percakapan |
| PUT | `/:id` | Update judul / pin |
| DELETE | `/:id` | Hapus percakapan |
| POST | `/:conversationId/messages` | Kirim pesan |
| POST | `/:conversationId/messages/stream` | Kirim pesan + stream respons AI (SSE) |
| GET | `/:conversationId/messages` | List pesan dalam percakapan |

### Messages — `/api/chat/messages`

| Method | Path | Keterangan |
|--------|------|------------|
| GET | `/:messageId` | Detail satu pesan |
| DELETE | `/:messageId` | Hapus pesan |

---

## Tipe data

### `ChatConversationResponse`

| Field | Tipe |
|-------|------|
| `id` | string (UUID) |
| `title` | string |
| `user_id` | string (UUID) |
| `is_pinned` | boolean |
| `pinned_at` | string \| null |
| `message_count` | number |
| `chat_messages` | array `ChatMessageResponse` (hanya pada GET `/:id`; urutan kronologis) |
| `created_at` | string (ISO) |
| `updated_at` | string (ISO) |

### `ChatMessageResponse`

| Field | Tipe |
|-------|------|
| `id` | string (UUID) |
| `conversation_id` | string (UUID) |
| `user_id` | string (UUID) |
| `role` | string |
| `content` | string |
| `model` | string \| null |
| `prompt_tokens` | number |
| `completion_tokens` | number |
| `total_tokens` | number |
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

## POST `/api/chat/conversations/stream`

Buat percakapan baru sekaligus kirim pesan pertama; respons AI di-stream via SSE jika AI tersedia.

**Body**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `content` | Ya | 1–10000 karakter |
| `title` | Tidak | 1–255 karakter |
| `model` | Tidak | max 100 karakter |
| `temperature` | Tidak | 0–2 |

**Respons**

- Jika streaming tidak tersedia: **201** — `data`: array berisi `ChatMessageResponse` user.
- Jika streaming aktif: **200** — `Content-Type: text/event-stream`.

Event SSE (urutan):

| `type` | `data` |
|--------|--------|
| `conversation_created` | `{ conversation_id, user_message }` |
| `ai_chunk` | string (potongan teks) |
| `ai_complete` | `ChatMessageResponse` (pesan AI) |
| `error` | pesan error generik |

Stream diakhiri dengan `data: [DONE]`.

---

## GET `/api/chat/conversations`

**Query:** `limit`, `offset` (default limit 10, max 100).

**Sukses — 200** — `data`: array percakapan, `meta`: paginasi.

---

## GET `/api/chat/conversations/:id`

**Sukses — 200** — `data`: satu `ChatConversationResponse` termasuk `chat_messages` (semua pesan, `created_at` naik).

Hanya percakapan milik user dari token.

---

## PUT `/api/chat/conversations/:id`

**Body**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `title` | Tidak | max 255 |
| `is_pinned` | Tidak | boolean |

**Sukses — 200** — `data`: percakapan terbaru.

---

## DELETE `/api/chat/conversations/:id`

**Sukses — 200** — `data`: `null`.

---

## POST `/api/chat/conversations/:conversationId/messages`

**Body**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `content` | Ya | 1–10000 karakter |
| `role` | Tidak | max 20 karakter |
| `model` | Tidak | max 100 karakter |
| `temperature` | Tidak | 0–2 |

**Sukses — 201** — `data`: array pesan (user + respons AI jika ada).

---

## POST `/api/chat/conversations/:conversationId/messages/stream`

Sama seperti endpoint non-stream, tetapi respons AI di-stream via SSE.

Event SSE:

| `type` | `data` |
|--------|--------|
| `user_message` | `ChatMessageResponse` |
| `ai_chunk` | string |
| `ai_complete` | `ChatMessageResponse` |
| `error` | pesan error generik |

---

## GET `/api/chat/conversations/:conversationId/messages`

**Sukses — 200** — `data`: array `ChatMessageResponse`.

---

## GET `/api/chat/messages/:messageId`

**Sukses — 200** — `data`: satu `ChatMessageResponse`.

---

## DELETE `/api/chat/messages/:messageId`

**Sukses — 200** — `data`: pesan yang dihapus.

---

## Error umum

| HTTP | Kondisi |
|------|---------|
| 404 | Percakapan / pesan tidak ditemukan |
| 403 | Bukan pemilik percakapan |
| 422 | Validasi body gagal |
| 500 | Error server lainnya |

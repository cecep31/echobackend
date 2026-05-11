# Dokumentasi API — echobackend

Referensi HTTP API untuk integrasi frontend. Semua route bisnis berada di bawah prefix `/api`, kecuali health check di root.

## Base URL

| Lingkungan | URL |
|------------|-----|
| Lokal | `http://localhost:<PORT>` (lihat `PORT` di `.env`) |
| Produksi | URL deploy Fly.io / reverse proxy Anda |

Contoh penuh: `GET /api/posts`, `POST /api/auth/login`.

## Autentikasi

Route yang membutuhkan login mengirim header:

```http
Authorization: Bearer <access_token>
```

Token didapat dari `POST /api/auth/login` atau `POST /api/auth/refresh`. Claim JWT memuat `user_id` (UUID).

**Middleware auth** yang gagal (token hilang, tidak valid, atau user bukan super admin pada route admin) mengembalikan JSON Echo `{"message":"..."}`, bukan envelope `success` di bawah.

| Situasi | HTTP |
|---------|------|
| Token tidak ada / tidak valid | 401 |
| Bukan super admin pada route admin | 403 |

## Format respons standar

Mayoritas handler memakai envelope dari `pkg/response`:

```json
{
  "success": true,
  "message": "Pesan human-readable",
  "data": {},
  "meta": {},
  "error": "",
  "errors": []
}
```

| Helper | HTTP | Catatan |
|--------|------|---------|
| Sukses | 200 | `success: true`, `data` opsional |
| Dibuat | 201 | Sama seperti sukses |
| Bad request | 400 | `success: false`, `error` berisi detail |
| Unauthorized | 401 | `error`: `"Unauthorized access"` |
| Forbidden | 403 | `error`: `"Access forbidden"` |
| Not found | 404 | |
| Conflict | 409 | Duplikat resource |
| Validasi | 422 | `errors` berisi array field (`field`, `message`, `value`, `tag`) |
| Server error | 500 | Pesan generik; detail hanya di log server |

### Paginasi (`meta`)

List yang dipaginasi memakai `SuccessWithMeta`:

```json
{
  "meta": {
    "total_items": 100,
    "offset": 0,
    "limit": 10,
    "total_pages": 10
  }
}
```

Query: `limit` (default bervariasi per endpoint, **maksimum 100**), `offset` (default `0`).

## Batasan global

- Ukuran body request: **10 MB** (lebih besar → **413**).
- CORS: `HTTP_ALLOW_ORIGINS` (default `*`).
- Rate limit global: aktif jika `HTTP_RATE_LIMIT_RPS` > 0.
- Rate limit khusus: `POST /api/auth/login` dan `POST /api/auth/forgot-password` — **5 request / 5 menit** per IP (burst 5).

## Health & root

| Method | Path | Auth | Respons |
|--------|------|------|---------|
| GET | `/` | Tidak | Envelope sukses (pesan welcome) |
| GET | `/health` | Tidak | `200` `{"status":"ok"}` atau `503` `{"status":"unhealthy","reason":"database unreachable"}` |

## Modul

| Modul | Base path | Dokumen |
|-------|-----------|---------|
| Auth | `/api/auth` | [auth.md](./auth.md) |
| Users & follow | `/api/users` | [users.md](./users.md) |
| Posts (komentar, view, like) | `/api/posts` | [posts.md](./posts.md) |
| Tags | `/api/tags` | [tags.md](./tags.md) |
| Chat | `/api/chat/conversations` | [chat.md](./chat.md) |
| Holdings | `/api/holdings`, `/api/holding-types` | [holdings.md](./holdings.md) |

Debug (`/api/debug/pprof/*`) hanya saat `APP_DEBUG=true` — tidak untuk frontend.

## Konvensi tipe

- **UUID**: string, primary key user/post/komentar/konversasi.
- **Waktu**: ISO 8601 / RFC3339 (`2026-05-12T08:00:00Z`).
- **Nullable**: field pointer di Go → `null` atau dihilangkan (`omitempty`).
- **Angka finansial (holdings)**: string desimal di JSON (mis. `"1500000.00"`), bukan number.

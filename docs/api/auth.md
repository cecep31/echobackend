# Modul Auth — `/api/auth`

Registrasi, login, refresh token, reset password, dan ubah password.

## Ringkasan endpoint

| Method | Path | Auth | Rate limit |
|--------|------|------|------------|
| POST | `/register` | Tidak | Global |
| POST | `/login` | Tidak | 5 / 5 menit |
| POST | `/check-username` | Tidak | Global |
| POST | `/forgot-password` | Tidak | 5 / 5 menit |
| POST | `/reset-password` | Tidak | Global |
| POST | `/refresh` | Tidak | Global |
| PATCH | `/password` | Bearer | Global |

---

## POST `/api/auth/register`

Membuat akun baru.

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `email` | string | Ya | format email |
| `username` | string | Ya | 3–30 karakter |
| `password` | string | Ya | min 6 karakter |

**Sukses — 201**

```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe"
  }
}
```

**Error**

| HTTP | Kondisi |
|------|---------|
| 400 | Body tidak valid |
| 422 | Validasi gagal |
| 409 | Email atau username sudah dipakai |
| 500 | Error server |

---

## POST `/api/auth/login`

Login dengan email **atau** username di field `identifier`.

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `identifier` | string | Ya | — |
| `password` | string | Ya | min 6 |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ...",
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "johndoe"
    }
  }
}
```

**Error**

| HTTP | Kondisi |
|------|---------|
| 400 | Body tidak valid |
| 401 | Kredensial salah |
| 429 | Rate limit (jika middleware memblokir) |
| 500 | Error server |

Simpan `access_token` untuk header `Authorization`; `refresh_token` untuk endpoint refresh.

---

## POST `/api/auth/check-username`

Cek ketersediaan username sebelum registrasi.

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `username` | string | Ya | 3–30 karakter |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Username availability checked",
  "data": {
    "username": "johndoe",
    "available": true
  }
}
```

---

## POST `/api/auth/forgot-password`

Meminta reset password. Respons **sama** untuk email terdaftar maupun tidak (anti-enumeration).

**Body**

| Field | Tipe | Wajib |
|-------|------|-------|
| `email` | string | Ya (email) |

**Sukses — 200**

```json
{
  "success": true,
  "message": "If the email exists, a password reset link has been sent",
  "data": null
}
```

---

## POST `/api/auth/reset-password`

Set password baru dengan token dari email reset.

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `token` | string | Ya | — |
| `password` | string | Ya | min 6 |

**Sukses — 200** — `data`: `null`, message sukses reset.

**Error**

| HTTP | Kondisi |
|------|---------|
| 400 | Token tidak valid / kedaluwarsa |

---

## POST `/api/auth/refresh`

Perpanjang sesi dengan refresh token.

**Body**

| Field | Tipe | Wajib |
|-------|------|-------|
| `refresh_token` | string | Ya |

**Sukses — 200** — Bentuk `data` sama seperti login (`access_token`, `refresh_token`, `user`).

**Error**

| HTTP | Kondisi |
|------|---------|
| 401 | Refresh token tidak valid / kedaluwarsa |

---

## PATCH `/api/auth/password`

Ubah password user yang sedang login.

**Header:** `Authorization: Bearer <access_token>`

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `current_password` | string | Ya | min 6 |
| `new_password` | string | Ya | min 6 |

**Sukses — 200** — `data`: `null`.

**Error**

| HTTP | Kondisi |
|------|---------|
| 401 | Tidak login atau password lama salah |

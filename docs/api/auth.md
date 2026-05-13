# Modul Auth — `/api/auth`

Registrasi, login, OAuth, refresh token, reset password, ubah password, logout, profil, dan log aktivitas.

## Ringkasan endpoint

| Method | Path | Auth | Rate limit |
|--------|------|------|------------|
| POST | `/register` | Tidak | Global |
| POST | `/login` | Tidak | 5 / 5 menit |
| POST | `/check-username` | Tidak | Global |
| GET | `/email/:email` | Tidak | Global |
| POST | `/forgot-password` | Tidak | 5 / 5 menit |
| POST | `/reset-password` | Tidak | Global |
| POST | `/refresh` | Tidak | Global |
| POST | `/logout` | Bearer | Global |
| GET | `/profile` | Bearer | Global |
| PATCH | `/password` | Bearer | Global |
| GET | `/activity-logs` | Bearer | Global |
| GET | `/activity-logs/recent` | Bearer | Global |
| GET | `/activity-logs/failed-logins` | Bearer | Global |
| GET | `/oauth/github` | Tidak | Global |
| GET | `/oauth/github/callback` | Tidak | Global |

---

## POST `/api/auth/register`

Membuat akun baru.

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `email` | string | Ya | format email |
| `username` | string | Ya | 3–30 karakter |
| `password` | string | Ya | min 8 karakter |

> **Catatan:** Disarankan min 8 karakter, huruf besar, huruf kecil, angka, dan karakter spesial.

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
    "refresh_token": "pl_...",
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
| 429 | Rate limit |
| 500 | Error server |

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

## GET `/api/auth/email/:email`

Cek ketersediaan email sebelum registrasi.

**Parameter Path**

| Param | Tipe | Wajib |
|-------|------|-------|
| `email` | string | Ya |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Email availability checked",
  "data": {
    "email": "user@example.com",
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

Set password baru dengan token reset.

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `token` | string | Ya | — |
| `password` | string | Ya | min 8 |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Password reset successful",
  "data": null
}
```

**Error**

| HTTP | Kondisi |
|------|---------|
| 400 | Token tidak valid / kedaluwarsa / sudah dipakai |

---

## POST `/api/auth/refresh`

Perpanjang sesi dengan refresh token. Menghasilkan access token dan refresh token baru (rotasi).

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

## POST `/api/auth/logout`

Logout user dengan menghapus session refresh token.

**Header:** `Authorization: Bearer <access_token>`

**Body**

| Field | Tipe | Wajib |
|-------|------|-------|
| `refresh_token` | string | Ya |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Logout successful",
  "data": null
}
```

---

## GET `/api/auth/profile`

Mendapatkan profil user yang sedang login.

**Header:** `Authorization: Bearer <access_token>`

**Sukses — 200**

```json
{
  "success": true,
  "message": "Profile retrieved successfully",
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "johndoe",
    "first_name": "John",
    "last_name": "Doe",
    "image": "https://...",
    "is_super_admin": false,
    "followers_count": 0,
    "following_count": 0
  }
}
```

---

## PATCH `/api/auth/password`

Ubah password user yang sedang login.

**Header:** `Authorization: Bearer <access_token>`

**Body**

| Field | Tipe | Wajib | Validasi |
|-------|------|-------|----------|
| `current_password` | string | Ya | min 8 |
| `new_password` | string | Ya | min 8 |

**Sukses — 200** — `data`: `null`.

**Error**

| HTTP | Kondisi |
|------|---------|
| 401 | Tidak login atau password lama salah |

---

## GET `/api/auth/activity-logs`

Daftar log aktivitas auth user yang sedang login (paginasi).

**Header:** `Authorization: Bearer <access_token>`

**Query Parameters**

| Param | Tipe | Default | Keterangan |
|-------|------|---------|------------|
| `limit` | int | 20 | Maks 100 |
| `offset` | int | 0 | — |
| `activity_type` | string | — | Filter: `login`, `login_failed`, `logout`, `register`, `password_change`, `password_reset_request`, `password_reset`, `token_refresh`, `oauth_login`, `oauth_login_failed` |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Activity logs retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "user_id": "uuid",
      "activity_type": "login",
      "ip_address": "127.0.0.1",
      "user_agent": "Mozilla/5.0...",
      "status": "success",
      "error_message": null,
      "metadata": null,
      "created_at": "2026-05-13T07:00:00Z"
    }
  ],
  "meta": {
    "total_items": 50,
    "offset": 0,
    "limit": 20,
    "total_pages": 3
  }
}
```

---

## GET `/api/auth/activity-logs/recent`

Log aktivitas terbaru user yang sedang login.

**Header:** `Authorization: Bearer <access_token>`

**Query Parameters**

| Param | Tipe | Default | Keterangan |
|-------|------|---------|------------|
| `limit` | int | 10 | Maks 50 |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Recent activity retrieved successfully",
  "data": [...]
}
```

---

## GET `/api/auth/activity-logs/failed-logins`

Daftar login gagal (semua user, untuk admin monitoring).

**Header:** `Authorization: Bearer <access_token>`

**Query Parameters**

| Param | Tipe | Default | Keterangan |
|-------|------|---------|------------|
| `limit` | int | 20 | Maks 100 |
| `offset` | int | 0 | — |
| `since_hours` | int | 24 | Jam terakhir |

**Sukses — 200**

```json
{
  "success": true,
  "message": "Failed logins retrieved successfully",
  "data": [...],
  "meta": { "total_items": 5, "offset": 0, "limit": 20, "total_pages": 1 }
}
```

---

## GET `/api/auth/oauth/github`

Redirect ke halaman otorisasi GitHub. Mengarahkan browser ke:

```
https://github.com/login/oauth/authorize?client_id=...&redirect_uri=...&scope=user:email
```

**Sukses — 307** Redirect ke GitHub.

---

## GET `/api/auth/oauth/github/callback`

Callback dari GitHub setelah user otorisasi. Menukar `code` dengan access token GitHub, mengambil profil GitHub, lalu membuat/mencari user lokal.

Jika email GitHub tidak tersedia, endpoint akan meminta email dari `https://api.github.com/user/emails`.

**Query Parameters**

| Param | Tipe | Keterangan |
|-------|------|------------|
| `code` | string | Kode otorisasi dari GitHub |

**Alur sukses — 307** Redirect ke:

```
{FRONTEND_URL}?access_token=eyJ...&refresh_token=pl_...
```

**Alur gagal — 307** Redirect ke:

```
{FRONTEND_URL}?error=<error_type>
```

| Error type | Kondisi |
|------------|---------|
| `missing_code` | Parameter `code` kosong |
| `github_token_failed` | Gagal menukar code dengan token GitHub |
| `github_user_failed` | Gagal mengambil profil GitHub |
| `oauth_login_failed` | Gagal membuat/login user |

---

## Log Aktivitas Auth

Semua operasi auth dicatat ke `auth_activity_logs` dengan tipe berikut:

| Tipe | Keterangan |
|------|------------|
| `login` | Login berhasil |
| `login_failed` | Login gagal |
| `logout` | Logout |
| `register` | Registrasi berhasil |
| `password_change` | Perubahan password |
| `password_reset_request` | Permintaan reset password |
| `password_reset` | Reset password berhasil |
| `token_refresh` | Refresh token |
| `oauth_login` | Login via OAuth (GitHub) |
| `oauth_login_failed` | Login OAuth gagal |

Setiap log menyimpan: `user_id`, `activity_type`, `ip_address`, `user_agent`, `status` (success/failure/pending), `error_message`, `metadata` (JSON), `created_at`.
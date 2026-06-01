# Modul Exchange Rates — `/api/exchange-rates`

Nilai tukar mata uang dari Yahoo Finance. **Semua route membutuhkan Bearer token.**

Hasil disimpan di Valkey/Redis selama 15 menit jika `VALKEY_URL` aktif. Jika cache tidak aktif, endpoint tetap berjalan dan selalu fetch ke provider.

## GET `/api/exchange-rates`

Ambil nilai tukar dari satu mata uang ke mata uang lain.

**Query**

| Param | Wajib | Keterangan |
|-------|-------|------------|
| `from` | Ya | Kode mata uang 3 huruf, mis. `USD` |
| `to` | Ya | Kode mata uang 3 huruf, mis. `IDR` |

Kode mata uang dinormalisasi ke uppercase. Input selain 3 huruf alfabet mengembalikan `400`.

**Contoh**

```http
GET /api/exchange-rates?from=USD&to=IDR
Authorization: Bearer <access_token>
```

**Sukses — 200**

```json
{
  "success": true,
  "message": "Exchange rate fetched successfully",
  "data": {
    "from": "USD",
    "to": "IDR",
    "symbol": "USDIDR=X",
    "rate": 16250,
    "source": "Yahoo Finance",
    "cached": false,
    "fetchedAt": "2026-06-01T10:00:00Z"
  }
}
```

| Field | Tipe | Keterangan |
|-------|------|------------|
| `from` | string | Mata uang asal |
| `to` | string | Mata uang tujuan |
| `symbol` | string | Symbol Yahoo Finance yang dipakai |
| `rate` | number | Nilai 1 `from` dalam `to` |
| `source` | string | Provider data |
| `cached` | boolean | `true` jika response berasal dari cache |
| `fetchedAt` | string | Waktu fetch/cache record dibuat, RFC3339 UTC |

**Catatan provider**

Service mencoba direct pair lebih dulu, misalnya `USDIDR=X`. Jika tidak ada, service mencoba inverse pair, misalnya `IDRUSD=X`, lalu menghitung `1 / rate`.

Untuk `from` dan `to` yang sama, response langsung `rate = 1`.

**Error**

| HTTP | Kondisi |
|------|---------|
| 400 | `from` atau `to` bukan kode mata uang 3 huruf |
| 500 | Provider gagal atau pair tidak ditemukan |

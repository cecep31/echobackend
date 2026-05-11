# Modul Holdings — `/api/holdings` & `/api/holding-types`

Portofolio investasi per user per bulan/tahun. **Semua route holdings dan holding-types membutuhkan Bearer token.**

## Tipe data: `Holding`

| Field | Tipe | Keterangan |
|-------|------|------------|
| `id` | number | int64 |
| `user_id` | string (UUID) | |
| `name` | string | |
| `symbol` | string \| null | |
| `platform` | string | |
| `holding_type_id` | number | |
| `holding_type` | object \| null | `HoldingType` |
| `currency` | string | 3 huruf, mis. `IDR` |
| `invested_amount` | string | Desimal |
| `current_value` | string | Desimal |
| `gain_amount` | string | Read-only (DB) |
| `gain_percent` | string | Read-only |
| `units` | string \| null | |
| `avg_buy_price` | string \| null | |
| `current_price` | string \| null | |
| `last_updated` | string \| null | |
| `notes` | string \| null | |
| `month` | number | 1–12 |
| `year` | number | |
| `created_at` | string | |
| `updated_at` | string | |

### `HoldingType`

```json
{ "id": 1, "code": "STOCK", "name": "Saham", "notes": null }
```

---

## Holdings — ringkasan route

| Method | Path | Keterangan |
|--------|------|------------|
| GET | `` | List + filter bulan/tahun |
| GET | `/summary` | Agregat portofolio |
| GET | `/trends` | Tren multi-tahun |
| GET | `/compare` | Banding dua bulan |
| GET | `/monthly` | Seri bulanan |
| POST | `` | Buat holding |
| POST | `/duplicate` | Salin bulan ke bulan |
| POST | `/sync` | Sinkron harga |
| GET | `/:id` | Detail by id |
| PUT | `/:id` | Update |
| DELETE | `/:id` | Hapus |

## GET `/api/holdings`

**Query**

| Param | Keterangan |
|-------|------------|
| `month` | 1–12 |
| `year` | Tahun |
| `sortBy` | Field sort (lihat service) |
| `order` | `asc` / `desc` |

**Sukses — 200** — `data`: `Holding[]`.

---

## GET `/api/holdings/summary`

**Query:** `month`, `year` (opsional).

**Sukses — 200** — `data` (`HoldingSummaryResponse`):

| Field | Tipe |
|-------|------|
| `totalInvested` | string |
| `totalCurrentValue` | string |
| `totalProfitLoss` | string |
| `totalProfitLossPercentage` | string |
| `holdingsCount` | number |
| `typeBreakdown` | array breakdown per tipe |
| `platformBreakdown` | array breakdown per platform |

---

## GET `/api/holdings/trends`

**Query:** `years` — tahun dipisah koma, mis. `2024,2025`.

**Sukses — 200** — `data`: `HoldingTrendResponse[]` (`date`, `invested`, `current`, `profitLoss`, `profitLossPercentage`).

---

## GET `/api/holdings/compare`

**Query:** `fromMonth`, `fromYear`, `toMonth`, `toYear` (default periode: bulan lalu → sekarang).

**Sukses — 200** — perbandingan `fromMonth` / `toMonth`, `summary`, `typeComparison`, `platformComparison`.

---

## GET `/api/holdings/monthly`

**Query:** `startMonth`, `startYear`, `endMonth`, `endYear`.

**Sukses — 200** — `data`: `HoldingMonthlyDataResponse[]`.

---

## POST `/api/holdings`

**Body (`CreateHoldingRequest`)**

| Field | Wajib | Validasi |
|-------|-------|----------|
| `name` | Ya | |
| `platform` | Ya | |
| `holding_type_id` | Ya | |
| `currency` | Ya | panjang 3 |
| `invested_amount` | Ya | string desimal |
| `current_value` | Ya | string desimal |
| `month` | Ya | 1–12 |
| `year` | Ya | min 2000 |
| `symbol`, `units`, `avg_buy_price`, `current_price`, `last_updated`, `notes` | Tidak | |

**Sukses — 201** — `data`: array `Holding` (satu elemen).

---

## POST `/api/holdings/duplicate`

**Body**

| Field | Wajib |
|-------|-------|
| `fromMonth`, `fromYear`, `toMonth`, `toYear` | Ya (1–12 / 1900–2100) |
| `overwrite` | boolean |

**Sukses — 201** — `data`: `DuplicateResultItem[]` (`id`, `name`, `month`, `year`).

**Error 400** jika bulan sumber = tujuan.

---

## POST `/api/holdings/sync`

Sinkron harga untuk periode aktif user.

**Sukses — 200** — `data`: `{ "syncedCount", "month", "year" }`.

---

## GET `/api/holdings/:id`

**Path:** `id` numerik.

| HTTP | Kondisi |
|------|---------|
| 404 | Tidak ada / bukan milik user |
| 403 | Bukan pemilik |

---

## PUT `/api/holdings/:id`

**Body:** `UpdateHoldingRequest` — field opsional; handler tidak selalu memanggil validator global pada body.

**Sukses — 200** — `data`: `[Holding]`.

---

## DELETE `/api/holdings/:id`

**Sukses — 200** — `data`: `null`.

---

## GET `/api/holding-types`

**Sukses — 200** — `data`: `HoldingType[]`.

Daftar tipe aset (saham, reksadana, dll.) untuk dropdown form.

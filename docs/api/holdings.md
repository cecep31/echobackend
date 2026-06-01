# Holdings Module - `/api/holdings` & `/api/holding-types`

Investment portfolio records per user per month/year. **All holdings and holding-types routes require a Bearer token.**

## Data Type: `Holding`

| Field | Type | Description |
|-------|------|-------------|
| `id` | number | int64 |
| `user_id` | string (UUID) | |
| `name` | string | |
| `symbol` | string \| null | |
| `platform` | string | |
| `holding_type_id` | number | |
| `holding_type` | object \| null | `HoldingType` |
| `currency` | string | 3 letters, for example `IDR` |
| `invested_amount` | string | Decimal |
| `current_value` | string | Decimal |
| `gain_amount` | string | Read-only (DB) |
| `gain_percent` | string | Read-only |
| `units` | string \| null | |
| `avg_buy_price` | string \| null | |
| `current_price` | string \| null | |
| `last_updated` | string \| null | |
| `notes` | string \| null | |
| `month` | number | 1-12 |
| `year` | number | |
| `created_at` | string | |
| `updated_at` | string | |

### `HoldingType`

```json
{ "id": 1, "code": "STOCK", "name": "Saham", "notes": null }
```

---

## Holdings - Route Summary

| Method | Path | Description |
|--------|------|-------------|
| GET | `` | List + month/year filters |
| GET | `/summary` | Portfolio aggregate |
| GET | `/trends` | Multi-year trends |
| GET | `/compare` | Compare two months |
| GET | `/monthly` | Monthly series |
| POST | `` | Create holding |
| POST | `/duplicate` | Copy one month to another |
| POST | `/sync` | Sync prices |
| GET | `/:id` | Detail by ID |
| PUT | `/:id` | Update |
| DELETE | `/:id` | Delete |

## GET `/api/holdings`

**Query**

| Param | Default | Description |
|-------|---------|-------------|
| `month` | Current month | 1-12 |
| `year` | Current year | |
| `sortBy` | `created_at` | `created_at`, `updated_at`, `name`, `platform`, `invested_amount`, `current_value`, `holding_type` |
| `order` | `desc` | `asc` / `desc` |

**Success - 200** - `data`: `Holding[]`.

---

## GET `/api/holdings/summary`

**Query:** `month`, `year` (optional).

**Success - 200** - `data` (`HoldingSummaryResponse`):

| Field | Type |
|-------|------|
| `totalInvested` | string |
| `totalCurrentValue` | string |
| `totalProfitLoss` | string |
| `totalProfitLossPercentage` | string |
| `holdingsCount` | number |
| `typeBreakdown` | breakdown array by type |
| `platformBreakdown` | breakdown array by platform |

---

## GET `/api/holdings/trends`

**Query:** `years` - comma-separated years, for example `2024,2025`.

**Success - 200** - `data`: `HoldingTrendResponse[]` (`date`, `invested`, `current`, `profitLoss`, `profitLossPercentage`).

---

## GET `/api/holdings/compare`

**Query:** `fromMonth`, `fromYear`, `toMonth`, `toYear` (default period: previous month -> current month).

**Success - 200** - comparison of `fromMonth` / `toMonth`, `summary`, `typeComparison`, `platformComparison`.

---

## GET `/api/holdings/monthly`

**Query**

| Param | Default | Description |
|-------|---------|-------------|
| `startMonth`, `startYear` | Current month/year | Range start |
| `endMonth`, `endYear` | 11 months before start | Range end when omitted |

**Success - 200** - `data`: `HoldingMonthlyDataResponse[]`.

---

## POST `/api/holdings`

**Body (`CreateHoldingRequest`)**

| Field | Required | Validation |
|-------|----------|------------|
| `name` | Yes | |
| `platform` | Yes | |
| `holding_type_id` | Yes | |
| `currency` | Yes | length 3 |
| `invested_amount` | Yes | decimal string |
| `current_value` | Yes | decimal string |
| `month` | Yes | 1-12 |
| `year` | Yes | min 2000 |
| `symbol`, `units`, `avg_buy_price`, `current_price`, `last_updated`, `notes` | No | |

**Success - 201** - `data`: `Holding[]` (one element).

---

## POST `/api/holdings/duplicate`

**Body**

| Field | Required |
|-------|----------|
| `fromMonth`, `fromYear`, `toMonth`, `toYear` | Yes (1-12 / 1900-2100) |
| `overwrite` | boolean |

**Success - 201** - `data`: `DuplicateResultItem[]` (`id`, `name`, `month`, `year`).

**Error 400** when source month = destination month.

---

## POST `/api/holdings/sync`

Sync prices for the user's active period.

**Success - 200** - `data`: `{ "syncedCount", "month", "year" }`.

---

## GET `/api/holdings/:id`

**Path:** numeric `id`.

| HTTP | Condition |
|------|-----------|
| 404 | Holding not found |
| 403 | Not the owner |

---

## PUT `/api/holdings/:id`

**Body:** `UpdateHoldingRequest` - optional fields; the handler does not always call global body validation.

**Success - 200** - `data`: `[Holding]`.

---

## DELETE `/api/holdings/:id`

**Success - 200** - `data`: `null`.

---

## GET `/api/holding-types`

**Success - 200** - `data`: `HoldingType[]`.

Asset type list (stocks, mutual funds, etc.) for form dropdowns.

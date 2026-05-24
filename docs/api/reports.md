# Modul Reports — `/api/reports`

Dashboard statistik admin. **Semua route membutuhkan Bearer token + super admin.**

## Ringkasan route

| Method | Path | Keterangan |
|--------|------|------------|
| GET | `/overview` | Ringkasan + engagement |
| GET | `/users` | Laporan user |
| GET | `/posts` | Laporan post |
| GET | `/engagement` | Metrik engagement |

---

## Query umum (opsional)

| Param | Keterangan |
|-------|------------|
| `startDate` | Filter awal periode (format string, lihat service) |
| `endDate` | Filter akhir periode |

Endpoint `/users` dan `/posts` juga menerima:

| Param | Default | Keterangan |
|-------|---------|------------|
| `limit` | 10 | max 100 |

Endpoint `/posts` tambahan:

| Param | Keterangan |
|-------|------------|
| `tagId` | Filter numerik by tag |

---

## GET `/api/reports/overview`

**Sukses — 200** — `data`:

```json
{
  "overview": { /* OverviewStatsResponse */ },
  "engagement": { /* EngagementMetricsResponse */ }
}
```

### `OverviewStatsResponse`

| Field | Tipe |
|-------|------|
| `totalUsers` | number |
| `totalPosts` | number |
| `totalViews` | number |
| `totalLikes` | number |
| `totalComments` | number |
| `newUsersToday` | number |
| `newPostsToday` | number |
| `activeUsersThisWeek` | number |

---

## GET `/api/reports/users`

**Sukses — 200** — `data`: `UserReportResponse`

| Field | Tipe |
|-------|------|
| `totalUsers` | number |
| `newUsersThisPeriod` | number |
| `activeUsers` | number |
| `topContributors` | array (`id`, `username`, `firstName`, `lastName`, `postCount`, `totalViews`, `totalLikes`) |
| `growthTrend` | array (`date`, `newUsers`, `cumulativeUsers`) |

---

## GET `/api/reports/posts`

**Sukses — 200** — `data`: `PostReportResponse`

| Field | Tipe |
|-------|------|
| `totalPosts` | number |
| `newPostsThisPeriod` | number |
| `totalViews` | number |
| `totalLikes` | number |
| `totalComments` | number |
| `avgEngagementRate` | number |
| `topPosts` | array performa post |
| `tagPerformance` | array performa tag |

---

## GET `/api/reports/engagement`

**Sukses — 200** — `data`: `EngagementMetricsResponse`

| Field | Tipe |
|-------|------|
| `totalEngagements` | number |
| `avgLikesPerPost` | number |
| `avgCommentsPerPost` | number |
| `avgViewsPerPost` | number |
| `periodComparison` | `{ current, previous, changePercent }` |

---

## Error umum

| HTTP | Kondisi |
|------|---------|
| 401 | Tidak login |
| 403 | Bukan super admin |
| 500 | Error server |

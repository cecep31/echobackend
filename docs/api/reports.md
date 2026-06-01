# Reports Module - `/api/reports`

Admin statistics dashboard. **All routes require a Bearer token and super admin access.**

## Route Summary

| Method | Path | Description |
|--------|------|-------------|
| GET | `/overview` | Overview + engagement |
| GET | `/users` | User report |
| GET | `/posts` | Post report |
| GET | `/engagement` | Engagement metrics |

---

## Common Query Parameters (Optional)

| Param | Description |
|-------|-------------|
| `startDate` | Start period filter (string format; see service behavior) |
| `endDate` | End period filter |

The `/users` and `/posts` endpoints also accept:

| Param | Default | Description |
|-------|---------|-------------|
| `limit` | 10 | max 100 |

The `/posts` endpoint also accepts:

| Param | Description |
|-------|-------------|
| `tagId` | Numeric tag filter |

---

## GET `/api/reports/overview`

**Success - 200** - `data`:

```json
{
  "overview": { /* OverviewStatsResponse */ },
  "engagement": { /* EngagementMetricsResponse */ }
}
```

### `OverviewStatsResponse`

| Field | Type |
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

**Success - 200** - `data`: `UserReportResponse`

| Field | Type |
|-------|------|
| `totalUsers` | number |
| `newUsersThisPeriod` | number |
| `activeUsers` | number |
| `topContributors` | array (`id`, `username`, `firstName`, `lastName`, `postCount`, `totalViews`, `totalLikes`) |
| `growthTrend` | array (`date`, `newUsers`, `cumulativeUsers`) |

---

## GET `/api/reports/posts`

**Success - 200** - `data`: `PostReportResponse`

| Field | Type |
|-------|------|
| `totalPosts` | number |
| `newPostsThisPeriod` | number |
| `totalViews` | number |
| `totalLikes` | number |
| `totalComments` | number |
| `avgEngagementRate` | number |
| `topPosts` | post performance array |
| `tagPerformance` | tag performance array |

---

## GET `/api/reports/engagement`

**Success - 200** - `data`: `EngagementMetricsResponse`

| Field | Type |
|-------|------|
| `totalEngagements` | number |
| `avgLikesPerPost` | number |
| `avgCommentsPerPost` | number |
| `avgViewsPerPost` | number |
| `periodComparison` | `{ current, previous, changePercent }` |

---

## Common Errors

| HTTP | Condition |
|------|-----------|
| 401 | Not authenticated |
| 403 | Not a super admin |
| 500 | Server error |

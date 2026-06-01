# Notifications Module - `/api/notifications`

In-app notifications for the logged-in user. **All routes require a Bearer token.**

Notifications are created automatically by the system, for example when a new comment or follow event occurs. There is no public HTTP endpoint for creating notifications.

## Route Summary

| Method | Path | Description |
|--------|------|-------------|
| GET | `` | List notifications |
| GET | `/unread-count` | Unread count |
| PATCH | `/:id/read` | Mark one notification as read |
| PATCH | `/read-all` | Mark all notifications as read |

---

## Data Type: `NotificationResponse`

| Field | Type |
|-------|------|
| `id` | string (UUID) |
| `user_id` | string (UUID) |
| `type` | string |
| `title` | string |
| `message` | string \| null |
| `read` | boolean |
| `data` | object \| null |
| `created_at` | string \| null |
| `updated_at` | string \| null |

### Notification Types Currently Used

| `type` | Trigger |
|--------|---------|
| `comment` | New comment on the user's post |
| `follow` | Another user follows the account |

---

## GET `/api/notifications`

**Query**

| Param | Default | Description |
|-------|---------|-------------|
| `unread` | - | `true` = unread only |
| `limit` | 20 | max 100 |
| `offset` | 0 | |

**Success - 200** - `data`: `NotificationResponse[]`, `meta`: pagination.

---

## GET `/api/notifications/unread-count`

**Success - 200** - `data`:

```json
{ "unread_count": 3 }
```

---

## PATCH `/api/notifications/:id/read`

**Success - 200** - `data`: `NotificationResponse` with `read: true`.

---

## PATCH `/api/notifications/read-all`

**Success - 200** - `data`:

```json
{ "updated_count": 5 }
```

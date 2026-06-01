# Chat Module - `/api/chat/conversations` & `/api/chat/messages`

Per-user chat conversation CRUD, messages, and streaming AI responses. All routes **require** `Authorization: Bearer <token>`.

## Route Summary

### Conversations - `/api/chat/conversations`

| Method | Path | Description |
|--------|------|-------------|
| POST | `` | Create conversation |
| POST | `/stream` | Create conversation + first message (SSE) |
| GET | `` | List user conversations |
| GET | `/:id` | Conversation detail |
| PUT | `/:id` | Update title / pin |
| DELETE | `/:id` | Delete conversation |
| POST | `/:conversationId/messages` | Send message |
| POST | `/:conversationId/messages/stream` | Send message + stream AI response (SSE) |
| GET | `/:conversationId/messages` | List messages in a conversation |

### Messages - `/api/chat/messages`

| Method | Path | Description |
|--------|------|-------------|
| GET | `/:messageId` | Message detail |
| DELETE | `/:messageId` | Delete message |

---

## Data Types

### `ChatConversationResponse`

| Field | Type |
|-------|------|
| `id` | string (UUID) |
| `title` | string |
| `user_id` | string (UUID) |
| `is_pinned` | boolean |
| `pinned_at` | string \| null |
| `message_count` | number |
| `chat_messages` | `ChatMessageResponse[]` (only on GET `/:id`; chronological order) |
| `created_at` | string (ISO) |
| `updated_at` | string (ISO) |

### `ChatMessageResponse`

| Field | Type |
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

| Field | Required | Validation |
|-------|----------|------------|
| `title` | Yes | max 255 characters |

**Success - 201** - `data`: `ChatConversationResponse`.

---

## POST `/api/chat/conversations/stream`

Create a new conversation and send the first message; the AI response is streamed via SSE when AI is available.

**Body**

| Field | Required | Validation |
|-------|----------|------------|
| `content` | Yes | 1-10000 characters |
| `title` | No | 1-255 characters |
| `model` | No | max 100 characters |
| `temperature` | No | 0-2 |

**Response**

- If streaming is unavailable: **201** - `data`: array containing the user `ChatMessageResponse`.
- If streaming is active: **200** - `Content-Type: text/event-stream`.

SSE event order:

| `type` | `data` |
|--------|--------|
| `conversation_created` | `{ conversation_id, user_message }` |
| `ai_chunk` | string (text chunk) |
| `ai_complete` | `ChatMessageResponse` (AI message) |
| `error` | generic error message |

The stream ends with `data: [DONE]`.

---

## GET `/api/chat/conversations`

**Query:** `limit`, `offset` (default limit 10, max 100).

**Success - 200** - `data`: conversation array, `meta`: pagination.

---

## GET `/api/chat/conversations/:id`

**Success - 200** - `data`: one `ChatConversationResponse` including `chat_messages` (all messages, ascending `created_at`).

Only conversations owned by the token user are returned.

---

## PUT `/api/chat/conversations/:id`

**Body**

| Field | Required | Validation |
|-------|----------|------------|
| `title` | No | max 255 |
| `is_pinned` | No | boolean |

**Success - 200** - `data`: updated conversation.

---

## DELETE `/api/chat/conversations/:id`

**Success - 200** - `data`: `null`.

---

## POST `/api/chat/conversations/:conversationId/messages`

**Body**

| Field | Required | Validation |
|-------|----------|------------|
| `content` | Yes | 1-10000 characters |
| `role` | No | max 20 characters |
| `model` | No | max 100 characters |
| `temperature` | No | 0-2 |

**Success - 201** - `data`: message array (user + AI response when available).

---

## POST `/api/chat/conversations/:conversationId/messages/stream`

Same as the non-stream endpoint, but the AI response is streamed via SSE.

SSE events:

| `type` | `data` |
|--------|--------|
| `user_message` | `ChatMessageResponse` |
| `ai_chunk` | string |
| `ai_complete` | `ChatMessageResponse` |
| `error` | generic error message |

---

## GET `/api/chat/conversations/:conversationId/messages`

**Success - 200** - `data`: `ChatMessageResponse[]`.

---

## GET `/api/chat/messages/:messageId`

**Success - 200** - `data`: one `ChatMessageResponse`.

---

## DELETE `/api/chat/messages/:messageId`

**Success - 200** - `data`: deleted message.

---

## Common Errors

| HTTP | Condition |
|------|-----------|
| 404 | Conversation / message not found |
| 403 | Not the conversation owner |
| 422 | Body validation failed |
| 500 | Other server error |

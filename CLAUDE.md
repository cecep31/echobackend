# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

Note: [AGENTS.md](AGENTS.md) covers the same ground — if you change conventions documented here, update both files.

## Commands

```bash
go run cmd/main.go              # Run server (requires .env with DATABASE_URL + JWT_SECRET)
air                             # Hot reload (reads .env automatically)
go test ./...                   # All tests — no DB required, no external test deps
go test ./internal/service/...  # Service tests only
go test -race ./...             # Race checker
go test -run TestName ./internal/service/...  # Single test
go vet ./...                    # Static analysis
golangci-lint run               # Lint
gosec ./...                     # Security scan

# Migrations (goose, config from .env GOOSE_* vars)
goose up / goose down / goose status
goose create <name> sql         # Always sql, never go
```

First-time migration setup: goose stores version history in the non-default `custom` schema (`GOOSE_TABLE=custom.goose_migrations`). Run once before the first `goose up`:
`psql "$DATABASE_URL" -c 'CREATE SCHEMA IF NOT EXISTS custom;'`

**CI** (`.github/workflows/main.yml`): on PRs/pushes run `go vet`, `go test`, and `golangci-lint`. On push to `main` only (after tests pass): build/push Docker image tags `latest` + `sha-*`.

## Architecture

- **Framework**: Echo **v5** (not v4). The API differs from v4 — handlers receive `*echo.Context` (a pointer).
- **Entry point**: `cmd/main.go` — loads config → builds DI container → registers routes → starts server with graceful shutdown (10 s) and resource cleanup.
- **DI**: Manual wiring in `internal/di/container.go` — no DI framework. Every new handler/service/repository must be instantiated there.
- **Layering**: `internal/handler` → `internal/service` → `internal/repository` (GORM). Models in `internal/model`, request/response structs in `internal/dto`, shared error sentinels in `internal/apperror`.
- **`internal/platform/`**: App-owned infrastructure adapters — `cache` (Valkey), `database`, `email`, `logger`, `queue` (asynq), `storage` (MinIO/S3).
- **`pkg/`**: Reusable helper packages — `applog`, `market`, `response`, `validator`.
- **Routes**: All under `/api/*`, defined per module in `internal/routes/*Routes.go` and wired in `routes.go`. Protected routes use `r.authMiddleware.Auth()`; admin routes chain `r.authMiddleware.AuthAdmin()` after it.
- **Health**: `GET /health` pings the DB (200/503). Used by Docker HEALTHCHECK and load balancers.
- **Pagination**: Use `handler.ParsePaginationParams(c, defaultLimit)` — returns `(limit, offset)`, capped at 100. Build meta with `response.CalculatePaginationMeta(total, offset, limit)` and return via `response.SuccessWithMeta`.

## Response Format

All handlers return JSON through `pkg/response`:

```go
response.Success(c, "message", data)         // 200
response.Created(c, "message", data)         // 201
response.BadRequest(c, "msg", err)           // 400
response.Unauthorized(c, "msg")              // 401
response.Forbidden(c, "msg")                 // 403
response.NotFound(c, "msg", err)             // 404
response.Conflict(c, "msg", "reason")        // 409 — takes a string reason, not an error
response.ValidationError(c, "msg", err)      // 422
response.FromValidateError(c, err)           // 422 with structured field errors
response.TooManyRequests(c, "msg")           // 429
response.InternalServerError(c, "msg", err)  // 500 — err logged server-side only, never sent to client
```

## Config & Env

- `config.Load()` reads `.env` (godotenv) then environment variables. See `.env.example` for optional keys.
- **Required**: `DATABASE_URL`, `JWT_SECRET` — the app panics without them.
- Many keys accept **fallback aliases** (legacy names); first-set key wins. Full list in `config/config.go`.
- Valkey/Redis caching is **optional** — leave `VALKEY_URL` empty to disable (app runs fail-open).

## Testing

- Tests live in `internal/service/`, `internal/handler/`, `internal/middleware/`, `config/`, and `pkg/` — no repository or DB integration tests, so `go test ./...` never needs PostgreSQL.
- Service tests use hand-written mocks in `internal/service/mocks_test.go` — no mockgen or code generation.
- Test files are white-box: `*_test.go` in the same package.

## Migrations

- Goose with **raw SQL** files in `migrations/`, numbered sequentially (`001_init_schema.sql`, …).
- Schema uses PostgreSQL triggers for count fields, UUID primary keys (v7 preferred), `ON DELETE CASCADE`, and soft deletes via `deleted_at`.

## Deployment

- GitHub Actions: push to `main` → test/lint → Docker build → push `cecep31/echobackend:latest` and `:sha-*` to Docker Hub. Deploy by pulling a pinned `sha-*` tag on your host.

## API Docs

Per-module HTTP API reference for frontend integration lives in `docs/api/` (auth, users, posts, tags, chat, holdings, exchange rates, bookmarks, notifications, reports). Update the relevant doc when changing an endpoint.

# AGENTS.md — echobackend

## Commands

```bash
# Local dev services (Postgres 18 + Valkey + MinIO via docker-compose.yml)
docker compose up -d --wait # Start services (or: make up). Creates the `custom` schema automatically.
docker compose down         # Stop (or: make down; make down-clean also wipes data)
make help                   # All shortcuts: up, dev, test, lint, check, migrate-*, ...

go run cmd/main.go          # Run server (requires .env with DATABASE_URL + JWT_SECRET)
air                         # Hot reload (reads .env automatically)
go test ./...               # All tests (service + pkg layers only; no DB integration tests)
go test ./internal/service/...  # Service tests only
go test ./pkg/...           # Reusable package tests only
go test -race ./...         # Race checker
go vet ./...                # Static analysis
go fmt ./...                # Format
golangci-lint run           # Lint
gosec ./...                 # Security scan

# Migrations (requires .env with GOOSE_* vars)
goose up                    # Apply pending
goose down                  # Rollback one
goose status                # Check current
goose create <name> sql     # New migration file

# First-time setup: goose stores version history in the `custom` schema.
# `docker compose up` already creates it via scripts/init-db.sql.
# For an external Postgres, run this once before the first `goose up`:
psql "$DATABASE_URL" -c 'CREATE SCHEMA IF NOT EXISTS custom;'
```

## Architecture

- **Framework**: Echo **v5** (not v4). API differs — handlers receive `*echo.Context` (pointer).
- **Entry point**: `cmd/main.go` — loads config → creates DI container → registers routes → starts server with graceful shutdown.
- **DI**: Manual wiring in `internal/di/container.go`. All handler/service/repo instances created there.
- **Layering**: `handler` → `service` → `repository`. No DI framework.
- **`internal/platform/`**: App-owned infrastructure adapters (`cache`, `database`, `email`, `queue`, `storage`).
- **`pkg/`**: Reusable helper packages (`applog`, `market`, `response`, `validator`).
- **`internal/model/`**: GORM entities (domain models).
- **`internal/repository/`**: Data access layer using GORM.
- **`internal/dto/`**: Request/response structs. `internal/apperror/` for shared app error sentinels.
- **API routes**: All under `/api/*`, defined in `internal/routes/*Routes.go`. Auth-protected routes use `r.authMiddleware.Auth()`.
- **Health**: `GET /health` — pings DB (200/503). Used by Docker HEALTHCHECK and load balancers.
- **Auth gates**: routes apply `r.authMiddleware.Auth()`, and admin routes chain `r.authMiddleware.AuthAdmin()` after it (e.g. `posts.PUT("/:id", h.UpdatePost, r.authMiddleware.Auth(), r.authMiddleware.AuthAdmin())`).
- **Pagination**: Use `handler.ParsePaginationParams(c, defaultLimit)` — returns `(limit, offset)`, max cap 100. Build response meta with `response.CalculatePaginationMeta(total, offset, limit)` and pass via `response.SuccessWithMeta`.

## Config & Env

- Config loaded via `config.Load()` — reads `.env` (godotenv) then environment variables.
- **Required**: `DATABASE_URL`, `JWT_SECRET`. App panics if missing.
- Many keys accept **fallback aliases** (legacy names). First-set key wins. See `config/config.go` for full list.
- `GOOSE_TABLE=custom.goose_migrations` — non-default goose table location; create the `custom` schema (`psql "$DATABASE_URL" -c 'CREATE SCHEMA IF NOT EXISTS custom;'`) once before the first `goose up`.
- Valkey/Redis caching is **optional** — leave `VALKEY_URL` empty to disable (app runs fail-open).

## Testing

- Tests exist mostly in `internal/service/`, `internal/handler/`, `internal/middleware/`, `config/`, and `pkg/`. No repository or DB integration tests.
- **No external test dependencies** — service tests use hand-written mocks (`internal/service/mocks_test.go`). No mockgen or code-gen.
- No testcontainers or integration test harness. Running `go test ./...` does not require PostgreSQL.
- Test file pattern: `*_test.go` in the same package (white-box).

## Response Format

All handlers use `pkg/response` for consistent JSON:

```go
response.Success(c, "message", data)        // 200
response.Created(c, "message", data)         // 201
response.ValidationError(c, "msg", err)       // 422
response.BadRequest(c, "msg", err)            // 400
response.NotFound(c, "msg", err)              // 404
response.Unauthorized(c, "msg")              // 401
response.Forbidden(c, "msg")                 // 403
response.Conflict(c, "msg", conflictErr)      // 409
response.InternalServerError(c, "msg", err)  // 500 — err logged server-side only, never sent to client
response.FromValidateError(c, err)            // 422 with structured field errors
```

Use `response.TooManyRequests(c, "msg")` for 429 (rate limit). `response.Conflict` takes a string reason, not an `error`.

## CI

`.github/workflows/main.yml` runs on PRs and pushes to `main`:

1. **test** — `go vet ./...`, `go test ./...`, `golangci-lint`
2. **docker** (push to `main` only, after test) — build & push `cecep31/echobackend:latest`, `:sha-<12-char>`, and `:sha-<full>`

Still useful locally before pushing: `go test ./...`, `go vet ./...`, `golangci-lint run`.

## Migrations

- Goose with **raw SQL** files in `migrations/`. Numbered `001_init_schema.sql` through `010_drop_uuid_ossp.sql`.
- Uses PostgreSQL features: triggers for count fields, `uuid_generate_v4()` (v7 preferred), `ON DELETE CASCADE`, soft deletes via `deleted_at`.
- **New migrations**: `goose create <name> sql` (always `sql`, never `go`).
- **First-time setup**: The local Postgres from `docker compose up` auto-creates the `custom` schema (`scripts/init-db.sql`). For an external Postgres, run `psql "$DATABASE_URL" -c 'CREATE SCHEMA IF NOT EXISTS custom;'` once before the first `goose up`.

## Deployment

- GitHub Actions on push to `main`: test/lint → Docker build → push to Docker Hub (`cecep31/echobackend:latest` and `:sha-*`).
- Docker image is built with Go 1.26, multi-stage, runs as non-root user. Pull a specific `sha-*` tag for reproducible deploys.

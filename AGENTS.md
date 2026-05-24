# AGENTS.md — echobackend

## Commands

```bash
go run cmd/main.go          # Run server (requires .env with DATABASE_URL + JWT_SECRET)
air                         # Hot reload (reads .env automatically)
go test ./...               # All tests (service + pkg layers only; no DB integration tests)
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
```

## Architecture

- **Framework**: Echo **v5** (not v4). API differs — handlers receive `*echo.Context` (pointer).
- **Entry point**: `cmd/main.go` — loads config → creates DI container → registers routes → starts server with graceful shutdown.
- **DI**: Manual wiring in `internal/di/container.go`. All handler/service/repo instances created there.
- **Layering**: `handler` → `service` → `repository`. No DI framework.
- **`pkg/`**: Infra-agnostic packages (`cache`, `database`, `market`, `response`, `storage`, `validator`).
- **`internal/dto/`**: Request/response structs. `internal/errors/` for shared error types.
- **API routes**: All under `/api/*`, defined in `internal/routes/*Routes.go`. Auth-protected routes use `r.authMiddleware.Auth()`.
- **Health**: `GET /health` — pings DB (200/503). Used by Fly.io and Docker HEALTHCHECK.
- **Debug routes**: `GET /api/debug/pprof/*` — only registered when `APP_DEBUG=true`.

## Config & Env

- Config loaded via `config.Load()` — reads `.env` (godotenv) then environment variables.
- **Required**: `DATABASE_URL`, `JWT_SECRET`. App panics if missing.
- Many keys accept **fallback aliases** (legacy names). First-set key wins. See `config/config.go` for full list.
- `GOOSE_TABLE=custom.goose_migrations` — non-default goose table location; migrations won't work without this.
- Valkey/Redis caching is **optional** — leave `VALKEY_URL` empty to disable (app runs fail-open).

## Testing

- Tests exist in `internal/service/` and `pkg/` only. No handler or repository tests.
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

Pagination via `handler.ParsePaginationParams(c, defaultLimit)` — max cap 100.

## Migrations

- Goose with **raw SQL** files in `migrations/`. Numbered `001_init_schema.sql` through `008_...`.
- Uses PostgreSQL features: triggers for count fields, `uuid_generate_v4()`, `ON DELETE CASCADE`, soft deletes via `deleted_at`.
- **New migrations**: `goose create <name> sql` (always `sql`, never `go`).

## Deployment

- **Fly.io** via GitHub Actions (push to `main` → Docker build → push to Docker Hub → `flyctl deploy`).
- Docker image: `cecep31/echobackend:latest`. Built with Go 1.26, runs as non-root user.
- Primary region: `sin` (Singapore).
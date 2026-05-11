# echobackend — Agent Notes

Go REST API (Echo v5 + GORM + PostgreSQL). Single binary, manual DI, deployed to Fly.io via Docker.

## Build & Run

- **Entry point:** `cmd/main.go`
- **Build:** `go build -o bin/main cmd/main.go`
- **Dev (hot reload):** `air` (config in `.air.toml`; builds to `tmp/main.exe` on Windows)
- **Run:** `go run cmd/main.go`
- **No Makefile.** All build/test/lint commands must be run directly via `go` tooling or external binaries.

## Testing & Quality

- **All tests:** `go test ./...` (no tests currently exist — `**/*_test.go` returns nothing)
- **Race check:** `go test -race ./...`
- **Coverage:** `go test -cover ./...` or `go test -coverprofile=coverage.out ./...`
- **Vet:** `go vet ./...`
- **Format:** `go fmt ./...` (or `gofmt -s -w .`)
- **Lint:** requires `golangci-lint` installed separately (`golangci-lint run`)
- **Security:** requires `gosec` installed separately (`gosec ./...`)

## Architecture

- **Manual DI** wired in `internal/di/container.go`. When adding a service/repo/handler:
  1. Create struct + constructor.
  2. Register in `Container` and `NewContainer`.
  3. Thread it through `routes.NewRoutes` and the route setup methods.
- **Graceful shutdown:** `CleanupManager` (LIFO) in `internal/di/cleanup.go`. Register closable resources in `container.go`.
- **Handlers MUST** use `pkg/response` helpers (`Success`, `Created`, `BadRequest`, `Conflict`, `Unauthorized`, `Forbidden`, `NotFound`, `ValidationError`, `InternalServerError`, `FromValidateError`, `SuccessWithMeta`). Never return raw `c.JSON` with ad-hoc maps.
- **Validation:** Echo uses a custom validator wrapping `go-playground/validator/v10`. Use `validate` struct tags.
- **Debug routes** (`/api/debug/pprof/*`) are registered only when `APP_DEBUG=true`.
- **Domain errors** are defined in `internal/errors/errors.go`. Handlers branch on these sentinel errors (e.g., `apperrors.ErrUserExists` → `response.Conflict`) — always add new domain errors here.
- **Shared handler helpers** in `internal/handler/helper.go`: `GetUserIDFromClaims` extracts user ID from 3 JWT claim formats (MapClaims, Token, map); `ParsePaginationParams` parses limit/offset and caps limit at 100.
- **S3Storage** (`pkg/storage/s3_storage.go`): `NewS3Storage` can return `nil` if the MinIO client fails to initialize. Services using it must nil-check.
- **ValkeyCache** (`pkg/cache/valkey.go`): `NewValkeyCache` returns `nil` when `VALKEY_URL` is empty or connection fails (fail-open). Services using it must nil-check. Setting `CACHE_TTL_SECONDS=0` disables cache writes even when connected.
- **Package layout:** `internal/dto/` contains request/response DTOs; `internal/model/` contains GORM entities.

## Echo v5 Quirks

- Handler signature uses **pointer receiver**: `func(c *echo.Context) error` — not value receiver like Echo v4.
- **All routes under `/api` prefix** (e.g., `POST /api/auth/register`, `GET /api/posts/:id`). Exception: `/health` and `/` are at root level.
- Auth middleware is passed as last arg to route registration, not chained separately: `posts.POST("", handler.Create, r.authMiddleware.Auth())`.
- **Admin middleware** (`r.authMiddleware.AuthAdmin()`) queries DB for `is_super_admin` — must chain after `Auth()`.
- **Per-route rate limits**: `/api/auth/login` and `/api/auth/forgot-password` have 5 req / 5 min (independent of global `HTTP_RATE_LIMIT_RPS`).

## Environment / Config

- Config package is at **root-level** `config/config.go` (imported as `echobackend/config`), not `internal/config`.
- Loaded from `.env` via `godotenv`, then env vars. See `.env.example` for the full schema.
- **Key env vars with fallback aliases** (first-match wins in `config/config.go`):
  - `S3_ENDPOINT` → `MINIO_ENDPOINT`
  - `S3_ACCESS_KEY` → `MINIO_ACCESS_KEY`
  - `S3_SECRET_KEY` → `MINIO_SECRET_KEY`
  - `S3_BUCKET` → `MINIO_BUCKET`
  - `S3_USE_SSL` → `MINIO_USE_SSL`
  - `DB_POOL_MAX_OPEN` → `MAX_OPEN_CONNS`
  - `DB_POOL_MAX_IDLE` → `MAX_IDLE_CONNS`
  - `DB_POOL_CONN_LIFETIME` → `CONN_MAX_LIFETIME`
  - `HTTP_RATE_LIMIT_RPS` → `RATE_LIMITER_MAX`
  - `HTTP_RATE_LIMIT_WINDOW_SEC` → `RATE_LIMITER_TTL`
  - `HTTP_TRUST_PROXY` → `TRUST_PROXY`
  - `APP_DEBUG` → `DEBUG`
- **`HTTP_ALLOW_ORIGINS`** controls CORS (`Access-Control-Allow-Origin`). Defaults to `"*"`. Comma-separated list or `"*"`.
- **Postgres default DSN requires TLS** (`sslmode=require`). For local dev without SSL, override `DATABASE_URL` (see `.env.example`).
- `HTTP_TRUST_PROXY=true` switches Echo IP extraction to `X-Forwarded-For` (use behind a trusted reverse proxy only).
- `APP_DEBUG=true` also enables GORM `Info` logging (vs `Error`) and registers pprof debug routes.
- **`VALKEY_URL`** must be set for caching to work. Empty string disables caching entirely.

## Database

- GORM with `pgx/v5` driver. Connection pooling and retry logic live in `pkg/database/setup.go`.
- **Uses the default pgx extended query protocol** (named statements, binary encoding). Prepared statements work as expected. If running behind PgBouncer in transaction-pooling mode, you'll need to switch to simple protocol.
- **Migrations are raw SQL** in `migrations/`, managed via [goose](https://github.com/pressly/goose) CLI. Goose env vars are in `.env` (`GOOSE_DRIVER`, `GOOSE_DBSTRING`, `GOOSE_MIGRATION_DIR`, `GOOSE_TABLE`). Use `goose up` / `goose down` / `goose status`. Create new migrations with `goose create <name> sql`.
- **UUID v7** is used for primary keys (see `pkg/validator.IsValidUUID` and migrations).
- **Health check**: `/health` endpoint pings the database (returns 200 or 503). Used by Fly.io and Docker.

## CI / Deploy

- GitHub Actions (`.github/workflows/main.yml`) triggers on push to `main` only.
- Builds a Docker image, pushes to Docker Hub (`cecep31/echobackend:latest`), then deploys to Fly.io.
- `fly.toml` references the pre-built image; it does not build from source. Primary region: `sin` (Singapore).
- **Dockerfile uses `golang:1.26`** while `go.mod` specifies `go 1.25.0`. This works because Go is backward-compatible, but be aware of the mismatch.

## Error Handling Convention

- **Repositories** return raw errors or custom domain errors.
- **Services** wrap/transform repository errors into business-meaningful errors.
- **Handlers** map errors to HTTP status codes using `pkg/response` helpers. Never map errors inline — always push that responsibility to the handler layer.

## Request Body Limit

- Echo middleware enforces **10 MB** body limit (`middleware.BodyLimit(10 * 1024 * 1024)` in `internal/middleware/setup.go:23`). Requests exceeding this return 413.
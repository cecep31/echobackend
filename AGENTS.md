# echobackend — Agent Notes

Go REST API (Echo v5 + GORM + PostgreSQL). Single binary, manual DI, deployed to Fly.io via Docker.

## Build & Run

- **Entry point:** `cmd/main.go`
- **Build:** `make build` → `bin/echobackend`
- **Dev (hot reload):** `make dev` (requires `air`)
- **Run:** `go run cmd/main.go` or `make run`
- **Quality gate:** `make check` runs `fmt → vet → test-race → test-coverage-func`

## Testing

- **All tests:** `make test` (timeout 30s)
- **Race check:** `make test-race`
- **Coverage:** `make test-coverage` or `make test-coverage-func`
- There are currently no unit tests in the project.

## Architecture

- **Manual DI** wired in `internal/di/container.go`. When adding a service/repo/handler:
  1. Create struct + constructor.
  2. Register in `Container` and `NewContainer`.
  3. Thread it through `routes.NewRoutes` and the route setup methods.
- **Graceful shutdown:** `CleanupManager` (LIFO) in `internal/di/cleanup.go`. Register closable resources in `container.go`.
- **Handlers MUST** use `pkg/response` helpers (`Success`, `BadRequest`, `ValidationError`, `NotFound`, etc.). Never return raw `c.JSON` with ad-hoc maps.
- **Validation:** Echo uses a custom validator wrapping `go-playground/validator/v10`. Use `validate` struct tags.
- **Debug routes** (`/v1/debug/pprof/*`) are registered only when `APP_DEBUG=true`.

## Environment / Config

- Loaded from `.env` via `godotenv`, then env vars. See `.env.example` for the full schema.
- **Key env vars with fallback aliases:**
  - `S3_ENDPOINT` → falls back to `MINIO_ENDPOINT`
  - `S3_ACCESS_KEY` → `MINIO_ACCESS_KEY`
  - `S3_SECRET_KEY` → `MINIO_SECRET_KEY`
  - `S3_BUCKET` → `MINIO_BUCKET`
  - `DB_POOL_MAX_OPEN` → `MAX_OPEN_CONNS`
  - `HTTP_TRUST_PROXY` → `TRUST_PROXY`
  - `APP_DEBUG` → `DEBUG`
- **Postgres default DSN requires TLS** (`sslmode=require`). For local dev without SSL, override `DATABASE_URL` (see `.env.example`).
- `HTTP_TRUST_PROXY=true` switches Echo IP extraction to `X-Forwarded-For` (use behind a trusted reverse proxy only).

## Database

- GORM with `pgx/v5` driver. Connection pooling and retry logic live in `pkg/database/setup.go`.
- **Migrations are raw SQL** in `migrations/`. There is no migration runner in the app; apply them manually or with an external tool (e.g., `migrate`, `psql`).

## CI / Deploy

- GitHub Actions (`.github/workflows/main.yml`) builds a Docker image, pushes to Docker Hub (`cecep31/echobackend:latest`), then deploys to Fly.io.
- `fly.toml` references the pre-built image; it does not build from source.

## Linting

- `make lint` requires `golangci-lint` installed separately.
- `make fmt` + `make vet` are always available.

# Echo Backend API

A modern REST API built with Go, Echo v5, GORM, and PostgreSQL.

## Tech Stack

- **Go** 1.25+
- **Echo** v5 (web framework)
- **GORM** v2 (ORM)
- **PostgreSQL** 14+ (via pgx/v5)
- **Valkey/Redis** (caching, optional)
- **MinIO/S3** (file storage, optional)
- **JWT** (authentication)
- **Docker** (deployment)

## Quick Start

```bash
# Clone and setup
cp .env.example .env

# Edit .env with your config

# Run with hot reload (requires air)
air

# Or run normally
go run cmd/main.go
```

Server starts at `http://localhost:8080`

## Commands

| Command | Description |
|---------|-------------|
| `go build -o bin/main cmd/main.go` | Build binary |
| `go run cmd/main.go` | Run the server |
| `air` | Run with hot reload |
| `go test ./...` | Run all tests |
| `go test -race ./...` | Run tests with race detection |
| `go test -cover ./...` | Run tests with coverage |
| `go vet ./...` | Run go vet |
| `go fmt ./...` | Format code |
| `golangci-lint run` | Run linter (requires install) |
| `gosec ./...` | Security scan (requires install) |

## Environment Variables

Key environment variables (with fallback aliases):

| Variable | Fallback | Description |
|----------|----------|-------------|
| `S3_ENDPOINT` | `MINIO_ENDPOINT` | S3/MinIO endpoint |
| `S3_ACCESS_KEY` | `MINIO_ACCESS_KEY` | Access key |
| `S3_SECRET_KEY` | `MINIO_SECRET_KEY` | Secret key |
| `S3_BUCKET` | `MINIO_BUCKET` | Bucket name |
| `S3_USE_SSL` | `MINIO_USE_SSL` | Use SSL for S3 |
| `DB_POOL_MAX_OPEN` | `MAX_OPEN_CONNS` | Max open DB connections |
| `DB_POOL_MAX_IDLE` | `MAX_IDLE_CONNS` | Max idle DB connections |
| `DB_POOL_CONN_LIFETIME` | `CONN_MAX_LIFETIME` | Connection max lifetime |
| `HTTP_RATE_LIMIT_RPS` | `RATE_LIMITER_MAX` | Global rate limit (req/sec) |
| `HTTP_RATE_LIMIT_WINDOW_SEC` | `RATE_LIMITER_TTL` | Rate limit window |
| `HTTP_TRUST_PROXY` | `TRUST_PROXY` | Trust X-Forwarded-For |
| `APP_DEBUG` | `DEBUG` | Enable debug mode & pprof routes |

For full config, see `.env.example`.

## Project Structure

```
├── cmd/main.go              # Entry point
├── config/config.go          # Configuration (env-based)
├── internal/
│   ├── di/                   # Dependency injection (container.go, cleanup.go)
│   ├── handler/              # HTTP handlers
│   ├── service/              # Business logic
│   ├── repository/           # Data access
│   ├── model/                # GORM models
│   ├── dto/                  # Request/response DTOs
│   ├── errors/               # Domain error sentinels
│   ├── middleware/            # Echo middleware
│   └── routes/               # Route definitions
├── pkg/
│   ├── response/             # API response helpers
│   ├── database/             # DB setup & pooling
│   ├── cache/                # Valkey/Redis cache
│   ├── storage/              # S3/MinIO storage
│   └── validator/            # Custom Echo validator
├── migrations/                # SQL migrations (goose)
├── .air.toml                 # Air hot reload config
└── Dockerfile                # Docker build
```

## Development

### Cross-Platform Support

`.air.toml` dikonfigurasi untuk multi-platform:
- **Windows**: output `tmp/main.exe`
- **Linux/macOS**: output `tmp/main`

### Database Migrations (Goose)

Migrations menggunakan [goose](https://github.com/pressly/goose) CLI. Konfigurasi sudah ada di `.env`:

```bash
# Install goose (sekali saja)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Apply semua migration
goose up

# Rollback 1 migration
goose down

# Cek status migration
goose status

# Buat migration baru
goose create nama_migration sql
```

Env vars di `.env` yang dipakai goose:

| Variable | Description |
|----------|-------------|
| `GOOSE_DRIVER` | Driver database (`postgres`) |
| `GOOSE_DBSTRING` | Koneksi string (bisa pakai `${DATABASE_URL}`) |
| `GOOSE_MIGRATION_DIR` | Folder migration (`./migrations`) |
| `GOOSE_TABLE` | Nama tabel tracking (`custom.goose_migrations`) |

> Jika `GOOSE_DBSTRING` tidak resolve `${DATABASE_URL}` secara otomatis, export manual:
> ```bash
> export GOOSE_DBSTRING="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
> ```

## Deployment

Docker image di-deploy ke Fly.io via GitHub Actions (`.github/workflows/main.yml`).

```bash
# Build locally
docker build -t echobackend .

# Run
docker run -p 8080:8080 echobackend
```

## License

MIT
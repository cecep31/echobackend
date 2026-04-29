# Echo Backend API

A modern REST API built with Go, Echo v5, GORM, and PostgreSQL.

## Tech Stack

- **Go** 1.25+
- **Echo** v5 (web framework)
- **GORM** v1 (ORM)
- **PostgreSQL** 14+
- **MinIO/S3** (file storage)
- **JWT** (authentication)
- **Docker** (deployment)

## Quick Start

```bash
# Clone and setup
cp .env.example .env

# Edit .env with your config

# Run with hot reload
make dev

# Or run normally
make run
```

Server starts at `http://localhost:8080`

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make build` | Build binary to `bin/echobackend` |
| `make run` | Run with `go run` |
| `make dev` | Run with **air** (hot reload) |
| `make test` | Run all tests |
| `make test-race` | Run with race detection |
| `make test-coverage` | Generate coverage HTML report |
| `make test-coverage-func` | Show coverage by function |
| `make fmt` | Format code |
| `make vet` | Run go vet |
| `make lint` | Run golangci-lint |
| `make check` | Run fmt → vet → test-race → test-coverage-func |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |

## Environment Variables

Key environment variables (with fallback aliases):

| Variable | Fallback | Description |
|----------|----------|-------------|
| `S3_ENDPOINT` | `MINIO_ENDPOINT` | S3/MinIO endpoint |
| `S3_ACCESS_KEY` | `MINIO_ACCESS_KEY` | Access key |
| `S3_SECRET_KEY` | `MINIO_SECRET_KEY` | Secret key |
| `S3_BUCKET` | `MINIO_BUCKET` | Bucket name |
| `DB_POOL_MAX_OPEN` | `MAX_OPEN_CONNS` | Max open DB connections |
| `HTTP_TRUST_PROXY` | `TRUST_PROXY` | Trust X-Forwarded-For |
| `APP_DEBUG` | `DEBUG` | Enable debug mode |

For full config, see `.env.example`.

## Project Structure

```
├── cmd/main.go           # Entry point
├── internal/
│   ├── di/               # Dependency injection (container.go)
│   ├── handler/          # HTTP handlers
│   ├── service/          # Business logic
│   ├── repository/       # Data access
│   ├── model/            # GORM models
│   ├── middleware/       # Echo middleware
│   └── routes/           # Route definitions
├── pkg/
│   ├── response/         # API response helpers
│   └── database/         # DB setup & pooling
├── migrations/           # Raw SQL migrations (manual)
├── .air.toml             # Air hot reload config
└── Makefile
```

## Development

### Cross-Platform Support

`.air.toml` sudah dikonfigurasi untuk:
- **Linux/macOS**: gunakan `[build]` section
- **Windows**: gunakan `[build.windows]` section

### Manual Migrations

Migrations adalah raw SQL di folder `migrations/`. Tidak ada migration runner bawaan. Apply manual:

```bash
psql -d your_database -f migrations/*.sql
```

## API Documentation

See `api_doc.md` for detailed endpoints.

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
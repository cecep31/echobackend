# Gemini CLI Project Context: echobackend

## Project Overview
`echobackend` is a modern REST API for a blog/content management system. It provides a robust backend for managing users, posts, comments, chat conversations, and file uploads.

### Core Technologies
- **Language:** Go 1.25+
- **Web Framework:** [Echo v5](https://github.com/labstack/echo)
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL driver
- **Database:** PostgreSQL 14+
- **Cache:** [Valkey](https://valkey.io/) (Redis-compatible)
- **Authentication:** JWT (via `github.com/golang-jwt/jwt/v5`)
- **Storage:** MinIO/S3 compatible storage (via `github.com/minio/minio-go/v7`)
- **Validation:** Go Playground Validator v10
- **Dev Tooling:** [Air](https://github.com/air-verse/air) for hot reloading

### Architecture
The project follows a layered "Clean Architecture" pattern with strict separation of concerns:
- **`cmd/`**: Entry point (`main.go`) where configuration is loaded and the server is started.
- **`internal/`**: Core business logic.
    - **`handler/`**: HTTP request handling, input validation, and response formatting.
    - **`service/`**: Domain logic and orchestration between repositories.
    - **`repository/`**: Data access layer using GORM.
    - **`model/`**: GORM database entities and common structs.
    - **`di/`**: Manual dependency injection container (`container.go`) and resource cleanup.
    - **`middleware/`**: Custom Echo middlewares (Auth, Logging, setup, etc.).
    - **`routes/`**: Centralized route definitions.
- **`pkg/`**: Shared utility packages.
    - **`response/`**: API response helpers.
    - **`database/`**: DB setup & pooling.
    - **`cache/`**: Valkey cache implementation.
    - **`storage/`**: S3 storage implementation.
    - **`validator/`**: Custom Echo validator wrapper.
- **`migrations/`**: Raw SQL-based database migrations (applied manually).

## Building and Running

### Development Commands
- **Install Dependencies:** `go mod download`
- **Hot Reload (Dev Mode):** `make dev` (requires `air`)
- **Run Manually:** `go run cmd/main.go` or `make run`
- **Build Binary:** `make build` (outputs to `bin/echobackend`)
- **Clean Artifacts:** `make clean`

### Testing
- **Run All Tests:** `make test`
- **Run with Race Detection:** `make test-race`
- **Generate Coverage:** `make test-coverage`
- **Show Coverage by Function:** `make test-coverage-func`

### Code Quality
- **Format Code:** `make fmt`
- **Run Vet:** `make vet`
- **Run Linter:** `make lint` (requires `golangci-lint`)
- **Security Audit:** `make security` (requires `gosec`)
- **Full Quality Check:** `make check`

### Docker
- **Build Image:** `make docker-build`
- **Run Container:** `make docker-run`

## Development Conventions

### 1. API Responses
All handlers MUST use the standard response utilities in `pkg/response`. This ensures a consistent JSON structure (`success`, `message`, `data`, `error`, `errors`, `meta`).
- `response.Success(c, message, data)`
- `response.Created(c, message, data)`
- `response.BadRequest(c, message, err)`
- `response.Unauthorized(c, message)`
- `response.Forbidden(c, message)`
- `response.NotFound(c, message, err)`
- `response.ValidationError(c, message, err)`
- `response.FromValidateError(c, err)`: Specialized helper for Echo validation errors.
- `response.InternalServerError(c, message, err)`
- `response.Conflict(c, message, conflictError)`

### 2. Dependency Injection
Dependencies are manually wired in `internal/di/container.go`. When adding a new service or repository:
1. Create the new struct and constructor.
2. Register it in the `Container` struct and `NewContainer` function.
3. Pass it through to the relevant handler and routes.
4. If the dependency requires cleanup (e.g., closing a connection), register it with the `CleanupManager`.

### 3. Resource Cleanup
The project uses a `CleanupManager` in `internal/di` to ensure database connections, cache connections, and other resources are closed gracefully during shutdown. Register new closable resources in `internal/di/container.go`.

### 4. Validation
Use struct tags with the `validate` key for input validation. The project uses `go-playground/validator`. Handlers should use `c.Bind` and `c.Validate` (which is automatically called if configured in Echo).

### 5. Error Handling
- Repositories should return raw database errors or custom domain errors.
- Services should wrap or transform repository errors into business-meaningful errors.
- Handlers are responsible for mapping errors to the appropriate HTTP status codes using `pkg/response`.

### 6. Environment Variables
Local configuration uses `.env`. Refer to `config/config.go` for the schema of available environment variables and their fallbacks (e.g., `S3_ENDPOINT` fallbacks to `MINIO_ENDPOINT`).
- `PORT`: HTTP server port (default: 8080).
- `DATABASE_URL`: Postgres DSN.
- `VALKEY_URL`: Connection URL for Valkey/Redis.
- `JWT_SECRET`: Secret for signing JWTs.
- `APP_DEBUG` or `DEBUG`: Enables verbose logging and GORM info logs.

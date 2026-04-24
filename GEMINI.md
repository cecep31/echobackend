# Gemini CLI Project Context: echobackend

## Project Overview
`echobackend` is a modern REST API for a blog/content management system. It provides a robust backend for managing users, posts, comments, chat conversations, and file uploads.

### Core Technologies
- **Language:** Go 1.25+
- **Web Framework:** [Echo v5](https://github.com/labstack/echo)
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL driver
- **Database:** PostgreSQL 14+
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
    - **`middleware/`**: Custom Echo middlewares (Auth, Logging, etc.).
    - **`routes/`**: Centralized route definitions.
- **`pkg/`**: Shared utility packages (database setup, custom validator, standard response formatters).
- **`migrations/`**: SQL-based database migrations.

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

### Code Quality
- **Format Code:** `make fmt`
- **Run Vet:** `make vet`
- **Run Linter:** `make lint` (requires `golangci-lint`)
- **Full Quality Check:** `make check`

### Docker
- **Build Image:** `make docker-build`
- **Run Container:** `make docker-run`

## Development Conventions

### 1. API Responses
All handlers MUST use the standard response utilities in `pkg/response`. Never return `c.JSON` directly with raw maps.
- `response.Success(c, message, data)`
- `response.BadRequest(c, message, err)`
- `response.ValidationError(c, message, err)`
- `response.NotFound(c, message, err)`

### 2. Dependency Injection
Dependencies are manually wired in `internal/di/container.go`. When adding a new service or repository:
1. Create the new struct and constructor.
2. Register it in the `Container` struct and `NewContainer` function.
3. Pass it through to the relevant handler and routes.

### 3. Resource Cleanup
The project uses a `CleanupManager` to ensure database connections and other resources are closed gracefully during shutdown. Register new closable resources in `internal/di/container.go`.

### 4. Validation
Use struct tags with the `validate` key for input validation. Echo is configured to use a custom validator that wraps `go-playground/validator`.

### 5. Error Handling
- Repositories should return raw database errors or custom domain errors.
- Services should wrap or transform repository errors into business-meaningful errors.
- Handlers are responsible for mapping errors to the appropriate HTTP status codes using `pkg/response`.

### 6. Environment Variables
Always use `.env` for local configuration. Reference `config/config.go` for the schema of available environment variables.

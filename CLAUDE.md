# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build, Lint, and Test Commands

### Essential Commands
- **Build application**: `make build` or `go build -o bin/echobackend cmd/main.go`
- **Run application**: `make run` or `go run cmd/main.go`
- **Development with hot reload**: `make dev` (uses air)
- **Run tests**: `make test` or `go test -timeout 30s -v ./...`
- **Run short tests**: `make test-short`
- **Run tests with race detection**: `make test-race`
- **Run tests with coverage**: `make test-coverage` (generates HTML report)
- **Run benchmarks**: `make bench`
- **Format code**: `make fmt` or `go fmt ./...`
- **Run linter**: `make lint` (uses golangci-lint)
- **Run vet**: `make vet` or `go vet ./...`
- **Run all quality checks**: `make check` (fmt + vet + test-race + test-coverage-func)
- **Clean build artifacts**: `make clean`

### Docker Commands
- **Build Docker image**: `make docker-build`
- **Run Docker container**: `make docker-run`

### Project Configuration
- **Module management**: `make mod-tidy` (go mod tidy)
- **Download dependencies**: `make mod-download` (go mod download)

## High-Level Architecture

### Project Structure
```
echobackend/
├── cmd/
│   └── main.go                    # Application entry point with dependency injection
├── config/
│   └── config.go                 # Configuration management
├── internal/
│   ├── di/                       # Dependency injection container (uber/dig)
│   │   ├── container.go          # DI container setup
│   │   └── cleanup.go            # Resource cleanup management
│   ├── handler/                  # HTTP handlers (11 handlers)
│   │   ├── auth_handler.go
│   │   ├── comment_handler.go
│   │   ├── page_handler.go
│   │   ├── post_handler.go
│   │   ├── post_like_handler.go
│   │   ├── post_view_handler.go
│   │   ├── tag_handler.go
│   │   ├── user_follow_handler.go
│   │   ├── user_handler.go
│   │   └── workspace_handler.go
│   ├── middleware/               # Custom middleware
│   │   ├── auth_middleware.go    # JWT authentication
│   │   └── setup.go              # Middleware initialization
│   ├── repository/               # Database repositories (11 repos)
│   │   ├── comment_repository.go
│   │   ├── post_repository.go
│   │   ├── post_like_repository.go
│   │   ├── post_view_repository.go
│   │   ├── session_repository.go
│   │   ├── tag_repository.go
│   │   ├── user_follow_repository.go
│   │   ├── user_repository.go
│   │   └── workspace_repository.go
│   ├── routes/                   # API routing
│   │   └── routes.go             # Route setup (v1 API)
│   ├── service/                  # Business logic (10 services)
│   │   ├── auth_service.go
│   │   ├── comment_service.go
│   │   ├── page_service.go
│   │   ├── post_service.go
│   │   ├── post_like_service.go
│   │   ├── post_view_service.go
│   │   ├── tag_service.go
│   │   ├── user_follow_service.go
│   │   ├── user_service.go
│   │   └── workspace_service.go
│   └── model/                    # Database models
│       ├── block.go
│       ├── comment.go
│       ├── page.go
│       ├── post.go
│       ├── post_comment.go
│       ├── post_like.go
│       ├── post_view.go
│       ├── session.go
│       ├── tag.go
│       ├── user.go
│       ├── user_follow.go
│       ├── workspace.go
│       └── workspace_user.go
├── pkg/
│   ├── database/                 # Database utilities
│   │   ├── setup.go              # Database initialization
│   │   └── wrapper.go            # Database wrapper with cleanup
│   ├── response/                 # HTTP response utilities
│   │   └── response.go
│   ├── storage/                  # Storage utilities
│   │   └── s3_storage.go         # S3 storage implementation
│   ├── utils/                    # Utility functions
│   │   ├── errors.go             # Error handling utilities
│   │   ├── pagination.go         # Pagination utilities
│   │   └── security.go           # Security utilities
│   ├── validator/                # Custom validator
│   │   └── validator.go          # Validator implementation
│   └── worker/                   # Worker pool
│       └── pool.go
├── migrations/                   # Database migrations
├── test/                         # Test files
└── scripts/                      # Utility scripts
```

### Key Architectural Patterns

1. **Dependency Injection**: Uses uber/dig for dependency injection with cleanup management
2. **Layered Architecture**: Clear separation between handlers, services, and repositories
3. **REST API**: v1 API with comprehensive endpoints for blog/content management
4. **JWT Authentication**: Custom auth middleware with token validation and admin checks
5. **GORM + PostgreSQL**: ORM with PostgreSQL database (upgraded to v1.31.0)
6. **MinIO Storage**: Migrated from AWS SDK to MinIO for S3 storage operations
7. **Graceful Shutdown**: Proper resource cleanup on application shutdown
8. **Input Validation**: Enhanced validation with custom validator and security utilities

### Main Application Flow
1. `cmd/main.go` initializes DI container and Echo server
2. Routes are set up in `internal/routes/routes.go` with v1 API group
3. Handlers process HTTP requests and delegate to services
4. Services contain business logic and call repositories
5. Repositories handle database operations with GORM
6. Models define database schema and relationships

### Configuration
- Environment-based configuration in `config/config.go`
- Database configuration with connection pooling
- JWT secret for authentication
- S3 storage configuration
- Debug mode support

### Testing
- Basic test structure in `test/memory_leak_test.go`
- Uses standard Go testing framework
- Dependency injection testing with cleanup manager
- Memory leak testing for resource management
- Run `make test-short` for quick validation during development

### Code Quality
- Comprehensive golangci-lint configuration
- 24 different linters enabled including staticcheck, govet, errcheck
- Custom linting rules for test files and exclusions
- Air for hot reloading in development

### Recent Changes (from git log)
- Enhanced security and validation (34bab2b)
- Comprehensive project documentation (a1b22b1)
- Go module dependencies update (438d724)
- JWT authentication middleware with token validation and admin checks (fc82649)
- Input validation and database connection reliability improvements (0cf408a)
- Migrate from AWS SDK to MinIO for S3 storage operations (a3ee6f9)
- Post publishing functionality (440c0f3)
- JWT claim key standardization to snake_case (db3df5c)
- GORM upgrade to v1.31.0 (204f3a2)
- Reduced maximum post limit from 50 to 20 in GetPostsRandom handler (647e63d)

## Common Development Tasks

### Adding a New API Endpoint
1. Create model in `internal/model/`
2. Create repository interface and implementation in `internal/repository/`
3. Create service in `internal/service/`
4. Create handler in `internal/handler/`
5. Register dependencies in `internal/di/container.go`
6. Add route in `internal/routes/routes.go`

### Running the Application
1. Set up PostgreSQL database
2. Configure environment variables in `.env`
3. Run `make build` or `make dev` for development
4. Access API at `http://localhost:8080`

### Testing
- Use `make test` for standard tests
- Use `make test-race` for race condition detection
- Use `make test-coverage` for coverage reports
- Current test structure is minimal, consider expanding test coverage

### Development Workflow
1. Use `make dev` for hot reloading during development
2. Run `make check` before commits (includes fmt, vet, race tests, coverage)
3. Use `make lint` for code quality checks
4. Run `make test` to ensure tests pass
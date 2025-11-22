# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build, Lint, and Test Commands

### Essential Commands
- **Build application**: `make build` or `go build -o bin/echobackend cmd/main.go`
- **Run application**: `make run` or `go run cmd/main.go`
- **Development with hot reload**: `make dev` (uses air)
- **Run tests**: `make test` or `go test -timeout 30s -v ./...`
- **Run short tests**: `make test-short` (excludes long-running tests)
- **Run tests with race detection**: `make test-race` (detects race conditions)
- **Run tests with coverage**: `make test-coverage` (generates HTML report in `coverage/` directory)
- **Run benchmarks**: `make bench` (runs benchmarks with memory stats)
- **Format code**: `make fmt` or `go fmt ./...`
- **Run linter**: `make lint` (uses golangci-lint with comprehensive rules)
- **Run vet**: `make vet` or `go vet ./...`
- **Run all quality checks**: `make check` (fmt + vet + test-race + test-coverage-func)
- **Clean build artifacts**: `make clean`
- **Security analysis**: `make security` (uses gosec for security scanning)
- **Module tidy**: `make mod-tidy` (clean up go.mod and go.sum)
- **Download dependencies**: `make mod-download` (download all dependencies)

### Test Commands (Specific)
- **Run specific test file**: `go test -v ./path/to/file_test.go`
- **Run specific test function**: `go test -v -run TestFunctionName ./...`
- **Run tests for specific package**: `go test -v ./internal/handler/`
- **Run tests with verbose output**: `go test -v -timeout 30s ./...`
- **Run tests in parallel**: `go test -parallel 4 ./...`
- **Debug specific test**: `go test -v -run TestFunctionName -debug ./...`

### Docker Commands
- **Build Docker image**: `make docker-build`
- **Run Docker container**: `make docker-run`
- **Build and run with environment**: `docker build -t echobackend . && docker run -p 8080:8080 echobackend`

### Project Configuration
- **Module management**: `make mod-tidy` (go mod tidy)
- **Download dependencies**: `make mod-download` (go mod download)
- **Show available commands**: `make help`

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
├── migrations/                   # Database migrations
├── test/                         # Test files
└── scripts/                      # Utility scripts
```

### Key Architectural Patterns

1. **Dependency Injection**: Uses uber/dig for dependency injection with cleanup management
2. **Layered Architecture**: Clear separation between handlers, services, and repositories
3. **REST API**: v1 API with comprehensive endpoints for blog/content management
4. **JWT Authentication**: Custom auth middleware with token validation and admin checks
5. **GORM + PostgreSQL**: ORM with PostgreSQL database (upgraded to v1.31.1)
6. **MinIO Storage**: Migrated from AWS SDK to MinIO for S3 storage operations
7. **Graceful Shutdown**: Proper resource cleanup on application shutdown
8. **Input Validation**: Enhanced validation with custom validator and security utilities
9. **Middleware Pipeline**: Custom middleware for logging, CORS, recovery, and auth
10. **Error Handling**: Consistent error response format across all endpoints

### Main Application Flow
1. `cmd/main.go` initializes DI container and Echo server
2. Routes are set up in `internal/routes/routes.go` with v1 API group
3. Handlers process HTTP requests and delegate to services
4. Services contain business logic and call repositories
5. Repositories handle database operations with GORM
6. Models define database schema and relationships
7. Middleware handles authentication, logging, and error handling

### Configuration
- Environment-based configuration in `config/config.go`
- Database configuration with connection pooling (max 30 open, 2 idle, 30min lifetime)
- JWT secret for authentication with configurable expiration
- S3/MinIO storage configuration for file uploads
- Debug mode support with pprof endpoints
- Rate limiting configuration (optional)
- Support for both MinIO and AWS S3

### Testing
- Test files in `test/` directory and alongside packages (`*_test.go`)
- Uses standard Go testing framework with testify for assertions
- Dependency injection testing with cleanup manager in `test/memory_leak_test.go`
- Security utilities testing in `pkg/utils/security_test.go`
- Post filtering and query testing in `test/post_filter_test.go`
- Memory leak testing for resource management
- Run `make test-short` for quick validation during development
- Run `make test-race` to detect race conditions in concurrent code

### Code Quality
- Comprehensive golangci-lint configuration with 24 linters
- Linters include staticcheck, govet, errcheck, gocyclo, gosec, and more
- Custom linting rules for test files and exclusions
- Line length limit: 140 characters
- Cyclomatic complexity limit: 15
- Air for hot reloading in development
- Security scanning with gosec
- Docker multi-stage builds for production

### Debugging and Profiling
- **Debug mode**: Set `DEBUG=true` in environment to enable pprof endpoints
- **Profiling endpoints**: Available at `/v1/debug/pprof/*` when debug mode is enabled
- **Memory profiling**: `go tool pprof http://localhost:8080/v1/debug/pprof/heap`
- **CPU profiling**: `go tool pprof http://localhost:8080/v1/debug/pprof/profile`
- **Goroutine analysis**: `go tool pprof http://localhost:8080/v1/debug/pprof/goroutine`
- **Trace analysis**: `go tool trace http://localhost:8080/v1/debug/pprof/trace`
- **Development server**: `make dev` for hot reloading with air

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

### Development Workflow

1. **Setup and Development**:
   ```bash
   # Clone and setup
   git clone <your-repo-url>
   cd echobackend
   cp .env.example .env
   # Edit .env with your database credentials

   # Install dependencies
   make mod-download

   # Start development server with hot reload
   make dev

   # Or run directly
   make run
   ```

2. **Before Committing**:
   ```bash
   # Run all quality checks
   make check

   # Run linter
   make lint

   # Run tests with coverage
   make test-coverage

   # Security scan
   make security
   ```

3. **Database Setup**:
   - Requires PostgreSQL 14+
   - Default connection: `postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable`
   - Connection pooling: 30 max open, 2 idle, 30min lifetime
   - Migrations are in `migrations/` directory

4. **Environment Variables** (see `.env.example`):
   - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
   - `JWT_SECRET` - Required for authentication
   - `MINIO_ENDPOINT`, `MINIO_ACCESS_KEY`, `MINIO_SECRET_KEY`, `MINIO_BUCKET`
   - `DEBUG` - Enable debug mode and pprof endpoints
   - `PORT` - Server port (default: 8080)

### API Endpoints Overview

The API provides comprehensive endpoints for:
- **Authentication**: Register, login, JWT token management
- **Users**: User management, profiles, admin operations
- **Posts**: Create, read, update, delete posts with pagination
- **Comments**: Post comments with real-time WebSocket support
- **Tags**: Tag management for posts
- **User Interactions**: Follow/unfollow users, like posts
- **Workspaces**: Collaborative workspace management
- **Pages**: Content pages within workspaces
- **Chat**: Conversations and messages
- **Analytics**: Post views, likes, and engagement metrics

### WebSocket Support
- **Endpoint**: `wss://nest.pilput.me/ws/posts`
- **Authentication**: JWT token required as query parameter
- **Features**: Real-time comments, typing indicators, read receipts
- **Events**: `sendComment`, `typing`, `markAsRead`, `newComment`, `userTyping`

### Deployment

1. **Local Development**:
   ```bash
   # Using Docker
   make docker-build
   make docker-run

   # Direct deployment
   make build
   ./bin/echobackend
   ```

2. **Production Considerations**:
   - Set `DEBUG=false`
   - Configure proper JWT secrets
   - Set up PostgreSQL with proper credentials
   - Configure MinIO/S3 for file storage
   - Set up reverse proxy (nginx) for SSL termination
   - Configure proper CORS settings
   - Set up monitoring and logging

### Adding a New API Endpoint
1. Create model in `internal/model/`
2. Create repository interface and implementation in `internal/repository/`
3. Create service in `internal/service/`
4. Create handler in `internal/handler/`
5. Register dependencies in `internal/di/container.go`
6. Add route in `internal/routes/routes.go`
7. Add comprehensive tests
8. Update API documentation

### Running the Application
1. Set up PostgreSQL database
2. Configure environment variables in `.env`
3. Run `make build` or `make dev` for development
4. Access API at `http://localhost:8080`
5. Access API documentation at `http://localhost:8080/v1` endpoints

### Testing
- Use `make test` for standard tests
- Use `make test-race` for race condition detection
- Use `make test-coverage` for coverage reports
- Use `make test-short` for quick validation during development
- Current test structure includes basic functionality tests
- Consider expanding test coverage for critical paths and edge cases

### Performance and Monitoring
- **Profiling**: Enable with `DEBUG=true` and access `/v1/debug/pprof/*`
- **Memory Usage**: Monitor with `go tool pprof` heap profiles
- **CPU Usage**: Monitor with `go tool pprof` CPU profiles
- **Goroutines**: Monitor with `go tool pprof` goroutine profiles
- **Response Times**: Built-in logging for request/response times
- **Database Performance**: GORM logging and connection pool monitoring

### Security Considerations
- **JWT Authentication**: All protected endpoints require valid JWT tokens
- **Input Validation**: Custom validator with security utilities
- **Rate Limiting**: Configurable rate limiting for sensitive endpoints
- **SQL Injection**: Protected by GORM ORM and parameterized queries
- **XSS Protection**: Input sanitization and validation
- **CORS**: Configurable CORS settings for frontend integration
- **Security Headers**: Built-in security middleware

## Dependencies and Tools

### Core Dependencies
- **Web Framework**: Echo v4.13.4
- **ORM**: GORM v1.31.1 with PostgreSQL driver
- **Dependency Injection**: Uber Dig v1.19.0
- **JWT Authentication**: github.com/golang-jwt/jwt/v5
- **Validation**: go-playground/validator/v10
- **Logging**: rs/zerolog
- **S3 Storage**: MinIO Go client v7.0.97
- **Testing**: github.com/stretchr/testify

### Development Tools
- **Hot Reload**: Air for development
- **Linting**: GolangCI-Lint with 24 linters
- **Security**: GoSec for security analysis
- **Database**: PostgreSQL 14+
- **Storage**: MinIO or AWS S3

### Code Quality Standards
- **Line Length**: Maximum 140 characters
- **Cyclomatic Complexity**: Maximum 15
- **Test Coverage**: Target 80%+ coverage
- **Security**: Regular security scans with gosec
- **Documentation**: Comprehensive API documentation

## Troubleshooting

### Common Issues
1. **Database Connection Failed**:
   - Check PostgreSQL is running
   - Verify database credentials in `.env`
   - Ensure database exists and user has permissions

2. **JWT Authentication Errors**:
   - Verify `JWT_SECRET` is set and non-empty
   - Check token expiration time
   - Ensure proper Authorization header format

3. **File Upload Errors**:
   - Verify MinIO/S3 configuration
   - Check bucket exists and has proper permissions
   - Ensure network connectivity to storage service

4. **Memory Issues**:
   - Monitor with `go tool pprof` heap profiles
   - Check for goroutine leaks
   - Verify proper resource cleanup

5. **Performance Issues**:
   - Use pprof CPU profiling
   - Check database query performance
   - Monitor connection pool usage

### Getting Help
- Check API documentation at `api_doc.md`
- Review error logs for detailed information
- Use debugging endpoints when `DEBUG=true`
- Check recent commits for breaking changes

## Contributing Guidelines

### Code Style
- Follow Go naming conventions
- Use meaningful variable and function names
- Keep functions focused and under 50 lines when possible
- Add comments for complex logic
- Use godoc format for public APIs

### Testing Requirements
- Write tests for new features
- Ensure existing tests pass
- Add integration tests for API endpoints
- Test edge cases and error conditions
- Maintain test coverage above 80%

### Commit Messages
- Use clear, descriptive commit messages
- Follow conventional commit format when possible
- Reference related issues or pull requests
- Include motivation for changes when not obvious

### Pull Request Guidelines
- Ensure all tests pass
- Include appropriate test coverage
- Update documentation for new features
- Squash commits for clean history
- Address all code review feedback
# Echo Backend API - Project Context

## Project Overview

Echo Backend API is a REST API for a blog/content management system built with Go, Echo framework, and PostgreSQL. It provides a full-featured backend for managing users, posts, comments, tags, workspaces, and pages with comprehensive authentication and authorization features.

### Key Technologies & Architecture

- **Go 1.24.0** with modern dependency management
- **Echo v4 framework** for HTTP routing
- **PostgreSQL** database with GORM ORM
- **JWT authentication** for secure API access
- **Dependency injection** using go.uber.org/dig
- **Docker support** for containerization
- **S3 storage** (MinIO compatible) for file uploads
- **Gin-gonic style** architecture with separation of concerns (handlers, services, repositories)

### Core Features

- User authentication & management
- Posts, comments, and tags
- User follows and post likes
- Workspaces and pages
- JWT authentication
- PostgreSQL database
- Rate limiting capabilities
- File storage with S3/MinIO
- Comprehensive logging with zerolog

### Project Structure

```
├── cmd/                    # Application entry point (main.go)
├── internal/               # Private application code
│   ├── di/                 # Dependency injection setup
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic
│   ├── repository/         # Data access layer
│   ├── model/              # Database models
│   └── middleware/         # Custom middleware
├── pkg/                    # Shared packages
│   ├── database/           # Database connection & wrapper
│   ├── storage/            # S3 storage interface
│   ├── validator/          # Custom validation
│   └── ...                 # Other shared utilities
├── config/                 # Configuration management
├── migrations/             # Database migrations (future)
├── api_doc.md              # API documentation
└── Makefile                # Build and development commands
```

## Building and Running

### Prerequisites
- Go 1.21+
- PostgreSQL 14+

### Setup and Development

1. **Clone and setup:**
   ```bash
   git clone <your-repo-url>
   cd echobackend
   cp .env.example .env
   ```

2. **Configure database:**
   Edit `.env` file with your PostgreSQL credentials:
   ```env
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=your_database_name
   ```

3. **Build and run:**
   ```bash
   go mod download
   go run cmd/main.go
   ```

### Make Commands

The project includes a comprehensive Makefile with the following commands:
- `make build` - Build the binary
- `make run` - Run the application
- `make test` - Run all tests
- `make test-race` - Run tests with race detection
- `make test-coverage` - Generate coverage report
- `make fmt` - Format code
- `make lint` - Run linter
- `make check` - Run all quality checks
- `make dev` - Run with hot reload using Air
- `make docker-build` - Build Docker image
- `make docker-run` - Run Docker container

### Development Conventions

- **Code formatting:** Follow standard Go formatting (use `make fmt`)
- **Testing:** Write tests for new features (use `make test`)
- **Dependencies:** Use Go modules, organize with dependency injection
- **Error handling:** Proper error handling and logging
- **Configuration:** Use environment variables with defaults
- **Authentication:** JWT-based authentication for protected routes

### Configuration

The application uses a comprehensive configuration system that reads from environment variables with sensible defaults:

- **App_Port** - Server port (default: 8080)
- **JWT_SECRET** - JWT signing key
- **DATABASE_URL** - PostgreSQL connection string
- **MaxOpenConns/MaxIdleConns** - Database connection pool settings
- **S3 Configuration** - For file storage (MinIO compatible)
- **Rate limiting settings** - For API protection

### Testing

The project supports multiple testing approaches:
- Unit tests: `make test`
- Short tests: `make test-short`
- Race condition testing: `make test-race`
- Coverage reports: `make test-coverage`

### Docker Support

The application includes a Dockerfile and can be built and run as a containerized application:
- Build: `make docker-build`
- Run: `make docker-run`

### API Documentation

Detailed API endpoints and examples are available in [api_doc.md](api_doc.md).

## Architecture Patterns

The application follows clean architecture principles with clear separation of concerns:

- **Models**: Define data structures
- **Repositories**: Handle database operations
- **Services**: Implement business logic
- **Handlers**: Handle HTTP requests/responses
- **Middleware**: Handle cross-cutting concerns (auth, logging, etc.)
- **Dependency Injection**: Wire components together using go.uber.org/dig

This architecture promotes testability, maintainability, and scalability of the application.
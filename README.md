# Echo Backend API

A RESTful API backend service built with Echo v4 framework and PostgreSQL database.

## Project Structure

```
.
├── cmd/                    # Main applications for this project
│   └── main.go            # Application entry point
├── internal/              # Private application and library code
│   ├── api/              # API handlers
│   ├── middleware/       # Custom middleware
│   ├── models/           # Database models
│   ├── repository/       # Data access layer
│   ├── service/         # Business logic layer
│   └── utils/           # Utility functions
├── pkg/                  # Library code that could be used by other projects
├── config/              # Configuration files
├── migrations/          # Database migration files
├── docs/               # Documentation files
└── test/               # Additional test files
```

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker (optional)

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your environment variables
3. Run `go mod download` to install dependencies
4. Run database migrations
5. Start the server with `go run cmd/main.go`

## Development

- Use `air` for hot reload during development
- Run tests with `go test ./...`
- Build with `go build -o bin/app cmd/main.go`

## Docker

```bash
# Build the image
docker build -t echobackend .

# Run the container
docker run -p 8080:8080 echobackend
```

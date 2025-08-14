# Echo Backend API

A REST API for a blog/content management system built with Go, Echo framework, and PostgreSQL.

## Features

- User authentication & management
- Posts, comments, and tags
- User follows and post likes
- Workspaces and pages
- JWT authentication
- PostgreSQL database
- Docker support

## Quick Start

### Prerequisites
- Go 1.21+
- PostgreSQL 14+

### Setup

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

3. **Run:**
   ```bash
   go mod download
   go run cmd/main.go
   ```

   Server starts at `http://localhost:8080`

### Development

- **Hot reload:** `air`
- **Tests:** `go test ./...`
- **Build:** `make build`

### Docker

```bash
docker build -t echobackend .
docker run -p 8080:8080 echobackend
```

## API Documentation

For detailed API endpoints and examples, see [api_doc.md](api_doc.md).

## Project Structure

```
├── cmd/                    # Application entry point
├── internal/               # Private application code
│   ├── handler/            # HTTP handlers
│   ├── service/            # Business logic
│   ├── repository/         # Data access
│   ├── model/              # Database models
│   └── middleware/         # Custom middleware
├── pkg/                    # Shared packages
├── migrations/             # Database migrations
└── config/                 # Configuration
```

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit changes: `git commit -am 'Add feature'`
4. Push to branch: `git push origin feature/my-feature`
5. Submit a pull request

Please follow existing code style and include tests for new features.

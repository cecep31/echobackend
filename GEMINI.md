# Gemini Project: Echo Backend API

## Project Overview

This is a REST API for a blog/content management system built with Go, the Echo framework, and PostgreSQL. It features user authentication, posts, comments, tags, user follows, post likes, workspaces, and pages.

**Key Technologies:**

*   **Backend:** Go, Echo
*   **Database:** PostgreSQL
*   **ORM:** GORM
*   **Authentication:** JWT
*   **Dependency Injection:** go.uber.org/dig
*   **Validation:** go-playground/validator
*   **Logging:** zerolog

**Architecture:**

The project follows a standard layered architecture:

*   `cmd/main.go`: Application entry point, responsible for initialization and startup.
*   `internal/`: Internal application logic, separated into:
    *   `handler`: HTTP handlers that receive requests and send responses.
    *   `service`: Business logic.
    *   `repository`: Data access layer that interacts with the database.
    *   `model`: GORM models that represent database tables.
    *   `middleware`: Custom middleware for tasks like authentication.
*   `pkg/`: Shared packages that can be used by other applications.
*   `config/`: Configuration management.
*   `migrations/`: Database migrations.
*   `routes/`: Route definitions.

## Building and Running

### Prerequisites

*   Go 1.21+
*   PostgreSQL 14+
*   Docker (optional)

### Setup

1.  **Clone the repository:**
    ```bash
    git clone <your-repo-url>
    cd echobackend
    ```

2.  **Configure environment variables:**
    Copy the `.env.example` file to `.env` and update it with your PostgreSQL database credentials.
    ```bash
    cp .env.example .env
    ```

### Commands

The following commands are available in the `Makefile`:

*   **Run the application:**
    ```bash
    make run
    ```
    This will start the server at `http://localhost:8080`.

*   **Run with hot-reloading:**
    ```bash
    make dev
    ```

*   **Build the application:**
    ```bash
    make build
    ```
    This will create a binary in the `bin/` directory.

*   **Run tests:**
    ```bash
    make test
    ```

*   **Run tests with coverage:**
    ```bash
    make test-coverage
    ```

*   **Lint the code:**
    ```bash
    make lint
    ```

*   **Format the code:**
    ```bash
    make fmt
    ```

### Docker

*   **Build the Docker image:**
    ```bash
    make docker-build
    ```

*   **Run the Docker container:**
    ```bash
    make docker-run
    ```

## Development Conventions

*   **Code Style:** Follow the standard Go formatting guidelines. Use `make fmt` to format the code.
*   **Testing:** Include tests for new features. Use `make test` to run all tests.
*   **Commits:** Follow conventional commit message standards.
*   **Dependencies:** Manage dependencies using Go modules. Use `go mod tidy` to clean up unused dependencies.
*   **API Documentation:** The API is documented in `api_doc.md`.

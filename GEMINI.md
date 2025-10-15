# GEMINI Project Analysis

## Project Overview

This project is a backend API for a blog/content management system, built with the Go programming language and the Echo web framework. It uses a PostgreSQL database for data storage and JWT for user authentication. The API provides a comprehensive set of endpoints for managing users, posts, comments, tags, and workspaces.

The project is well-structured, with a clear separation of concerns between handlers, services, repositories, and models. It also includes support for Docker, making it easy to build and deploy the application in a containerized environment.

## Building and Running

### Prerequisites

*   Go 1.24+
*   PostgreSQL 16+
*   Docker (optional)

### Key Commands

*   **Run the application:**
    ```bash
    go run cmd/main.go
    ```
    The server will start on `http://localhost:8080`.

*   **Run with hot-reloading (requires `air`):**
    ```bash
    air
    ```

*   **Run tests:**
    ```bash
    go test ./...
    ```

*   **Build the application:**
    ```bash
    make build
    ```

*   **Run with Docker:**
    ```bash
    docker build -t echobackend .
    docker run -p 8080:8080 echobackend
    ```

## Development Conventions

*   **Code Style:** The codebase follows standard Go conventions.
*   **Testing:** The project includes a `test` directory, and the `go test ./...` command is used to run tests. New features should include corresponding tests.
*   **API Documentation:** The `api_doc.md` file provides detailed documentation for all API endpoints. This file should be kept up-to-date as the API evolves.
*   **Dependency Management:** The project uses Go modules for dependency management. The `go.mod` file lists all the project's dependencies.
*   **Configuration:** The application is configured using environment variables. A `.env.example` file is provided to show the required variables.

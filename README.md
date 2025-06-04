# Echo Backend API

Echo Backend API is a RESTful API backend service built with the Echo v4 framework and PostgreSQL database. It provides a robust foundation for building scalable and maintainable web applications.

Project structure:

```
.
├── cmd/                    # Main applications for this project
│   └── main.go            # Application entry point
├── internal/              # Private application and library code
│   ├── handler/          # API handlers
│   ├── middleware/       # Custom middleware
│   ├── model/            # Database models
│   ├── repository/       # Data access layer
│   ├── service/          # Business logic layer
│   └── utils/            # Utility functions
├── pkg/                  # Library code that could be used by other projects
├── config/               # Configuration files
├── migrations/           # Database migration files
├── docs/                 # Documentation files
└── test/                 # Additional test files
```

Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Docker (optional)

Set up the Echo Backend API by following these steps:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/echobackend.git
   cd echobackend
   ```

2. **Configure environment variables:**
   Copy the `.env.example` file to `.env` and update the environment variables with your configuration:
   ```bash
   cp .env.example .env
   ```

3. **Install dependencies:**
   ```bash
   go mod download
   ```

4. **Set up the database:**
   - Create a PostgreSQL database.
   - Update the `DATABASE_URL` in the `.env` file with your database connection string.
   - Run the database migrations:
     ```bash
     go run cmd/migrate.go
     ```

5. **Start the server:**
   ```bash
   go run cmd/main.go
   ```

Enhance your development experience with these tools and commands:

- **Hot Reload:**
  Use `air` for hot reload during development:
  ```bash
  air
  ```

- **Run Tests:**
  Execute tests with:
  ```bash
  go test ./...
  ```

- **Build the Application:**
  Build the application with:
  ```bash
  go build -o bin/app cmd/main.go
  ```

Containerize the Echo Backend API with Docker.

```bash
# Build the Docker image
docker build -t echobackend .

# Run the Docker container
docker run -p 8080:8080 echobackend
```

The API provides the following endpoints:

### Authentication

- **POST /api/auth/login**
  - Logs in a user.
  - Request Body:
    ```json
    {
      "username": "string",
      "password": "string"
    }
    ```
  - Response:
    ```json
    {
      "token": "string"
    }
    ```

- **POST /api/auth/register**
  - Registers a new user.
  - Request Body:
    ```json
    {
      "username": "string",
      "password": "string",
      "email": "string"
    }
    ```
  - Response:
    ```json
    {
      "message": "User registered successfully"
    }
    ```

### Users

- **GET /api/users**
  - Retrieves a list of users.
  - Response:
    ```json
    [
      {
        "id": "string",
        "username": "string",
        "email": "string"
      }
    ]
    ```

- **GET /api/users/{id}**
  - Retrieves a user by ID.
  - Response:
    ```json
    {
      "id": "string",
      "username": "string",
      "email": "string"
    }
    ```

### Posts

- **GET /api/posts**
  - Retrieves a list of posts.
  - Response:
    ```json
    [
      {
        "id": "string",
        "title": "string",
        "content": "string",
        "author": "string"
      }
    ]
    ```

- **POST /api/posts**
  - Creates a new post.
  - Request Body:
    ```json
    {
      "title": "string",
      "content": "string",
      "author": "string"
    }
    ```
  - Response:
    ```json
    {
      "id": "string",
      "title": "string",
      "content": "string",
      "author": "string"
    }
    ```

### Tags

- **GET /api/tags**
  - Retrieves a list of tags.
  - Response:
    ```json
    [
      {
        "id": "string",
        "name": "string"
      }
    ]
    ```

- **POST /api/tags**
  - Creates a new tag.
  - Request Body:
    ```json
    {
      "name": "string"
    }
    ```
  - Response:
    ```json
    {
      "id": "string",
      "name": "string"
    }
    ```

Example requests to the API:

### Login

```bash
curl -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d '{"username": "user", "password": "pass"}'
```

### Register

```bash
curl -X POST http://localhost:8080/api/auth/register -H "Content-Type: application/json" -d '{"username": "newuser", "password": "newpass", "email": "newuser@example.com"}'
```

### Get Users

```bash
curl -X GET http://localhost:8080/api/users
```

### Create Post

```bash
curl -X POST http://localhost:8080/api/posts -H "Content-Type: application/json" -d '{"title": "New Post", "content": "This is a new post.", "author": "user"}'
```

Contribute to the Echo Backend API by following these steps:

1. **Fork the repository:**
   - Click the "Fork" button at the top right of the repository page.

2. **Clone your fork:**
   ```bash
   git clone https://github.com/yourusername/echobackend.git
   cd echobackend
   ```

3. **Create a new branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Make your changes:**
   - Implement your feature or bug fix.

5. **Commit your changes:**
   ```bash
   git commit -m "Add your commit message"
   ```

6. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a pull request:**
   - Go to the original repository and click the "New Pull Request" button.
   - Select your fork and the branch you pushed to.
   - Provide a clear description of your changes and submit the pull request.

Follow the existing code style and conventions. Write clear and concise commit messages. Include tests for new features and bug fixes.

Open an issue on the GitHub repository if you find a bug or have a feature request. Provide as much detail as possible, including steps to reproduce the issue or a description of the feature you'd like to see.

By contributing to the Echo Backend API, you agree that your contributions will be licensed under the MIT License.

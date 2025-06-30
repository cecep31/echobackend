# Go Echo Backend Makefile

# Variables
BINARY_NAME=echobackend
MAIN_PATH=./cmd/main.go
TEST_TIMEOUT=30s
COVERAGE_DIR=coverage

# Build commands
.PHONY: build
build:
	go build -o bin/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: run
run:
	go run $(MAIN_PATH)

.PHONY: clean
clean:
	go clean
	rm -rf bin/
	rm -rf $(COVERAGE_DIR)/

# Test commands
.PHONY: test
test:
	go test -timeout $(TEST_TIMEOUT) -v ./...

.PHONY: test-short
test-short:
	go test -timeout $(TEST_TIMEOUT) -short -v ./...

.PHONY: test-race
test-race:
	go test -timeout $(TEST_TIMEOUT) -race -v ./...

.PHONY: test-coverage
test-coverage:
	mkdir -p $(COVERAGE_DIR)
	go test -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated: $(COVERAGE_DIR)/coverage.html"

.PHONY: test-coverage-func
test-coverage-func:
	mkdir -p $(COVERAGE_DIR)
	go test -timeout $(TEST_TIMEOUT) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	go tool cover -func=$(COVERAGE_DIR)/coverage.out

# Benchmarks
.PHONY: bench
bench:
	go test -bench=. -benchmem ./...

# Code quality
.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: mod-download
mod-download:
	go mod download

# Security
.PHONY: security
security:
	gosec ./...

# All quality checks
.PHONY: check
check: fmt vet test-race test-coverage-func

# Development
.PHONY: dev
dev:
	air

# Docker commands
.PHONY: docker-build
docker-build:
	docker build -t $(BINARY_NAME) .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 $(BINARY_NAME)

# Database commands (if you add migrations later)
.PHONY: migrate-up
migrate-up:
	@echo "Add your migration up command here"

.PHONY: migrate-down
migrate-down:
	@echo "Add your migration down command here"

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application"
	@echo "  clean           - Clean build artifacts"
	@echo "  test            - Run all tests"
	@echo "  test-short      - Run short tests"
	@echo "  test-race       - Run tests with race detection"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  bench           - Run benchmarks"
	@echo "  fmt             - Format code"
	@echo "  vet             - Run go vet"
	@echo "  lint            - Run golangci-lint"
	@echo "  mod-tidy        - Tidy go modules"
	@echo "  check           - Run all quality checks"
	@echo "  dev             - Run with air for hot reload"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  help            - Show this help message"
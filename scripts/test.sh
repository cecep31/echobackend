#!/bin/bash

# Test script for Echo Backend
set -e

echo "🧪 Running Echo Backend Tests"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}✓${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

# Check if Go is installed
if ! command -v go &> /dev/null; then
    print_error "Go is not installed or not in PATH"
    exit 1
fi

print_status "Go version: $(go version)"

# Download dependencies
print_status "Downloading dependencies..."
go mod download

# Install test dependencies
go mod tidy

# Run tests with different configurations
echo ""
echo "🔍 Running Unit Tests..."

# Basic test run
print_status "Running all tests..."
go test -v ./...

echo ""
echo "🏃 Running tests with race detection..."
go test -race -v ./...

echo ""
echo "📊 Running tests with coverage..."
mkdir -p coverage
go test -coverprofile=coverage/coverage.out ./...
go tool cover -func=coverage/coverage.out

echo ""
echo "📈 Generating HTML coverage report..."
go tool cover -html=coverage/coverage.out -o coverage/coverage.html
print_status "Coverage report generated: coverage/coverage.html"

echo ""
echo "🚀 Running benchmarks..."
go test -bench=. -benchmem ./...

echo ""
echo "🔍 Running go vet..."
go vet ./...

echo ""
echo "🎯 Running gofmt check..."
if [ $(gofmt -l . | wc -l) -gt 0 ]; then
    print_warning "Code formatting issues found:"
    gofmt -l .
    print_warning "Run 'gofmt -w .' to fix formatting"
else
    print_status "Code formatting is correct"
fi

echo ""
print_status "All tests completed successfully!"
print_status "Coverage report available at: coverage/coverage.html"
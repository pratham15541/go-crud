.PHONY: build run dev test test-unit test-integration test-coverage clean migrate-up migrate-down docker-build docker-run help

# Variables
APP_NAME=go-crud
BIN_DIR=bin
DOCKER_IMAGE=go-crud:latest
GO_VERSION=1.21

# Default target
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  dev            - Run in development mode with hot reload"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  clean          - Clean build artifacts"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  lint           - Run golangci-lint"
	@echo "  fmt            - Format code"
	@echo "  deps           - Download dependencies"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/server cmd/server/main.go
	@echo "Build complete: $(BIN_DIR)/server"

# Run the application
run: build
	@echo "Starting $(APP_NAME)..."
	@./$(BIN_DIR)/server

# Development mode with hot reload
dev:
	@echo "Starting development server..."
	@go run cmd/server/main.go

# Run all tests
test:
	@echo "Running all tests..."
	@go test -v ./...

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v ./tests/unit/...

# Run integration tests only
test-integration:
	@echo "Running integration tests..."
	@go test -v ./tests/integration/...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	@./scripts/migrate.sh up

migrate-down:
	@echo "Rolling back database migrations..."
	@./scripts/migrate.sh down

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE)

# Code quality
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy

# Install tools
tools:
	@echo "Installing development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/server/main.go -o docs/swagger

# Run development environment
dev-env:
	@echo "Starting development environment..."
	@docker-compose -f docker/docker-compose.dev.yml up -d

# Stop development environment
dev-env-down:
	@echo "Stopping development environment..."
	@docker-compose -f docker/docker-compose.dev.yml down
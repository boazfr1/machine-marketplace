# Makefile for machine-marketplace server

# Variables
BINARY_NAME=machine-marketplace
BUILD_DIR=build
GO_FILES=$(shell find . -name "*.go" -type f)
MAIN_FILE=main.go

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Build flags
LDFLAGS=-ldflags "-s -w"
BUILD_FLAGS=-trimpath

.PHONY: all build clean test deps fmt vet run-auth run-order run-process help

# Default target
all: clean deps fmt vet test build

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for production (optimized)
build-prod:
	@echo "Building $(BINARY_NAME) for production..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)
	@echo "Production build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Formatting complete"

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...
	@echo "Vet complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Tests complete"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Quick build and run auth service
run-auth: build
	@echo "Starting auth service..."
	./$(BUILD_DIR)/$(BINARY_NAME) auth-service --db-host=localhost --db-user=boazfrid --db-password=""

# Quick build and run order service
run-order: build
	@echo "Starting order service..."
	./$(BUILD_DIR)/$(BINARY_NAME) order-service --db-host=localhost --db-user=boazfrid --db-password=""

# Quick build and run process service
run-process: build
	@echo "Starting process service..."
	./$(BUILD_DIR)/$(BINARY_NAME) process-service --db-host=localhost --db-user=boazfrid --db-password=""

# Development mode - build and run with auto-restart (requires air)
dev:
	@if command -v air > /dev/null; then \
		echo "Starting development server with air..."; \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to regular build and run..."; \
		make run-auth; \
	fi

# Install development tools
install-tools:
	@echo "Installing development tools..."
	$(GOGET) github.com/cosmtrek/air@latest
	@echo "Development tools installed"

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t machine-marketplace-server .
	@echo "Docker build complete"

# Docker run
docker-run:
	@echo "Running Docker container..."
	docker run -p 8080:8080 machine-marketplace-server

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Clean, deps, fmt, vet, test, and build"
	@echo "  build        - Build the binary"
	@echo "  build-prod   - Build optimized binary for production"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  fmt          - Format Go code"
	@echo "  vet          - Run go vet"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  run-auth     - Build and run auth service (with local DB config)"
	@echo "  run-order    - Build and run order service (with local DB config)"
	@echo "  run-process  - Build and run process service (with local DB config)"
	@echo "  run-auth-custom   - Run auth service with custom DB settings"
	@echo "  run-order-custom  - Run order service with custom DB settings"
	@echo "  run-process-custom- Run process service with custom DB settings"
	@echo "  dev          - Start development server with auto-restart"
	@echo "  install-tools- Install development tools (air)"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  help         - Show this help message"

# Local database configuration
DB_HOST=localhost
DB_USER=boazfrid
DB_PASSWORD=
DB_NAME=machine_market
DB_FLAGS=--db-host=$(DB_HOST) --db-user=$(DB_USER) --db-password=$(DB_PASSWORD) --db-name=$(DB_NAME)

# Run with custom database settings
run-auth-custom: build
	@echo "Starting auth service with custom DB settings..."
	./$(BUILD_DIR)/$(BINARY_NAME) auth-service $(DB_FLAGS)

run-order-custom: build
	@echo "Starting order service with custom DB settings..."
	./$(BUILD_DIR)/$(BINARY_NAME) order-service $(DB_FLAGS)

run-process-custom: build
	@echo "Starting process service with custom DB settings..."
	./$(BUILD_DIR)/$(BINARY_NAME) process-service $(DB_FLAGS)

# Quick shortcuts
auth: run-auth
order: run-order
process: run-process

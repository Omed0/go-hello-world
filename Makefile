# Go Task Management API - Makefile
# This file provides convenient commands for development

.PHONY: help build run test clean dev fmt vet deps sqlc

# Default target
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application binary"
	@echo "  make run      - Run the application"
	@echo "  make dev      - Run in development mode with hot reload"
	@echo "  make test     - Run all tests"
	@echo "  make fmt      - Format all Go files"
	@echo "  make vet      - Run go vet"
	@echo "  make deps     - Download dependencies"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make sqlc     - Generate database code (requires sqlc)"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server main.go

# Run the application
run:
	@echo "Starting server..."
	go run main.go

# Development mode (install air first: go install github.com/cosmtrek/air@latest)
dev:
	@echo "Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "Falling back to regular run..."; \
		go run main.go; \
	fi

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out

# Generate database code (requires sqlc to be installed)
sqlc:
	@echo "Generating database code..."
	@if command -v sqlc > /dev/null; then \
		sqlc generate; \
	else \
		echo "sqlc not installed. Install from https://docs.sqlc.dev/en/latest/overview/install.html"; \
	fi

# Lint code (requires golangci-lint)
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install from https://golangci-lint.run/usage/install/"; \
	fi

# Check all (format, vet, test)
check: fmt vet test
	@echo "All checks passed!"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	@echo "Tools installed. You may also want to install:"
	@echo "  - sqlc: https://docs.sqlc.dev/en/latest/overview/install.html"
	@echo "  - golangci-lint: https://golangci-lint.run/usage/install/"

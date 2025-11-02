# OpenSCAD HTTP API Server - Task Runner

# List all available recipes
default:
    @just --list

# Download dependencies
deps:
    go mod download
    go mod verify

# Generate swagger documentation
swagger:
    swag init

# Build the application
build: swagger
    #!/bin/bash
    set -e
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
    go build -ldflags "-X github.com/stevexciv/scad-server/version.commit=$COMMIT -X github.com/stevexciv/scad-server/version.tag=$TAG" -o bin/scad-server .

# Run tests
test:
    go test -v -cover ./...

# Run tests with coverage report
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out -o coverage.html
    @echo "Coverage report generated at coverage.html"

# Run the application
run: swagger
    #!/bin/bash
    set -e
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "unknown")
    go run -ldflags "-X github.com/stevexciv/scad-server/version.commit=$COMMIT -X github.com/stevexciv/scad-server/version.tag=$TAG" main.go

# Clean build artifacts
clean:
    rm -rf bin/
    rm -rf tmp/
    rm -f coverage.out coverage.html
    rm -f *.stl *.png *.svg *.pdf *.scad

# Build Docker image
docker-build:
    docker build -t stevexciv/scad-server:latest .

# Run Docker container
docker-run:
    docker run -p 8000:8000 stevexciv/scad-server:latest

# Run linter
lint:
    golangci-lint run ./... || true

# Format code
fmt:
    go fmt ./...

# Run go vet
vet:
    go vet ./...

# Run all checks and build
all: deps fmt vet swagger test build

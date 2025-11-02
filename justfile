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
    go build -o bin/scad-server .

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
    go run main.go

# Clean build artifacts
clean:
    rm -rf bin/
    rm -rf tmp/
    rm -f coverage.out coverage.html
    rm -f *.stl *.png *.svg *.pdf *.scad

# Build Docker image
docker-build:
    docker build -t scad-server:latest .

# Run Docker container
docker-run:
    docker run -p 8000:8000 scad-server:latest

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

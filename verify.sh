#!/bin/bash

# Verification script for scad-server project

set -e

echo "========================================"
echo "OpenSCAD HTTP API Server - Verification"
echo "========================================"
echo ""

# Check Go version
echo "1. Checking Go version..."
go version
echo ""

# Check dependencies
echo "2. Downloading dependencies..."
go mod download
go mod verify
echo "✓ Dependencies verified"
echo ""

# Format check
echo "3. Checking code formatting..."
if [ -n "$(gofmt -l .)" ]; then
    echo "✗ Code is not formatted. Run: just fmt"
    exit 1
fi
echo "✓ Code is formatted"
echo ""

# Vet check
echo "4. Running go vet..."
go vet ./...
echo "✓ go vet passed"
echo ""

# Generate Swagger docs
echo "5. Generating Swagger documentation..."
if command -v swag &> /dev/null; then
    swag init
    echo "✓ Swagger docs generated"
else
    echo "⚠ swag not found. Install with: go install github.com/swaggo/swag/cmd/swag@latest"
fi
echo ""

# Run tests
echo "6. Running tests..."
go test -v -race -cover ./...
echo "✓ All tests passed"
echo ""

# Build application
echo "7. Building application..."
go build -o bin/scad-server .
echo "✓ Build successful"
echo ""

# Check binary
echo "8. Checking binary..."
if [ -f "bin/scad-server" ]; then
    ls -lh bin/scad-server
    echo "✓ Binary created successfully"
else
    echo "✗ Binary not found"
    exit 1
fi
echo ""

echo "========================================"
echo "✓ All verification checks passed!"
echo "========================================"
echo ""
echo "Next steps:"
echo "  - Run the server: just run"
echo "  - Build: just build"
echo "  - Build Docker: just docker-build"
echo "  - View docs: http://localhost:8080/swagger/index.html"
echo ""

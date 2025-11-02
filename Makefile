.PHONY: help build test clean run docker-build docker-run swagger deps

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

deps: ## Download dependencies
	go mod download
	go mod verify

swagger: ## Generate swagger documentation
	swag init

build: swagger ## Build the application
	go build -o bin/scad-server .

test: ## Run tests
	go test -v -cover ./...

test-coverage: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

run: swagger ## Run the application
	go run main.go

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	rm -f *.stl *.png *.svg *.pdf *.scad

docker-build: ## Build Docker image
	docker build -t scad-server:latest .

docker-run: ## Run Docker container
	docker run -p 8080:8080 scad-server:latest

lint: ## Run linter
	golangci-lint run ./... || true

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

all: deps fmt vet swagger test build ## Run all checks and build

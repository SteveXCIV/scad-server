# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for version info
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install swag for generating Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code (including .git for version info)
COPY . .

# Generate Swagger documentation
RUN swag init

# Build the application with Git metadata
RUN set -e && \
    COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    TAG=$(git describe --tags --abbrev=0 2>/dev/null || echo "unknown") && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
      -ldflags "-X github.com/stevexciv/scad-server/version.commit=$COMMIT -X github.com/stevexciv/scad-server/version.tag=$TAG" \
      -o scad-server .

# Final stage
FROM openscad/openscad:trixie

# Install ca-certificates and wget for health checks
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates wget && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/scad-server .

# Expose port
EXPOSE 8000

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8000

# Run the application
CMD ["./scad-server"]

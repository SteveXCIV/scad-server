# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install swag for generating Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate Swagger documentation
RUN swag init

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scad-server .

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

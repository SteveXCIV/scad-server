# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o scad-server .

# Final stage
FROM openscad/openscad:latest

# Install ca-certificates for HTTPS
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/scad-server .

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release
ENV PORT=8080

# Run the application
CMD ["./scad-server"]

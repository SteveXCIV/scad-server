# OpenSCAD HTTP API Server

An HTTP API that provides access to core OpenSCAD functionality, specifically file export and summarization. The API accepts SCAD file content via JSON payloads and returns processed results.

I decided to build this because I wanted more-or-less a "CI" server for generating OpenSCAD models on my home server.
Hopefully others will find this useful as well.

## Features

- **Export to Multiple Formats**: PNG, STL (binary + ASCII), SVG, PDF, and 3MF
- **Summary Generation**: Get diagnostics about SCAD models
- **Format-Specific Options**: Supports a subset of format-specific parameters from the OpenSCAD CLI
- **OpenAPI Documentation**: Interactive API docs
- **Docker/OCI Support**: Fits into existing homelabs already using k8s/compose/etc 

## API Endpoints

### Base URL
```
http://localhost:8000/openscad/v1
```

### Endpoints

#### 1. Export SCAD to Various Formats
```
POST /openscad/v1/export
```

Exports OpenSCAD content to PNG, STL (binary/ASCII), SVG, PDF, or 3MF format.

**Supported Formats:**
- `png` - Image export for visualization
- `stl_binary` - Binary STL (useful for older slicer software)
- `stl_ascii` - ASCII STL (useful for older slicer softwate)
- `svg` - Vector graphics
- `pdf` - Document export
- `3mf` - 3D Manufacturing Format (good option for more modern slicers)

#### 2. Generate Summary Information
```
POST /openscad/v1/summary
```

Generates summary information for OpenSCAD content.

**Summary Types:**
- `all` - All available summary information (default)
- `cache` - Cache statistics
- `time` - Timing information
- `camera` - Camera position
- `geometry` - Geometry statistics
- `bounding-box` - Bounding box dimensions
- `area` - Surface area

#### 3. Health Check
```
GET /health
```

Returns the health status of the API.

## Installation

### Prerequisites
- Go 1.23 or later
- OpenSCAD on your PATH (for local development)
- Docker (for containerized deployment)
- just (task runner):  https://github.com/casey/just (optional but **very** helpful)

### Local Development

1. Clone the repository:
```bash
git clone https://github.com/stevexciv/scad-server.git
cd scad-server
```

2. Install dependencies:
```bash
go mod download
```

3. Generate Swagger documentation:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init
```

4. Run the server:
```bash
just run
```

The server will start on `http://localhost:8000`.

### Docker Deployment

This project uses the OpenSCAD Trixie development build, which provides built-in EGL support for headless rendering and supports both AMD64 and ARM64 architectures.

1. Build the Docker image:
```bash
just docker-build
```

2. Run the container:
```bash
just docker-run
```

For a specific platform:
```bash
docker build --platform linux/amd64 -t scad-server:latest .
docker run --platform linux/amd64 -p 8000:8000 scad-server:latest
```

## Usage Examples

### Export to PNG

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "cube([10,10,10]);",
    "format": "png",
    "options": {
      "png": {
        "width": 800,
        "height": 600
      }
    }
  }' \
  --output cube.png
```

### Export to STL (Binary)

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "sphere(r=5);",
    "format": "stl_binary",
    "options": {
      "stl": {
        "decimal_precision": 6
      }
    }
  }' \
  --output sphere.stl
```

### Export to SVG

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "circle(r=10);",
    "format": "svg",
    "options": {
      "svg": {
        "fill": true,
        "fill_color": "blue",
        "stroke": true,
        "stroke_color": "black",
        "stroke_width": 0.5
      }
    }
  }' \
  --output circle.svg
```

### Export to PDF

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "square([20,20]);",
    "format": "pdf",
    "options": {
      "pdf": {
        "paper_size": "a4",
        "orientation": "landscape",
        "show_grid": true,
        "grid_size": 10
      }
    }
  }' \
  --output square.pdf
```

### Generate Summary

```bash
curl -X POST http://localhost:8000/openscad/v1/summary \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "cube([10,10,10]);",
    "summary_type": "geometry"
  }'
```

## API Documentation

### Interactive/Human-readable

Once the server is running, you can access the interactive Swagger documentation at:

```
http://localhost:8000/swagger/index.html
```

### Spec/Machine-readable

Additionally, the JSON API spec can be found at:

```
http://localhost:8000/swagger/doc.json
```

Tools able to dynamically generate API clients can use this spec file to generate a client for scad-server.

## Configuration

The server can be configured using environment variables:

- `SCADSRV_PORT` - Server port (default: 8000)
- `SCADSRV_GIN_MODE` - Gin framework mode: `debug`, `release`, or `test` (default: release)

Example:
```bash
SCADSRV_PORT=3000 SCADSRV_GIN_MODE=debug just run
```

## Testing

Run all tests:
```bash
just test
```

Run tests with coverage report:
```bash
just test-coverage
```

## Available Commands

List all available tasks:
```bash
just --list
```

Common commands:
- `just build` - Build the application
- `just test` - Run tests
- `just run` - Run the server
- `just clean` - Clean build artifacts
- `just docker-build` - Build Docker image
- `just docker-run` - Run Docker container

## Project Structure

```
.
├── main.go                 # Application entry point
├── models/                 # Data models
│   └── models.go
├── handlers/               # HTTP handlers
│   ├── handlers.go
│   └── handlers_test.go
├── services/               # Business logic
│   ├── openscad.go
│   └── openscad_test.go
├── docs/                   # Swagger documentation (generated)
├── Dockerfile              # Docker configuration
├── justfile                # Task runner configuration
├── .gitignore              # Git ignore file
├── .dockerignore           # Docker ignore file
├── go.mod                  # Go module file
└── README.md               # This file
```

## Security Considerations

- Input validation for SCAD content
- Request size limits to prevent resource exhaustion
- Timeout mechanisms for long-running renders (default: 5 minutes)

## Performance

- OpenSCAD processes run in isolated temporary directories
- Automatic cleanup of temporary files
- Configurable timeout for processing
- Resource limits enforced by Docker container

## License

GPL-3.0 License - See COPYING file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues and questions, please open an issue on the GitHub repository:
https://github.com/stevexciv/scad-server/issues


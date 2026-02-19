# API Documentation

## Overview

The OpenSCAD HTTP API provides RESTful endpoints for converting OpenSCAD content to various formats and generating summary information.

## Base URL

```
http://localhost:8000
```

## Authentication

Currently, no authentication is required. This should be added in production environments.

## Endpoints

### 1. Health Check

Check if the API is running and healthy.

**Endpoint:** `GET /health`

**Response:**
```json
{
  "status": "ok"
}
```

**Status Codes:**
- `200 OK` - Service is healthy

---

### 2. Export SCAD Content

Export OpenSCAD content to various formats.

**Endpoint:** `POST /openscad/v1/export`

**Content-Type:** `application/json`

**Request Body:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| scad_content | string | Yes | The OpenSCAD code to export |
| format | string | Yes | Output format: `png`, `stl_binary`, `stl_ascii`, `svg`, `pdf`, `3mf`, `webp`, `avif` |
| options | object | No | Format-specific options (see below) |

#### Format-Specific Options

##### PNG Options (`options.png`)

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| width | integer | No | 800 | Image width in pixels |
| height | integer | No | 600 | Image height in pixels |

> **Note:** WebP and AVIF formats reuse `options.png` for dimension customization. OpenSCAD renders to PNG first, then the server converts to the requested format.

##### STL Options (`options.stl`)

| Field | Type | Required | Default | Range | Description |
|-------|------|----------|---------|-------|-------------|
| decimal_precision | integer | No | 6 | 1-16 | Decimal precision for coordinates |

##### SVG Options (`options.svg`)

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| fill | boolean | No | false | Fill shapes |
| fill_color | string | No | "white" | Fill color |
| stroke | boolean | No | true | Draw strokes |
| stroke_color | string | No | "black" | Stroke color |
| stroke_width | float | No | 0.35 | Stroke width |

##### PDF Options (`options.pdf`)

| Field | Type | Required | Default | Valid Values | Description |
|-------|------|----------|---------|--------------|-------------|
| paper_size | string | No | "a4" | a6, a5, a4, a3, letter, legal, tabloid | Paper size |
| orientation | string | No | "portrait" | portrait, landscape, auto | Page orientation |
| show_grid | boolean | No | false | - | Show grid |
| grid_size | float | No | 10 | - | Grid size |
| fill | boolean | No | false | - | Fill shapes |
| fill_color | string | No | "black" | - | Fill color |
| stroke | boolean | No | true | - | Draw strokes |
| stroke_color | string | No | "black" | - | Stroke color |
| stroke_width | float | No | 0.35 | - | Stroke width |

##### 3MF Options (`options.3mf`)

| Field | Type | Required | Default | Valid Values | Description |
|-------|------|----------|---------|--------------|-------------|
| unit | string | No | "millimeter" | micron, millimeter, centimeter, meter, inch, foot | Unit of measurement |
| decimal_precision | integer | No | 6 | 1-16 | Decimal precision for coordinates |
| color | string | No | "#f9d72c" | - | Color in hex format |
| color_mode | string | No | "model" | model, none, selected-only | Color mode |
| material_type | string | No | "color" | color, basematerial | Material type |
| add_metadata | boolean | No | true | - | Include metadata in file |
| metadata_title | string | No | "" | - | Model title metadata |
| metadata_designer | string | No | "" | - | Designer metadata |
| metadata_description | string | No | "" | - | Description metadata |
| metadata_copyright | string | No | "" | - | Copyright metadata |

**Example Request - PNG:**
```json
{
  "scad_content": "cube([10,10,10]);",
  "format": "png",
  "options": {
    "png": {
      "width": 800,
      "height": 600
    }
  }
}
```

**Example Request - STL:**
```json
{
  "scad_content": "sphere(r=5);",
  "format": "stl_binary",
  "options": {
    "stl": {
      "decimal_precision": 6
    }
  }
}
```

**Response:**
- Binary data in the requested format

**Status Codes:**
- `200 OK` - Export successful, returns binary data
- `400 Bad Request` - Invalid request parameters
- `500 Internal Server Error` - Export failed

**Error Response:**
```json
{
  "error": "export failed",
  "message": "detailed error message"
}
```

---

### 3. Generate Summary

Generate summary information for OpenSCAD content.

**Endpoint:** `POST /openscad/v1/summary`

**Content-Type:** `application/json`

**Request Body:**

| Field | Type | Required | Default | Description |
|-------|------|----------|---------|-------------|
| scad_content | string | Yes | - | The OpenSCAD code to analyze |
| summary_type | string | No | "all" | Type of summary: `all`, `cache`, `time`, `camera`, `geometry`, `bounding-box`, `area` |

**Example Request:**
```json
{
  "scad_content": "cube([10,10,10]);",
  "summary_type": "all"
}
```

**Response:**
```json
{
  "summary": {
    "cache": { ... },
    "time": { ... },
    "camera": { ... },
    "geometry": { ... },
    "bounding_box": { ... },
    "area": { ... }
  }
}
```

**Status Codes:**
- `200 OK` - Summary generated successfully
- `400 Bad Request` - Invalid request parameters
- `500 Internal Server Error` - Summary generation failed

**Error Response:**
```json
{
  "error": "summary generation failed",
  "message": "detailed error message"
}
```

---

## Complete Examples

### Export a Cube to PNG

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "cube([10,10,10]);",
    "format": "png",
    "options": {
      "png": {
        "width": 1024,
        "height": 768
      }
    }
  }' \
  --output cube.png
```

### Export a Sphere to Binary STL

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "sphere(r=5);",
    "format": "stl_binary"
  }' \
  --output sphere.stl
```

### Export a Circle to SVG with Custom Colors

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
        "stroke_color": "red",
        "stroke_width": 1.0
      }
    }
  }' \
  --output circle.svg
```

### Export to PDF with Grid

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "square([20,20]);",
    "format": "pdf",
    "options": {
      "pdf": {
        "paper_size": "a3",
        "orientation": "landscape",
        "show_grid": true,
        "grid_size": 5
      }
    }
  }' \
  --output square.pdf
```

### Export to 3MF with Metadata

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "cube([10,10,10]);",
    "format": "3mf",
    "options": {
      "3mf": {
        "unit": "centimeter",
        "decimal_precision": 6,
        "color": "#ff0000",
        "add_metadata": true,
        "metadata_title": "Red Cube",
        "metadata_designer": "OpenSCAD",
        "metadata_description": "A simple red cube model"
      }
    }
  }' \
  --output cube.3mf
```

### Export a Cube to WebP

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "cube([10,10,10]);",
    "format": "webp",
    "options": {
      "png": {
        "width": 800,
        "height": 600
      }
    }
  }' \
  --output cube.webp
```

### Export a Cube to AVIF

```bash
curl -X POST http://localhost:8000/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "cube([10,10,10]);",
    "format": "avif",
    "options": {
      "png": {
        "width": 800,
        "height": 600
      }
    }
  }' \
  --output cube.avif
```

### Generate Complete Summary

```bash
curl -X POST http://localhost:8000/openscad/v1/summary \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "difference() { cube([20,20,20], center=true); sphere(r=12); }",
    "summary_type": "all"
  }'
```

### Generate Geometry Summary Only

```bash
curl -X POST http://localhost:8000/openscad/v1/summary \
  -H "Content-Type: application/json" \
  -d '{
    "scad_content": "cube([10,10,10]);",
    "summary_type": "geometry"
  }'
```

---

## Error Handling

All endpoints return appropriate HTTP status codes and JSON error responses when errors occur.

### Common Error Codes

- `400 Bad Request` - Invalid request parameters or SCAD syntax
- `500 Internal Server Error` - Processing failed (OpenSCAD error, timeout, etc.)

### Error Response Format

```json
{
  "error": "error_type",
  "message": "detailed error description"
}
```

---

## Rate Limiting

Currently, no rate limiting is implemented. Consider adding rate limiting in production environments.

---

## Timeouts

- Default processing timeout: 5 minutes
- Can be adjusted in the service configuration

---

## Content Type Headers

### Export Endpoint Response Content Types

| Format | Content-Type |
|--------|--------------|
| png | `image/png` |
| stl_binary | `application/octet-stream` |
| stl_ascii | `application/octet-stream` |
| svg | `image/svg+xml` |
| pdf | `application/pdf` |
| 3mf | `application/vnd.ms-package.3dmodel+xml` |
| webp | `image/webp` |
| avif | `image/avif` |

---

## OpenAPI/Swagger Documentation

Interactive API documentation is available at:

```
http://localhost:8000/swagger/index.html
```

This provides a complete, interactive interface to explore and test all API endpoints.

---

## Notes

1. **SCAD Content**: Must be valid OpenSCAD code
2. **File Size**: Large or complex SCAD files may take longer to process
3. **Temporary Files**: All temporary files are automatically cleaned up after processing
4. **Security**: Input is validated to prevent code injection
5. **Docker**: When running in Docker, OpenSCAD is pre-installed in the container

---

## Support

For issues, questions, or contributions, please visit:
https://github.com/stevexciv/scad-server

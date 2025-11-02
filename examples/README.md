# Example Requests

This directory contains example JSON payloads for testing the API.

## Usage

### Export to PNG
```bash
curl -X POST http://localhost:8080/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d @examples/export-png.json \
  --output cube.png
```

### Export to STL
```bash
curl -X POST http://localhost:8080/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d @examples/export-stl.json \
  --output sphere.stl
```

### Export to SVG
```bash
curl -X POST http://localhost:8080/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d @examples/export-svg.json \
  --output circle.svg
```

### Export to PDF
```bash
curl -X POST http://localhost:8080/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d @examples/export-pdf.json \
  --output square.pdf
```

### Generate Summary
```bash
curl -X POST http://localhost:8080/openscad/v1/summary \
  -H "Content-Type: application/json" \
  -d @examples/summary.json
```

## Complex Example

Create a more complex SCAD model:

```json
{
  "scad_content": "difference() { cube([20,20,20], center=true); sphere(r=12); }",
  "format": "png",
  "options": {
    "png": {
      "width": 1024,
      "height": 768
    }
  }
}
```

Save this to a file (e.g., `complex.json`) and run:

```bash
curl -X POST http://localhost:8080/openscad/v1/export \
  -H "Content-Type: application/json" \
  -d @complex.json \
  --output complex.png
```

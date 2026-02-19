package services

import (
	"testing"

	"github.com/stevexciv/scad-server/models"
)

func TestValidateFormat(t *testing.T) {
	service := NewOpenSCADService()

	tests := []struct {
		name    string
		format  string
		wantErr bool
	}{
		{"Valid PNG", "png", false},
		{"Valid STL Binary", "stl_binary", false},
		{"Valid STL ASCII", "stl_ascii", false},
		{"Valid SVG", "svg", false},
		{"Valid PDF", "pdf", false},
		{"Valid 3MF", "3mf", false},
		{"Valid WebP", "webp", false},
		{"Valid AVIF", "avif", false},
		{"Invalid format", "invalid", true},
		{"Empty format", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetOutputExtension(t *testing.T) {
	service := NewOpenSCADService()

	tests := []struct {
		name          string
		format        string
		wantExt       string
		wantExportFmt string
	}{
		{"PNG", "png", "png", ""},
		{"STL Binary", "stl_binary", "stl", "binstl"},
		{"STL ASCII", "stl_ascii", "stl", "asciistl"},
		{"SVG", "svg", "svg", ""},
		{"PDF", "pdf", "pdf", ""},
		{"3MF", "3mf", "3mf", ""},
		{"WebP", "webp", "png", ""},
		{"AVIF", "avif", "png", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext, exportFmt := service.getOutputExtension(tt.format)
			if ext != tt.wantExt {
				t.Errorf("getOutputExtension() ext = %v, want %v", ext, tt.wantExt)
			}
			if exportFmt != tt.wantExportFmt {
				t.Errorf("getOutputExtension() exportFmt = %v, want %v", exportFmt, tt.wantExportFmt)
			}
		})
	}
}

func TestGetContentType(t *testing.T) {
	service := NewOpenSCADService()

	tests := []struct {
		name   string
		format string
		want   string
	}{
		{"PNG", "png", "image/png"},
		{"STL Binary", "stl_binary", "application/octet-stream"},
		{"STL ASCII", "stl_ascii", "application/octet-stream"},
		{"SVG", "svg", "image/svg+xml"},
		{"PDF", "pdf", "application/pdf"},
		{"3MF", "3mf", "application/vnd.ms-package.3dmodel+xml"},
		{"WebP", "webp", "image/webp"},
		{"AVIF", "avif", "image/avif"},
		{"Unknown", "unknown", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.getContentType(tt.format)
			if got != tt.want {
				t.Errorf("getContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildExportOptions(t *testing.T) {
	service := NewOpenSCADService()

	t.Run("PNG options", func(t *testing.T) {
		width := 1024
		height := 768
		req := &models.ExportRequest{
			Format: "png",
			Options: models.ExportOptions{
				PNG: &models.PNGOptions{
					Width:  &width,
					Height: &height,
				},
			},
		}
		args := service.buildExportOptions(req)
		if len(args) != 2 {
			t.Errorf("Expected 2 args, got %d", len(args))
		}
		if args[0] != "--imgsize" {
			t.Errorf("Expected --imgsize, got %s", args[0])
		}
		if args[1] != "1024,768" {
			t.Errorf("Expected 1024,768, got %s", args[1])
		}
	})

	t.Run("SVG options", func(t *testing.T) {
		fill := true
		fillColor := "red"
		req := &models.ExportRequest{
			Format: "svg",
			Options: models.ExportOptions{
				SVG: &models.SVGOptions{
					Fill:      &fill,
					FillColor: &fillColor,
				},
			},
		}
		args := service.buildExportOptions(req)
		if len(args) != 4 {
			t.Errorf("Expected 4 args, got %d", len(args))
		}
	})

	t.Run("PDF options", func(t *testing.T) {
		paperSize := "a3"
		orientation := "landscape"
		showGrid := true
		req := &models.ExportRequest{
			Format: "pdf",
			Options: models.ExportOptions{
				PDF: &models.PDFOptions{
					PaperSize:   &paperSize,
					Orientation: &orientation,
					ShowGrid:    &showGrid,
				},
			},
		}
		args := service.buildExportOptions(req)
		if len(args) != 6 {
			t.Errorf("Expected 6 args, got %d", len(args))
		}
	})

	t.Run("WebP options (reuses PNG)", func(t *testing.T) {
		width := 800
		height := 600
		req := &models.ExportRequest{
			Format: "webp",
			Options: models.ExportOptions{
				PNG: &models.PNGOptions{
					Width:  &width,
					Height: &height,
				},
			},
		}
		args := service.buildExportOptions(req)
		if len(args) != 2 {
			t.Errorf("Expected 2 args, got %d", len(args))
		}
		if args[0] != "--imgsize" {
			t.Errorf("Expected --imgsize, got %s", args[0])
		}
		if args[1] != "800,600" {
			t.Errorf("Expected 800,600, got %s", args[1])
		}
	})

	t.Run("AVIF options (reuses PNG)", func(t *testing.T) {
		width := 640
		height := 480
		req := &models.ExportRequest{
			Format: "avif",
			Options: models.ExportOptions{
				PNG: &models.PNGOptions{
					Width:  &width,
					Height: &height,
				},
			},
		}
		args := service.buildExportOptions(req)
		if len(args) != 2 {
			t.Errorf("Expected 2 args, got %d", len(args))
		}
		if args[0] != "--imgsize" {
			t.Errorf("Expected --imgsize, got %s", args[0])
		}
		if args[1] != "640,480" {
			t.Errorf("Expected 640,480, got %s", args[1])
		}
	})

	t.Run("No options", func(t *testing.T) {
		req := &models.ExportRequest{
			Format: "png",
		}
		args := service.buildExportOptions(req)
		if len(args) != 0 {
			t.Errorf("Expected 0 args, got %d", len(args))
		}
	})

	t.Run("3MF options", func(t *testing.T) {
		unit := "centimeter"
		precision := 8
		color := "#ff0000"
		colorMode := "model"
		addMetadata := true
		title := "Test Model"
		req := &models.ExportRequest{
			Format: "3mf",
			Options: models.ExportOptions{
				ThreeMF: &models.ThreeMFOptions{
					Unit:             &unit,
					DecimalPrecision: &precision,
					Color:            &color,
					ColorMode:        &colorMode,
					AddMetadata:      &addMetadata,
					MetadataTitle:    &title,
				},
			},
		}
		args := service.buildExportOptions(req)
		if len(args) < 10 {
			t.Errorf("Expected at least 10 args for 3MF options, got %d", len(args))
		}
		// Check that 3MF specific options are present
		hasUnit := false
		hasPrecision := false
		for i := 0; i < len(args)-1; i++ {
			if args[i] == "-O" {
				if args[i+1] == "export-3mf/unit=centimeter" {
					hasUnit = true
				}
				if args[i+1] == "export-3mf/decimal-precision=8" {
					hasPrecision = true
				}
			}
		}
		if !hasUnit {
			t.Errorf("Expected unit option in args")
		}
		if !hasPrecision {
			t.Errorf("Expected precision option in args")
		}
	})
}

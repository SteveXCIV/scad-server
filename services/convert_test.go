package services

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
)

// createTestPNG creates a minimal valid PNG image for testing.
func createTestPNG(width, height int) ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
		}
	}
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func TestConvertPNGToWebP(t *testing.T) {
	t.Run("Valid PNG", func(t *testing.T) {
		pngData, err := createTestPNG(4, 4)
		if err != nil {
			t.Fatalf("Failed to create test PNG: %v", err)
		}

		webpData, err := convertPNGToWebP(pngData)
		if err != nil {
			t.Fatalf("convertPNGToWebP() error = %v", err)
		}

		if len(webpData) == 0 {
			t.Error("Expected non-empty WebP output")
		}

		// WebP files start with "RIFF" magic bytes
		if len(webpData) < 12 {
			t.Fatalf("WebP output too short: %d bytes", len(webpData))
		}
		if string(webpData[0:4]) != "RIFF" {
			t.Errorf("Expected RIFF magic bytes, got %q", webpData[0:4])
		}
		if string(webpData[8:12]) != "WEBP" {
			t.Errorf("Expected WEBP signature, got %q", webpData[8:12])
		}
	})

	t.Run("Invalid input", func(t *testing.T) {
		_, err := convertPNGToWebP([]byte("not a png"))
		if err == nil {
			t.Error("Expected error for invalid PNG input")
		}
	})
}

func TestConvertPNGToAVIF(t *testing.T) {
	t.Run("Valid PNG", func(t *testing.T) {
		pngData, err := createTestPNG(4, 4)
		if err != nil {
			t.Fatalf("Failed to create test PNG: %v", err)
		}

		avifData, err := convertPNGToAVIF(pngData)
		if err != nil {
			t.Fatalf("convertPNGToAVIF() error = %v", err)
		}

		if len(avifData) == 0 {
			t.Error("Expected non-empty AVIF output")
		}

		// AVIF files contain an "ftyp" box near the start
		if len(avifData) < 12 {
			t.Fatalf("AVIF output too short: %d bytes", len(avifData))
		}
		if string(avifData[4:8]) != "ftyp" {
			t.Errorf("Expected ftyp box, got %q", avifData[4:8])
		}
	})

	t.Run("Invalid input", func(t *testing.T) {
		_, err := convertPNGToAVIF([]byte("not a png"))
		if err == nil {
			t.Error("Expected error for invalid PNG input")
		}
	})
}

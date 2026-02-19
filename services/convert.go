package services

import (
	"bytes"
	"fmt"
	"image/png"
	"log"

	avif "github.com/Kagami/go-avif"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

const (
	defaultWebPQuality float32 = 80
)

// convertPNGToWebP takes raw PNG bytes and returns WebP-encoded bytes.
func convertPNGToWebP(pngData []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, defaultWebPQuality)
	if err != nil {
		return nil, fmt.Errorf("failed to create WebP encoder options: %w", err)
	}

	var buf bytes.Buffer
	if err := webp.Encode(&buf, img, options); err != nil {
		return nil, fmt.Errorf("failed to encode WebP: %w", err)
	}

	log.Printf("[Convert] PNG (%d bytes) -> WebP (%d bytes)", len(pngData), buf.Len())
	return buf.Bytes(), nil
}

// convertPNGToAVIF takes raw PNG bytes and returns AVIF-encoded bytes.
func convertPNGToAVIF(pngData []byte) ([]byte, error) {
	img, err := png.Decode(bytes.NewReader(pngData))
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}

	var buf bytes.Buffer
	if err := avif.Encode(&buf, img, nil); err != nil {
		return nil, fmt.Errorf("failed to encode AVIF: %w", err)
	}

	log.Printf("[Convert] PNG (%d bytes) -> AVIF (%d bytes)", len(pngData), buf.Len())
	return buf.Bytes(), nil
}

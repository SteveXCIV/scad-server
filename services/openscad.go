package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/stevexciv/scad-server/models"
)

const (
	defaultTimeout = 5 * time.Minute
	openscadCmd    = "openscad"
)

// OpenSCADService provides OpenSCAD operations
type OpenSCADService struct {
	timeout time.Duration
}

// NewOpenSCADService creates a new OpenSCAD service
func NewOpenSCADService() *OpenSCADService {
	return &OpenSCADService{
		timeout: defaultTimeout,
	}
}

// Export exports SCAD content to the specified format
func (s *OpenSCADService) Export(req *models.ExportRequest) ([]byte, string, error) {
	// Validate format
	if err := s.validateFormat(req.Format); err != nil {
		return nil, "", err
	}

	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "scad-export-*")
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			// Log error but don't fail the operation
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp directory %s: %v\n", tmpDir, err)
		}
	}()

	// Write SCAD content to temporary file
	scadFile := filepath.Join(tmpDir, "input.scad")
	if err := os.WriteFile(scadFile, []byte(req.ScadContent), 0644); err != nil {
		return nil, "", fmt.Errorf("failed to write SCAD file: %w", err)
	}

	// Determine output file extension
	outputExt, exportFormat := s.getOutputExtension(req.Format)
	outputFile := filepath.Join(tmpDir, "output."+outputExt)

	// Build OpenSCAD command arguments
	args := []string{"-o", outputFile}

	// Add export format if needed
	if exportFormat != "" {
		args = append(args, "--export-format", exportFormat)
	}

	// Add format-specific options
	args = append(args, s.buildExportOptions(req)...)

	// Add input file
	args = append(args, scadFile)

	// Execute OpenSCAD command
	if err := s.executeCommand(args); err != nil {
		return nil, "", err
	}

	// Read output file
	data, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read output file: %w", err)
	}

	// Get content type
	contentType := s.getContentType(req.Format)

	return data, contentType, nil
}

// Summary generates summary information for SCAD content
func (s *OpenSCADService) Summary(req *models.SummaryRequest) (*models.SummaryResponse, error) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "scad-summary-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			// Log error but don't fail the operation
			fmt.Fprintf(os.Stderr, "warning: failed to remove temp directory %s: %v\n", tmpDir, err)
		}
	}()

	// Write SCAD content to temporary file
	scadFile := filepath.Join(tmpDir, "input.scad")
	if err := os.WriteFile(scadFile, []byte(req.ScadContent), 0644); err != nil {
		return nil, fmt.Errorf("failed to write SCAD file: %w", err)
	}

	// Create summary output file
	summaryFile := filepath.Join(tmpDir, "summary.json")

	// Build OpenSCAD command arguments
	summaryType := req.SummaryType
	if summaryType == "" {
		summaryType = "all"
	}

	args := []string{
		"--summary", summaryType,
		"--summary-file", summaryFile,
		"-o", filepath.Join(tmpDir, "dummy.stl"),
		scadFile,
	}

	// Execute OpenSCAD command
	if err := s.executeCommand(args); err != nil {
		return nil, err
	}

	// Read summary file
	data, err := os.ReadFile(summaryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read summary file: %w", err)
	}

	// Parse summary JSON
	var summary map[string]interface{}
	if err := json.Unmarshal(data, &summary); err != nil {
		return nil, fmt.Errorf("failed to parse summary JSON: %w", err)
	}

	return &models.SummaryResponse{Summary: summary}, nil
}

func (s *OpenSCADService) validateFormat(format string) error {
	validFormats := map[string]bool{
		"png":        true,
		"stl_binary": true,
		"stl_ascii":  true,
		"svg":        true,
		"pdf":        true,
		"3mf":        true,
	}

	if !validFormats[format] {
		return fmt.Errorf("unsupported format: %s", format)
	}

	return nil
}

func (s *OpenSCADService) getOutputExtension(format string) (string, string) {
	switch format {
	case "png":
		return "png", ""
	case "stl_binary":
		return "stl", "binstl"
	case "stl_ascii":
		return "stl", "asciistl"
	case "svg":
		return "svg", ""
	case "pdf":
		return "pdf", ""
	case "3mf":
		return "3mf", ""
	default:
		return "", ""
	}
}

func (s *OpenSCADService) buildExportOptions(req *models.ExportRequest) []string {
	var args []string

	switch req.Format {
	case "png":
		if req.Options.PNG != nil {
			if req.Options.PNG.Width != nil || req.Options.PNG.Height != nil {
				width := 800
				height := 600
				if req.Options.PNG.Width != nil {
					width = *req.Options.PNG.Width
				}
				if req.Options.PNG.Height != nil {
					height = *req.Options.PNG.Height
				}
				args = append(args, "--imgsize", fmt.Sprintf("%d,%d", width, height))
			}
		}

	case "stl_binary", "stl_ascii":
		if req.Options.STL != nil && req.Options.STL.DecimalPrecision != nil {
			precision := *req.Options.STL.DecimalPrecision
			if precision >= 1 && precision <= 16 {
				args = append(args, "-O", fmt.Sprintf("export-stl/decimal-precision=%d", precision))
			}
		}

	case "svg":
		if req.Options.SVG != nil {
			if req.Options.SVG.Fill != nil {
				args = append(args, "-O", fmt.Sprintf("export-svg/fill=%t", *req.Options.SVG.Fill))
			}
			if req.Options.SVG.FillColor != nil {
				args = append(args, "-O", fmt.Sprintf("export-svg/fill-color=%s", *req.Options.SVG.FillColor))
			}
			if req.Options.SVG.Stroke != nil {
				args = append(args, "-O", fmt.Sprintf("export-svg/stroke=%t", *req.Options.SVG.Stroke))
			}
			if req.Options.SVG.StrokeColor != nil {
				args = append(args, "-O", fmt.Sprintf("export-svg/stroke-color=%s", *req.Options.SVG.StrokeColor))
			}
			if req.Options.SVG.StrokeWidth != nil {
				args = append(args, "-O", fmt.Sprintf("export-svg/stroke-width=%s", strconv.FormatFloat(*req.Options.SVG.StrokeWidth, 'f', -1, 64)))
			}
		}

	case "pdf":
		if req.Options.PDF != nil {
			if req.Options.PDF.PaperSize != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/paper-size=%s", *req.Options.PDF.PaperSize))
			}
			if req.Options.PDF.Orientation != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/orientation=%s", *req.Options.PDF.Orientation))
			}
			if req.Options.PDF.ShowGrid != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/show-grid=%t", *req.Options.PDF.ShowGrid))
			}
			if req.Options.PDF.GridSize != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/grid-size=%s", strconv.FormatFloat(*req.Options.PDF.GridSize, 'f', -1, 64)))
			}
			if req.Options.PDF.Fill != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/fill=%t", *req.Options.PDF.Fill))
			}
			if req.Options.PDF.FillColor != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/fill-color=%s", *req.Options.PDF.FillColor))
			}
			if req.Options.PDF.Stroke != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/stroke=%t", *req.Options.PDF.Stroke))
			}
			if req.Options.PDF.StrokeColor != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/stroke-color=%s", *req.Options.PDF.StrokeColor))
			}
			if req.Options.PDF.StrokeWidth != nil {
				args = append(args, "-O", fmt.Sprintf("export-pdf/stroke-width=%s", strconv.FormatFloat(*req.Options.PDF.StrokeWidth, 'f', -1, 64)))
			}
		}

	case "3mf":
		if req.Options.ThreeMF != nil {
			if req.Options.ThreeMF.Unit != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/unit=%s", *req.Options.ThreeMF.Unit))
			}
			if req.Options.ThreeMF.DecimalPrecision != nil {
				precision := *req.Options.ThreeMF.DecimalPrecision
				if precision >= 1 && precision <= 16 {
					args = append(args, "-O", fmt.Sprintf("export-3mf/decimal-precision=%d", precision))
				}
			}
			if req.Options.ThreeMF.Color != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/color=%s", *req.Options.ThreeMF.Color))
			}
			if req.Options.ThreeMF.ColorMode != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/color-mode=%s", *req.Options.ThreeMF.ColorMode))
			}
			if req.Options.ThreeMF.MaterialType != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/material-type=%s", *req.Options.ThreeMF.MaterialType))
			}
			if req.Options.ThreeMF.AddMetadata != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/add-meta-data=%t", *req.Options.ThreeMF.AddMetadata))
			}
			if req.Options.ThreeMF.MetadataTitle != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/meta-data-title=%s", *req.Options.ThreeMF.MetadataTitle))
			}
			if req.Options.ThreeMF.MetadataDesigner != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/meta-data-designer=%s", *req.Options.ThreeMF.MetadataDesigner))
			}
			if req.Options.ThreeMF.MetadataDesc != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/meta-data-description=%s", *req.Options.ThreeMF.MetadataDesc))
			}
			if req.Options.ThreeMF.MetadataCopyright != nil {
				args = append(args, "-O", fmt.Sprintf("export-3mf/meta-data-copyright=%s", *req.Options.ThreeMF.MetadataCopyright))
			}
		}
	}

	return args
}

func (s *OpenSCADService) getContentType(format string) string {
	switch format {
	case "png":
		return "image/png"
	case "stl_binary", "stl_ascii":
		return "application/octet-stream"
	case "svg":
		return "image/svg+xml"
	case "pdf":
		return "application/pdf"
	case "3mf":
		return "application/vnd.ms-package.3dmodel+xml"
	default:
		return "application/octet-stream"
	}
}

func (s *OpenSCADService) executeCommand(args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, openscadCmd, args...)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("openscad command timed out")
		}
		return fmt.Errorf("openscad command failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

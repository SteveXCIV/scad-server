package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stevexciv/scad-server/models"
	"github.com/stevexciv/scad-server/services"
	"github.com/stevexciv/scad-server/version"
)

// Handler provides HTTP handlers
type Handler struct {
	openscadService services.OpenSCADExporter
}

// NewHandler creates a new handler with the default OpenSCAD service
func NewHandler() *Handler {
	return &Handler{
		openscadService: services.NewOpenSCADService(),
	}
}

// NewHandlerWithService creates a new handler with a custom OpenSCAD exporter
func NewHandlerWithService(exporter services.OpenSCADExporter) *Handler {
	return &Handler{
		openscadService: exporter,
	}
}

// Export handles the export endpoint
// @Summary Export SCAD to various formats
// @Description Exports OpenSCAD content to PNG, STL (binary/ASCII), SVG, PDF, 3MF, WebP, or AVIF format
// @Tags export
// @Accept json
// @Produce octet-stream
// @Param request body models.ExportRequest true "Export request"
// @Success 200 {file} binary "Exported file"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /openscad/v1/export [post]
func (h *Handler) Export(c *gin.Context) {
	var req models.ExportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	data, contentType, err := h.openscadService.Export(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "unsupported format: "+req.Format {
			statusCode = http.StatusBadRequest
		}
		log.Printf("OpenSCAD export error: %v", err)
		c.JSON(statusCode, models.ErrorResponse{
			Error:   "export failed",
			Message: err.Error(),
		})
		return
	}

	c.Data(http.StatusOK, contentType, data)
}

// Summary handles the summary endpoint
// @Summary Generate summary information
// @Description Generates summary information for OpenSCAD content
// @Tags summary
// @Accept json
// @Produce json
// @Param request body models.SummaryRequest true "Summary request"
// @Success 200 {object} models.SummaryResponse "Summary information"
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /openscad/v1/summary [post]
func (h *Handler) Summary(c *gin.Context) {
	var req models.SummaryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "invalid request",
			Message: err.Error(),
		})
		return
	}

	response, err := h.openscadService.Summary(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "summary generation failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck handles the health check endpoint
// @Summary Health check
// @Description Checks if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	info := version.GetInfo()
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"commit": info.Commit,
		"tag":    info.Tag,
	})
}

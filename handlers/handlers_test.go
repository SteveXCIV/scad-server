package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stevexciv/scad-server/models"
	"github.com/stevexciv/scad-server/services"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockOpenSCADExporter is a mock implementation of OpenSCADExporter for testing
type MockOpenSCADExporter struct {
	ExportFunc func(req *models.ExportRequest) ([]byte, string, error)
	SummaryFunc func(req *models.SummaryRequest) (*models.SummaryResponse, error)
}

func (m *MockOpenSCADExporter) Export(req *models.ExportRequest) ([]byte, string, error) {
	if m.ExportFunc != nil {
		return m.ExportFunc(req)
	}
	// Default behavior: return mock data
	return []byte("mock export data"), "application/octet-stream", nil
}

func (m *MockOpenSCADExporter) Summary(req *models.SummaryRequest) (*models.SummaryResponse, error) {
	if m.SummaryFunc != nil {
		return m.SummaryFunc(req)
	}
	// Default behavior: return mock summary
	return &models.SummaryResponse{
		Summary: map[string]interface{}{
			"facets": 6,
		},
	}, nil
}

func TestHealthCheck(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response["status"] != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response["status"])
	}
}

func TestExportEndpoint_InvalidJSON(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/openscad/v1/export", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestExportEndpoint_MissingRequiredFields(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name string
		body string
	}{
		{"Missing scad_content", `{"format":"png"}`},
		{"Missing format", `{"scad_content":"cube([10,10,10]);"}`},
		{"Empty scad_content", `{"scad_content":"","format":"png"}`},
		{"Empty format", `{"scad_content":"cube([10,10,10]);","format":""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/openscad/v1/export", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestExportEndpoint_InvalidFormat(t *testing.T) {
	router := setupRouter()

	reqBody := models.ExportRequest{
		ScadContent: "cube([10,10,10]);",
		Format:      "invalid_format",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/openscad/v1/export", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var errResp models.ErrorResponse
	if err := json.Unmarshal(w.Body.Bytes(), &errResp); err != nil {
		t.Errorf("Failed to parse error response: %v", err)
	}

	if errResp.Error != "export failed" {
		t.Errorf("Expected error 'export failed', got '%s'", errResp.Error)
	}
}

func TestSummaryEndpoint_InvalidJSON(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/openscad/v1/summary", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestSummaryEndpoint_MissingRequiredFields(t *testing.T) {
	router := setupRouter()

	tests := []struct {
		name string
		body string
	}{
		{"Missing scad_content", `{"summary_type":"all"}`},
		{"Empty scad_content", `{"scad_content":"","summary_type":"all"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/openscad/v1/summary", bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestSummaryEndpoint_ValidRequest(t *testing.T) {
	mock := &MockOpenSCADExporter{
		SummaryFunc: func(req *models.SummaryRequest) (*models.SummaryResponse, error) {
			return &models.SummaryResponse{
				Summary: map[string]interface{}{
					"facets": 6,
				},
			}, nil
		},
	}
	router := setupRouterWithMock(mock)

	reqBody := models.SummaryRequest{
		ScadContent: "cube([10,10,10]);",
		SummaryType: "all",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/openscad/v1/summary", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response models.SummaryResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response: %v", err)
	}

	if response.Summary == nil {
		t.Errorf("Expected summary data, got nil")
	}
}

func TestExportEndpoint_ValidFormats(t *testing.T) {
	mock := &MockOpenSCADExporter{
		ExportFunc: func(req *models.ExportRequest) ([]byte, string, error) {
			contentType := "application/octet-stream"
			switch req.Format {
			case "png":
				contentType = "image/png"
			case "svg":
				contentType = "image/svg+xml"
			case "pdf":
				contentType = "application/pdf"
			case "stl_binary", "stl_ascii":
				contentType = "application/octet-stream"
			}
			return []byte("mock export data"), contentType, nil
		},
	}
	router := setupRouterWithMock(mock)

	formats := []string{"png", "stl_binary", "stl_ascii", "svg", "pdf"}

	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			reqBody := models.ExportRequest{
				ScadContent: "cube([10,10,10]);",
				Format:      format,
			}

			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request: %v", err)
			}
			w := httptest.NewRecorder()
			req, err := http.NewRequest("POST", "/openscad/v1/export", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			router.ServeHTTP(w, req)

			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d: %s", w.Code, w.Body.String())
			}
		})
	}
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	h := NewHandler()

	router.GET("/health", h.HealthCheck)

	v1 := router.Group("/openscad/v1")
	{
		v1.POST("/export", h.Export)
		v1.POST("/summary", h.Summary)
	}

	return router
}

func setupRouterWithMock(exporter services.OpenSCADExporter) *gin.Engine {
	router := gin.Default()
	h := NewHandlerWithService(exporter)

	router.GET("/health", h.HealthCheck)

	v1 := router.Group("/openscad/v1")
	{
		v1.POST("/export", h.Export)
		v1.POST("/summary", h.Summary)
	}

	return router
}

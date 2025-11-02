package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stevexciv/scad-server/models"
)

func init() {
	gin.SetMode(gin.TestMode)
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
	router := setupRouter()

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

	// May fail if OpenSCAD is not installed, but should return proper error
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status 200 or 500, got %d", w.Code)
	}
}

func TestExportEndpoint_ValidFormats(t *testing.T) {
	router := setupRouter()

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

			// May fail if OpenSCAD is not installed, but should return proper error
			if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
				t.Errorf("Expected status 200 or 500, got %d", w.Code)
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

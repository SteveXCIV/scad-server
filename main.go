package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	_ "github.com/stevexciv/scad-server/docs"
	"github.com/stevexciv/scad-server/handlers"
	"github.com/stevexciv/scad-server/version"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title OpenSCAD HTTP API
// @version 1.0
// @description RESTful HTTP API that provides headless access to core OpenSCAD functionality
// @description Supports file export and summary generation capabilities
//
// @contact.name API Support
// @contact.url https://github.com/stevexciv/scad-server
//
// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
//
// @BasePath /
func main() {
	// Log version information
	info := version.GetInfo()
	log.Printf("Starting scad-server (commit: %s, tag: %s)", info.Commit, info.Tag)

	// Check if OpenSCAD is available
	if err := checkOpenSCAD(); err != nil {
		log.Fatalf("OpenSCAD not available: %v", err)
	}

	// Set Gin mode from environment variable
	mode := os.Getenv("SCADSRV_GIN_MODE")
	if mode == "" {
		mode = gin.ReleaseMode
	}
	gin.SetMode(mode)

	router := gin.Default()

	// Create handler
	h := handlers.NewHandler()

	// Health check endpoint
	router.GET("/health", h.HealthCheck)

	// API v1 routes
	v1 := router.Group("/openscad/v1")
	{
		v1.POST("/export", h.Export)
		v1.POST("/summary", h.Summary)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Get port from environment variable
	port := os.Getenv("SCADSRV_PORT")
	if port == "" {
		port = "8000"
	}

	// Check if port is available
	addr := ":" + port
	if err := checkPortAvailable(addr); err != nil {
		log.Fatalf("Port %s is not available: %v", port, err)
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// checkOpenSCAD verifies that the openscad binary is available
func checkOpenSCAD() error {
	cmd := exec.Command("openscad", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("openscad binary not found or not executable: %w", err)
	}
	return nil
}

// checkPortAvailable checks if the specified port is available for binding
func checkPortAvailable(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return listener.Close()
}

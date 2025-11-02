package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/stevexciv/scad-server/docs"
	"github.com/stevexciv/scad-server/handlers"
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
	// Set Gin mode from environment variable
	mode := os.Getenv("GIN_MODE")
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

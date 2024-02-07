package api

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// InitServer starts and initializes the API and its routes
func InitServer() {

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{AllowOrigins: []string{"*"}}))

	// Register routes for the server
	registerRoutes(e)

	// Get the port/address to start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "80" // Default port if not specified
	}
	address := ":" + port
	e.Logger.Fatal(e.Start(address))
}

// registerRoutes registers all routes for all versions of the API
func registerRoutes(e *echo.Echo) {
	// Serve the data collector

	// Start the API groups
	// api := e.Group("/api")
	// v1 := api.Group("/v1")
	{
	}
}

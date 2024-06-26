package api

import (
	"net/http"
	"os"

	"github.com/easymirror/easymirror-backend/internal/api/v1/router"
	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/easymirror/easymirror-backend/internal/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// InitServer starts and initializes the API and its routes
func InitServer(db *db.Database) {

	e := echo.New()
	e.Use(log.NewMiddlewareLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost*", "https://easymirror.io", "https://www.easymirror.io"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))

	// Register routes for the server
	router.Register(e, db)

	// Get the port/address to start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	address := ":" + port

	// Start the server
	e.Logger.Fatal(e.Start(address))
}

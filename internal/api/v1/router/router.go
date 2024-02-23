package router

import (
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/history"
	"github.com/easymirror/easymirror-backend/internal/db"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Register registers all routes for all versions of the API
func Register(e *echo.Echo, db *db.Database) {
	// Start the API groups
	api := e.Group("/api")
	v1 := api.Group("/v1", echojwt.WithConfig(jwtConfig(db)))
	{

		// TODO endpoint to get new session
		// TODO endpoint to refresh token

		// Endpoint to get a user's history
		history := &history.Handler{Database: db}
		v1.GET("/history", history.GetHistory)
	}
}

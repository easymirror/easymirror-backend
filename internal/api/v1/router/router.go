package router

import (
	"net/http"

	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/account"
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/auth"
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/history"
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/mirrors"
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/upload"
	"github.com/easymirror/easymirror-backend/internal/build"
	"github.com/easymirror/easymirror-backend/internal/db"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// Register registers all routes for all versions of the API
func Register(e *echo.Echo, db *db.Database) {
	// Start the API groups
	api := e.Group("/api")

	// Add health-check endpoint
	{
		api.GET("/health-check", healthcheck)
		api.GET("/build-info", buildInfo)
	}

	v1 := api.Group("/v1", echojwt.WithConfig(jwtConfig()))
	{
		// Auth endpounts
		auth := auth.Handler{Database: db}
		api.GET("/v1/auth/init", auth.NewJWT)
		api.GET("/v1/auth/refresh", auth.RefreshJWT)

		// Upload endpoints
		upload := upload.NewHandler(db)
		v1.GET("/mirror/new", upload.Init)
		v1.GET("/mirror", upload.PresignUri)
		v1.PUT("/mirror", upload.Mirror)

		// Account endpoints
		account := &account.Handler{Database: db}
		v1.GET("/user", account.GetUserInfo)
		v1.PATCH("/user/update", account.UpdateUser)

		// Mirrors endpoints
		mirrors := mirrors.Handler{Database: db}
		api.GET("/v1/mirror/:id", mirrors.GetMirror)

		// History Endpoints
		history := &history.Handler{Database: db}
		v1.GET("/history", history.GetHistory)
		v1.GET("/history/:id", history.GetFiles)
		v1.PATCH("/history/:id", history.UpdateHistoryItem)
		v1.DELETE("/history/:id", history.DeleteHistoryItem)

	}
}

// healthcheck is a handlder for incoming requests to the `/health-check` endpoint.
// It returns a 200 status code to indicate everything is OK
func healthcheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}

// buildInfo is a handlder for incoming requests to the `/build-info` endpoint.
func buildInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, build.GetInfo())
}

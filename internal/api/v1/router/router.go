package router

import (
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/account"
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/history"
	"github.com/easymirror/easymirror-backend/internal/api/v1/handlers/upload"
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

		// Upload endpoints
		upload := upload.NewHandler(db)
		v1.POST("/upload", upload.Upload) // TODO: delete this endpoint?
		v1.GET("/mirror", upload.Init)
		v1.PUT("/mirror", upload.Mirror)

		// Account endpoints
		account := &account.Handler{Database: db}
		v1.GET("/user", account.GetUserInfo)
		v1.PATCH("/user/update", account.UpdateUser)

		// History Endpoints
		history := &history.Handler{Database: db}
		v1.GET("/history", history.GetHistory)
		v1.GET("/history/:id", history.GetFiles)
		v1.PATCH("/history/:id", history.UpdateHistoryItem)
		v1.DELETE("/history/:id", history.DeleteHistoryItem)

	}
}

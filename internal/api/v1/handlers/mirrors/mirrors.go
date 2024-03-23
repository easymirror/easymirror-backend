package mirrors

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/easymirror/easymirror-backend/internal/mirrorlink"
	"github.com/labstack/echo/v4"
)

// GetMirrorLink is a handler for incoming `/mirrors/:id` requests
//
// It returns data about the mirror with a given id
func (h *Handler) GetMirror(c echo.Context) error {
	// Get the ID from the url
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		response := map[string]any{
			"success": false,
			"error":   "empty_id",
		}
		return c.JSON(http.StatusNotFound, response)
	}

	// Get sharelink from database
	type Response struct {
		Success               bool   `json:"success"`
		Error                 string `json:"error,omitempty"`
		*mirrorlink.ShareLink        // Embed everything from the ShareLink
	}

	ctx, cancel := context.WithTimeoutCause(context.Background(), 30*time.Second, errors.New("database took too long"))
	defer cancel()
	sl, err := mirrorlink.GetMirror(ctx, h.Database, id)
	if err != nil {
		log.Println("Error getting mirror: ", err)
		// TODO: Refactor this
		switch err.Error() {
		case "nothing found", "not found":
			return c.JSON(http.StatusNotFound, Response{Error: "not_found"})
		case "db took too long":
			return c.JSON(http.StatusNotFound, Response{Error: "db_took_long"})
		}
	}

	// Return
	return c.JSON(http.StatusOK, Response{Success: true, ShareLink: sl})
}

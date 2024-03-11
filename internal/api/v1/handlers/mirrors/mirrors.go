package mirrors

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success    bool      `json:"success"`
	ID         *string   `json:"id"`
	Name       *string   `json:"name"`
	UploadDate time.Time `json:"upload_date"`
	Links      HostLinks `json:"links"`
	Status     string    `json:"status"`
}
type HostLinks struct {
	Bunkr      *string `json:"bunkr"`
	Gofile     *string `json:"gofile"`
	Pixeldrain *string `json:"pixeldrain"`
	Cyberfile  *string `json:"cyberfile"`
	SaintTo    *string `json:"saint_to"`
	Cyberdrop  *string `json:"cyberdrop"`
}

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

	// Get info about the mirror link
	// TODO: Refactor this
	response := &Response{}
	ctx, cancel := context.WithTimeoutCause(context.Background(), 30*time.Second, errors.New("database took too long"))
	defer cancel()
	query := `
		SELECT mirroring_links.nickname, mirroring_links.upload_date, host_links.*
		FROM mirroring_links
		RIGHT JOIN host_links ON mirroring_links.id = host_links.mirror_id
		WHERE mirroring_links.id=($1);
	`
	row := h.PostgresConn.QueryRowContext(ctx, query, id)
	var err error
	if err = row.Scan(
		&response.Name,
		&response.UploadDate,
		&response.ID,
		&response.Links.Bunkr,
		&response.Links.Gofile,
		&response.Links.Pixeldrain,
		&response.Links.Cyberfile,
		&response.Links.SaintTo,
		&response.Links.Cyberdrop,
	); err == sql.ErrNoRows {
		response := map[string]any{
			"success": false,
			"error":   "not_found",
		}
		return c.JSON(http.StatusNotFound, response)
	}
	if err != nil {
		var response map[string]any
		switch err {
		// case sql.ErrNoRows:
		// 	response = map[string]any{"success": false, "error": "not_found"}
		case context.DeadlineExceeded:
			response = map[string]any{"success": false, "error": "db_too_long"}
		default:
			response = map[string]any{"success": false, "error": "not_found"}
		}
		return c.JSON(http.StatusNotFound, response)
	}

	// Return
	response.Success = true
	response.Status = "complete"
	return c.JSON(http.StatusOK, response)
}

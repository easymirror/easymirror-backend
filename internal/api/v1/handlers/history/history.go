package history

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/labstack/echo/v4"
)

// GetHistory returns a list of items a user has uploaded
func (h *Handler) GetHistory(c echo.Context) error {
	// Get the user-id from the JWT token
	user, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Get list of uploads from database
	page := c.QueryParam("page")
	var pageNum int
	pageNum, err = strconv.Atoi(page)
	if err != nil {
		pageNum = 0
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	links, err := user.MirrorLinks(ctx, h.Database, pageNum)
	if err != nil {
		log.Println("Error getting links:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Return list of items
	return c.JSON(http.StatusOK, links)
}

// UpdateHistoryItem updates the name of a given history item
func (h *Handler) UpdateHistoryItem(c echo.Context) error {
	// Get the user-id from the JWT token
	user, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Get the ID of the history item
	id := c.Param("id")

	// Get the new name of the link and validate
	name := c.FormValue("name")
	if strings.TrimSpace(name) == "" {
		resp := map[string]any{}
		resp["success"] = false
		resp["error"] = "name cannot be empty"
		return c.JSON(http.StatusBadRequest, resp)
	}

	// Update backend
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err = user.UpdateMirrorLinkName(ctx, h.Database, id, name); err != nil {
		// TODO return response based on error
		log.Println("Failed to update mirror link name:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Return response
	resp := map[string]any{}
	resp["success"] = true
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) DeleteHistoryItem(c echo.Context) error {
	// Get the user-id from the JWT token
	user, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Get the ID of the history item
	id := c.Param("id")

	// Delete the item
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err = user.DeleteMirrorLink(ctx, h.Database, id); err != nil {
		// TODO return response based on error
		log.Println("Failed to delete mirror link:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Return response
	resp := map[string]any{}
	resp["success"] = true
	return c.JSON(http.StatusOK, resp)
}

func (h *Handler) GetFiles(c echo.Context) error {
	// Get the user-id from the JWT token
	user, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Get the ID of the mirror link item
	id := c.Param("id")

	// Get the files
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	files, err := user.GetFiles(ctx, h.Database, id)
	if err != nil {
		// TODO return response based on error
		log.Println("Failed to get files in link:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Return
	return c.JSON(http.StatusOK, files)
}

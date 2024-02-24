package history

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetHistory returns a list of items a user has uploaded
func (h *Handler) GetHistory(c echo.Context) error {
	// Get the user-id from the JWT token
	token, ok := c.Get("jwt-token").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	user, err := user.FromJWT(token)
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
	token, ok := c.Get("jwt-token").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	user, err := user.FromJWT(token)
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

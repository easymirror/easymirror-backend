package history

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	*db.Database
}

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

package account

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUserInfo(c echo.Context) error {
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

	// Get the user's info
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	info, err := user.Info(ctx, h.Database)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	return c.JSON(http.StatusOK, info)
}

package account

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetUserInfo(c echo.Context) error {
	// Get the user-id from the JWT token
	user, err := user.FromEcho(c)
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

func (h *Handler) UpdateUser(c echo.Context) error {
	// Get the user-id from the JWT token
	u, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Update the user's info based on the key
	payload := make(map[string]string)
	if err = (&echo.DefaultBinder{}).BindBody(c, &payload); err != nil {
		log.Println("Error binding payload:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	var key user.InfoKey
	var val string
	switch {
	case payload["first_name"] != "":
		key = user.FirstNameKey
		val = payload["first_name"]
	case payload["last_name"] != "":
		key = user.LastNameKey
		val = payload["last_name"]
	case payload["phone"] != "":
		key = user.PhoneKey
		val = payload["phone"]
	case payload["username"] != "":
		key = user.UsernameKey
		val = payload["username"]
	default:
		resp := map[string]any{}
		resp["success"] = false
		resp["error"] = "bad request"
		return c.JSON(http.StatusBadRequest, resp)
	}

	if strings.TrimSpace(val) == "" {
		resp := map[string]any{}
		resp["success"] = false
		resp["error"] = "bad request"
		return c.JSON(http.StatusBadRequest, resp)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err = u.Update(ctx, h.Database, key, val); err != nil {
		log.Println("Failed to update user:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	resp := map[string]any{}
	resp["success"] = true
	return c.JSON(http.StatusOK, resp)
}

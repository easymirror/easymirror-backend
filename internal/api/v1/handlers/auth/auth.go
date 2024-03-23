package auth

import (
	"log"
	"net/http"

	"github.com/easymirror/easymirror-backend/internal/auth"
	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/labstack/echo/v4"
)

// NewJWT is a handler to issue new JWT Tokens
func (h *Handler) NewJWT(c echo.Context) error {
	// Create new user
	u, err := user.Create(h.Database)
	if err != nil {
		log.Println("Error creating user:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
	jwt, err := auth.GenerateJWT(u.ID().String())
	if err != nil {
		log.Println("Error generating JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	c.SetCookie(&http.Cookie{Name: auth.RefreshCookieName, Value: jwt.RefreshToken, HttpOnly: true, Path: "/"})
	response := map[string]any{
		"success":      true,
		"access_token": jwt.AccessToken,
	}
	return c.JSON(http.StatusOK, response)
}

// RefreshJWT is a handler to refresh expired access tokens
func (h *Handler) RefreshJWT(c echo.Context) error {
	// Get refresh token from cookie
	cookie, err := c.Cookie(auth.RefreshCookieName)
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	// Refresh token
	newAccess, err := auth.RefreshAccessToken(cookie.Value)
	if err != nil {
		// TODO: Handle this error better
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// Return new token
	response := map[string]any{
		"success":      true,
		"access_token": newAccess,
	}
	return c.JSON(http.StatusOK, response)
}

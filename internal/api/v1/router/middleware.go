package router

import (
	"errors"
	"net/http"
	"os"

	"github.com/easymirror/easymirror-backend/internal/auth"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func generateUnauthorizedResponse(c echo.Context, action string) error {
	type Response struct {
		Success bool   `json:"success"`
		Action  string `json:"action"`
	}
	response := Response{
		Success: false,
		Action:  action,
	}
	return c.JSON(http.StatusUnauthorized, response)
}

// jwtConfig provides a config middleware for authenticating JWT tokens
func jwtConfig() echojwt.Config {
	return echojwt.Config{
		SigningKey:    []byte(os.Getenv("JWT_ACCESS_SECRET")),
		SigningMethod: echojwt.AlgorithmHS256,
		TokenLookup:   "header:Authorization:Bearer ,cookie:user_session",
		ContextKey:    "jwt-token",
		// ContinueOnIgnoredError: true, // Set this to `true` so it can go to the correct handler
		ErrorHandler: func(c echo.Context, err error) error {
			if errors.Is(err, echojwt.ErrJWTInvalid) {
				// Check cookies to see if a refresh token is present. If present, tell them to refresh
				if _, err := c.Cookie(auth.RefreshCookieName); err == nil {
					return generateUnauthorizedResponse(c, "refresh_token")
				}
			}
			return generateUnauthorizedResponse(c, "new_token")
		},
	}
}

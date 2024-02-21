package router

import (
	"errors"
	"net/http"
	"os"

	"github.com/easymirror/easymirror-backend/internal/auth"
	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/easymirror/easymirror-backend/internal/user"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// jwtConfig provides a config middleware for authenticating JWT tokens
func jwtConfig(db *db.Database) echojwt.Config {
	return echojwt.Config{
		SigningKey:    []byte(os.Getenv("JWT_ACCESS_SECRET")),
		SigningMethod: echojwt.AlgorithmHS256,
		TokenLookup:   "header:Authorization:Bearer ,cookie:user_session",
		ContextKey:    "jwt-token",
		// ContinueOnIgnoredError: true, // Set this to `true` so it can go to the correct handler
		ErrorHandler: func(c echo.Context, err error) error {
			if errors.Is(err, echojwt.ErrJWTInvalid) {
				// Check cookies to see if a refresh token is present. If present, refresh acess token & return
				cookie, err := c.Cookie(auth.RefreshCookieName)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Internal Server Error")
				}

				// refresh access token
				newAccess, err := auth.RefreshAccessToken(cookie.Value)
				if err != nil {
					return c.String(http.StatusInternalServerError, "Internal Server Error")
				}

				// Return new access token
				c.Response().Header().Set("Authorization", newAccess)
				return c.String(http.StatusNoContent, `"action":"refresh"`)
			}

			// Create and set new JWT Pair (access & refresh token)
			user, err := user.Create(db)
			if err != nil {
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}

			t, err := auth.GenerateJWT(user.ID().String())
			if err != nil {
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}
			c.Response().Header().Set("Authorization", t.AccessToken)
			c.SetCookie(&http.Cookie{Name: auth.RefreshCookieName, Value: t.RefreshToken, HttpOnly: true})

			// Return nil so it can continue to the appropriate handler
			return nil
		},
	}
}

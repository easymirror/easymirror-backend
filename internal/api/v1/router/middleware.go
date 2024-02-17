package router

import (
	"errors"
	"log"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

// jwtConfig provides a config middleware for authenticating JWT tokens
func jwtConfig() echojwt.Config {
	return echojwt.Config{
		SigningKey:             []byte(os.Getenv("JWT_ACCESS_SECRET")),
		SigningMethod:          echojwt.AlgorithmHS256,
		TokenLookup:            "header:Authorization:Bearer ,cookie:user_session",
		ContextKey:             "jwt-token",
		ContinueOnIgnoredError: true, // Set this to `true` so it can go to the correct handle
		ErrorHandler: func(c echo.Context, err error) error {
			log.Println("There was an error with the JWT:", err)
			switch {
			case errors.Is(err, echojwt.ErrJWTMissing):
				// TODO Create a new JWT Pair (access & refresh token)
				log.Println("JWT is missing.")
			case errors.Is(err, echojwt.ErrJWTInvalid):
				// TODO check cookies to see if a refresh token is present
				log.Println("JWT is invalid.")
			}

			// Return nil so it can continue to the appropriate handler
			return nil
		},
	}
}

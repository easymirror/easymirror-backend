package router

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/easymirror/easymirror-backend/internal/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Load .env file
	if err := godotenv.Load("../../../../.env"); err != nil {
		log.Println("no env file loaded.")
	}
}

// go test -v -timeout 30s -run ^TestJWTConfig$ github.com/easymirror/easymirror-backend/internal/api/v1/router
func TestJWTConfig(t *testing.T) {
	// Setup server
	e := echo.New()
	e.Use(echojwt.WithConfig(jwtConfig()))
	e.GET("/", func(c echo.Context) error {
		token, ok := c.Get("jwt-token").(*jwt.Token) // by default token is stored under `user` key
		if !ok {
			t.Fatal("failed to retrieve JWT Token")
		}
		return c.JSON(http.StatusOK, token.Claims)
	})

	// Test with a valid auth header
	t.Run("Valid Auth Header", func(t *testing.T) {
		// Generate a valid JWT Token
		token, err := auth.GenerateJWT("test_user_id")
		if err != nil {
			t.Fatalf("Error generating JWT: %v", err)
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", token.AccessToken))
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		// Assertions
		assert.Equal(t, http.StatusOK, res.Code)
	})

	// Test with an expired auth header
	t.Run("Expired Auth Header", func(t *testing.T) {
		expiredToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDg0ODA0NzMsImlhdCI6MTcwODQ3OTU3MywiaXNzIjoiZWFzeW1pcnJvci5pbyIsInN1YiI6InRlc3RfdXNlcl9pZCJ9.EFji46DaEDEHCdI-A7RWUQzDijZIZDxlNzny98vGoe0"

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("Bearer %v", expiredToken))

		// Set refresh cookies
		req.AddCookie(&http.Cookie{Name: auth.RefreshCookieName, Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE3MDg0Nzk1NzMsImlzcyI6ImVhc3ltaXJyb3IuaW8ifQ.PWKEbin5Sq47BkGe_JwJlXCbH4UEcR6b3T42Q7pdnY8"})

		// Server HTTP
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		// Assertions
		assert.NotEmpty(t, res.Header().Get("Authorization"), "Was expecting new access token to be set") // assert that a new access token / `Authorization` header is set
		assert.Equal(t, http.StatusNoContent, res.Code, "Expected: %v || Received: %v", http.StatusNoContent, res.Code)
	})

	// Test with missing auth header
	t.Run("No Auth Header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		// Server HTTP
		res := httptest.NewRecorder()
		e.ServeHTTP(res, req)

		// Assertions
		assert.NotEmpty(t, res.Header().Get("set-cookie"), "no cookies were set")                 // assert that refresh_token is set in cookies
		assert.NotEmpty(t, res.Header().Get("Authorization"), "Authorization header was not set") // assert that `Authorization` is set in header
	})
}

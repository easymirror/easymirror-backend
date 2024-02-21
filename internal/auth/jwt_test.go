package auth

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no env file loaded.")
	}
}

// go test -v -timeout 30s -run ^TestGenerateJWT$ github.com/easymirror/easymirror-backend/internal/auth
func TestGenerateJWT(t *testing.T) {
	accessSecret := []byte(os.Getenv("JWT_ACCESS_SECRET"))
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(10 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    issuer,
	}
	token, err := generateJWT(accessSecret, claims)
	if err != nil {
		panic(err)
	}
	fmt.Println(token)
	assert.NotEmpty(t, token)
}

// go test -v -timeout 30s -run ^TestValidateJWT$ github.com/easymirror/easymirror-backend/internal/auth
func TestValidateJWT(t *testing.T) {
	// Generate a JWT Token
	token, err := GenerateJWT("some_id")
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Validate the JWT Token
	tkn, err := ValidateJWT(token.AccessToken)
	if err != nil {
		t.Fatalf("Error validating JWT: %v", err)
	}

	assert.Equal(t, true, tkn.Valid)
}

// go test -v -timeout 30s -run ^TestRefreshAccessToken$ github.com/easymirror/easymirror-backend/internal/auth
func TestRefreshAccessToken(t *testing.T) {
	// Generate a JWT
	auth, err := GenerateJWT("some_id")
	if err != nil {
		t.Fatalf("Error generating JWT: %v", err)
	}

	// Refresh the token
	newAuth, err := RefreshAccessToken(auth.RefreshToken)
	if err != nil {
		t.Fatalf("Error refreshing JWT: %v", err)
	}
	assert.NotEmpty(t, newAuth)
}

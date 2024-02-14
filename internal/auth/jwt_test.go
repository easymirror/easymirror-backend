package auth

import (
	"log"
	"testing"

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
	token, err := GenerateJWT("some_id")
	if err != nil {
		panic(err)
	}
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

	tkn.Claims.Valid()
	assert.Equal(t, true, tkn.Valid)
}

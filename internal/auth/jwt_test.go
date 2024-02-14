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

package cyberfile

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Load Env
	if err := godotenv.Load("../../../.env"); err != nil {
		panic(err)
	}
}

// go test -v -timeout 30s -run ^TestGetAuthToken$ github.com/easymirror/easymirror-backend/internal/hosts/cyberfile
func TestGetAuthToken(t *testing.T) {
	// Create new account
	a := account{
		username: os.Getenv("CYBERFILE_USERNAME"),
		password: os.Getenv("CYBERFILE_PASSWORD"),
	}

	// Run test
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	token, err := a.GetAccessToken(ctx)
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	assert.NotEmpty(t, token)
}

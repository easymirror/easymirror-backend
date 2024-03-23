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

// go test -v -timeout 30s -run ^TestCreateFolder$ github.com/easymirror/easymirror-backend/internal/hosts/cyberfile
func TestCreateFolder(t *testing.T) {
	// Create account
	a := newAccount(os.Getenv("CYBERFILE_USERNAME"), os.Getenv("CYBERFILE_PASSWORD"))
	if _, err := a.GetAccessToken(context.Background()); err != nil {
		t.Fatalf("Failed to get access token: %v", err)
	}

	// Run Test
	id, err := createFolder(context.Background(), a, "some_mirror_id")
	if err != nil {
		t.Fatalf("Failed to create folder: %v", err)
	}
	assert.NotEmpty(t, id)
}

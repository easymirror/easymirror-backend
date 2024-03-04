package bunkr

import (
	"context"
	"fmt"
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

// go test -v -timeout 30s -run ^TestGetUploadLink$ github.com/easymirror/easymirror-backend/internal/hosts/bunkr
func TestGetUploadLink(t *testing.T) {
	link, err := getUploadLink(context.Background())
	if err != nil {
		t.Fatalf("Error getting upload link: %v", err)
	}
	assert.NotEmpty(t, link)
	fmt.Println(link)

}


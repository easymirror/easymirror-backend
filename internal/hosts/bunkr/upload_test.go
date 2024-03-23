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

// go test -v -timeout 30s -run ^TestUpload$ github.com/easymirror/easymirror-backend/internal/hosts/bunkr
func TestUpload(t *testing.T) {
	presignURL := "https://easymirror.s3.us-east-1.amazonaws.com/48751a99-dc86-4e08-a885-ef0f75337779/mattaio.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA4LVGZIJF45Z3KRFR%2F20240304%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20240304T192514Z&X-Amz-Expires=86400&X-Amz-SignedHeaders=host&x-id=GetObject&X-Amz-Signature=c731217155e560a1afc22cbfb25a6480e4f48936091613aa78bc66febef001d2"
	albumID := "307806"

	something, err := upload(context.Background(), albumID, presignURL)
	if err != nil {
		t.Fatalf("Error uploading: %v", err)
	}
	assert.NotEmpty(t, something)
}

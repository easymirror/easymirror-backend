package bunkr

import (
	"context"
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

func TestCreateFolder(t *testing.T) {
	id, err := createFolder(context.Background(), "Some Name", true, true)
	if err != nil {
		t.Fatalf("error creating folder: %v", err)
	}
	assert.NotEmpty(t, id)
}

package upload

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/easymirror/easymirror-backend/internal/db"
)

type Handler struct {
	*db.Database
	S3Client *s3.Client
}

// NewHandler returns a new upload handler with a S3Client
func NewHandler(db *db.Database) *Handler {

	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println("S3 Region:", cfg.Region)

	return &Handler{
		Database: db,
		S3Client: s3.NewFromConfig(cfg),
	}
}

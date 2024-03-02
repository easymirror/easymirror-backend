package upload

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type mirrorHost string

const (
	GofileHost     mirrorHost = "gofile"
	BunkrHost      mirrorHost = "bunkr"
	PixelDrainHost mirrorHost = "pixeldrain"
	CyberfileHost  mirrorHost = "cyberfile"
)

// Mirror handles incoming PUT requests for mirroring sites.
func (h *Handler) Mirror(c echo.Context) error {
	// Get user data from the JWT token
	token, ok := c.Get("jwt-token").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	user, err := user.FromJWT(token)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Parse the body
	body := &struct {
		MirrorID string       `json:"id"`
		Sites    []mirrorHost `json:"sites"`
	}{}
	err = (&echo.DefaultBinder{}).BindBody(c, &body)
	if err != nil {
		log.Println("Error binding body: ", err)
		return err
	}
	fmt.Println(user) // TODO: Delete this

	// Get files from AWS S3 bucket
	result, err := h.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Prefix: aws.String(body.MirrorID + "/"),
	})
	if err != nil {
		log.Println("Error getting folder:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	presignedLinks := make([]string, *result.KeyCount)
	for i, file := range result.Contents {
		// Create presigned URLs for each file in bucket
		if file.Key == nil {
			log.Println("Cannot create presign url. Name is empty")
			continue
		}
		url, err := getPresignURL(h.S3Client, file.Key)
		if err != nil {
			log.Println("Error creating presigned url:", err)
			continue
		}
		presignedLinks[i] = url
	}

	// Parse which sites to mirror to
	for _, host := range body.Sites {
		// TODO: Begin mirroring process
		switch host {
		case BunkrHost:
		case GofileHost:
		case PixelDrainHost:
		case CyberfileHost:
		}
	}
	// TODO: Save mirror links to `host_links` table
	return nil
}

// getPresignURL creates a presigned URL for a given file key so users can make GET requests to
func getPresignURL(s3client *s3.Client, fileKey *string) (string, error) {
	presignClient := s3.NewPresignClient(s3client)
	presignedUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    fileKey,
		})
	if err != nil {
		return "", fmt.Errorf("getPresignURL error: %w", err)
	}
	return presignedUrl.URL, nil
}

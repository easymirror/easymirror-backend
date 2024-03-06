package upload

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/easymirror/easymirror-backend/internal/hosts/bunkr"
	"github.com/easymirror/easymirror-backend/internal/hosts/pixeldrain"
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
	go func() {
		// Make sure sites are unique so we only upload once to the host
		siteMap := map[mirrorHost]bool{}
		for _, chosen := range body.Sites {
			siteMap[chosen] = true
		}

		// Start TX
		tx, err := h.Database.PostgresConn.Begin()
		if err != nil {
			log.Println("Error creating transaction:", err)
			return
		}

		// TODO: Defer delete from AWS S3
		for host := range siteMap {
			switch host {
			case BunkrHost:
				log.Println("Uploading file to Bunkr. Presigned link:", presignedLinks)
				folder, err := bunkr.UploadTx(context.TODO(), tx, body.MirrorID, presignedLinks)
				if err != nil {
					log.Println("Error uploading to bunk:", err)
					continue
				}

				// TODO: Do something with the ID
				fmt.Printf("Folder Link: %q\n", folder)
			case GofileHost:

			case PixelDrainHost:
				log.Println("Uploading file to pixeldrain. Presigned link:", presignedLinks)
				folder, err := pixeldrain.UploadTX(context.TODO(), tx, body.MirrorID, presignedLinks)
				if err != nil {
					log.Println("Error uploading to pixel drain:", err)
					continue
				}

				// TODO: Do something with the ID
				fmt.Printf("Folder Link: %q\n", folder)
			case CyberfileHost:
			}
		}

		// Save mirror links to `host_links` table
		log.Println("Saving to database...")
		if err = tx.Commit(); err != nil {
			log.Println("Error committing tx:", err)
			tx.Rollback()
			return
		}
	}()

	// Return Response
	response := map[string]any{
		"success":   true,
		"mirror_id": body.MirrorID,
	}
	return c.JSON(http.StatusOK, response)
}

// getPresignURL creates a presigned URL for a given file key so users can make GET requests to
func getPresignURL(s3client *s3.Client, fileKey *string) (string, error) {
	presignClient := s3.NewPresignClient(s3client)
	presignedUrl, err := presignClient.PresignGetObject(context.Background(),
		&s3.GetObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    fileKey,
		},
		s3.WithPresignExpires(24*time.Hour))
	if err != nil {
		return "", fmt.Errorf("getPresignURL error: %w", err)
	}
	return presignedUrl.URL, nil
}

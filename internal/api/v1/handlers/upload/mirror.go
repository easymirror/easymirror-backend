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
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/easymirror/easymirror-backend/internal/hosts/bunkr"
	"github.com/easymirror/easymirror-backend/internal/hosts/pixeldrain"
	"github.com/easymirror/easymirror-backend/internal/user"
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
	user, err := user.FromEcho(c)
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
	// TODO Validate user has access to mirror ID

	// Get files from AWS S3 bucket
	files, err := getFilesInS3Dir(h.S3Client, body.MirrorID)
	if err != nil {
		log.Println("Error getting folder:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	var presignedLinks []string
	for _, file := range files {
		// Create presigned URLs for each file in bucket
		if file.Key == nil {
			log.Println("Cannot create presign url. Name is empty")
			continue
		}
		if *file.Key == body.MirrorID+"/" { // Skip the folder itself
			continue
		}

		url, err := getPresignURL(h.S3Client, file.Key)
		if err != nil {
			log.Println("Error creating presigned url:", err)
			continue
		}
		presignedLinks = append(presignedLinks, url)
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

		// Delete from AWS S3 when done
		defer deleteFromS3(h.S3Client, body.MirrorID)

		for host := range siteMap {
			switch host {
			case BunkrHost:
				log.Println("Uploading file to Bunkr. Presigned link:", presignedLinks)
				_, err := bunkr.UploadTx(context.TODO(), tx, body.MirrorID, presignedLinks)
				if err != nil {
					log.Println("Error uploading to bunk:", err)
					continue
				}
			case GofileHost:

			case PixelDrainHost:
				log.Println("Uploading file to pixeldrain. Presigned link:", presignedLinks)
				_, err := pixeldrain.UploadTX(context.TODO(), tx, body.MirrorID, presignedLinks)
				if err != nil {
					log.Println("Error uploading to pixel drain:", err)
					continue
				}
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

// deleteFromS3 deletes a given object in AWS S3
func deleteFromS3(s3client *s3.Client, mirrorID string) error {
	// List all the files in the directory
	objects, err := getFilesInS3Dir(s3client, mirrorID)
	if err != nil {
		return fmt.Errorf("failed to get files: %w", err)
	}
	ids := make([]types.ObjectIdentifier, len(objects))
	for i, obj := range objects {
		ids[i] = types.ObjectIdentifier{Key: obj.Key}
	}

	// Delete everything
	_, err = s3client.DeleteObjects(context.TODO(), &s3.DeleteObjectsInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Delete: &types.Delete{Objects: ids},
	})
	if err != nil {
		log.Printf("Couldn't delete objects from bucket. Here's why: %v\n", err)
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// getFilesInS3Dir returns a list of items in a given directory in a S3 bucket
func getFilesInS3Dir(s3client *s3.Client, mirrorID string) ([]types.Object, error) {
	result, err := s3client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Prefix: aws.String(mirrorID + "/"),
	})
	if err != nil {
		log.Printf("Couldn't list objects in bucket. Here's why: %v\n", err)
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}
	return result.Contents, err
}

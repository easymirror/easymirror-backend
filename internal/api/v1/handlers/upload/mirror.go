package upload

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/easymirror/easymirror-backend/internal/db"
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

const (
	maxMirrorTasks = 3 // The max number of gorountines when mirroring
	taskTimeout    = 1 * time.Hour
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

	// Generate presigned URLs for each file
	presignedLinks := genPresignURIs(h.S3Client, files, body.MirrorID)

	// Mirror the files
	go mirrorFiles(h.Database, h.S3Client, body.MirrorID, body.Sites, presignedLinks)

	// Return Response
	response := map[string]any{
		"success":   true,
		"mirror_id": body.MirrorID,
	}
	return c.JSON(http.StatusOK, response)
}

// mirrorFiles uploads files to the users other sites
func mirrorFiles(db *db.Database, s3client *s3.Client, mirrorID string, sites []mirrorHost, sourceURIs []string) {
	// Make sure sites are unique so we only upload once to the host
	siteMap := map[mirrorHost]bool{}
	for _, chosen := range sites {
		siteMap[chosen] = true
	}

	// Start TX
	tx, err := db.PostgresConn.Begin()
	if err != nil {
		log.Println("Error creating transaction:", err)
		return
	}

	// Delete from AWS S3 when done
	defer deleteFromS3(s3client, mirrorID)

	// Begin the mirroring process
	ctx, cancel := context.WithTimeout(context.Background(), taskTimeout)
	defer cancel()
	var wg sync.WaitGroup
	sem := make(chan int, maxMirrorTasks)
	for host := range siteMap {
		wg.Add(1)
		sem <- 1 // will block if there is MAX ints in sem / until

		switch host {
		case BunkrHost:
			go func() {
				defer wg.Done()
				defer func() { <-sem }() // removes an int from sem, allowing another to proceed
				_, err := bunkr.UploadTx(ctx, tx, mirrorID, sourceURIs)
				if err != nil {
					log.Println("Error uploading to bunk:", err)
				}
			}()
		case PixelDrainHost:
			go func() {
				defer wg.Done()
				defer func() { <-sem }() // removes an int from sem, allowing another to proceed
				_, err := pixeldrain.UploadTX(ctx, tx, mirrorID, sourceURIs)
				if err != nil {
					log.Println("Error uploading to pixel drain:", err)
				}
			}()
		case GofileHost:
		case CyberfileHost:
		}
	}
	wg.Wait() // Wait for all tasks to be finished

	// Save/Commit mirror links to the `host_links` table
	if err = tx.Commit(); err != nil {
		log.Println("Error committing tx:", err)
		tx.Rollback()
		return
	}
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

// genPresignURIs generates presigned URIs for objects in a given mirror ID
func genPresignURIs(s3client *s3.Client, files []types.Object, mirrorID string) []string {
	var presignedLinks []string
	for _, file := range files {
		// Create presigned URLs for each file in bucket
		if file.Key == nil {
			log.Println("Cannot create presign url. Name is empty")
			continue
		}
		if *file.Key == mirrorID+"/" { // Skip the folder itself
			continue
		}

		url, err := getPresignURL(s3client, file.Key)
		if err != nil {
			log.Println("Error creating presigned url:", err)
			continue
		}
		presignedLinks = append(presignedLinks, url)
	}
	return presignedLinks
}

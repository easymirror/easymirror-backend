package upload

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	// The buffer size (in bytes) to use when buffering data into chunks and sending them as parts to S3.
	// The minimum allowed part size is 5MB, and if this value is set to zero,
	// the DefaultUploadPartSize value will be used.
	partMiBs int64 = 50 * megabyte // 50MB
	megabyte       = 1024 * 1024   //  1 megabyte
)

const (
	presignExp = 6 * time.Hour
)

// Init is a handler for incoming GET requests.
// It returns a valid mirror link ID along with a URL for users to upload their content to.
func (h *Handler) Init(c echo.Context) error {
	// Get user data from JWT token
	user, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Generate a new UUID
	mirrorID := uuid.New()

	// Generate a new mirror link in the database
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	tx, err := h.PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error creating transaction:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	_, err = tx.Exec(`
	INSERT INTO mirroring_links (id, created_by_id, upload_date)
	VALUES
	(($1), ($2), ($3));
`, mirrorID, user.ID(), time.Now().UTC())
	if err != nil {
		log.Println("Error creating new mirror link:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		log.Println("Error committing to database:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Return to user
	response := map[string]any{"success": true, "id": mirrorID}
	return c.JSON(http.StatusOK, response)

}

// PresignUri is a handler for incoming GET requests.
// It returns a valid presigned uri for the user can upload their files to
func (h *Handler) PresignUri(c echo.Context) error {
	// Get user data from JWT token
	user, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Get the name of the file & mirror ID, if any
	filename := c.QueryParam("n")
	mirrorID := strings.TrimSpace(c.QueryParam("id"))
	if mirrorID == "" {
		// Generate a new mirror link if there is none
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		tx, err := h.PostgresConn.BeginTx(ctx, nil)
		if err != nil {
			log.Println("Error creating transaction:", err)
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
		mirrorID = uuid.NewString()
		_, err = tx.Exec(`
		INSERT INTO mirroring_links (id, created_by_id, upload_date)
		VALUES
		(($1), ($2), ($3));
	`, mirrorID, user.ID(), time.Now().UTC())
		if err != nil {
			log.Println("Error creating new mirror link:", err)
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
		if err = tx.Commit(); err != nil {
			tx.Rollback()
			log.Println("Error committing to database:", err)
			return c.String(http.StatusInternalServerError, "Internal server error")
		}
	}

	// TODO: Validate if the ID exists?

	// Generate a presign URL
	presignURL, err := putPresignURL(h.S3Client, presignExp, mirrorID, filename)
	if err != nil {
		log.Println("Error creating new presign url:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Return data in response
	type Upload struct {
		URI        string `json:"uri"`
		ValidUntil string `json:"valid_until"`
	}
	type Response struct {
		Success  bool   `json:"success"`
		MirrorID string `json:"mirror_id"`
		Upload   Upload `json:"upload"`
	}
	r := &Response{
		Success:  true,
		MirrorID: mirrorID,
		Upload: Upload{
			URI:        presignURL,
			ValidUntil: time.Now().Add(presignExp).Format(time.RFC3339),
		},
	}
	return c.JSON(http.StatusOK, r)
}

// Upload is a handler for incoming POST requests.
// It takes in files and mirrors to other hosts.
//
// Deprecated: This flow has been moved to a different flow.
// Users should call the Init method.
func (h *Handler) Upload(c echo.Context) error {
	// Get user data from the JWT
	user, err := user.FromEcho(c)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Read files form the upload
	form, err := c.MultipartForm()
	if err != nil {
		return err
	}
	files, ok := form.File["files"] // This files key is from the client, we named the input "files"
	if !ok {
		resp := map[string]any{"success": false, "error": "no files"}
		return c.JSON(http.StatusBadRequest, resp)
	}

	// Begin TX
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	tx, err := h.Database.PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Error beggning TX:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// Create a new mirror link
	mirrorID := uuid.New()
	_, err = tx.Exec(`
		INSERT INTO mirroring_links (id, created_by_id, upload_date)
		VALUES
		(($1), ($2), ($3));
	`, mirrorID, user.ID(), time.Now().UTC())
	if err != nil {
		log.Println("Error creating new mirror link:", err)
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// TODO: Add Goroutines to quickly process files
	var src multipart.File
	var srcBytes []byte
	for _, file := range files {
		src, err = file.Open()
		if err != nil {
			return err
		}
		defer src.Close() // TODO unmake this defer?
		srcBytes, err = io.ReadAll(src)
		if err != nil {
			log.Println("Error converting to bytes:", err)
			continue
		}

		// Upload to AWS S3
		if err = uploadToBucket(h.S3Client, srcBytes, mirrorID.String(), file.Filename); err != nil {
			log.Println("Could not upload file to bucket:", err)
			continue
		}

		// Upload file data to database
		if err = addFileData(tx, file, mirrorID.String()); err != nil {
			log.Println("Error uploading to database:", err)
		}
	}

	// Commit the changes to SQ
	if err = tx.Commit(); err != nil {
		log.Println("Error comitting tx:", err)
		tx.Rollback()
		return c.String(http.StatusInternalServerError, "Internal server error")
	}

	// TODO: Upload to other hosts
	// TODO: add a update statement for duration?

	resp := map[string]any{
		"success":   true,
		"mirror_id": mirrorID,
	}
	return c.JSON(http.StatusOK, resp)
}

// addFileData adds a file's data to the database
func addFileData(tx *sql.Tx, file *multipart.FileHeader, mirrorID string) error {
	_, err := tx.Exec(`
		INSERT INTO files (id, name, size_bytes, upload_date, mirror_link_id)
		VALUES
		(($1), ($2), ($3), ($4), ($5));
		`, uuid.NewString(), file.Filename, file.Size, time.Now().UTC(), mirrorID)
	if err != nil {
		return fmt.Errorf("tx.Exec error: %v", err)
	}
	return nil
}

// uploadToBucket uploads a given file to the AWS S3 bucket
func uploadToBucket(c *s3.Client, srcBytes []byte, mirrorID, fileName string) error {
	// We use a manager to upload data to an object in a bucket.
	// The upload manager breaks large data into parts and uploads the parts concurrently.
	contentBuffer := bytes.NewReader(srcBytes)
	uploader := manager.NewUploader(c, func(u *manager.Uploader) {
		u.PartSize = partMiBs
	})

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(filepath.Join(mirrorID, fileName)),
		Body:   contentBuffer,
	})
	if err != nil {
		return fmt.Errorf("uploader error: %w", err)
	}
	return nil
}

func putPresignURL(s3client *s3.Client, expiration time.Duration, mirrorID, filename string) (string, error) {
	presignClient := s3.NewPresignClient(s3client)
	presignedUrl, err := presignClient.PresignPutObject(context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    aws.String(filepath.Join(mirrorID, filename)),
		},
		s3.WithPresignExpires(expiration))
	if err != nil {
		return "", fmt.Errorf("presignPutObject error: %w", err)
	}
	return presignedUrl.URL, nil
}

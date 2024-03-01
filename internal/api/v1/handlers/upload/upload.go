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
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/easymirror/easymirror-backend/internal/user"
	"github.com/golang-jwt/jwt/v5"
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

func (h *Handler) Upload(c echo.Context) error {
	// Get user data from the JWT
	token, ok := c.Get("jwt-token").(*jwt.Token) // by default token is stored under `user` key
	if !ok {
		return c.String(http.StatusInternalServerError, "Internal server error")
	}
	user, err := user.FromJWT(token)
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

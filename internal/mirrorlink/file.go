package mirrorlink

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/google/uuid"
)

type File struct {
	ID         uuid.UUID `json:"id"`          // ID of the file
	Name       string    `json:"name"`        // Name of the file
	SizeBytes  int64     `json:"size"`        // Size of the file in bytes
	UploadDate time.Time `json:"upload_date"` // Date the file was uploaded
}

// GetFilesFromMirror returns a list of files from a given mirror link
func GetFilesFromMirror(ctx context.Context, db *db.Database, mirrorId, userID string) ([]File, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	// Begin TX
	tx, err := db.PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("BeginTx error: %w", err)
	}

	// Get files
	// We use a join here to ensure that we only select files that belongs to the user
	query := `
		SELECT files.id, files.name, files.size_bytes, files.upload_date
		FROM mirroring_links INNER JOIN files
		ON mirroring_links.id = files.mirror_link_id
		WHERE mirroring_links.created_by_id=($1)
		AND mirroring_links.id=($2);
	`
	rows, err := tx.Query(query, userID, mirrorId)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Parse the files
	files := []File{}
	defer rows.Close()
	for rows.Next() {

		// Scan the results into the appropriate variables.
		f := File{}
		if err := rows.Scan(&f.ID, &f.Name, &f.SizeBytes, &f.UploadDate); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// If there are no errors, update the temp values and append to array
		files = append(files, f)
	}
	return files, nil
}

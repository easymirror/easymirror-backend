package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/easymirror/easymirror-backend/internal/common"
	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/google/uuid"
)

const (
	query_limit = 25
)

type MirrorLink struct {
	ID         uuid.UUID `json:"id"`
	Nickname   string    `json:"name"`
	UploadDate time.Time `json:"upload_date"`
	DurationMS int64     `json:"duration"`
}

type File struct {
	ID         uuid.UUID `json:"id"`          // ID of the file
	Name       string    `json:"name"`        // Name of the file
	SizeBytes  int64     `json:"size"`        // Size of the file in bytes
	UploadDate time.Time `json:"upload_date"` // Date the file was uploaded
}

// Returns a list of items a user has uploaded
func (u user) MirrorLinks(ctx context.Context, db *db.Database, pageNum int) ([]MirrorLink, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	// Get links from database
	tx, err := db.PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("BeginTx error: %w", err)
	}
	rows, err := tx.Query(`
		SELECT "id","nickname", "upload_date", "duration_ms" from mirroring_links
		WHERE created_by_id=($1)
		LIMIT ($2)
		OFFSET ($3)
	`, u.ID(), query_limit, common.GetPageOffset(query_limit, pageNum))
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Parse links
	links := []MirrorLink{}
	defer rows.Close()
	for rows.Next() {

		// Scan the results into the appropriate variables.
		// We want to use temporary sql null variables since some values can be null.
		var link MirrorLink
		var tempName sql.NullString
		var tempDate sql.NullTime
		var tempDuration sql.NullInt64
		if err := rows.Scan(&link.ID, &tempName, &tempDate, &tempDuration); err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		// If there are no errors, update the temp values and append to array
		link.UploadDate = tempDate.Time
		link.DurationMS = tempDuration.Int64
		link.Nickname = tempName.String
		links = append(links, link)
	}

	// Return
	return links, nil
}

func (u user) UpdateMirrorLinkName(
	ctx context.Context,
	db *db.Database,
	linkID, name string,
) error {
	if db == nil {
		return errors.New("database is nil")
	}

	// Begin tx
	tx, err := db.PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("BeginTx error: %w", err)
	}
	_, err = tx.Exec(`
		UPDATE mirroring_links
		SET nickname = ($1)
		WHERE id = ($2)
		AND created_by_id = ($3);	
	`, name, linkID, u.ID())

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("exec error: %w", err)
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("Commit error: %w", err)
	}
	return nil
}

func (u user) DeleteMirrorLink(ctx context.Context, db *db.Database, linkID string) error {
	if db == nil {
		return errors.New("database is nil")
	}

	// Begin TX
	tx, err := db.PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("BeginTx error: %w", err)
	}

	// Delete Link
	statement := `
		DELETE FROM mirroring_links
		WHERE id = ($1)
		AND
		created_by_id = ($2);
	`
	_, err = tx.Exec(statement, linkID, u.ID())
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error executing tx: %w", err)
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("error comitting tx: %w", err)
	}
	return nil
}

func (u user) GetFiles(ctx context.Context, db *db.Database, mirrorLinkID string) ([]File, error) {
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
	rows, err := tx.Query(query, u.ID(), mirrorLinkID)
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

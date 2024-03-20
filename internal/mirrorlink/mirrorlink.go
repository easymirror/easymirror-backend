package mirrorlink

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

// MirrorLink represents a mirror link, but with many details omitted.
type MirrorLink struct {
	ID         uuid.UUID `json:"id"`
	Nickname   string    `json:"name"`
	UploadDate time.Time `json:"upload_date"`
	DurationMS int64     `json:"duration"`
}

type ShareLink struct {
	MirrorLink           // Embed everything from the `MirrorLink`
	Links      HostLinks `json:"links"`
	Status     string    `json:"status"`
}

type HostLinks struct {
	Bunkr      *string `json:"bunkr"`
	Gofile     *string `json:"gofile"`
	Pixeldrain *string `json:"pixeldrain"`
	Cyberfile  *string `json:"cyberfile"`
	SaintTo    *string `json:"saint_to"`
	Cyberdrop  *string `json:"cyberdrop"`
}

// Returns a list of items a user has uploaded
func GetUserLinks(ctx context.Context, db *db.Database, userID string, pageNum int) ([]MirrorLink, error) {
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
	`, userID, query_limit, common.GetPageOffset(query_limit, pageNum))
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

// UpdateName update's the name of a given mirror link
func UpdateName(ctx context.Context, db *db.Database, mirrorID, userID, newName string) error {
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
	`, newName, mirrorID, userID)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("exec error: %w", err)
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return fmt.Errorf("commit error: %w", err)
	}
	return nil
}

// Delete delete's a given mirror link
func Delete(ctx context.Context, db *db.Database, userID, mirrorID string) error {
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
	_, err = tx.Exec(statement, mirrorID, userID)
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

// GetMirror returns all related data about a given mirror link.
func GetMirror(ctx context.Context, db *db.Database, mirrorID string) (*ShareLink, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	// Query the database
	query := `
		SELECT mirroring_links.nickname, mirroring_links.upload_date, host_links.*
		FROM mirroring_links
		RIGHT JOIN host_links ON mirroring_links.id = host_links.mirror_id
		WHERE mirroring_links.id=($1);
		`
	row := db.PostgresConn.QueryRowContext(ctx, query, mirrorID)

	// Parse the results
	// Because some values can be null, define temp null strings
	var (
		tmpName sql.NullString
		tmpDate sql.NullTime
	)

	sl := &ShareLink{}
	err := row.Scan(
		&tmpName,
		&tmpDate,
		&sl.ID,
		&sl.Links.Bunkr,
		&sl.Links.Gofile,
		&sl.Links.Pixeldrain,
		&sl.Links.Cyberfile,
		&sl.Links.SaintTo,
		&sl.Links.Cyberdrop,
	)

	// Validate
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("nothing found")
		case context.DeadlineExceeded:
			return nil, errors.New("db took too long")
		default:
			return nil, errors.New("not found")
			// response = map[string]any{"success": false, "error": "not_found"}
			// return c.JSON(http.StatusNotFound, response)
		}
	}
	sl.Nickname = tmpName.String
	sl.UploadDate = tmpDate.Time

	// Return
	return sl, nil
}

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
	var links []MirrorLink
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

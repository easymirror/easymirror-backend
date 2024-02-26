package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/easymirror/easymirror-backend/internal/db"
)

// Struct containing user info
type Info struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Username    string    `json:"username"`
	MemberSince time.Time `json:"member_since"`
	NextRenew   time.Time `json:"next_renewal"`
}

// Returns user info
func (u user) Info(ctx context.Context, db *db.Database) (*Info, error) {
	if db == nil {
		return nil, errors.New("database is nil")
	}

	// Get info fdrom database
	tx, err := db.PostgresConn.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("BeginTx error: %w", err)
	}
	rows, err := tx.Query(`
		SELECT "first_name", "last_name", "email", "phone", "username", "member_since", "next_renewal" from users
		WHERE id=($1);
	`, u.ID())
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	// Parse the info
	defer rows.Close()
	info := &Info{ID: u.ID().String()}
	for rows.Next() {
		if err = rows.Scan(&info.FirstName, &info.LastName, &info.Email, &info.Phone, &info.Username, &info.MemberSince, &info.NextRenew); err != nil {
			log.Println("Error scanning row:", err)
		}
	}
	return info, nil
}

package user

import (
	"fmt"

	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/google/uuid"
)

type User interface {
	ID() uuid.UUID // returns the ID of the user
}

type user struct {
	id uuid.UUID
}

func (u user) ID() uuid.UUID {
	return u.id
}

// Create creates and registers a new user in the database and returns a User object
func Create(db *db.Database) (User, error) {
	// Create a new user
	user := newUser()

	// Save to database
	tx, err := db.PostgresConn.Begin()
	if err != nil {
		return nil, fmt.Errorf("db.PostgresConn.Begin error: %w", err)
	}
	_, err = tx.Exec(`
	INSERT INTO users (id)
	VALUES
	($1);
	`, user.ID())
	if err != nil {
		return nil, fmt.Errorf("error executing tx: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing tx: %w", err)
	}

	// Return
	return user, nil
}

// New creates and returns a new user object
func New() User {
	return newUser()
}

func newUser() User {
	return user{id: uuid.New()}
}

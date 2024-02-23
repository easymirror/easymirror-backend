package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type User interface {
	ID() uuid.UUID                                                                       // returns the ID of the user
	MirrorLinks(ctx context.Context, db *db.Database, pageNum int) ([]MirrorLink, error) // Returns a list of items a user has uploaded
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

// FromJWT converts JWT token to a User
func FromJWT(t *jwt.Token) (User, error) {
	if !t.Valid {
		return nil, errors.New("jwt token not valid")
	}

	// Get the user-id from the JWT
	userID, err := t.Claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("GetSubject error: %w", err)
	}

	// Convert to type User
	uId, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("parse uuid error: %w", err)
	}
	return &user{id: uId}, nil
}

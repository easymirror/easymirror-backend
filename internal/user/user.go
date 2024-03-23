package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/easymirror/easymirror-backend/internal/mirrorlink"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type User interface {
	ID() uuid.UUID                                                                                  // returns the ID of the user
	Info(ctx context.Context, db *db.Database) (*Info, error)                                       // Returns the user info
	MirrorLinks(ctx context.Context, db *db.Database, pageNum int) ([]mirrorlink.MirrorLink, error) // Returns a list of items a user has uploaded
	UpdateMirrorLinkName(ctx context.Context, db *db.Database, linkID, name string) error
	DeleteMirrorLink(ctx context.Context, db *db.Database, linkID string) error
	GetFiles(ctx context.Context, db *db.Database, linkID string) ([]mirrorlink.File, error)
	Update(ctx context.Context, db *db.Database, k InfoKey, newVal string) error
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
	INSERT INTO users (id, member_since)
	VALUES
	(($1), ($2));
	`, user.ID(), time.Now().UTC())
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

// FromEcho pulls a JWT token from an echo context.
func FromEcho(c echo.Context) (User, error) {
	token, ok := c.Get("jwt-token").(*jwt.Token) // by default token is stored under `jwt-token` key
	if !ok {
		log.Println("Error with JWT token.")
		return nil, c.String(http.StatusInternalServerError, "Internal server error")
	}

	user, err := FromJWT(token)
	if err != nil {
		log.Println("Error getting user from JWT:", err)
		return nil, c.String(http.StatusInternalServerError, "Internal server error")
	}
	return user, err
}

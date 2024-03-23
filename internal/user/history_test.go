package user

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no env file loaded.")
	}
}

func insertMirrorLink(u User, database *db.Database) {
	tx, err := database.PostgresConn.Begin()
	if err != nil {
		panic(err)
	}

	if _, err := tx.Exec(`INSERT INTO users (id) values (($1))`, u.ID()); err != nil {
		panic(err)
	}

	for i := 0; i < 100; i++ {
		_, err := tx.Exec(`
		INSERT INTO mirroring_links values (($1), ($2), ($3), ($4), ($5));
		`, uuid.NewString(), u.ID(), fmt.Sprintf("Link #%v", i), time.Now(), 100)
		if err != nil {
			panic(err)
		}
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		panic(err)
	}
}

// go test -v -timeout 30s -run ^TestGetMirrorLinks$ github.com/easymirror/easymirror-backend/internal/user
func TestGetMirrorLinks(t *testing.T) {

	// start database
	database, err := db.InitDB()
	if err != nil {
		t.Fatalf("Error starting database: %v", err)
	}
	defer database.CloseConnections()

	user := newUser()

	// Insert 100 documents into `mirroring_links`
	insertMirrorLink(user, database)
	defer func() {
		// Delete everything from database
		tx, err := database.PostgresConn.Begin()
		if err != nil {
			t.Fatal(err)
		}

		tx.Exec(`DELETE FROM mirroring_links;`)
		tx.Commit()
	}()

	// run tests
	links, err := user.MirrorLinks(context.Background(), database, 2)
	if err != nil {
		t.Fatalf("Error getting mirror links: %v", err)
	}
	if len(links) < 1 {
		t.Fatal("No links")
	}
	for i, link := range links {
		fmt.Printf("#%v = %v\n", i+1, link)
	}

}

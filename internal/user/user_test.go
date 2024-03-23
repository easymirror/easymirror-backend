package user

import (
	"log"
	"testing"

	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no env file loaded.")
	}
}

// go test -v -timeout 30s -run ^TestCreate$ github.com/easymirror/easymirror-backend/internal/user
func TestCreate(t *testing.T) {
	database, err := db.InitDB()
	if err != nil {
		t.Fatal(err)
	}

	user, err := Create(database)
	if err != nil {
		t.Fatalf("Error creating user: %v", err)
	}

	// Drop values from table
	_, err = database.PostgresConn.Exec(`
	DELETE FROM users
WHERE id=($1);`, user.ID())
	if err != nil {
		t.Fatal(err)
	}

	database.PostgresConn.Close()
}

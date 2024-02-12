package main

import (
	"log"

	easymirrorbackend "github.com/easymirror/easymirror-backend/internal/api"
	"github.com/easymirror/easymirror-backend/internal/db"
)

func main() {

	// TODO initialize environment file

	// Initialize database(s)
	database, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	// Defer a function that will recover from any panic
	defer func() {
		// Perform graceful shutdowns here
		if r := recover(); r != nil {
			log.Println("[main] Recovered from panic:", r)
			database.CloseConnections()
		}
	}()

	// initialize API server
	easymirrorbackend.InitServer()
}

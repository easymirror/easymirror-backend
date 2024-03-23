package main

import (
	"log"

	easymirrorbackend "github.com/easymirror/easymirror-backend/internal/api"
	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/joho/godotenv"
)

func main() {
	// Load the env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("no env file loaded.")
	}

	// Initialize database(s)
	log.Println("Initializing database...")
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
	log.Println("Starting api...")
	easymirrorbackend.InitServer(database)
}

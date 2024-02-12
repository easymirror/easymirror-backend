package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

// initPostgres initializes a connection to out postgres database
func initPostgres() (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("PGSQL_HOST"),
		os.Getenv("PGSQL_PORT"),
		os.Getenv("PGSQL_USERNAME"),
		os.Getenv("PGSQL_PASSWORD"),
		os.Getenv("PGSQL_DB_NAME"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("initPostgres error: %w", err)
	}

	db.SetMaxIdleConns(25)
	db.SetMaxOpenConns(25)
	// db.SetConnMaxLifetime(5 * time.Minute) // TODO debate on using this
	return db, err
}

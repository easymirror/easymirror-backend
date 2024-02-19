package db

import (
	"database/sql"
	"fmt"
)

type Database struct {
	PostgresConn *sql.DB
}

// InitDB initializes all of the databases
func InitDB() (*Database, error) {
	// Initialize PostgreSQL
	pqsql, err := initPostgres()
	if err != nil {
		return nil, fmt.Errorf("initPostgres error: %w", err)
	}
	// TODO Initialize MongoDB

	return &Database{PostgresConn: pqsql}, nil
}

// CloseConnections closes all underlying connections to the database
func (db *Database) CloseConnections() {
	// Close PostgreSQL connection
	db.PostgresConn.Close()

	// TODO Close MongoDB connection
}

package db

import (
	"database/sql"
	"fmt"
)

type Database struct {
	postgresConn *sql.DB
}

// InitDB initializes all of the databases
func InitDB() (*Database, error) {
	// Initialize PostgreSQL
	pqsql, err := initPostgres()
	if err != nil {
		return nil, fmt.Errorf("initPostgres error: %w", err)
	}
	// TODO Initialize MongoDB

	return &Database{postgresConn: pqsql}, nil
}

// CloseConnections closes all underlying connections to the database
func (db *Database) CloseConnections() {
	// Close PostgreSQL connection
	db.postgresConn.Close()

	// TODO Close MongoDB connection
}

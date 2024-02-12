package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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

	// Create tables
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("initPostgres error: %w", err)
	}

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users ( id uuid NOT NULL, first_name character varying(30), last_name character varying(30), email text, phone character varying(15), password text, username character varying(60), member_since date, next_renewal date, PRIMARY KEY (id), CONSTRAINT username UNIQUE (username) );`,
		`CREATE TABLE IF NOT EXISTS mirroring_links ( id uuid NOT NULL, created_by_id uuid NOT NULL, nickname character varying(60), upload_date date, duration_ms bigint, PRIMARY KEY (id), CONSTRAINT created_by_id FOREIGN KEY (created_by_id) REFERENCES public.users (id) MATCH SIMPLE ON UPDATE CASCADE ON DELETE NO ACTION NOT VALID );`,
		`CREATE TABLE IF NOT EXISTS files ( id uuid NOT NULL, name text NOT NULL, size_bytes bigint NOT NULL, upload_date date NOT NULL, mirror_link_id uuid, PRIMARY KEY (id), CONSTRAINT mirror_link_id FOREIGN KEY (mirror_link_id) REFERENCES public.mirroring_links (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE SET NULL NOT VALID );`,
		`CREATE TABLE IF NOT EXISTS host_links ( mirror_id uuid NOT NULL, bunkr text, gofile text, pixeldrain text, cyberfile text, saint_to text, cyberdrop text, CONSTRAINT mirror_id FOREIGN KEY (mirror_id) REFERENCES public.mirroring_links (id) MATCH SIMPLE ON UPDATE NO ACTION ON DELETE NO ACTION NOT VALID );`,
	}
	for _, query := range queries {
		_, err = tx.ExecContext(ctx, query)
		if err != nil {
			tx.Rollback()
			log.Println("Error creating tables:", err)
			return nil, err
		}
	}
	return db, tx.Commit()
}

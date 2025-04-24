package storage

import (
	"database/sql"
	"errors"
	"fmt"
)

// add more errors
var (
	ErrURLNotFound = errors.New("url not fount")
	ErrURLExists = errors.New("url exists")
)

type Storage struct {
	db *sql.DB 
}

func New(DSN string) (*Storage, error){
	const op = "storage.New"
	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil , fmt.Errorf("%s: %w", op, err)
	}
	EsureSchema(db)
	return &Storage{db}, nil
}

func EsureSchema(db *sql.DB) error {
 stmts := []string{
    `CREATE TABLE IF NOT EXISTS urls (
       id UUID PRIMARY KEY,
       address TEXT NOT NULL UNIQUE,
       check_interval INTERVAL NOT NULL,
       created_at TIMESTAMPTZ NOT NULL DEFAULT now()
     );`,
    `CREATE TABLE IF NOT EXISTS checks (
       id UUID PRIMARY KEY,
       url_id UUID NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
       status BOOLEAN NOT NULL,
       code INT,
       duration INTERVAL NOT NULL,
       checked_at TIMESTAMPTZ NOT NULL DEFAULT now()
     );`,
  }
  for i, stmt := range stmts {
    if _, err := db.Exec(stmt); err != nil {
      return fmt.Errorf("schema init [%d] failed: %w", i, err)
    }
  }
  return nil
}
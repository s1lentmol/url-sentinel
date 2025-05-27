package postgres

import (
	"context"
	"database/sql"
	"fmt"
)

// migrations contains all database migrations
var migrations = []string{
	// 001_init.sql
	`
-- Create URLs table
CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY,
    address TEXT NOT NULL UNIQUE,
    check_interval INTERVAL NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index on address for faster lookups
CREATE INDEX IF NOT EXISTS idx_urls_address ON urls(address);

-- Create checks table
CREATE TABLE IF NOT EXISTS checks (
    id UUID PRIMARY KEY,
    url_id UUID NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
    status BOOLEAN NOT NULL,
    code INT,
    duration INTERVAL NOT NULL,
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create indexes for checks
CREATE INDEX IF NOT EXISTS idx_checks_url_id ON checks(url_id);
CREATE INDEX IF NOT EXISTS idx_checks_checked_at ON checks(checked_at DESC);
	`,
}

// RunMigrations executes all SQL migrations in order
func RunMigrations(db *sql.DB) error {
	ctx := context.Background()

	for i, migration := range migrations {
		if _, err := db.ExecContext(ctx, migration); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", i+1, err)
		}
	}

	return nil
}

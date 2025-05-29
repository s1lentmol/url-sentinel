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

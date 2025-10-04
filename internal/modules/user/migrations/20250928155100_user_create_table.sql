-- +goose Up
-- +goose StatementBegin
CREATE TABLE users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Email normalization (critical for security)
    email VARCHAR(255) NOT NULL UNIQUE,
    email_normalized VARCHAR(255) GENERATED ALWAYS AS (LOWER(email)) STORED UNIQUE,
    -- Password security
    password_hash VARCHAR(255) NOT NULL,
    -- Rename for clarity
    password_changed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- Personal info
    name VARCHAR(255) NOT NULL DEFAULT '',
    date_of_birth DATE,
    phone VARCHAR(20),
    -- Keep NOT NULL for consistency
    -- Verification & security
    verified_email BOOLEAN NOT NULL DEFAULT FALSE,
    -- Account status (prevent hard deletes)
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    onboarded BOOLEAN NOT NULL DEFAULT FALSE,
    deactivated_at TIMESTAMP,
    -- Timestamps (with ON UPDATE for updated_at)
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Optimistic locking
    version INT NOT NULL DEFAULT 1,
    -- Security/compliance
    last_login_at TIMESTAMP,
    failed_login_attempts INT NOT NULL DEFAULT 0,
    locked_until TIMESTAMP
);
-- Index for case-insensitive email lookups
CREATE INDEX idx_users_email_normalized ON users (email_normalized);
-- Index for active users (common query pattern)
CREATE INDEX idx_users_active ON users (is_active)
WHERE is_active = TRUE;
-- Trigger to auto-update updated_at
CREATE OR REPLACE FUNCTION update_users_version_and_timestamp() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = CURRENT_TIMESTAMP;
NEW.version = OLD.version + 1;
IF OLD.password_hash IS DISTINCT
FROM NEW.password_hash THEN NEW.password_changed_at = CURRENT_TIMESTAMP;
END IF;
RETURN NEW;
END;
$$ LANGUAGE 'plpgsql';
-- 2. Attach the trigger to the users table
DROP TRIGGER IF EXISTS update_users_version ON users;
CREATE TRIGGER update_users_version BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_users_version_and_timestamp();
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP FUNCTION IF EXISTS update_users_version_and_timestamp();
-- +goose StatementEnd
-- +goose Up
-- +goose StatementBegin
CREATE TYPE token_type AS ENUM (
    'email_verify',
    'phone_verify',
    'password_reset'
);
CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(60) NOT NULL,
    -- 6-digit OTP or UUID for email links
    type token_type NOT NULL DEFAULT 'email_verify',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Ensure only one active token per type per user
    CONSTRAINT unique_active_token_per_user_type UNIQUE (user_id, type) DEFERRABLE INITIALLY DEFERRED
);
-- Indexes for performance
CREATE INDEX tokens_user_id ON tokens (user_id);
CREATE INDEX tokens_token ON tokens (token);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tokens;
DROP TYPE IF EXISTS token_type;
-- +goose StatementEnd
-- name: GetUserById :one
SELECT id,
    email,
    name,
    verified_email,
    is_active,
    password_hash,
    updated_at,
    last_login_at,
    locked_until,
    version
FROM users
WHERE id = $1
LIMIT 1;
-- name: GetUserByEmail :one
SELECT id,
    email,
    name,
    verified_email,
    password_hash,
    is_active,
    updated_at,
    last_login_at,
    locked_until,
    version
FROM users
WHERE email_normalized = LOWER($1)
LIMIT 1;
-- name: GetAllUsers :many
SELECT *
FROM users
ORDER BY name;
-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES ($1, $2)
RETURNING id,
    email;
-- name: DeleteUser :exec
UPDATE users
SET is_active = FALSE,
    deactivated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING id,
    is_active,
    deactivated_at,
    version;
-- name: GetToken :one
SELECT id,
    user_id,
    token,
    type,
    expires_at,
    created_at
FROM tokens
WHERE token = $1
    AND type = $2
    AND user_id = $3
LIMIT 1;
-- name: GetTokenByUserId :one
SELECT id,
    user_id,
    token,
    type,
    expires_at,
    created_at
FROM tokens
WHERE user_id = $1;
-- name: GetTokenByCode :one
SELECT id,
    user_id,
    token,
    type,
    expires_at,
    created_at
FROM tokens
WHERE token = $1;
-- name: CreateToken :one
INSERT INTO tokens (user_id, token, type, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING id,
    user_id,
    token,
    type,
    expires_at,
    created_at;
-- name: DeleteToken :exec
DELETE FROM tokens
WHERE id = $1;
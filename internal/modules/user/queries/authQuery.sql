-- name: VerifyUserEmail :exec
UPDATE users
SET verified_email = true
WHERE id = $1;
-- name: UpdatePassword :exec
UPDATE users
SET password_hash = $1
WHERE id = $2;
-- name: ResetFailedLogin :exec
UPDATE users
SET last_login_at = CURRENT_TIMESTAMP,
    failed_login_attempts = 0,
    locked_until = NULL
WHERE id = $1;
-- name: IncrementFailedLogin :exec
UPDATE users
SET failed_login_attempts = failed_login_attempts + 1,
    locked_until = CASE
        WHEN failed_login_attempts + 1 >= 5 THEN CURRENT_TIMESTAMP + INTERVAL '15 minutes'
        ELSE locked_until
    END
WHERE id = $1;
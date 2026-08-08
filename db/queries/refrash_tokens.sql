-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id,
    token_hash,
    expires_at
)
VALUES ($1, $2, $3)
RETURNING id;

-- name: GetTokenHash :one
SELECT
    id,
    user_id,
    token_hash,
    created_at,
    expires_at,
    revoked
FROM refresh_tokens
WHERE token_hash = $1
  AND revoked = FALSE
  AND expires_at > NOW();

-- name: GetTokenUserID :one
SELECT
    id,
    token_hash,
    created_at,
    expires_at,
    revoked
FROM refresh_tokens
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteTokenHash :exec
DELETE FROM refresh_tokens
WHERE token_hash = $1;

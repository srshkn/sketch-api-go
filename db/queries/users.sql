-- name: CreateUser :one
INSERT INTO users (
    username,
    email,
    password_hash
)
VALUES (
    $1,
    $2,
    $3
)
RETURNING id, username, email;

-- name: GetUserByEmail :one
SELECT id, username, email, password_hash
FROM users
WHERE LOWER(email) = LOWER($1);

-- name: GetUserByEmailOrUsername :one
SELECT id, username, email
FROM users
WHERE username = $1 OR LOWER(email) = LOWER($2)
LIMIT 1;

-- name: GetUserByID :one
SELECT id, username, email
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, username, email
FROM users
ORDER BY username;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

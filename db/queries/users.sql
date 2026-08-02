-- name: CreateUser :one
INSERT INTO users (
    name,
    password_hash
)
VALUES ($1, $2)
RETURNING id, name;

-- name: GetUser :one
SELECT id, name
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, name
FROM users
ORDER BY id;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

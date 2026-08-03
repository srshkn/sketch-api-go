-- name: CreateUser :one
INSERT INTO users (
    name,
    email,
    password_hash
)
VALUES ($1, $2, $3)
RETURNING id, name, email;

-- name: GetUserByEmail :one
SELECT id, name, email, password_hash
FROM users
WHERE LOWER(email) = LOWER($1);

-- name: ByName :one
SELECT id, name
FROM users
WHERE name = $1;

-- name: ByEmail :one
SELECT id, email
FROM users
WHERE LOWER(email) = LOWER($1);

-- name: ListUsers :many
SELECT id, name, email
FROM users
ORDER BY name;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

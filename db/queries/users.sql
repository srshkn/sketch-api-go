-- name: CheckUserExists :one
SELECT id, username, email
FROM users
WHERE username = $1 OR email = $2
LIMIT 1;

-- name: CreateUser :one
WITH new_user AS (
    INSERT INTO users (
        username,
        email,
        password_hash,
        role_id
    )
    VALUES (
        $1,
        $2,
        $3,
        (SELECT id FROM roles WHERE name = 'user')
    )
    RETURNING id, username, email, role_id
)
SELECT
    u.id,
    u.username,
    u.email,
    r.name AS role
FROM new_user u
JOIN roles r ON r.id = u.role_id;

-- name: GetUserByLogin :one
SELECT
    u.id,
    u.email,
    u.password_hash,
    r.name AS role
FROM users AS u
LEFT JOIN roles AS r ON r.id = u.role_id
WHERE u.email = $1;

-- name: GetUserByID :one
SELECT
    u.id,
    u.username,
    u.email,
    r.name AS role
FROM users AS u
LEFT JOIN roles AS r ON r.id = u.role_id
WHERE u.id = $1;
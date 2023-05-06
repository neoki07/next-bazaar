-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  hashed_password
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: TruncateUsersTable :exec
TRUNCATE TABLE users CASCADE;

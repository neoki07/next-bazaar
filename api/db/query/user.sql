-- name: CreateUser :one
INSERT INTO users (
  name,
  email,
  hashed_password
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
  name = sqlc.arg('name'),
  email = sqlc.arg('email'),
  hashed_password = sqlc.arg('hashed_password')
WHERE id = $1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: GetUsersByIDs :many
SELECT * FROM users
WHERE id = ANY((sqlc.arg('ids'))::uuid[])
ORDER BY id;

-- name: TruncateUsersTable :exec
TRUNCATE TABLE users CASCADE;

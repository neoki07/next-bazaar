-- name: CreateSession :one
INSERT INTO sessions (
  user_id,
  session_token,
  expired_at
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE session_token = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_token = $1;

-- name: TruncateSessionsTable :exec
TRUNCATE TABLE sessions CASCADE;

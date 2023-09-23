-- name: CreateSession :one
INSERT INTO sessions (
  user_id,
  session_token,
  session_token_expired_at,
  refresh_token,
  refresh_token_expired_at
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE session_token = $1 LIMIT 1;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE session_token = $1;

-- name: TruncateSessionsTable :exec
TRUNCATE TABLE sessions CASCADE;

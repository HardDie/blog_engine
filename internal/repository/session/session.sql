-- name: CreateOrUpdate :one
INSERT INTO sessions (user_id, session_hash)
VALUES (?, ?)
ON CONFLICT (user_id) DO UPDATE
SET session_hash = excluded.session_hash, updated_at = datetime('now'), deleted_at = NULL
RETURNING *;

-- name: GetByUserID :one
SELECT *
FROM sessions
WHERE session_hash = ?
  AND deleted_at IS NULL;

-- name: DeleteByID :exec
UPDATE sessions
SET deleted_at = datetime('now')
WHERE id = ?
  AND deleted_at IS NULL;
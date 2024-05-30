-- name: Create :one
INSERT INTO passwords (user_id, password_hash)
VALUES (?, ?)
RETURNING *;

-- name: GetByUserID :one
SELECT *
FROM passwords
WHERE user_id = ?
  AND deleted_at IS NULL;

-- name: Update :one
UPDATE passwords
SET password_hash = ?, updated_at = datetime('now')
WHERE id = ?
  AND deleted_at IS NULL
RETURNING *;

-- name: IncreaseFailedAttempts :one
UPDATE passwords
SET failed_attempts = failed_attempts + 1, updated_at = datetime('now')
WHERE id = ?
  AND deleted_at IS NULL
RETURNING *;

-- name: ResetFailedAttempts :one
UPDATE passwords
SET failed_attempts = 0, updated_at = datetime('now')
WHERE id = ?
  AND deleted_at IS NULL
RETURNING *;

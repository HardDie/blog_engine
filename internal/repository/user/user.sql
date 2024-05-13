-- name: GetByIDPublic :one
SELECT id, displayed_name, invited_by_user, created_at, updated_at, deleted_at
FROM users
WHERE id = ?
  AND deleted_at IS NULL;

-- name: GetByIDPrivate :one
SELECT *
FROM users
WHERE id = ?
  AND deleted_at IS NULL;

-- name: GetByName :one
SELECT *
FROM users
WHERE username = ?
  AND deleted_at IS NULL;

-- name: Create :one
INSERT INTO users (username, displayed_name, invited_by_user)
VALUES (?, ?, ?)
RETURNING *;

-- name: Update :one
UPDATE users
SET displayed_name = ?, email = ?, updated_at = datetime('now')
WHERE id = ?
  AND deleted_at IS NULL
RETURNING *;

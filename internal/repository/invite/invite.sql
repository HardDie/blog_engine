-- name: GetByID :one
SELECT *
FROM invites
WHERE id = ?
  AND deleted_at IS NULL;

-- name: GetActiveByUserID :one
SELECT *
FROM invites
WHERE id = ?
  AND is_activated IS FALSE
  AND deleted_at IS NULL;

-- name: GetAllByUserID :many
SELECT *
FROM invites
WHERE user_id = ?
  AND deleted_at IS NULL;

-- name: GetByInviteHash :one
SELECT *
FROM invites
WHERE invite_hash = ?
  AND is_activated IS FALSE
  AND deleted_at IS NULL;

-- name: CreateOrUpdate :one
INSERT INTO invites (user_id, invite_hash, is_activated)
VALUES (?, ?, false)
ON CONFLICT (user_id, is_activated) WHERE is_activated IS FALSE DO UPDATE
SET invite_hash = excluded.invite_hash, updated_at = datetime('now')
RETURNING *;

-- name: Delete :exec
UPDATE invites
SET deleted_at = datetime('now'), is_activated = true
WHERE id = ?
  AND deleted_at IS NULL;

-- name: Activate :one
UPDATE invites
SET is_activated = true
WHERE id = ?
  AND is_activated IS FALSE
  AND deleted_at IS NULL
RETURNING *;
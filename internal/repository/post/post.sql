-- name: List :many
SELECT sqlc.embed(posts), count(*) over()
FROM posts
WHERE deleted_at IS NULL
  AND CASE WHEN CAST(sqlc.arg(display_only_published) AS boolean) IS TRUE THEN is_published IS true ELSE true END
  AND CASE WHEN CAST(sqlc.arg(query) AS text) <> '' THEN lower(title) LIKE sqlc.arg(query) ELSE true END
  AND CASE WHEN CAST(sqlc.arg(related_to_user) AS int) > 0 THEN user_id = sqlc.arg(related_to_user) ELSE true END
ORDER BY id DESC
LIMIT CASE WHEN CAST(sqlc.arg(limit) AS int) > 0 THEN sqlc.arg(limit) ELSE 10 END
OFFSET sqlc.arg(offset);

-- name: Create :one
INSERT INTO posts (user_id, title, short, body, tags, is_published)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: Edit :one
UPDATE posts
SET title = ?, short = ?, body = ?, tags = ?, is_published = ?, updated_at = datetime('now')
WHERE id = ?
  AND deleted_at IS NULL
  AND user_id = ?
RETURNING *;

-- name: GetByID :one
SELECT *
FROM posts
WHERE deleted_at IS NULL
  AND id = ?
  AND CASE WHEN CAST(sqlc.narg(user_id) AS int) IS NULL THEN is_published IS TRUE ELSE user_id = sqlc.narg(user_id) END;
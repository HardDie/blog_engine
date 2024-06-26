// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: post.sql

package post

import (
	"context"
	"database/sql"
)

const create = `-- name: Create :one
INSERT INTO posts (user_id, title, short, body, tags, is_published)
VALUES (?, ?, ?, ?, ?, ?)
RETURNING id, user_id, title, short, body, tags, is_published, created_at, updated_at, deleted_at
`

type CreateParams struct {
	UserID      int64          `json:"userId"`
	Title       string         `json:"title"`
	Short       string         `json:"short"`
	Body        string         `json:"body"`
	Tags        sql.NullString `json:"tags"`
	IsPublished bool           `json:"isPublished"`
}

// Create
//
//	INSERT INTO posts (user_id, title, short, body, tags, is_published)
//	VALUES (?, ?, ?, ?, ?, ?)
//	RETURNING id, user_id, title, short, body, tags, is_published, created_at, updated_at, deleted_at
func (q *Queries) Create(ctx context.Context, arg CreateParams) (*Post, error) {
	row := q.queryRow(ctx, q.createStmt, create,
		arg.UserID,
		arg.Title,
		arg.Short,
		arg.Body,
		arg.Tags,
		arg.IsPublished,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Short,
		&i.Body,
		&i.Tags,
		&i.IsPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const edit = `-- name: Edit :one
UPDATE posts
SET title = ?, short = ?, body = ?, tags = ?, is_published = ?, updated_at = datetime('now')
WHERE id = ?
  AND deleted_at IS NULL
  AND user_id = ?
RETURNING id, user_id, title, short, body, tags, is_published, created_at, updated_at, deleted_at
`

type EditParams struct {
	Title       string         `json:"title"`
	Short       string         `json:"short"`
	Body        string         `json:"body"`
	Tags        sql.NullString `json:"tags"`
	IsPublished bool           `json:"isPublished"`
	ID          int64          `json:"id"`
	UserID      int64          `json:"userId"`
}

// Edit
//
//	UPDATE posts
//	SET title = ?, short = ?, body = ?, tags = ?, is_published = ?, updated_at = datetime('now')
//	WHERE id = ?
//	  AND deleted_at IS NULL
//	  AND user_id = ?
//	RETURNING id, user_id, title, short, body, tags, is_published, created_at, updated_at, deleted_at
func (q *Queries) Edit(ctx context.Context, arg EditParams) (*Post, error) {
	row := q.queryRow(ctx, q.editStmt, edit,
		arg.Title,
		arg.Short,
		arg.Body,
		arg.Tags,
		arg.IsPublished,
		arg.ID,
		arg.UserID,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Short,
		&i.Body,
		&i.Tags,
		&i.IsPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const getByID = `-- name: GetByID :one
SELECT id, user_id, title, short, body, tags, is_published, created_at, updated_at, deleted_at
FROM posts
WHERE deleted_at IS NULL
  AND id = ?
  AND CASE WHEN CAST(?2 AS int) IS NULL THEN is_published IS TRUE ELSE user_id = ?2 END
`

type GetByIDParams struct {
	ID     int64         `json:"id"`
	UserID sql.NullInt64 `json:"userId"`
}

// GetByID
//
//	SELECT id, user_id, title, short, body, tags, is_published, created_at, updated_at, deleted_at
//	FROM posts
//	WHERE deleted_at IS NULL
//	  AND id = ?
//	  AND CASE WHEN CAST(?2 AS int) IS NULL THEN is_published IS TRUE ELSE user_id = ?2 END
func (q *Queries) GetByID(ctx context.Context, arg GetByIDParams) (*Post, error) {
	row := q.queryRow(ctx, q.getByIDStmt, getByID, arg.ID, arg.UserID)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Short,
		&i.Body,
		&i.Tags,
		&i.IsPublished,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return &i, err
}

const list = `-- name: List :many
SELECT posts.id, posts.user_id, posts.title, posts.short, posts.body, posts.tags, posts.is_published, posts.created_at, posts.updated_at, posts.deleted_at, count(*) over()
FROM posts
WHERE deleted_at IS NULL
  AND CASE WHEN CAST(?1 AS boolean) IS TRUE THEN is_published IS true ELSE true END
  AND CASE WHEN CAST(?2 AS text) <> '' THEN lower(title) LIKE ?2 ELSE true END
  AND CASE WHEN CAST(?3 AS int) > 0 THEN user_id = ?3 ELSE true END
ORDER BY id DESC
LIMIT CASE WHEN CAST(?5 AS int) > 0 THEN ?5 ELSE 10 END
OFFSET ?4
`

type ListParams struct {
	DisplayOnlyPublished bool   `json:"displayOnlyPublished"`
	Query                string `json:"query"`
	RelatedToUser        int64  `json:"relatedToUser"`
	Offset               int64  `json:"offset"`
	Limit                int64  `json:"limit"`
}

type ListRow struct {
	Post  Post  `json:"post"`
	Count int64 `json:"count"`
}

// List
//
//	SELECT posts.id, posts.user_id, posts.title, posts.short, posts.body, posts.tags, posts.is_published, posts.created_at, posts.updated_at, posts.deleted_at, count(*) over()
//	FROM posts
//	WHERE deleted_at IS NULL
//	  AND CASE WHEN CAST(?1 AS boolean) IS TRUE THEN is_published IS true ELSE true END
//	  AND CASE WHEN CAST(?2 AS text) <> '' THEN lower(title) LIKE ?2 ELSE true END
//	  AND CASE WHEN CAST(?3 AS int) > 0 THEN user_id = ?3 ELSE true END
//	ORDER BY id DESC
//	LIMIT CASE WHEN CAST(?5 AS int) > 0 THEN ?5 ELSE 10 END
//	OFFSET ?4
func (q *Queries) List(ctx context.Context, arg ListParams) ([]*ListRow, error) {
	rows, err := q.query(ctx, q.listStmt, list,
		arg.DisplayOnlyPublished,
		arg.Query,
		arg.RelatedToUser,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*ListRow{}
	for rows.Next() {
		var i ListRow
		if err := rows.Scan(
			&i.Post.ID,
			&i.Post.UserID,
			&i.Post.Title,
			&i.Post.Short,
			&i.Post.Body,
			&i.Post.Tags,
			&i.Post.IsPublished,
			&i.Post.CreatedAt,
			&i.Post.UpdatedAt,
			&i.Post.DeletedAt,
			&i.Count,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

package repository

import (
	"strings"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IPost interface {
	List(filter *dto.ListPostFilter) ([]*entity.Post, int32, error)
	Create(req *dto.CreatePostDTO, userID int32) (*entity.Post, error)
	Edit(req *dto.EditPostDTO, userID int32) (*entity.Post, error)
}

type Post struct {
	db *db.DB
}

func NewPost(db *db.DB) *Post {
	return &Post{
		db: db,
	}
}

func (r *Post) List(filter *dto.ListPostFilter) ([]*entity.Post, int32, error) {
	var res []*entity.Post
	var total int32

	q := gosql.NewSelect().From("posts")
	q.Columns().Add("id", "user_id", "title", "short", "tags", "created_at", "updated_at", "count(*) over()")
	q.Where().AddExpression("deleted_at IS NULL")
	if filter.DisplayOnlyPublished {
		q.Where().AddExpression("is_published IS true")
	}
	if filter.Query != "" {
		q.Where().AddExpression("lower(title) LIKE ?", utils.PrepareStringToLike(filter.Query))
	}
	if filter.RelatedToUser > 0 {
		q.Where().AddExpression("user_id = ?", filter.RelatedToUser)
	}
	if filter.Limit > 0 {
		q.SetPagination(utils.GetPagination(filter.Limit, filter.Page))
	}
	q.AddOrder("id DESC")
	rows, err := r.db.DB.Query(q.String(), q.GetArguments()...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		post := &entity.Post{}
		var tags string

		err = rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Short, &tags, &post.CreatedAt, &post.UpdatedAt, &total)
		if err != nil {
			return nil, 0, err
		}
		post.Tags = strings.Split(tags, ";")
		res = append(res, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}
func (r *Post) Create(req *dto.CreatePostDTO, userID int32) (*entity.Post, error) {
	post := &entity.Post{
		UserID:      userID,
		Title:       req.Title,
		Short:       req.Short,
		Body:        req.Body,
		Tags:        req.Tags,
		IsPublished: req.IsPublished,
	}
	tags := strings.Join(req.Tags, ";")

	q := gosql.NewInsert().Into("posts")
	q.Columns().Add("user_id", "title", "short", "body", "tags", "is_published")
	q.Columns().Arg(userID, req.Title, req.Short, req.Body, tags, req.IsPublished)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}
func (r *Post) Edit(req *dto.EditPostDTO, userID int32) (*entity.Post, error) {
	post := &entity.Post{
		ID:          req.ID,
		UserID:      userID,
		Title:       req.Title,
		Short:       req.Short,
		Body:        req.Body,
		Tags:        req.Tags,
		IsPublished: req.IsPublished,
	}
	tags := strings.Join(req.Tags, ";")

	q := gosql.NewUpdate().Table("posts")
	q.Set().Append("title = ?", req.Title)
	q.Set().Append("short = ?", req.Short)
	q.Set().Append("body = ?", req.Body)
	q.Set().Append("tags = ?", tags)
	q.Set().Append("is_published = ?", req.IsPublished)
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", req.ID)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Where().AddExpression("user_id = ?", userID)
	q.Returning().Add("created_at", "updated_at")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return post, nil
}

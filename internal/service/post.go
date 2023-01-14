package service

import (
	"context"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/repository"
)

type IPost interface {
	Feed(ctx context.Context, req *dto.FeedPostDTO) ([]*entity.Post, int32, error)
	PublicGet(ctx context.Context, id int32) (*entity.Post, error)

	Create(ctx context.Context, req *dto.CreatePostDTO, userID int32) (*entity.Post, error)
	Edit(ctx context.Context, req *dto.EditPostDTO, userID int32) (*entity.Post, error)
	List(ctx context.Context, req *dto.ListPostDTO, userID int32) ([]*entity.Post, int32, error)
}

type Post struct {
	repository repository.IPost
}

func NewPost(repository repository.IPost) *Post {
	return &Post{
		repository: repository,
	}
}

func (p *Post) Feed(ctx context.Context, req *dto.FeedPostDTO) ([]*entity.Post, int32, error) {
	return p.repository.List(ctx, &dto.ListPostFilter{
		Limit:                req.Limit,
		Page:                 req.Page,
		Query:                req.Query,
		DisplayOnlyPublished: true,
	})
}
func (p *Post) PublicGet(ctx context.Context, id int32) (*entity.Post, error) {
	return p.repository.GetByID(ctx, id, nil)
}

func (p *Post) Create(ctx context.Context, req *dto.CreatePostDTO, userID int32) (*entity.Post, error) {
	return p.repository.Create(ctx, req, userID)
}
func (p *Post) Edit(ctx context.Context, req *dto.EditPostDTO, userID int32) (*entity.Post, error) {
	return p.repository.Edit(ctx, req, userID)
}
func (p *Post) List(ctx context.Context, req *dto.ListPostDTO, userID int32) ([]*entity.Post, int32, error) {
	return p.repository.List(ctx, &dto.ListPostFilter{
		Limit:                req.Limit,
		Page:                 req.Page,
		Query:                req.Query,
		RelatedToUser:        userID,
		DisplayOnlyPublished: false,
	})
}

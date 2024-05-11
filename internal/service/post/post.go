package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	repositoryPost "github.com/HardDie/blog_engine/internal/repository/post"
)

type IPost interface {
	Feed(ctx context.Context, req *dto.FeedPostDTO) ([]*entity.Post, int32, error)
	PublicGet(ctx context.Context, id int32) (*entity.Post, error)

	Create(ctx context.Context, req *dto.CreatePostDTO, userID int32) (*entity.Post, error)
	Edit(ctx context.Context, req *dto.EditPostDTO, userID int32) (*entity.Post, error)
	List(ctx context.Context, req *dto.ListPostDTO, userID int32) ([]*entity.Post, int32, error)
}

type Post struct {
	postRepository repositoryPost.IPost
}

func New(post repositoryPost.IPost) *Post {
	return &Post{
		postRepository: post,
	}
}

func (p *Post) Feed(ctx context.Context, req *dto.FeedPostDTO) ([]*entity.Post, int32, error) {
	posts, count, err := p.postRepository.List(ctx, &dto.ListPostFilter{
		Limit:                req.Limit,
		Page:                 req.Page,
		Query:                req.Query,
		DisplayOnlyPublished: true,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("Post.Feed() List: %w", err)
	}
	return posts, count, nil
}
func (p *Post) PublicGet(ctx context.Context, id int32) (*entity.Post, error) {
	resp, err := p.postRepository.GetByID(ctx, id, nil)
	if err != nil {
		switch {
		case errors.Is(err, repositoryPost.ErrorNotFound):
			return nil, ErrorPostNotFound
		}
		return nil, fmt.Errorf("Post.PublicGet() GetByID: %w", err)
	}
	return resp, nil
}
func (p *Post) Create(ctx context.Context, req *dto.CreatePostDTO, userID int32) (*entity.Post, error) {
	resp, err := p.postRepository.Create(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("Post.Create() Create: %w", err)
	}
	return resp, nil
}
func (p *Post) Edit(ctx context.Context, req *dto.EditPostDTO, userID int32) (*entity.Post, error) {
	resp, err := p.postRepository.Edit(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("Post.Edit() Edit: %w", err)
	}
	return resp, nil
}
func (p *Post) List(ctx context.Context, req *dto.ListPostDTO, userID int32) ([]*entity.Post, int32, error) {
	posts, count, err := p.postRepository.List(ctx, &dto.ListPostFilter{
		Limit:                req.Limit,
		Page:                 req.Page,
		Query:                req.Query,
		RelatedToUser:        userID,
		DisplayOnlyPublished: false,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("Post.List() List: %w", err)
	}
	return posts, count, nil
}

var (
	ErrorPostNotFound = errors.New("post not found")
)

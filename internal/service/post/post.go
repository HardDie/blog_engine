package post

import (
	"context"
	"errors"
	"fmt"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	repositoryPost "github.com/HardDie/blog_engine/internal/repository/post"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/user"
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
	userRepository repositoryUser.IUser
}

func New(post repositoryPost.IPost, user repositoryUser.IUser) *Post {
	return &Post{
		postRepository: post,
		userRepository: user,
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
	if posts == nil {
		posts = []*entity.Post{}
	}
	// Enrich every post with user info
	users := make(map[int32]*entity.User)
	for i, post := range posts {
		user, ok := users[post.UserID]
		if ok {
			posts[i].User = user
			continue
		}
		user, err := p.userRepository.GetByID(ctx, post.UserID, false)
		if err != nil {
			return nil, 0, fmt.Errorf("Post.Feed() user.GetByID: %w", err)
		}
		users[post.UserID] = user
		posts[i].User = user
	}
	return posts, count, nil
}
func (p *Post) PublicGet(ctx context.Context, id int32) (*entity.Post, error) {
	post, err := p.postRepository.GetByID(ctx, id, nil)
	if err != nil {
		switch {
		case errors.Is(err, repositoryPost.ErrorNotFound):
			return nil, ErrorPostNotFound
		}
		return nil, fmt.Errorf("Post.PublicGet() GetByID: %w", err)
	}
	user, err := p.userRepository.GetByID(ctx, post.UserID, false)
	if err != nil {
		return nil, fmt.Errorf("Post.PublicGet() user.GetByID: %w", err)
	}
	post.User = user
	return post, nil
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
	// Enrich every post with user info
	users := make(map[int32]*entity.User)
	for i, post := range posts {
		user, ok := users[post.UserID]
		if ok {
			posts[i].User = user
			continue
		}
		user, err := p.userRepository.GetByID(ctx, post.UserID, false)
		if err != nil {
			return nil, 0, fmt.Errorf("Post.List() user.GetByID: %w", err)
		}
		users[post.UserID] = user
		posts[i].User = user
	}
	return posts, count, nil
}

var (
	ErrorPostNotFound = errors.New("post not found")
)

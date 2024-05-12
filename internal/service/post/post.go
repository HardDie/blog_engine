package post

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	repositoryPost "github.com/HardDie/blog_engine/internal/repository/post"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/user"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IPost interface {
	Feed(ctx context.Context, req *dto.FeedPostDTO) ([]*entity.Post, int64, error)
	PublicGet(ctx context.Context, id int64) (*entity.Post, error)

	Create(ctx context.Context, req *dto.CreatePostDTO, userID int64) (*entity.Post, error)
	Edit(ctx context.Context, req *dto.EditPostDTO, userID int64) (*entity.Post, error)
	List(ctx context.Context, req *dto.ListPostDTO, userID int64) ([]*entity.Post, int64, error)
}

type Post struct {
	postRepository repositoryPost.Querier
	userRepository repositoryUser.IUser
}

func New(post repositoryPost.Querier, user repositoryUser.IUser) *Post {
	return &Post{
		postRepository: post,
		userRepository: user,
	}
}

func (p *Post) Feed(ctx context.Context, req *dto.FeedPostDTO) ([]*entity.Post, int64, error) {
	limit, offset := utils.GetPagination(req.Limit, req.Page)
	resp, err := p.postRepository.List(ctx, repositoryPost.ListParams{
		Limit:                limit,
		Offset:               offset,
		Query:                req.Query,
		DisplayOnlyPublished: true,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("Post.Feed() List: %w", err)
	}
	if len(resp) == 0 {
		return []*entity.Post{}, 0, nil
	}

	posts := make([]*entity.Post, 0, len(resp))
	users := make(map[int64]*entity.User)
	for _, el := range resp {
		post := &entity.Post{
			ID:          el.Post.ID,
			UserID:      el.Post.UserID,
			Title:       el.Post.Title,
			Short:       el.Post.Short,
			Body:        el.Post.Body,
			Tags:        strings.Split(el.Post.Tags.String, ";"),
			IsPublished: el.Post.IsPublished,
			CreatedAt:   el.Post.CreatedAt,
			UpdatedAt:   el.Post.UpdatedAt,
		}

		user, ok := users[post.UserID]
		if ok {
			post.User = user
			posts = append(posts, post)
			continue
		}

		user, err := p.userRepository.GetByID(ctx, post.UserID, false)
		if err != nil {
			return nil, 0, fmt.Errorf("Post.Feed() user.GetByID: %w", err)
		}
		users[post.UserID] = user
		post.User = user
		posts = append(posts, post)
	}
	return posts, resp[0].Count, nil
}
func (p *Post) PublicGet(ctx context.Context, id int64) (*entity.Post, error) {
	resp, err := p.postRepository.GetByID(ctx, repositoryPost.GetByIDParams{
		ID:     id,
		UserID: utils.NewSqlInt64(nil),
	})
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorPostNotFound
		}
		return nil, fmt.Errorf("Post.PublicGet() GetByID: %w", err)
	}
	post := &entity.Post{
		ID:          resp.ID,
		UserID:      resp.UserID,
		Title:       resp.Title,
		Short:       resp.Short,
		Body:        resp.Body,
		Tags:        strings.Split(resp.Tags.String, ";"),
		IsPublished: resp.IsPublished,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}

	user, err := p.userRepository.GetByID(ctx, post.UserID, false)
	if err != nil {
		return nil, fmt.Errorf("Post.PublicGet() user.GetByID: %w", err)
	}
	post.User = user
	return post, nil
}
func (p *Post) Create(ctx context.Context, req *dto.CreatePostDTO, userID int64) (*entity.Post, error) {
	resp, err := p.postRepository.Create(ctx, repositoryPost.CreateParams{
		UserID: userID,
		Title:  req.Title,
		Short:  req.Short,
		Body:   req.Body,
		Tags: sql.NullString{
			String: strings.Join(req.Tags, ";"),
			Valid:  true,
		},
		IsPublished: req.IsPublished,
	})
	if err != nil {
		return nil, fmt.Errorf("Post.Create() Create: %w", err)
	}
	post := &entity.Post{
		ID:          resp.ID,
		UserID:      resp.UserID,
		Title:       resp.Title,
		Short:       resp.Short,
		Body:        resp.Body,
		Tags:        strings.Split(resp.Tags.String, ";"),
		IsPublished: resp.IsPublished,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}
	return post, nil
}
func (p *Post) Edit(ctx context.Context, req *dto.EditPostDTO, userID int64) (*entity.Post, error) {
	resp, err := p.postRepository.Edit(ctx, repositoryPost.EditParams{
		Title: req.Title,
		Short: req.Short,
		Body:  req.Body,
		//Tags        sql.NullString `json:"tags"`
		IsPublished: req.IsPublished,
		ID:          req.ID,
		UserID:      userID,
	})
	if err != nil {
		return nil, fmt.Errorf("Post.Edit() Edit: %w", err)
	}
	post := &entity.Post{
		ID:          resp.ID,
		UserID:      resp.UserID,
		Title:       resp.Title,
		Short:       resp.Short,
		Body:        resp.Body,
		Tags:        strings.Split(resp.Tags.String, ";"),
		IsPublished: resp.IsPublished,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}
	return post, nil
}
func (p *Post) List(ctx context.Context, req *dto.ListPostDTO, userID int64) ([]*entity.Post, int64, error) {
	limit, offset := utils.GetPagination(req.Limit, req.Page)
	resp, err := p.postRepository.List(ctx, repositoryPost.ListParams{
		Limit:                limit,
		Offset:               offset,
		Query:                req.Query,
		RelatedToUser:        userID,
		DisplayOnlyPublished: false,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("Post.List() List: %w", err)
	}
	if len(resp) == 0 {
		return []*entity.Post{}, 0, nil
	}

	posts := make([]*entity.Post, 0, len(resp))
	users := make(map[int64]*entity.User)
	for _, el := range resp {
		post := &entity.Post{
			ID:          el.Post.ID,
			UserID:      el.Post.UserID,
			Title:       el.Post.Title,
			Short:       el.Post.Short,
			Body:        el.Post.Body,
			Tags:        strings.Split(el.Post.Tags.String, ";"),
			IsPublished: el.Post.IsPublished,
			CreatedAt:   el.Post.CreatedAt,
			UpdatedAt:   el.Post.UpdatedAt,
		}

		user, ok := users[post.UserID]
		if ok {
			post.User = user
			posts = append(posts, post)
			continue
		}

		user, err := p.userRepository.GetByID(ctx, post.UserID, false)
		if err != nil {
			return nil, 0, fmt.Errorf("Post.List() user.GetByID: %w", err)
		}
		users[post.UserID] = user
		post.User = user
		posts = append(posts, post)
	}
	return posts, resp[0].Count, nil
}

var (
	ErrorPostNotFound = errors.New("post not found")
)

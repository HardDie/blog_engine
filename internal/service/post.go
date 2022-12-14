package service

import (
	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/repository"
)

type IPost interface {
	Feed(req *dto.FeedPostDTO) ([]*entity.Post, int32, error)

	Create(req *dto.CreatePostDTO, userID int32) (*entity.Post, error)
	Edit(req *dto.EditPostDTO, userID int32) (*entity.Post, error)
	List(req *dto.ListPostDTO, userID int32) ([]*entity.Post, int32, error)
}

type Post struct {
	repository repository.IPost
}

func NewPost(repository repository.IPost) *Post {
	return &Post{
		repository: repository,
	}
}

func (p *Post) Feed(req *dto.FeedPostDTO) ([]*entity.Post, int32, error) {
	return p.repository.List(&dto.ListPostFilter{
		Limit:                req.Limit,
		Page:                 req.Page,
		Query:                req.Query,
		DisplayOnlyPublished: true,
	})
}

func (p *Post) Create(req *dto.CreatePostDTO, userID int32) (*entity.Post, error) {
	return p.repository.Create(req, userID)
}
func (p *Post) Edit(req *dto.EditPostDTO, userID int32) (*entity.Post, error) {
	return p.repository.Edit(req, userID)
}
func (p *Post) List(req *dto.ListPostDTO, userID int32) ([]*entity.Post, int32, error) {
	return p.repository.List(&dto.ListPostFilter{
		Limit:                req.Limit,
		Page:                 req.Page,
		Query:                req.Query,
		RelatedToUser:        userID,
		DisplayOnlyPublished: false,
	})
}

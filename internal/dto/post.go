package dto

type CreatePostDTO struct {
	Title       string   `json:"title" validate:"required"`
	Short       string   `json:"short" validate:"required"`
	Body        string   `json:"body" validate:"required"`
	Tags        []string `json:"tags" validate:"omitempty,dive,alphanum"`
	IsPublished bool     `json:"isPublished"`
}

type FeedPostDTO struct {
	Limit int32  `json:"limit" validate:"omitempty,gt=0"`
	Page  int32  `json:"page" validate:"omitempty,gt=0"`
	Query string `json:"query"`
}

type EditPostDTO struct {
	ID          int32    `json:"-"`
	Title       string   `json:"title" validate:"required"`
	Short       string   `json:"short" validate:"required"`
	Body        string   `json:"body" validate:"required"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"isPublished"`
}

type ListPostDTO struct {
	Limit int32  `json:"limit" validate:"omitempty,gt=0"`
	Page  int32  `json:"page" validate:"omitempty,gt=0"`
	Query string `json:"query"`
}

/*
 * internal
 */

type ListPostFilter struct {
	Limit                int32
	Page                 int32
	Query                string
	RelatedToUser        int32
	DisplayOnlyPublished bool
}

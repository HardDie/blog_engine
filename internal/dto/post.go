package dto

import "fmt"

type CreatePostDTO struct {
	Title       string   `json:"title"`
	Short       string   `json:"short"`
	Body        string   `json:"body"`
	Tags        []string `json:"tags"`
	IsPublished bool     `json:"isPublished"`
}

func (p *CreatePostDTO) Validate() error {
	if p == nil {
		return nil
	}
	if p.Title == "" {
		return fmt.Errorf("title can't be empty")
	}
	if p.Short == "" {
		return fmt.Errorf("short can't be empty")
	}
	if p.Body == "" {
		return fmt.Errorf("body can't be empty")
	}
	return nil
}

type FeedPostDTO struct {
	Limit int32  `json:"limit"`
	Page  int32  `json:"page"`
	Query string `json:"query"`
}

func (d *FeedPostDTO) Validate() error {
	if d == nil {
		return nil
	}
	if d.Limit < 0 {
		return fmt.Errorf("limit can't be less than 0")
	}
	if d.Page < 0 {
		return fmt.Errorf("page can't be less than 0")
	}
	return nil
}

type ListPostFilter struct {
	Limit                int32
	Page                 int32
	Query                string
	DisplayOnlyPublished bool
}

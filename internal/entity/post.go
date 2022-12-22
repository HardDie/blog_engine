package entity

import "time"

type Post struct {
	ID          int32      `json:"id"`
	UserID      int32      `json:"userId"`
	Title       string     `json:"title"`
	Short       string     `json:"short"`
	Body        string     `json:"body"`
	Tags        string     `json:"tags"`
	IsPublished bool       `json:"isPublished"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

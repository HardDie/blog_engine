package entity

import "time"

type Session struct {
	ID          int64      `json:"id"`
	UserID      int64      `json:"userId"`
	SessionHash string     `json:"sessionHash"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

package entity

import "time"

type Session struct {
	UserID      int64     `json:"userId"`
	SessionHash string    `json:"sessionHash"`
	CreatedAt   time.Time `json:"createdAt"`
}

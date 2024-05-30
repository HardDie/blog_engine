package models

import "time"

type Session struct {
	UserID      int64
	SessionHash string
	CreatedAt   time.Time
}

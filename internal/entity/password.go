package entity

import "time"

type Password struct {
	ID             int64      `json:"id"`
	UserID         int64      `json:"userId"`
	PasswordHash   string     `json:"passwordHash"`
	FailedAttempts int64      `json:"failedAttempts"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	BlockedAt      *time.Time `json:"blockedAt"`
}

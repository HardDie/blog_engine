package entity

import "time"

type Invite struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"userId"`
	InvitedHash string    `json:"passwordHash"`
	IsActivated bool      `json:"isActivated"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	DeletedAt   time.Time `json:"deletedAt"`
	BlockedAt   time.Time `json:"blockedAt"`
}

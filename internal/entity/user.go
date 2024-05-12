package entity

import "time"

type User struct {
	ID              int32      `json:"id"`
	Username        string     `json:"username,omitempty"`
	DisplayedName   string     `json:"displayedName"`
	Email           *string    `json:"email,omitempty"`
	InvitedByUserID int32      `json:"invitedByUserId"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt"`
}

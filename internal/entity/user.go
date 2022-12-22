package entity

import "time"

type User struct {
	ID              int32      `json:"id"`
	Username        string     `json:"username"`
	DisplayedName   string     `json:"displayedName"`
	Email           string     `json:"email"`
	InvitedByUserID int32      `json:"invitedByUserId"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
	DeletedAt       *time.Time `json:"deletedAt"`
}

package dto

type RegisterDTO struct {
	Username      string `json:"username" validate:"required"`
	Password      string `json:"password" validate:"required"`
	DisplayedName string `json:"displayedName" validate:"required"`
	Invite        string `json:"invite" validate:"required,uuid"`
}

type LoginDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

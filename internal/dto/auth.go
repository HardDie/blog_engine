package dto

import (
	"fmt"

	"github.com/google/uuid"
)

type RegisterDTO struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
	Invite   *string `json:"invite"`
}

func (d *RegisterDTO) Validate() error {
	if d == nil {
		return nil
	}
	if d.Username == nil || *d.Username == "" {
		return fmt.Errorf("username can't be empty")
	}
	if d.Password == nil || *d.Password == "" {
		return fmt.Errorf("password can't be empty")
	}
	if d.Invite == nil || *d.Invite == "" {
		return fmt.Errorf("invite can't be empty")
	}
	if _, err := uuid.Parse(*d.Invite); err != nil {
		return fmt.Errorf("invalid invite")
	}
	return nil
}

type LoginDTO struct {
	Username *string `json:"username"`
	Password *string `json:"password"`
}

func (d *LoginDTO) Validate() error {
	if d == nil {
		return nil
	}
	if d.Username == nil || *d.Username == "" {
		return fmt.Errorf("username can't be empty")
	}
	if d.Password == nil || *d.Password == "" {
		return fmt.Errorf("password can't be empty")
	}
	return nil
}

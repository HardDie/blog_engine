package dto

import "fmt"

type UpdatePasswordDTO struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (d *UpdatePasswordDTO) Validate() error {
	if d == nil {
		return nil
	}
	if d.OldPassword == "" {
		return fmt.Errorf("oldPassword can't be empty")
	}
	if d.NewPassword == "" {
		return fmt.Errorf("newPassword can't be empty")
	}
	if d.OldPassword == d.NewPassword {
		return fmt.Errorf("old and new passwords can't be same")
	}
	return nil
}

type UpdateProfileDTO struct {
	DisplayedName string `json:"displayedName"`
	Email         string `json:"email"`
}

func (d *UpdateProfileDTO) Validate() error {
	if d == nil {
		return nil
	}
	if d.DisplayedName == "" {
		return fmt.Errorf("displayedName can't be empty")
	}
	if d.Email != "" {
		//return fmt.Errorf("newPassword can't be empty")
	}
	return nil
}

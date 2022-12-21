package service

import (
	"fmt"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/repository"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IUser interface {
	Get(id int32) (*entity.User, error)

	Password(req *dto.UpdatePasswordDTO, userID int32) error
	Profile(req *dto.UpdateProfileDTO, userID int32) (*entity.User, error)
}

type User struct {
	userRepository     repository.IUser
	passwordRepository repository.IPassword
}

func NewUser(repository repository.IUser, password repository.IPassword) *User {
	return &User{
		userRepository:     repository,
		passwordRepository: password,
	}
}

func (s *User) Get(id int32) (*entity.User, error) {
	return s.userRepository.GetByID(id, false)
}

func (s *User) Password(req *dto.UpdatePasswordDTO, userID int32) error {
	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(userID)
	if err != nil {
		return err
	}
	if password == nil {
		return fmt.Errorf("password for user not exist")
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.OldPassword, password.PasswordHash) {
		return fmt.Errorf("invalid old password")
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	password, err = s.passwordRepository.Update(userID, hashPassword)
	if err != nil {
		return err
	}
	return nil
}
func (s *User) Profile(req *dto.UpdateProfileDTO, userID int32) (*entity.User, error) {
	return s.userRepository.Update(req, userID)
}

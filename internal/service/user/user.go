package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	repositoryPassword "github.com/HardDie/blog_engine/internal/repository/password"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/user"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IUser interface {
	Get(ctx context.Context, id int32) (*entity.User, error)

	Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int32) error
	Profile(ctx context.Context, req *dto.UpdateProfileDTO, userID int32) (*entity.User, error)
}

type User struct {
	userRepository     repositoryUser.IUser
	passwordRepository repositoryPassword.IPassword
}

func New(user repositoryUser.IUser, password repositoryPassword.IPassword) *User {
	return &User{
		userRepository:     user,
		passwordRepository: password,
	}
}

func (s *User) Get(ctx context.Context, id int32) (*entity.User, error) {
	resp, err := s.userRepository.GetByID(ctx, id, false)
	if err != nil {
		switch {
		case errors.Is(err, repositoryUser.ErrorNotFound):
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("User.Get() GetByID: %w", err)
	}
	return resp, nil
}

func (s *User) Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int32) error {
	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("User.Password() GetByUserID: %w", err)
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.OldPassword, password.PasswordHash) {
		return ErrorInvalidPassword
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.NewPassword)
	if err != nil {
		return fmt.Errorf("User.Password() HashBcrypt: %w", err)
	}

	// Update password
	password, err = s.passwordRepository.Update(ctx, userID, hashPassword)
	if err != nil {
		return fmt.Errorf("User.Password() Update: %w", err)
	}
	return nil
}
func (s *User) Profile(ctx context.Context, req *dto.UpdateProfileDTO, userID int32) (*entity.User, error) {
	resp, err := s.userRepository.Update(ctx, req, userID)
	if err != nil {
		return nil, fmt.Errorf("User.Profile() Update: %w", err)
	}
	return resp, nil
}

var (
	ErrorUserNotFound    = errors.New("user not found")
	ErrorInvalidPassword = errors.New("invalid password")
)

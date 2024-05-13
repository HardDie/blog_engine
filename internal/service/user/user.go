package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	repositoryPassword "github.com/HardDie/blog_engine/internal/repository/password"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/user"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IUser interface {
	Get(ctx context.Context, userID int64) (*entity.User, error)

	Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int64) error
	Profile(ctx context.Context, req *dto.UpdateProfileDTO, userID int64) (*entity.User, error)
}

type User struct {
	userRepository     repositoryUser.Querier
	passwordRepository repositoryPassword.Querier
}

func New(user repositoryUser.Querier, password repositoryPassword.Querier) *User {
	return &User{
		userRepository:     user,
		passwordRepository: password,
	}
}

func (s *User) Get(ctx context.Context, userID int64) (*entity.User, error) {
	resp, err := s.userRepository.GetByIDPublic(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("User.Get() GetByID: %w", err)
	}
	user := &entity.User{
		ID:              resp.ID,
		DisplayedName:   resp.DisplayedName,
		InvitedByUserID: resp.InvitedByUser,
		CreatedAt:       resp.CreatedAt,
		UpdatedAt:       resp.UpdatedAt,
	}
	return user, nil
}
func (s *User) Password(ctx context.Context, req *dto.UpdatePasswordDTO, userID int64) error {
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
	password, err = s.passwordRepository.Update(ctx, repositoryPassword.UpdateParams{
		ID:           password.ID,
		PasswordHash: hashPassword,
	})
	if err != nil {
		return fmt.Errorf("User.Password() Update: %w", err)
	}
	return nil
}
func (s *User) Profile(ctx context.Context, req *dto.UpdateProfileDTO, userID int64) (*entity.User, error) {
	resp, err := s.userRepository.Update(ctx, repositoryUser.UpdateParams{
		ID:            userID,
		DisplayedName: req.DisplayedName,
		Email:         utils.NewSqlString(req.Email),
	})
	if err != nil {
		return nil, fmt.Errorf("User.Profile() Update: %w", err)
	}
	user := &entity.User{
		ID:              resp.ID,
		Username:        resp.Username,
		DisplayedName:   resp.DisplayedName,
		Email:           utils.SqlStringToString(resp.Email),
		InvitedByUserID: resp.InvitedByUser,
		CreatedAt:       resp.CreatedAt,
		UpdatedAt:       resp.UpdatedAt,
	}
	return user, nil
}

var (
	ErrorUserNotFound    = errors.New("user not found")
	ErrorInvalidPassword = errors.New("invalid password")
)

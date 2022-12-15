package service

import (
	"fmt"
	"time"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/logger"
	"github.com/HardDie/blog_engine/internal/repository"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IAuth interface {
	Register(req *dto.RegisterDTO) (*entity.User, error)
	Login(req *dto.LoginDTO) (*entity.User, error)
	GenerateCookie(userID int32) (string, error)
	ValidateCookie(session string) (*int32, error)
}

type Auth struct {
	userRepository     repository.IUser
	passwordRepository repository.IPassword
	sessionRepository  repository.ISession
	inviteRepository   repository.IInvite
}

func NewAuth(user repository.IUser, password repository.IPassword,
	session repository.ISession, invite repository.IInvite) *Auth {
	return &Auth{
		userRepository:     user,
		passwordRepository: password,
		sessionRepository:  session,
		inviteRepository:   invite,
	}
}

func (s *Auth) Register(req *dto.RegisterDTO) (*entity.User, error) {
	// Check if such user already exist
	user, err := s.userRepository.GetByName(*req.Username)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("user already exist")
	}

	// Hashing invite
	hashInvite := utils.HashSha256(*req.Invite)

	// Validating invite
	invite, err := s.inviteRepository.CheckIfExistAndDisable(hashInvite)
	if err != nil {
		return nil, err
	}
	if invite == nil {
		return nil, fmt.Errorf("can't find such invite")
	}

	// Check if invite is not expired
	if time.Now().Sub(invite.UpdatedAt) > time.Hour*24 {
		err = s.inviteRepository.Delete(invite.ID)
		if err != nil {
			logger.Error.Printf("Can't delete expired invite [%d]: %v", invite.ID, err.Error())
		}
		return nil, fmt.Errorf("invite has expired")
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(*req.Password)
	if err != nil {
		return nil, err
	}

	// Create a user
	user, err = s.userRepository.Create(*req.Username, invite.UserID)
	if err != nil {
		return nil, err
	}

	// Create a password
	_, err = s.passwordRepository.Create(*user.ID, hashPassword)
	if err != nil {
		return nil, err
	}

	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		return nil, fmt.Errorf("generate session key: %w", err)
	}

	// Write session to DB
	_, err = s.sessionRepository.CreateOrUpdate(*user.ID, utils.HashSha256(sessionKey))
	if err != nil {
		return nil, fmt.Errorf("write session to DB: %w", err)
	}

	return user, nil
}
func (s *Auth) Login(req *dto.LoginDTO) (*entity.User, error) {
	// Check if such user exist
	user, err := s.userRepository.GetByName(*req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not exist")
	}

	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(*user.ID)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return nil, fmt.Errorf("password for user not exist")
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(*req.Password, password.PasswordHash) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
func (s *Auth) GenerateCookie(userID int32) (string, error) {
	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		return "", fmt.Errorf("generate session key: %w", err)
	}

	// Write session to DB
	_, err = s.sessionRepository.CreateOrUpdate(userID, utils.HashSha256(sessionKey))
	if err != nil {
		return "", fmt.Errorf("write session to DB: %w", err)
	}

	return sessionKey, nil
}
func (s *Auth) ValidateCookie(session string) (*int32, error) {
	sessionHash := utils.HashSha256(session)
	userID, err := s.sessionRepository.GetUserID(sessionHash)
	if err != nil {
		return nil, err
	}
	return userID, nil
}

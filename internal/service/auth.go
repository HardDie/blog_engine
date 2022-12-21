package service

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/logger"
	"github.com/HardDie/blog_engine/internal/repository"
	"github.com/HardDie/blog_engine/internal/utils"
)

var (
	ErrorSessionHasExpired = errors.New("session has expired")
)

type IAuth interface {
	Register(req *dto.RegisterDTO) (*entity.User, error)
	Login(req *dto.LoginDTO) (*entity.User, error)
	Logout(sessionID int32) error
	GenerateCookie(userID int32) (string, error)
	ValidateCookie(session string) (*entity.Session, error)
	GetUserInfo(userID int32) (*entity.User, error)
}

type Auth struct {
	userRepository     repository.IUser
	passwordRepository repository.IPassword
	sessionRepository  repository.ISession
	inviteRepository   repository.IInvite

	mutex sync.Mutex
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
	s.mutex.Lock()
	defer func() {
		s.mutex.Unlock()
	}()

	// Hashing invite
	hashInvite := utils.HashSha256(req.Invite)

	// Check if invite exist
	invite, err := s.inviteRepository.GetByInviteHash(hashInvite)
	if err != nil {
		return nil, fmt.Errorf("error while trying get invite: %w", err)
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

	// Check if username is not busy
	user, err := s.userRepository.GetByName(req.Username)
	if err != nil {
		return nil, fmt.Errorf("error while trying get user: %w", err)
	}
	if user != nil {
		return nil, fmt.Errorf("username already exist")
	}

	// Activating invite
	invite, err = s.inviteRepository.Activate(invite.ID)
	if err != nil || invite == nil {
		return nil, fmt.Errorf("error while activating invite: %w", err)
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.Password)
	if err != nil {
		return nil, err
	}

	// Create a user
	user, err = s.userRepository.Create(req.Username, req.DisplayedName, invite.UserID)
	if err != nil {
		return nil, err
	}

	// Create a password
	_, err = s.passwordRepository.Create(*user.ID, hashPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *Auth) Login(req *dto.LoginDTO) (*entity.User, error) {
	// Check if such user exist
	user, err := s.userRepository.GetByName(req.Username)
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
	if !utils.HashBcryptCompare(req.Password, password.PasswordHash) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}
func (s *Auth) Logout(sessionID int32) error {
	return s.sessionRepository.DeleteByID(sessionID)
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
func (s *Auth) ValidateCookie(sessionToken string) (*entity.Session, error) {
	// Check if session exist
	sessionHash := utils.HashSha256(sessionToken)
	session, err := s.sessionRepository.GetByUserID(sessionHash)
	if err != nil {
		return nil, err
	}

	// Check if session is not expired
	if time.Now().Sub(session.UpdatedAt) > time.Hour*24 {
		return nil, ErrorSessionHasExpired
	}
	return session, nil
}
func (s *Auth) GetUserInfo(userID int32) (*entity.User, error) {
	return s.userRepository.GetByID(userID, true)
}

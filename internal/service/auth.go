package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HardDie/blog_engine/internal/config"
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
	Register(ctx context.Context, req *dto.RegisterDTO) (*entity.User, error)
	Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error)
	Logout(ctx context.Context, sessionID int32) error
	GenerateCookie(ctx context.Context, userID int32) (string, error)
	ValidateCookie(ctx context.Context, session string) (*entity.Session, error)
	GetUserInfo(ctx context.Context, userID int32) (*entity.User, error)
}

type Auth struct {
	userRepository     repository.IUser
	passwordRepository repository.IPassword
	sessionRepository  repository.ISession
	inviteRepository   repository.IInvite

	cfg   *config.Config
	mutex sync.Mutex
}

func NewAuth(cfg *config.Config, user repository.IUser, password repository.IPassword,
	session repository.ISession, invite repository.IInvite) *Auth {
	return &Auth{
		cfg:                cfg,
		userRepository:     user,
		passwordRepository: password,
		sessionRepository:  session,
		inviteRepository:   invite,
	}
}

func (s *Auth) Register(ctx context.Context, req *dto.RegisterDTO) (*entity.User, error) {
	s.mutex.Lock()
	defer func() {
		s.mutex.Unlock()
	}()

	// Hashing invite
	hashInvite := utils.HashSha256(req.Invite)

	// Check if invite exist
	invite, err := s.inviteRepository.GetByInviteHash(ctx, hashInvite)
	if err != nil {
		return nil, fmt.Errorf("error while trying get invite: %w", err)
	}
	if invite == nil {
		return nil, fmt.Errorf("can't find such invite")
	}

	// Check if invite is not expired
	if time.Now().Sub(invite.UpdatedAt) > time.Hour*24 {
		err = s.inviteRepository.Delete(ctx, invite.ID)
		if err != nil {
			logger.Error.Printf("Can't delete expired invite [%d]: %v", invite.ID, err.Error())
		}
		return nil, fmt.Errorf("invite has expired")
	}

	// Check if username is not busy
	user, err := s.userRepository.GetByName(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("error while trying get user: %w", err)
	}
	if user != nil {
		return nil, fmt.Errorf("username already exist")
	}

	// Activating invite
	invite, err = s.inviteRepository.Activate(ctx, invite.ID)
	if err != nil || invite == nil {
		return nil, fmt.Errorf("error while activating invite: %w", err)
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.Password)
	if err != nil {
		return nil, err
	}

	// Create a user
	user, err = s.userRepository.Create(ctx, req.Username, req.DisplayedName, invite.UserID)
	if err != nil {
		return nil, err
	}

	// Create a password
	_, err = s.passwordRepository.Create(ctx, user.ID, hashPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (s *Auth) Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error) {
	// Check if such user exist
	user, err := s.userRepository.GetByName(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not exist")
	}

	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, err
	}
	if password == nil {
		return nil, fmt.Errorf("password for user not exist")
	}

	// Check if the password is locked after failed attempts
	if password.FailedAttempts >= int32(s.cfg.PwdMaxAttempts) {
		// Check if the password block time has expired
		if time.Now().Sub(password.UpdatedAt) <= time.Hour*time.Duration(s.cfg.PwdBlockTime) {
			return nil, fmt.Errorf("user was blocked after failed attempts")
		}
		// If the blocking time has expired, reset the counter of failed attempts
		password, err = s.passwordRepository.ResetFailedAttempts(ctx, password.ID)
		if err != nil {
			return nil, fmt.Errorf("error resetting the counter of failed attempts: %w", err)
		}
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.Password, password.PasswordHash) {
		// Increased number of failed attempts
		_, err = s.passwordRepository.IncreaseFailedAttempts(ctx, password.ID)
		if err != nil {
			logger.Error.Println("Error increasing failed attempts:", err.Error())
		}
		return nil, fmt.Errorf("invalid password")
	}

	// Reset the failed attempts counter after the first successful attempt
	if password.FailedAttempts > 0 {
		_, err = s.passwordRepository.ResetFailedAttempts(ctx, password.ID)
		if err != nil {
			logger.Error.Println("Error flushing failed attempts:", err.Error())
		}
	}
	return user, nil
}
func (s *Auth) Logout(ctx context.Context, sessionID int32) error {
	return s.sessionRepository.DeleteByID(ctx, sessionID)
}
func (s *Auth) GenerateCookie(ctx context.Context, userID int32) (string, error) {
	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		return "", fmt.Errorf("generate session key: %w", err)
	}

	// Write session to DB
	_, err = s.sessionRepository.CreateOrUpdate(ctx, userID, utils.HashSha256(sessionKey))
	if err != nil {
		return "", fmt.Errorf("write session to DB: %w", err)
	}

	return sessionKey, nil
}
func (s *Auth) ValidateCookie(ctx context.Context, sessionToken string) (*entity.Session, error) {
	// Check if session exist
	sessionHash := utils.HashSha256(sessionToken)
	session, err := s.sessionRepository.GetByUserID(ctx, sessionHash)
	if err != nil {
		return nil, err
	}

	// Check if session is not expired
	if time.Now().Sub(session.UpdatedAt) > time.Hour*24 {
		return nil, ErrorSessionHasExpired
	}
	return session, nil
}
func (s *Auth) GetUserInfo(ctx context.Context, userID int32) (*entity.User, error) {
	return s.userRepository.GetByID(ctx, userID, true)
}

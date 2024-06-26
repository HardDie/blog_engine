package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/HardDie/blog_engine/internal/config"
	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
	"github.com/HardDie/blog_engine/internal/logger"
	"github.com/HardDie/blog_engine/internal/models"
	repositorySession "github.com/HardDie/blog_engine/internal/repository/boltdb/session"
	repositoryInvite "github.com/HardDie/blog_engine/internal/repository/sqlite/invite"
	repositoryPassword "github.com/HardDie/blog_engine/internal/repository/sqlite/password"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/sqlite/user"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IAuth interface {
	Register(ctx context.Context, req *dto.RegisterDTO) (*entity.User, error)
	Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error)
	Logout(ctx context.Context, sessionHash string) error
	GenerateCookie(ctx context.Context, userID int64) (string, error)
	ValidateCookie(ctx context.Context, session string) (*entity.Session, error)
	GetUserInfo(ctx context.Context, userID int64) (*entity.User, error)
}

type Session interface {
	CreateOrUpdate(_ context.Context, userID int64, sessionHash string) (*models.Session, error)
	DeleteBySessionHash(_ context.Context, sessionHash string) error
	GetBySessionHash(_ context.Context, sessionHash string) (*models.Session, error)
}

type Auth struct {
	userRepository     repositoryUser.Querier
	passwordRepository repositoryPassword.Querier
	sessionRepository  Session
	inviteRepository   repositoryInvite.Querier

	cfg   *config.Config
	mutex sync.Mutex
}

func New(
	cfg *config.Config,
	user repositoryUser.Querier,
	password repositoryPassword.Querier,
	session Session,
	invite repositoryInvite.Querier,
) *Auth {
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorInviteNotFound
		}
		return nil, fmt.Errorf("Auth.Register() GetByInviteHash: %w", err)
	}

	// Check if invite is not expired
	if time.Now().Sub(invite.UpdatedAt) > time.Hour*24 {
		err = s.inviteRepository.Delete(ctx, invite.ID)
		if err != nil {
			logger.Error.Printf("Auth.Register(): Can't delete expired invite [%d]: %v", invite.ID, err.Error())
		}
		return nil, ErrorInviteExpired
	}

	// Check if username is not busy
	u, err := s.userRepository.GetByName(ctx, req.Username)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		u = nil
		// continue
	default:
		return nil, fmt.Errorf("Auth.Register() FindByName: %w", err)
	}
	if u != nil {
		return nil, ErrorUserExist
	}

	// Activating invite
	invite, err = s.inviteRepository.Activate(ctx, invite.ID)
	if err != nil {
		return nil, fmt.Errorf("Auth.Register() Activate: %w", err)
	}

	// Hashing password
	hashPassword, err := utils.HashBcrypt(req.Password)
	if err != nil {
		return nil, fmt.Errorf("Auth.Register() HashBcrypt: %w", err)
	}

	// Create a user
	resp, err := s.userRepository.Create(ctx, repositoryUser.CreateParams{
		Username:      req.Username,
		DisplayedName: req.DisplayedName,
		InvitedByUser: invite.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("Auth.Register() user.Create: %w", err)
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

	// Create a password
	_, err = s.passwordRepository.Create(ctx, repositoryPassword.CreateParams{
		UserID:       user.ID,
		PasswordHash: hashPassword,
	})
	if err != nil {
		return nil, fmt.Errorf("Auth.Register() password.Create: %w", err)
	}

	return user, nil
}
func (s *Auth) Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error) {
	// Check if such user exist
	resp, err := s.userRepository.GetByName(ctx, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("Auth.Login() user.GetByName: %w", err)
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

	// Get password from DB
	password, err := s.passwordRepository.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("Auth.Login() password.GetByUserID: %w", err)
	}

	// Check if the password is locked after failed attempts
	if password.FailedAttempts >= int64(s.cfg.PwdMaxAttempts) {
		// Check if the password block time has expired
		if time.Now().Sub(password.UpdatedAt) <= time.Hour*time.Duration(s.cfg.PwdBlockTime) {
			return nil, ErrorUserBlocked
		}
		// If the blocking time has expired, reset the counter of failed attempts
		password, err = s.passwordRepository.ResetFailedAttempts(ctx, password.ID)
		if err != nil {
			return nil, fmt.Errorf("Auth.Login() ResetFailedAttempts: %w", err)
		}
	}

	// Check if password is correct
	if !utils.HashBcryptCompare(req.Password, password.PasswordHash) {
		// Increased number of failed attempts
		_, err = s.passwordRepository.IncreaseFailedAttempts(ctx, password.ID)
		if err != nil {
			logger.Error.Println("Auth.Login() IncreaseFailedAttempts:", err.Error())
		}
		return nil, ErrorInvalidPassword
	}

	// Reset the failed attempts counter after the first successful attempt
	if password.FailedAttempts > 0 {
		_, err = s.passwordRepository.ResetFailedAttempts(ctx, password.ID)
		if err != nil {
			logger.Error.Println("Auth.Login() ResetFailedAttempts:", err.Error())
		}
	}
	return user, nil
}
func (s *Auth) Logout(ctx context.Context, sessionHash string) error {
	err := s.sessionRepository.DeleteBySessionHash(ctx, sessionHash)
	if err != nil {
		return fmt.Errorf("Auth.Logout() DeleteByID: %w", err)
	}
	return nil
}
func (s *Auth) GenerateCookie(ctx context.Context, userID int64) (string, error) {
	// Generate session key
	sessionKey, err := utils.GenerateSessionKey()
	if err != nil {
		return "", fmt.Errorf("Auth.GenerateCookie() GenerateSessionKey: %w", err)
	}

	// Write session to DB
	_, err = s.sessionRepository.CreateOrUpdate(ctx, userID, utils.HashSha256(sessionKey))
	if err != nil {
		return "", fmt.Errorf("Auth.GenerateCookie() CreateOrUpdate: %w", err)
	}

	return sessionKey, nil
}
func (s *Auth) ValidateCookie(ctx context.Context, sessionToken string) (*entity.Session, error) {
	// Check if session exist
	sessionHash := utils.HashSha256(sessionToken)
	resp, err := s.sessionRepository.GetBySessionHash(ctx, sessionHash)
	if err != nil {
		switch {
		case errors.Is(err, repositorySession.ErrorNotFound):
			return nil, ErrorSessionNotFound
		}
		return nil, fmt.Errorf("Auth.ValidateCookie() GetyByUserID: %w", err)
	}
	session := &entity.Session{
		UserID:      resp.UserID,
		SessionHash: resp.SessionHash,
		CreatedAt:   resp.CreatedAt,
	}

	// Check if session is not expired
	if time.Now().Sub(session.CreatedAt) > time.Hour*24 {
		return nil, ErrorSessionHasExpired
	}
	return session, nil
}
func (s *Auth) GetUserInfo(ctx context.Context, userID int64) (*entity.User, error) {
	resp, err := s.userRepository.GetByIDPrivate(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("Auth.GetUserInfo() GetByID: %w", err)
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
	ErrorInviteNotFound    = errors.New("invite not found")
	ErrorInviteExpired     = errors.New("invite has expired")
	ErrorUserExist         = errors.New("user exist")
	ErrorUserNotFound      = errors.New("user not found")
	ErrorUserBlocked       = errors.New("user blocked")
	ErrorInvalidPassword   = errors.New("invalid password")
	ErrorSessionNotFound   = errors.New("session not found")
	ErrorSessionHasExpired = errors.New("session has expired")
)

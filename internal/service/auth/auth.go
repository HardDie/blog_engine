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
	repositoryInvite "github.com/HardDie/blog_engine/internal/repository/invite"
	repositoryPassword "github.com/HardDie/blog_engine/internal/repository/password"
	repositorySession "github.com/HardDie/blog_engine/internal/repository/session"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/user"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IAuth interface {
	Register(ctx context.Context, req *dto.RegisterDTO) (*entity.User, error)
	Login(ctx context.Context, req *dto.LoginDTO) (*entity.User, error)
	Logout(ctx context.Context, sessionID int64) error
	GenerateCookie(ctx context.Context, userID int64) (string, error)
	ValidateCookie(ctx context.Context, session string) (*entity.Session, error)
	GetUserInfo(ctx context.Context, userID int64) (*entity.User, error)
}

type Auth struct {
	userRepository     repositoryUser.IUser
	passwordRepository repositoryPassword.Querier
	sessionRepository  repositorySession.Querier
	inviteRepository   repositoryInvite.Querier

	cfg   *config.Config
	mutex sync.Mutex
}

func New(
	cfg *config.Config,
	user repositoryUser.IUser,
	password repositoryPassword.Querier,
	session repositorySession.Querier,
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
	user, err := s.userRepository.FindByName(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("Auth.Register() FindByName: %w", err)
	}
	if user != nil {
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
	user, err = s.userRepository.Create(ctx, req.Username, req.DisplayedName, invite.UserID)
	if err != nil {
		return nil, fmt.Errorf("Auth.Register() user.Create: %w", err)
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
	user, err := s.userRepository.GetByName(ctx, req.Username)
	if err != nil {
		switch {
		case errors.Is(err, repositoryUser.ErrorNotFound):
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("Auth.Login() user.GetByName: %w", err)
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
func (s *Auth) Logout(ctx context.Context, sessionID int64) error {
	err := s.sessionRepository.DeleteByID(ctx, sessionID)
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
	_, err = s.sessionRepository.CreateOrUpdate(ctx, repositorySession.CreateOrUpdateParams{
		UserID:      userID,
		SessionHash: utils.HashSha256(sessionKey),
	})
	if err != nil {
		return "", fmt.Errorf("Auth.GenerateCookie() CreateOrUpdate: %w", err)
	}

	return sessionKey, nil
}
func (s *Auth) ValidateCookie(ctx context.Context, sessionToken string) (*entity.Session, error) {
	// Check if session exist
	sessionHash := utils.HashSha256(sessionToken)
	resp, err := s.sessionRepository.GetByUserID(ctx, sessionHash)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrorSessionNotFound
		}
		return nil, fmt.Errorf("Auth.ValidateCookie() GetyByUserID: %w", err)
	}
	session := &entity.Session{
		ID:          resp.ID,
		UserID:      resp.UserID,
		SessionHash: resp.SessionHash,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}

	// Check if session is not expired
	if time.Now().Sub(session.UpdatedAt) > time.Hour*24 {
		return nil, ErrorSessionHasExpired
	}
	return session, nil
}
func (s *Auth) GetUserInfo(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := s.userRepository.GetByID(ctx, userID, true)
	if err != nil {
		switch {
		case errors.Is(err, repositoryUser.ErrorNotFound):
			return nil, ErrorUserNotFound
		}
		return nil, fmt.Errorf("Auth.GetUserInfo() GetByID: %w", err)
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

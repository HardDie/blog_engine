package invite

import (
	"context"
	"errors"
	"fmt"

	repositoryInvite "github.com/HardDie/blog_engine/internal/repository/invite"
	repositoryUser "github.com/HardDie/blog_engine/internal/repository/user"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IInvite interface {
	Generate(ctx context.Context, userID int32) (string, error)
	Revoke(ctx context.Context, userID int32) error
}

type Invite struct {
	userRepository   repositoryUser.IUser
	inviteRepository repositoryInvite.IInvite
}

func New(user repositoryUser.IUser, invite repositoryInvite.IInvite) *Invite {
	return &Invite{
		userRepository:   user,
		inviteRepository: invite,
	}
}

func (s *Invite) Generate(ctx context.Context, userID int32) (string, error) {
	// Generate invite
	inviteCode, err := utils.UUIDGenerate()
	if err != nil {
		return "", fmt.Errorf("Invite.Generate() UUIDGenerate: %w", err)
	}
	// Hashing invite for DB
	inviteHash := utils.HashSha256(inviteCode)
	// Write hash of invite into DB
	_, err = s.inviteRepository.CreateOrUpdate(ctx, userID, inviteHash)
	if err != nil {
		return "", fmt.Errorf("Invite.Generate() CreateOrUpdate: %w", err)
	}
	return inviteCode, nil
}
func (s *Invite) Revoke(ctx context.Context, userID int32) error {
	invite, err := s.inviteRepository.GetActiveByUserID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, repositoryInvite.ErrorNotFound):
			return ErrorInviteNotFound
		}
		return err
	}
	err = s.inviteRepository.Delete(ctx, invite.ID)
	if err != nil {
		return fmt.Errorf("Invite.Revoke() Delete: %w", err)
	}
	return nil
}

var (
	ErrorInviteNotFound = errors.New("invite not found")
)
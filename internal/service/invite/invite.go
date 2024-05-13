package invite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	repositoryInvite "github.com/HardDie/blog_engine/internal/repository/invite"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IInvite interface {
	Generate(ctx context.Context, userID int64) (string, error)
	Revoke(ctx context.Context, userID int64) error
}

type Invite struct {
	inviteRepository repositoryInvite.Querier
}

func New(invite repositoryInvite.Querier) *Invite {
	return &Invite{
		inviteRepository: invite,
	}
}

func (s *Invite) Generate(ctx context.Context, userID int64) (string, error) {
	// Generate invite
	inviteCode, err := utils.UUIDGenerate()
	if err != nil {
		return "", fmt.Errorf("Invite.Generate() UUIDGenerate: %w", err)
	}
	// Hashing invite for DB
	inviteHash := utils.HashSha256(inviteCode)
	// Write hash of invite into DB
	_, err = s.inviteRepository.CreateOrUpdate(ctx, repositoryInvite.CreateOrUpdateParams{
		UserID:     userID,
		InviteHash: inviteHash,
	})
	if err != nil {
		return "", fmt.Errorf("Invite.Generate() CreateOrUpdate: %w", err)
	}
	return inviteCode, nil
}
func (s *Invite) Revoke(ctx context.Context, userID int64) error {
	invite, err := s.inviteRepository.GetActiveByUserID(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
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

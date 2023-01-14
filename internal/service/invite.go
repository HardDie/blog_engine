package service

import (
	"context"
	"fmt"

	"github.com/HardDie/blog_engine/internal/repository"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IInvite interface {
	Generate(ctx context.Context, userID int32) (string, error)
	Revoke(ctx context.Context, userID int32) error
}

type Invite struct {
	userRepository   repository.IUser
	inviteRepository repository.IInvite
}

func NewInvite(user repository.IUser, invite repository.IInvite) *Invite {
	return &Invite{
		userRepository:   user,
		inviteRepository: invite,
	}
}

func (s *Invite) Generate(ctx context.Context, userID int32) (string, error) {
	// Generate invite
	inviteCode, err := utils.UUIDGenerate()
	if err != nil {
		return "", err
	}
	// Hashing invite for DB
	inviteHash := utils.HashSha256(inviteCode)
	// Write hash of invite into DB
	_, err = s.inviteRepository.CreateOrUpdate(ctx, userID, inviteHash)
	if err != nil {
		return "", err
	}
	return inviteCode, nil
}
func (s *Invite) Revoke(ctx context.Context, userID int32) error {
	invite, err := s.inviteRepository.GetActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}
	if invite == nil {
		return fmt.Errorf("no active invites")
	}
	err = s.inviteRepository.Delete(ctx, invite.ID)
	if err != nil {
		return err
	}
	return nil
}

package service

import (
	"fmt"

	"github.com/HardDie/blog_engine/internal/repository"
	"github.com/HardDie/blog_engine/internal/utils"
)

type IInvite interface {
	Generate(userID int32) (string, error)
	Revoke(userID int32) error
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

func (s *Invite) Generate(userID int32) (string, error) {
	// Generate invite
	inviteCode, err := utils.UUIDGenerate()
	if err != nil {
		return "", err
	}
	// Hashing invite for DB
	inviteHash := utils.HashSha256(inviteCode)
	// Write hash of invite into DB
	_, err = s.inviteRepository.CreateOrUpdate(userID, inviteHash)
	if err != nil {
		return "", err
	}
	return inviteCode, nil
}
func (s *Invite) Revoke(userID int32) error {
	invite, err := s.inviteRepository.GetActiveByUserID(userID)
	if err != nil {
		return err
	}
	if invite == nil {
		return fmt.Errorf("no active invites")
	}
	err = s.inviteRepository.Delete(invite.ID)
	if err != nil {
		return err
	}
	return nil
}

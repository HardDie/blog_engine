package repository

import (
	"database/sql"
	"errors"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IInvite interface {
	GetByID(id int32) (*entity.Invite, error)
	GetActiveByUserID(userID int32) (*entity.Invite, error)
	GetAllByUserID(userID int32) ([]*entity.Invite, error)
	CreateOrUpdate(userID int32, inviteHash string) (*entity.Invite, error)
	Delete(id int32) error
}

type Invite struct {
	db *db.DB
}

func NewInvite(db *db.DB) *Invite {
	return &Invite{
		db: db,
	}
}

func (r *Invite) GetByID(id int32) (*entity.Invite, error) {
	invite := &entity.Invite{
		ID: id,
	}

	row := r.db.DB.QueryRow(`
SELECT user_id, invite_hash, is_activated, created_at, updated_at
FROM invites
WHERE id = $1 AND deleted_at IS NULL`, id)

	err := row.Scan(&invite.UserID, &invite.InvitedHash, &invite.IsActivated, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return invite, nil
}
func (r *Invite) GetActiveByUserID(userID int32) (*entity.Invite, error) {
	invite := &entity.Invite{
		UserID:      userID,
		IsActivated: false,
	}

	row := r.db.DB.QueryRow(`
SELECT id, invite_hash, created_at, updated_at
FROM invites
WHERE user_id = $1 AND is_activated IS false AND deleted_at IS NULL`, userID)

	err := row.Scan(&invite.ID, &invite.InvitedHash, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return invite, nil
}
func (r *Invite) GetAllByUserID(userID int32) ([]*entity.Invite, error) {
	rows, err := r.db.DB.Query(`
SELECT id, invite_hash, created_at, updated_at
FROM invites
WHERE user_id = $1 AND deleted_at IS NULL`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resp []*entity.Invite
	for rows.Next() {
		invite := &entity.Invite{
			UserID: userID,
		}
		err = rows.Scan(&invite.ID, &invite.InvitedHash, &invite.CreatedAt, &invite.UpdatedAt)
		if err != nil {
			return nil, err
		}
		resp = append(resp, invite)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return resp, nil
}
func (r *Invite) CreateOrUpdate(userID int32, inviteHash string) (*entity.Invite, error) {
	invite := &entity.Invite{
		UserID:      userID,
		InvitedHash: inviteHash,
	}

	row := r.db.DB.QueryRow(`
	INSERT INTO invites (user_id, invite_hash, is_activated)
	VALUES ($1, $2, false)
	ON CONFLICT (user_id, is_activated) WHERE is_activated IS FALSE DO UPDATE
	SET invite_hash = $2, updated_at = datetime('now')
	RETURNING id, created_at, updated_at`, userID, inviteHash)

	err := row.Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return invite, nil
}
func (r *Invite) Delete(id int32) error {
	row := r.db.DB.QueryRow(`
UPDATE invites
SET deleted_at = datetime('now'), is_activated = true
WHERE id = $1 AND deleted_at IS NULL
RETURNING id`, id)

	err := row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("invite not exist")
		}
		return err
	}
	return nil
}

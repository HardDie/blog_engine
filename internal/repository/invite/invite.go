package invite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IInvite interface {
	GetByID(ctx context.Context, id int32) (*entity.Invite, error)
	GetActiveByUserID(ctx context.Context, userID int32) (*entity.Invite, error)
	GetByInviteHash(ctx context.Context, inviteHash string) (*entity.Invite, error)
	GetAllByUserID(ctx context.Context, userID int32) ([]*entity.Invite, error)
	CreateOrUpdate(ctx context.Context, userID int32, inviteHash string) (*entity.Invite, error)
	Delete(ctx context.Context, id int32) error
	Activate(ctx context.Context, id int32) (*entity.Invite, error)
}

type Invite struct {
	db *db.DB
}

func New(db *db.DB) *Invite {
	return &Invite{
		db: db,
	}
}

func (r *Invite) GetByID(ctx context.Context, id int32) (*entity.Invite, error) {
	invite := &entity.Invite{
		ID: id,
	}

	q := gosql.NewSelect().From("invites")
	q.Columns().Add("user_id", "invite_hash", "is_activated", "created_at", "updated_at")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.UserID, &invite.InviteHash, &invite.IsActivated, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("Invite.GetByID() Scan: %w", err)
	}
	return invite, nil
}
func (r *Invite) GetActiveByUserID(ctx context.Context, userID int32) (*entity.Invite, error) {
	invite := &entity.Invite{
		UserID:      userID,
		IsActivated: false,
	}

	q := gosql.NewSelect().From("invites")
	q.Columns().Add("id", "invite_hash", "created_at", "updated_at")
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("is_activated IS false")
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.ID, &invite.InviteHash, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("Invite.GetActiveByUserID() Scan: %w", err)
	}
	return invite, nil
}
func (r *Invite) GetAllByUserID(ctx context.Context, userID int32) ([]*entity.Invite, error) {
	q := gosql.NewSelect().From("invites")
	q.Columns().Add("id", "invite_hash", "created_at", "updated_at")
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")

	rows, err := r.db.DB.QueryContext(ctx, q.String(), q.GetArguments()...)
	if err != nil {
		return nil, fmt.Errorf("Invite.GetAllByUserID() QueryContext: %w", err)
	}
	defer rows.Close()

	var resp []*entity.Invite
	for rows.Next() {
		invite := &entity.Invite{
			UserID: userID,
		}
		err = rows.Scan(&invite.ID, &invite.InviteHash, &invite.CreatedAt, &invite.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("Invite.GetAllByUserID() Scan: %w", err)
		}
		resp = append(resp, invite)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Invite.GetAllByUserID() Err: %w", err)
	}

	return resp, nil
}
func (r *Invite) GetByInviteHash(ctx context.Context, inviteHash string) (*entity.Invite, error) {
	invite := &entity.Invite{
		InviteHash: inviteHash,
	}

	q := gosql.NewSelect().From("invites")
	q.Columns().Add("id", "user_id", "is_activated", "created_at", "updated_at")
	q.Where().AddExpression("invite_hash = ?", inviteHash)
	q.Where().AddExpression("is_activated IS false")
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.ID, &invite.UserID, &invite.IsActivated, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("Invite.GetByInviteHash() Scan: %w", err)
	}
	return invite, nil
}
func (r *Invite) CreateOrUpdate(ctx context.Context, userID int32, inviteHash string) (*entity.Invite, error) {
	invite := &entity.Invite{
		UserID:     userID,
		InviteHash: inviteHash,
	}

	row := r.db.DB.QueryRowContext(ctx, `
	INSERT INTO invites (user_id, invite_hash, is_activated)
	VALUES ($1, $2, false)
	ON CONFLICT (user_id, is_activated) WHERE is_activated IS FALSE DO UPDATE
	SET invite_hash = $2, updated_at = datetime('now')
	RETURNING id, created_at, updated_at`, userID, inviteHash)

	err := row.Scan(&invite.ID, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Invite.CreateOrUpdate() Scan: %w", err)
	}
	return invite, nil
}
func (r *Invite) Delete(ctx context.Context, id int32) error {
	q := gosql.NewUpdate().Table("invites")
	q.Set().Add("deleted_at = datetime('now')", "is_activated = true")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("id")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&id)
	if err != nil {
		return fmt.Errorf("Invite.Delete() Scan: %w", err)
	}
	return nil
}
func (r *Invite) Activate(ctx context.Context, id int32) (*entity.Invite, error) {
	invite := &entity.Invite{
		ID: id,
	}

	q := gosql.NewUpdate().Table("invites")
	q.Set().Add("is_activated = true")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("is_activated IS false")
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "invite_hash", "is_activated", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&invite.UserID, &invite.InviteHash, &invite.IsActivated, &invite.CreatedAt, &invite.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Invite.Activate() Scan: %w", err)
	}
	return invite, nil
}

var (
	ErrorNotFound = errors.New("invite not found")
)

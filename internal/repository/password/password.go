package password

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IPassword interface {
	Create(ctx context.Context, userID int64, passwordHash string) (*entity.Password, error)
	GetByUserID(ctx context.Context, userID int64) (*entity.Password, error)
	Update(ctx context.Context, id int32, passwordHash string) (*entity.Password, error)
	IncreaseFailedAttempts(ctx context.Context, id int32) (*entity.Password, error)
	ResetFailedAttempts(ctx context.Context, id int32) (*entity.Password, error)
}

type Password struct {
	db *db.DB
}

func New(db *db.DB) *Password {
	return &Password{
		db: db,
	}
}

func (r *Password) Create(ctx context.Context, userID int64, passwordHash string) (*entity.Password, error) {
	password := &entity.Password{
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	q := gosql.NewInsert().Into("passwords")
	q.Columns().Add("user_id", "password_hash")
	q.Columns().Arg(userID, passwordHash)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.ID, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Password.Create() Scan: %w", err)
	}
	return password, nil
}
func (r *Password) GetByUserID(ctx context.Context, userID int64) (*entity.Password, error) {
	password := &entity.Password{
		UserID: userID,
	}

	q := gosql.NewSelect().From("passwords")
	q.Columns().Add("id", "password_hash", "failed_attempts", "created_at", "updated_at")
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.ID, &password.PasswordHash, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrorNotFound
		}
		return nil, fmt.Errorf("Password.GetByUserID() Scan: %w", err)
	}
	return password, nil
}
func (r *Password) Update(ctx context.Context, id int32, passwordHash string) (*entity.Password, error) {
	password := &entity.Password{
		ID:           id,
		PasswordHash: passwordHash,
	}

	q := gosql.NewUpdate().Table("passwords")
	q.Set().Append("password_hash = ?", passwordHash)
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "failed_attempts", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.UserID, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Password.Update() Scan: %w", err)
	}
	return password, nil
}
func (r *Password) IncreaseFailedAttempts(ctx context.Context, id int32) (*entity.Password, error) {
	password := &entity.Password{
		ID: id,
	}

	q := gosql.NewUpdate().Table("passwords")
	q.Set().Add("failed_attempts = failed_attempts + 1")
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "password_hash", "failed_attempts", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.UserID, &password.PasswordHash, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Password.IncreaseFailedAttempts() Scan: %w", err)
	}
	return password, nil
}
func (r *Password) ResetFailedAttempts(ctx context.Context, id int32) (*entity.Password, error) {
	password := &entity.Password{
		ID: id,
	}

	q := gosql.NewUpdate().Table("passwords")
	q.Set().Add("failed_attempts = 0")
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "password_hash", "failed_attempts", "created_at", "updated_at")
	row := r.db.DB.QueryRowContext(ctx, q.String(), q.GetArguments()...)

	err := row.Scan(&password.UserID, &password.PasswordHash, &password.FailedAttempts, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("Passowrd.ResetFailedAttempts() Scan: %w", err)
	}
	return password, nil
}

var (
	ErrorNotFound = errors.New("password not found")
)

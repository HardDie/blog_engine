package repository

import (
	"database/sql"
	"errors"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IPassword interface {
	Create(userID int32, passwordHash string) (*entity.Password, error)
	GetByUserID(userID int32) (*entity.Password, error)
	Update(id int32, passwordHash string) (*entity.Password, error)
}

type Password struct {
	db *db.DB
}

func NewPassword(db *db.DB) *Password {
	return &Password{
		db: db,
	}
}

func (r *Password) Create(userID int32, passwordHash string) (*entity.Password, error) {
	password := &entity.Password{
		UserID:       userID,
		PasswordHash: passwordHash,
	}

	q := gosql.NewInsert().Into("passwords")
	q.Columns().Add("user_id", "password_hash")
	q.Columns().Arg(userID, passwordHash)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&password.ID, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return password, nil
}
func (r *Password) GetByUserID(userID int32) (*entity.Password, error) {
	password := &entity.Password{
		UserID: userID,
	}

	q := gosql.NewSelect().From("passwords")
	q.Columns().Add("id", "password_hash", "created_at", "updated_at")
	q.Where().AddExpression("user_id = ?", userID)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&password.ID, &password.PasswordHash, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return password, nil
}
func (r *Password) Update(id int32, passwordHash string) (*entity.Password, error) {
	password := &entity.Password{
		ID:           id,
		PasswordHash: passwordHash,
	}

	q := gosql.NewUpdate().Table("passwords")
	q.Set().Append("password_hash = ?", passwordHash)
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("user_id", "created_at", "updated_at")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&password.UserID, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return password, nil
}

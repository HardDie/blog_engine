package repository

import (
	"database/sql"
	"errors"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IUser interface {
	GetByName(name string) (*entity.User, error)
	Create(name string, invitedByUserID int32) (*entity.User, error)
}

type User struct {
	db *db.DB
}

func NewUser(db *db.DB) *User {
	return &User{
		db: db,
	}
}

func (r *User) GetByName(name string) (*entity.User, error) {
	user := &entity.User{
		Username: &name,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("id", "displayed_name", "email", "invited_by_user", "created_at", "updated_at", "deleted_at")
	q.Where().AddExpression("username = ?", name)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&user.ID, &user.DisplayedName, &user.Email, &user.InvitedByUserID, &user.CreatedAt,
		&user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (r *User) Create(name string, invitedByUserID int32) (*entity.User, error) {
	user := &entity.User{
		Username:        &name,
		InvitedByUserID: &invitedByUserID,
	}

	q := gosql.NewInsert().Into("users")
	q.Columns().Add("username", "invited_by_user")
	q.Columns().Arg(name, invitedByUserID)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

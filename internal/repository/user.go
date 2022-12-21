package repository

import (
	"database/sql"
	"errors"

	"github.com/dimonrus/gosql"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/dto"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IUser interface {
	GetByID(id int32, showPrivateInfo bool) (*entity.User, error)
	GetByName(name string) (*entity.User, error)
	Create(name string, invitedByUserID int32) (*entity.User, error)
	Update(req *dto.UpdateProfileDTO, id int32) (*entity.User, error)
}

type User struct {
	db *db.DB
}

func NewUser(db *db.DB) *User {
	return &User{
		db: db,
	}
}

func (r *User) GetByID(id int32, showPrivateInfo bool) (*entity.User, error) {
	user := &entity.User{
		ID: &id,
	}

	q := gosql.NewSelect().From("users")
	q.Columns().Add("displayed_name", "invited_by_user", "created_at", "updated_at", "deleted_at")
	if showPrivateInfo {
		q.Columns().Add("username", "email")
	}
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	var err error
	if showPrivateInfo {
		err = row.Scan(&user.DisplayedName, &user.InvitedByUserID, &user.CreatedAt,
			&user.UpdatedAt, &user.DeletedAt, &user.Username, &user.Email)
	} else {
		err = row.Scan(&user.DisplayedName, &user.InvitedByUserID, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt)
	}
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil

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
	q.Columns().Add("username", "displayed_name", "invited_by_user")
	q.Columns().Arg(name, name, invitedByUserID)
	q.Returning().Add("id", "created_at", "updated_at")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (r *User) Update(req *dto.UpdateProfileDTO, id int32) (*entity.User, error) {
	user := &entity.User{
		ID:            &id,
		DisplayedName: &req.DisplayedName,
		Email:         &req.Email,
	}

	q := gosql.NewUpdate().Table("users")
	q.Set().Append("displayed_name = ?", req.DisplayedName)
	q.Set().Append("email = ?", req.Email)
	q.Set().Append("updated_at = datetime('now')")
	q.Where().AddExpression("id = ?", id)
	q.Where().AddExpression("deleted_at IS NULL")
	q.Returning().Add("username", "invited_by_user", "created_at", "updated_at")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&user.Username, &user.InvitedByUserID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

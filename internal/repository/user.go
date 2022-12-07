package repository

import (
	"database/sql"
	"errors"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IUser interface {
	GetByName(name string) (*entity.User, error)
	Create(name string) (*entity.User, error)
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
	user := &entity.User{}

	row := r.db.DB.QueryRow(`
SELECT id, username, displayed_name, email, invited_by_user, created_at, updated_at, deleted_at
FROM users
WHERE username = $1 AND deleted_at IS NULL`, name)

	err := row.Scan(&user.ID, &user.Username, &user.DisplayedName, &user.Email, &user.InvitedByUserID, &user.CreatedAt,
		&user.UpdatedAt, &user.DeletedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}
func (r *User) Create(name string) (*entity.User, error) {
	user := &entity.User{
		Username: &name,
	}

	row := r.db.DB.QueryRow(`
INSERT INTO users (username)
VALUES ($1)
RETURNING id, created_at, updated_at`, name)

	err := row.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

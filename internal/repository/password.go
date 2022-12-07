package repository

import (
	"database/sql"
	"errors"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type IPassword interface {
	Create(userID int32, passwordHash string) (*entity.Password, error)
	GetByUserID(userID int32) (*entity.Password, error)
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

	row := r.db.DB.QueryRow(`
INSERT INTO passwords (user_id, password_hash)
VALUES ($1, $2)
RETURNING id, created_at, updated_at`, userID, passwordHash)

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

	row := r.db.DB.QueryRow(`
SELECT id, password_hash, created_at, updated_at
FROM passwords
WHERE user_id = $1 AND deleted_at IS NULL`, userID)

	err := row.Scan(&password.ID, &password.PasswordHash, &password.CreatedAt, &password.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return password, nil
}

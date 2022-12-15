package repository

import (
	"github.com/dimonrus/gosql"

	"github.com/HardDie/blog_engine/internal/db"
	"github.com/HardDie/blog_engine/internal/entity"
)

type ISession interface {
	CreateOrUpdate(userID int32, sessionHash string) (*entity.Session, error)
	GetUserID(sessionHash string) (*int32, error)
}

type Session struct {
	db *db.DB
}

func NewSession(db *db.DB) *Session {
	return &Session{
		db: db,
	}
}

func (r *Session) CreateOrUpdate(userID int32, sessionHash string) (*entity.Session, error) {
	session := &entity.Session{
		UserID:      userID,
		SessionHash: sessionHash,
	}

	row := r.db.DB.QueryRow(`
INSERT INTO sessions (user_id, session_hash)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
SET session_hash = $2, updated_at = datetime('now')
RETURNING id, created_at, updated_at`, userID, sessionHash)

	err := row.Scan(&session.ID, &session.CreatedAt, &session.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return session, nil
}
func (r *Session) GetUserID(sessionHash string) (*int32, error) {
	var userID int32

	q := gosql.NewSelect().From("sessions")
	q.Columns().Add("user_id")
	q.Where().AddExpression("session_hash = ?", sessionHash)
	q.Where().AddExpression("deleted_at IS NULL")
	row := r.db.DB.QueryRow(q.String(), q.GetArguments()...)

	err := row.Scan(&userID)
	if err != nil {
		return nil, err
	}
	return &userID, nil
}

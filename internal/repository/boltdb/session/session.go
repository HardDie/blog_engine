package session

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/boltdb/bolt"

	"github.com/HardDie/blog_engine/internal/boltdb"
	"github.com/HardDie/blog_engine/internal/models"
)

var (
	ErrorNotFound = errors.New("session not found")
)

type Session struct {
	db *boltdb.DB
}

func New(db *boltdb.DB) *Session {
	return &Session{
		db: db,
	}
}

func (s *Session) CreateOrUpdate(_ context.Context, userID int64, sessionHash string) (*models.Session, error) {
	ses := models.Session{
		UserID:      userID,
		SessionHash: sessionHash,
		CreatedAt:   time.Now(),
	}

	var buf bytes.Buffer
	err := gob.NewEncoder(&buf).Encode(ses)
	if err != nil {
		return nil, fmt.Errorf("Session.CreateOrUpdate() Encode: %w", err)
	}

	err = s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltdb.BucketSessions))
		if b == nil {
			return fmt.Errorf("Session.CreateOrUpdate() Bucket: b == nil")
		}
		err := b.Put([]byte(sessionHash), buf.Bytes())
		if err != nil {
			return fmt.Errorf("Session.CreateOrUpdate() Put: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &ses, nil
}

func (s *Session) DeleteBySessionHash(_ context.Context, sessionHash string) error {
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltdb.BucketSessions))
		if b == nil {
			return fmt.Errorf("Session.DeleteBySessionHash() Bucket: b == nil")
		}
		err := b.Delete([]byte(sessionHash))
		if err != nil {
			return fmt.Errorf("Session.DeleteBySessionHash() Delete: %w", err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) GetBySessionHash(_ context.Context, sessionHash string) (*models.Session, error) {
	var ses models.Session
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(boltdb.BucketSessions))
		if b == nil {
			return fmt.Errorf("Session.GetBySessionHash() Bucket: b == nil")
		}
		data := b.Get([]byte(sessionHash))
		if data == nil {
			return ErrorNotFound
		}
		err := gob.NewDecoder(bytes.NewReader(data)).Decode(&ses)
		if err != nil {
			return fmt.Errorf("Session.GetBySessionHash() Decode: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ses, nil
}

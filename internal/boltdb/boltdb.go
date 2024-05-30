package boltdb

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

const BucketSessions = "sessions"

type DB struct {
	*bolt.DB
}

func Get(dbpath string) (*DB, error) {
	db, err := bolt.Open(dbpath, 0644, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("error init boltdb: %w", err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(BucketSessions))
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("error creating bucket %q: %w", BucketSessions, err)
	}
	return &DB{
		DB: db,
	}, nil
}

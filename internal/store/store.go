package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.etcd.io/bbolt"
)

var (
	recordBucket = []byte("records")
	userBucket   = []byte("users")
	eventBucket  = []byte("events")
	auditBucket  = []byte("audits")
)

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path, 0o600, &bbolt.Options{Timeout: 2 * time.Second})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	if err := s.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, bucket := range [][]byte{recordBucket, userBucket, eventBucket, auditBucket} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) Path() string { return s.path }

func encode(value any) ([]byte, error) { return json.Marshal(value) }

func decode(data []byte, target any) error { return json.Unmarshal(data, target) }

func (s *Store) withRead(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.View(fn)
}

func (s *Store) withWrite(fn func(*bbolt.Tx) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store is closed")
	}
	return s.db.Update(fn)
}

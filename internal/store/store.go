package store

import (
	"errors"
	"sync"

	"go.etcd.io/bbolt"
)

var ErrNotFound = errors.New("entity not found")

var buckets = [][]byte{[]byte("records"), []byte("audits"), []byte("workflows"), []byte("attachments"), []byte("claim_codes"), []byte("meta")}

type Store struct {
	db   *bbolt.DB
	mu   sync.RWMutex
	path string
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(path, 0600, &bbolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range buckets {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
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

func (s *Store) Put(bucket, key, value []byte) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Put(key, value) })
}

func (s *Store) Get(bucket, key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store closed")
	}
	var result []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		v := tx.Bucket(bucket).Get(key)
		if v == nil {
			return ErrNotFound
		}
		result = append([]byte(nil), v...)
		return nil
	})
	return result, err
}

func (s *Store) Delete(bucket, key []byte) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket(bucket).Delete(key) })
}

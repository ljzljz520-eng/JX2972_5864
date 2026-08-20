package store

import "go.etcd.io/bbolt"

type Txn struct{ tx *bbolt.Tx }

func (s *Store) Update(fn func(*Txn) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return ErrNotFound
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return fn(&Txn{tx: tx}) })
}

func (t *Txn) Put(bucket, key, value []byte) error { return t.tx.Bucket(bucket).Put(key, value) }

func (t *Txn) Delete(bucket, key []byte) error { return t.tx.Bucket(bucket).Delete(key) }

func (t *Txn) Get(bucket, key []byte) []byte {
	v := t.tx.Bucket(bucket).Get(key)
	if v == nil {
		return nil
	}
	return append([]byte(nil), v...)
}

func (t *Txn) ForEach(bucket []byte, fn func([]byte, []byte) error) error {
	return t.tx.Bucket(bucket).ForEach(fn)
}

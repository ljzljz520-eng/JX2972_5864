package store

import (
	"bytes"
	"go.etcd.io/bbolt"
	"sort"
	"strings"
)

type Entry struct{ Key, Value []byte }

func (s *Store) List(bucket []byte) ([]Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, ErrNotFound
	}
	entries := make([]Entry, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(bucket).ForEach(func(k, v []byte) error {
			entries = append(entries, Entry{append([]byte(nil), k...), append([]byte(nil), v...)})
			return nil
		})
	})
	return entries, err
}

func SortEntries(entries []Entry) []Entry {
	sort.Slice(entries, func(i, j int) bool { return bytes.Compare(entries[i].Key, entries[j].Key) < 0 })
	return entries
}

func Prefix(entries []Entry, prefix string) []Entry {
	filtered := make([]Entry, 0, len(entries))
	for _, e := range entries {
		if strings.HasPrefix(string(e.Key), prefix) {
			filtered = append(filtered, e)
		}
	}
	return filtered
}

func Keys(entries []Entry) []string {
	result := make([]string, 0, len(entries))
	for _, e := range entries {
		result = append(result, string(e.Key))
	}
	return result
}

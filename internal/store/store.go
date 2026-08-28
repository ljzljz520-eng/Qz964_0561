package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"os"
	"path/filepath"
	"time"
)

var buckets = [][]byte{[]byte("records"), []byte("users"), []byte("events"), []byte("audits")}

type Store struct {
	db   *bbolt.DB
	path string
}

func Open(path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	db, err := bbolt.Open(path, 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, path: path}
	if err := s.init(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) init() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, err := tx.CreateBucketIfNotExists(b); err != nil {
				return err
			}
		}
		return nil
	})
}
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
func (s *Store) Path() string { return s.path }
func putJSON(tx *bbolt.Tx, bucket, key string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(bucket)).Put([]byte(key), data)
}
func getJSON(tx *bbolt.Tx, bucket, key string, v any) error {
	b := tx.Bucket([]byte(bucket)).Get([]byte(key))
	if b == nil {
		return fmt.Errorf("%s %s not found", bucket, key)
	}
	return json.Unmarshal(append([]byte(nil), b...), v)
}
func deleteKey(tx *bbolt.Tx, bucket, key string) error {
	return tx.Bucket([]byte(bucket)).Delete([]byte(key))
}

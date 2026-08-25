package store

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"timber-safety/internal/domain"
)

func (s *Store) SaveRecord(r domain.Record) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "records", r.ID, r) })
}
func (s *Store) GetRecord(id string) (domain.Record, error) {
	var r domain.Record
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, "records", id, &r) })
	return r, err
}
func (s *Store) DeleteRecord(id string) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return deleteKey(tx, "records", id) })
}
func (s *Store) ListRecords() ([]domain.Record, error) {
	out := []domain.Record{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, v []byte) error {
			var r domain.Record
			if err := json.Unmarshal(v, &r); err != nil {
				return err
			}
			out = append(out, r)
			return nil
		})
	})
	return out, err
}
func (s *Store) UpdateRecord(id string, fn func(*domain.Record) error) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		var r domain.Record
		if err := getJSON(tx, "records", id, &r); err != nil {
			return err
		}
		if err := fn(&r); err != nil {
			return err
		}
		return putJSON(tx, "records", id, r)
	})
}
func (s *Store) RequireRecord(id string) (domain.Record, error) {
	r, err := s.GetRecord(id)
	if err != nil {
		return r, fmt.Errorf("record %s: %w", id, err)
	}
	return r, nil
}

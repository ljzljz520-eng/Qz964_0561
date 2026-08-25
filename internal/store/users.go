package store

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"timber-safety/internal/domain"
)

func (s *Store) SaveUser(u domain.User) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "users", u.ID, u) })
}
func (s *Store) GetUser(id string) (domain.User, error) {
	var u domain.User
	err := s.db.View(func(tx *bbolt.Tx) error { return getJSON(tx, "users", id, &u) })
	return u, err
}
func (s *Store) ListUsers() ([]domain.User, error) {
	out := []domain.User{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("users")).ForEach(func(_, v []byte) error {
			var u domain.User
			if err := json.Unmarshal(v, &u); err != nil {
				return err
			}
			out = append(out, u)
			return nil
		})
	})
	return out, err
}
func (s *Store) ActiveUsers() ([]domain.User, error) {
	all, err := s.ListUsers()
	if err != nil {
		return nil, err
	}
	out := all[:0]
	for _, u := range all {
		if u.Active {
			out = append(out, u)
		}
	}
	return out, nil
}

package store

import (
	"encoding/json"
	"go.etcd.io/bbolt"
	"timber-safety/internal/domain"
)

func (s *Store) SaveEvent(e domain.Event) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "events", e.ID, e) })
}
func (s *Store) ListEvents(recordID string) ([]domain.Event, error) {
	out := []domain.Event{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("events")).ForEach(func(_, v []byte) error {
			var e domain.Event
			if err := json.Unmarshal(v, &e); err != nil {
				return err
			}
			if recordID == "" || e.RecordID == recordID {
				out = append(out, e)
			}
			return nil
		})
	})
	return out, err
}
func (s *Store) SaveAudit(a domain.Audit) error {
	return s.db.Update(func(tx *bbolt.Tx) error { return putJSON(tx, "audits", a.ID, a) })
}
func (s *Store) ListAudits(targetID string) ([]domain.Audit, error) {
	out := []domain.Audit{}
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("audits")).ForEach(func(_, v []byte) error {
			var a domain.Audit
			if err := json.Unmarshal(v, &a); err != nil {
				return err
			}
			if targetID == "" || a.TargetID == targetID {
				out = append(out, a)
			}
			return nil
		})
	})
	return out, err
}
func (s *Store) Count(bucket string) (int, error) {
	n := 0
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, _ []byte) error { n++; return nil })
	})
	return n, err
}

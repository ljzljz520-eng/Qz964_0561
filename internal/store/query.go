package store

import (
	"encoding/json"
	"strings"
	"timber-safety/internal/domain"
	"timber-safety/internal/validation"
)

func (s *Store) FilterRecords(q validation.Query) ([]domain.Record, error) {
	items, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	q = q.Clean()
	out := []domain.Record{}
	for _, r := range items {
		if q.State != "" && string(r.State) != q.State {
			continue
		}
		if q.Risk != "" && string(r.RiskLevel) != q.Risk {
			continue
		}
		if q.Source != "" && !strings.Contains(strings.ToLower(r.Source), strings.ToLower(q.Source)) {
			continue
		}
		out = append(out, r)
	}
	return out, nil
}
func (s *Store) SaveConfirmations(id string, items []domain.Confirmation) error {
	return s.UpdateRecord(id, func(r *domain.Record) error {
		r.Confirmations = append([]domain.Confirmation(nil), items...)
		return nil
	})
}
func (s *Store) ExportRecords() ([]byte, error) {
	items, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	return json.Marshal(items)
}
func (s *Store) ImportRecords(data []byte) error {
	var items []domain.Record
	if err := json.Unmarshal(data, &items); err != nil {
		return err
	}
	for _, r := range items {
		if err := s.SaveRecord(r); err != nil {
			return err
		}
	}
	return nil
}
func (s *Store) HasRecord(id string) bool { _, err := s.GetRecord(id); return err == nil }

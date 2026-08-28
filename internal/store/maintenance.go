package store

import "timber-safety/internal/domain"

func (s *Store) PurgeArchived() (int, error) {
	items, e := s.ListRecords()
	if e != nil {
		return 0, e
	}
	n := 0
	for _, r := range items {
		if r.State == domain.StateArchived {
			if e := s.DeleteRecord(r.ID); e != nil {
				return n, e
			}
			n++
		}
	}
	return n, nil
}
func (s *Store) RecordsByState(state domain.RecordState) ([]domain.Record, error) {
	items, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range items {
		if r.State == state {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) RecordsByRisk(level domain.RiskLevel) ([]domain.Record, error) {
	items, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range items {
		if r.RiskLevel == level {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Store) CountConfirmations(id string) (int, error) {
	r, e := s.GetRecord(id)
	if e != nil {
		return 0, e
	}
	return len(r.Confirmations), nil
}
func (s *Store) UpdateState(id string, state domain.RecordState) error {
	return s.UpdateRecord(id, func(r *domain.Record) error { r.State = state; return nil })
}
func (s *Store) UpdateRisk(id string, risk domain.RiskLevel) error {
	return s.UpdateRecord(id, func(r *domain.Record) error { r.RiskLevel = risk; return nil })
}
func (s *Store) IDs() ([]string, error) {
	items, e := s.ListRecords()
	if e != nil {
		return nil, e
	}
	out := []string{}
	for _, r := range items {
		out = append(out, r.ID)
	}
	return out, nil
}

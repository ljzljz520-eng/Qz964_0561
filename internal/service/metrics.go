package service

import (
	"timber-safety/internal/domain"
	"timber-safety/internal/validation"
	"time"
)

type Metrics struct {
	Total, Open, Reviewed, Rejected int
	Risk                            map[domain.RiskLevel]int
	Sources                         map[string]int
	UpdatedAt                       time.Time
}

func (s *Service) Metrics(q validation.Query) (Metrics, error) {
	items, e := s.Search(q)
	if e != nil {
		return Metrics{}, e
	}
	m := Metrics{Risk: map[domain.RiskLevel]int{}, Sources: map[string]int{}, UpdatedAt: s.now()}
	for _, r := range items {
		m.Total++
		m.Risk[r.RiskLevel]++
		m.Sources[r.Source]++
		if r.IsOpen() {
			m.Open++
		}
		if r.State == domain.StateReviewed || r.State == domain.StateArchived {
			m.Reviewed++
		}
		if r.State == domain.StateRejected {
			m.Rejected++
		}
	}
	return m, nil
}
func (s *Service) StateCounts(q validation.Query) (map[domain.RecordState]int, error) {
	items, e := s.Search(q)
	if e != nil {
		return nil, e
	}
	out := map[domain.RecordState]int{}
	for _, r := range items {
		out[r.State]++
	}
	return out, nil
}
func (s *Service) SourceCounts(q validation.Query) (map[string]int, error) {
	items, e := s.Search(q)
	if e != nil {
		return nil, e
	}
	out := map[string]int{}
	for _, r := range items {
		out[r.Source]++
	}
	return out, nil
}
func (s *Service) RiskLevels(q validation.Query) ([]domain.RiskLevel, error) {
	items, e := s.Search(q)
	if e != nil {
		return nil, e
	}
	seen := map[domain.RiskLevel]bool{}
	out := []domain.RiskLevel{}
	for _, r := range items {
		if !seen[r.RiskLevel] {
			seen[r.RiskLevel] = true
			out = append(out, r.RiskLevel)
		}
	}
	return out, nil
}
func (s *Service) RecordsNeedingReview(q validation.Query) ([]domain.Record, error) {
	items, e := s.Search(q)
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range items {
		if domain.RequiresSecondReview(r) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Service) RecordsForArchive(q validation.Query) ([]domain.Record, error) {
	items, e := s.Search(q)
	if e != nil {
		return nil, e
	}
	out := []domain.Record{}
	for _, r := range items {
		if domain.CanArchive(r) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (s *Service) ApplyChecklist(id string, c domain.Checklist) error {
	if !c.Complete() {
		return validation.ValidateRequiredReview(c.PassedCount())
	}
	return s.Touch(id)
}

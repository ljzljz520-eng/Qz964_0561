package service

import (
	"sort"
	"timber-safety/internal/domain"
	"timber-safety/internal/validation"
)

type Summary struct {
	Total   int
	ByState map[domain.RecordState]int
	ByRisk  map[domain.RiskLevel]int
	Records []domain.Record
}

func (s *Service) Summary(q validation.Query) (Summary, error) {
	items, err := s.Search(q)
	if err != nil {
		return Summary{}, err
	}
	sum := Summary{ByState: map[domain.RecordState]int{}, ByRisk: map[domain.RiskLevel]int{}, Records: items}
	for _, r := range items {
		sum.Total++
		sum.ByState[r.State]++
		sum.ByRisk[r.RiskLevel]++
	}
	return sum, nil
}
func (s *Service) TopRisks(limit int) ([]domain.Record, error) {
	items, err := s.Search(validation.Query{})
	if err != nil {
		return nil, err
	}
	sort.Slice(items, func(i, j int) bool {
		return domain.DefectSeverity(items[i].Measurements.Defects) > domain.DefectSeverity(items[j].Measurements.Defects)
	})
	if limit < 1 {
		limit = 5
	}
	if len(items) > limit {
		items = items[:limit]
	}
	return items, nil
}
func (s *Service) Recalculate(id string) (domain.Record, error) {
	var out domain.Record
	err := s.Store.UpdateRecord(id, func(r *domain.Record) error { r.RefreshRisk(); out = *r; return nil })
	return out, err
}

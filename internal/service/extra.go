package service

import (
	"timber-safety/internal/domain"
	"timber-safety/internal/validation"
	"time"
)

func structQuery() validation.Query { return validation.Query{} }
func (s *Service) Archive(id string, actor domain.User) error {
	return s.Transition(id, domain.StateArchived, actor)
}
func (s *Service) Reject(id string, actor domain.User) error {
	return s.Transition(id, domain.StateRejected, actor)
}
func (s *Service) SetClock(clock func() time.Time) { s.Clock = clock }
func (s *Service) ReviewLabel(id string) (string, error) {
	r, e := s.Get(id)
	if e != nil {
		return "", e
	}
	return domain.ReviewLabel(r), nil
}
func (s *Service) OpenCount() (int, error) {
	items, e := s.Search(validation.Query{})
	if e != nil {
		return 0, e
	}
	n := 0
	for _, r := range items {
		if r.IsOpen() {
			n++
		}
	}
	return n, nil
}
func (s *Service) RiskCount(level domain.RiskLevel) (int, error) {
	items, e := s.Search(validation.Query{Risk: string(level)})
	return len(items), e
}

package service

import (
	"timber-safety/internal/domain"
	"timber-safety/internal/validation"
	"time"
)

func (s *Service) BuildLifecycle(id string) (domain.Lifecycle, error) {
	r, e := s.Get(id)
	if e != nil {
		return domain.Lifecycle{}, e
	}
	l := domain.Lifecycle{RecordID: id, Created: &r.CreatedAt}
	l.Mark(r.State, r.UpdatedAt)
	return l, nil
}
func (s *Service) MarkChecklist(id string, c domain.Checklist) error {
	if e := validation.ValidateChecklist(c); e != nil {
		return e
	}
	if !c.Complete() {
		return validation.ValidateNote("checklist incomplete")
	}
	return s.Touch(id)
}
func (s *Service) Age(id string) (time.Duration, error) {
	l, e := s.BuildLifecycle(id)
	if e != nil {
		return 0, e
	}
	return l.TotalDuration(), nil
}
func (s *Service) IsStale(id string, ttl time.Duration) (bool, error) {
	r, e := s.Get(id)
	if e != nil {
		return false, e
	}
	return time.Since(r.UpdatedAt) > ttl, nil
}

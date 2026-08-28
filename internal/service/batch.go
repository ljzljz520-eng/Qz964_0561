package service

import (
	"fmt"
	"timber-safety/internal/domain"
	"timber-safety/internal/notify"
	"time"
)

func (s *Service) RegisterBatch(records []domain.Record, actor domain.User) (int, []error) {
	ok := 0
	errs := []error{}
	for i, r := range records {
		if err := s.Register(r, actor); err != nil {
			errs = append(errs, fmt.Errorf("record %d: %w", i, err))
		} else {
			ok++
		}
	}
	return ok, errs
}
func (s *Service) ProcessBatch(ids []string, actor domain.User) (int, []error) {
	ok := 0
	errs := []error{}
	for _, id := range ids {
		if err := s.Transition(id, domain.StateProcessing, actor); err != nil {
			errs = append(errs, err)
		} else {
			ok++
		}
	}
	return ok, errs
}
func (s *Service) ConfirmBatch(id string, inspectors []domain.User, decision string) (int, []error) {
	ok := 0
	errs := []error{}
	for _, u := range inspectors {
		if err := s.Confirm(id, u, decision, ""); err != nil {
			errs = append(errs, err)
		} else {
			ok++
		}
	}
	return ok, errs
}
func (s *Service) NotifyOpen(n *notify.Notifier) (int, error) {
	items, err := s.Search(structQuery())
	if err != nil {
		return 0, err
	}
	count := 0
	for _, r := range items {
		if r.IsOpen() && domain.ShouldNotify(r) {
			if err := n.NotifyRecord(r); err != nil {
				return count, err
			}
			count++
		}
	}
	return count, nil
}
func (s *Service) Touch(id string) error {
	return s.Store.UpdateRecord(id, func(r *domain.Record) error { r.UpdatedAt = time.Now().UTC(); return nil })
}

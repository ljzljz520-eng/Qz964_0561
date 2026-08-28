package service

import (
	"fmt"
	"sync"
	"timber-safety/internal/domain"
	"timber-safety/internal/store"
	"timber-safety/internal/validation"
	"time"
)

type Service struct {
	Store       *store.Store
	Clock       func() time.Time
	mu          sync.Mutex
	reviewDelay time.Duration
}

func New(s *store.Store) *Service {
	return &Service{Store: s, Clock: func() time.Time { return time.Now().UTC() }}
}
func (s *Service) now() time.Time {
	if s.Clock == nil {
		return time.Now().UTC()
	}
	return s.Clock()
}
func (s *Service) SetReviewDelay(d time.Duration) { s.reviewDelay = d }
func (s *Service) Register(r domain.Record, actor domain.User) error {
	if !validation.CanRegister(actor) {
		return fmt.Errorf("actor cannot register")
	}
	r.Normalize()
	if err := validation.ValidateRecord(r); err != nil {
		return err
	}
	r.State = domain.StateRegistered
	r.RiskLevel = r.CalculateRisk()
	r.UpdatedAt = s.now()
	if err := s.Store.SaveRecord(r); err != nil {
		return err
	}
	return s.Store.SaveEvent(domain.NewEvent(r.ID+"-registered", r.ID, "registered", actor.ID, "record accepted", s.now()))
}
func (s *Service) Transition(id string, next domain.RecordState, actor domain.User) error {
	if !validation.CanReview(actor) {
		return fmt.Errorf("actor cannot process")
	}
	var moved domain.Record
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.Store.UpdateRecord(id, func(r *domain.Record) error {
		if err := validation.ValidateTransition(*r, next); err != nil {
			return err
		}
		if !r.Move(next, s.now()) {
			return fmt.Errorf("transition refused")
		}
		moved = *r
		return nil
	})
	if err != nil {
		return err
	}
	return s.Store.SaveEvent(domain.NewEvent(fmt.Sprintf("%s-%d", id, moved.Version), id, "transition", actor.ID, string(next), s.now()))
}
func (s *Service) Confirm(id string, inspector domain.User, decision, notes string) error {
	if !validation.CanReview(inspector) {
		return fmt.Errorf("actor cannot review")
	}
	if err := validation.ValidateDecision(decision); err != nil {
		return err
	}
	if s.reviewDelay > 0 {
		time.Sleep(s.reviewDelay)
	}
	snapshot, err := s.Store.GetRecord(id)
	if err != nil {
		return err
	}
	snapshot.Confirmations = append(snapshot.Confirmations, domain.Confirmation{RecordID: id, InspectorID: inspector.ID, Decision: domain.NormalizeDecision(decision), Notes: notes, ConfirmedAt: s.now()})
	snapshot.State = snapshot.ReviewState(2)
	snapshot.Version++
	snapshot.UpdatedAt = s.now()
	if err := s.Store.SaveRecord(snapshot); err != nil {
		return err
	}
	return s.Store.SaveAudit(domain.NewAudit(id+"-"+inspector.ID, inspector.ID, "confirm", id, decision, s.now()))
}
func (s *Service) ProcessAndConfirm(id string, inspector domain.User, decision string) error {
	if err := s.Transition(id, domain.StateProcessing, inspector); err != nil {
		return err
	}
	return s.Confirm(id, inspector, decision, "")
}
func (s *Service) Get(id string) (domain.Record, error) { return s.Store.RequireRecord(id) }
func (s *Service) Search(q validation.Query) ([]domain.Record, error) {
	if err := validation.ValidateQuery(q); err != nil {
		return nil, err
	}
	items, err := s.Store.ListRecords()
	if err != nil {
		return nil, err
	}
	q = q.Clean()
	out := make([]domain.Record, 0, len(items))
	for _, r := range items {
		if q.State != "" && string(r.State) != q.State {
			continue
		}
		if q.Risk != "" && string(r.RiskLevel) != q.Risk {
			continue
		}
		if q.Source != "" && !validation.MatchText(r.Source, q.Source) {
			continue
		}
		out = append(out, r)
	}
	return domain.SortRecords(out), nil
}
func (s *Service) AuditTrail(id string) ([]domain.Audit, []domain.Event, error) {
	a, err := s.Store.ListAudits(id)
	if err != nil {
		return nil, nil, err
	}
	e, err := s.Store.ListEvents(id)
	return a, e, err
}

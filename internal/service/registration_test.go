package service

import (
	"testing"
	"timber-safety/internal/domain"
	"timber-safety/internal/store"
	"time"
)

func TestRegistrationRejectsInvalid(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/db")
	defer s.Close()
	svc := New(s)
	r := domain.NewRecord("r", "bad", "松", "场", domain.Measurements{Length: 1, Width: 1, Height: 1}, time.Now())
	if svc.Register(r, domain.User{ID: "u", Role: "registrar", Active: true}) == nil {
		t.Fatal("expected error")
	}
}

package service

import (
	"testing"
	"timber-safety/internal/domain"
	"timber-safety/internal/store"
	"timber-safety/internal/validation"
	"time"
)

func TestQuerySummary(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/db")
	defer s.Close()
	svc := New(s)
	r := domain.NewRecord("r", "木材35", "松", "场", domain.Measurements{Length: 1, Width: 1, Height: 1}, time.Now())
	_ = svc.Register(r, domain.User{ID: "u", Role: "registrar", Active: true})
	sum, e := svc.Summary(validation.Query{})
	if e != nil || sum.Total != 1 {
		t.Fatal(sum, e)
	}
}

package store

import (
	"testing"
	"timber-safety/internal/domain"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/db"
	s, e := Open(path)
	if e != nil {
		t.Fatal(e)
	}
	r := domain.NewRecord("r", "木材35", "松", "场", domain.Measurements{Length: 1, Width: 1, Height: 1}, time.Now())
	if e = s.SaveRecord(r); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(path)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.GetRecord("r")
	if e != nil || got.TimberCode != "木材35" {
		t.Fatal(got, e)
	}
}

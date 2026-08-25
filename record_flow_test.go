package main

import (
	"sync"
	"testing"
	"timber-safety/internal/domain"
	"timber-safety/internal/service"
	"timber-safety/internal/store"
	"timber-safety/internal/validation"
	"time"
)

func testService(t *testing.T) *service.Service {
	t.Helper()
	s, e := store.Open(t.TempDir() + "/db")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { s.Close() })
	return service.New(s)
}
func baseRecord() domain.Record {
	return domain.NewRecord("r35", "木材35", "松木", "林场A", domain.Measurements{Length: 300, Width: 20, Height: 20, Moisture: 12}, time.Now().UTC())
}
func actor(id, role string) domain.User {
	return domain.User{ID: id, Name: id, Role: role, Active: true}
}
func TestRecordFlow35(t *testing.T) {
	svc := testService(t)
	if err := svc.Register(baseRecord(), actor("u1", "registrar")); err != nil {
		t.Fatal(err)
	}
	ins := actor("i1", "inspector")
	if err := svc.Confirm("r35", ins, "approve", ""); err != nil {
		t.Fatal(err)
	}
	r, _ := svc.Get("r35")
	if len(r.Confirmations) != 1 {
		t.Fatalf("confirmations=%d", len(r.Confirmations))
	}
}
func TestWorkflowOne(t *testing.T) {
	svc := testService(t)
	if err := svc.Register(baseRecord(), actor("u", "registrar")); err != nil {
		t.Fatal(err)
	}
	r, e := svc.Get("r35")
	if e != nil || r.State != domain.StateRegistered {
		t.Fatal(r, e)
	}
}
func TestWorkflowTwo(t *testing.T) {
	svc := testService(t)
	if err := svc.Register(baseRecord(), actor("u", "registrar")); err != nil {
		t.Fatal(err)
	}
	if err := svc.Transition("r35", domain.StateProcessing, actor("i", "inspector")); err != nil {
		t.Fatal(err)
	}
}
func TestWorkflowThree(t *testing.T) {
	svc := testService(t)
	if err := svc.Register(baseRecord(), actor("u", "registrar")); err != nil {
		t.Fatal(err)
	}
	sum, e := svc.Summary(validation.Query{})
	if e != nil || sum.Total != 1 {
		t.Fatal(sum, e)
	}
}
func TestConcurrentConfirmationsExposeLostUpdate(t *testing.T) {
	svc := testService(t)
	svc.SetReviewDelay(time.Millisecond)
	if err := svc.Register(baseRecord(), actor("u", "registrar")); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	for _, id := range []string{"i1", "i2"} {
		wg.Add(1)
		go func(id string) { defer wg.Done(); _ = svc.Confirm("r35", actor(id, "inspector"), "approve", "") }(id)
	}
	wg.Wait()
	r, _ := svc.Get("r35")
	if len(r.Confirmations) != 2 {
		t.Fatalf("parallel confirmations lost: got %d", len(r.Confirmations))
	}
}

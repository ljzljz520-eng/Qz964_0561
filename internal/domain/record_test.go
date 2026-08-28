package domain

import (
	"testing"
	"time"
)

func TestRiskAndTransitions(t *testing.T) {
	r := NewRecord("r", "木材35", "松", "场", Measurements{Length: 700, Width: 1, Height: 1, Moisture: 30, Defects: []string{"裂纹", "腐朽", "虫蛀"}}, time.Now())
	if r.CalculateRisk() != RiskCritical {
		t.Fatal(r.CalculateRisk())
	}
	if !r.Move(StateRegistered, time.Now()) {
		t.Fatal("move")
	}
	if r.Move(StateArchived, time.Now()) {
		t.Fatal("invalid move")
	}
}

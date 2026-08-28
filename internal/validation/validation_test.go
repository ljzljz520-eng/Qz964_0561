package validation

import (
	"testing"
	"timber-safety/internal/domain"
	"time"
)

func TestValidation(t *testing.T) {
	r := domain.NewRecord("r", "木材35", "松", "场", domain.Measurements{Length: 1, Width: 1, Height: 1}, time.Now())
	if e := ValidateRecord(r); e != nil {
		t.Fatal(e)
	}
	if ValidateCode("bad") {
		t.Fatal("code")
	}
}

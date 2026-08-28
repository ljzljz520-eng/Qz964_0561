package validation

import (
	"fmt"
	"strings"
	"timber-safety/internal/domain"
)

func ValidateTransition(r domain.Record, next domain.RecordState) error {
	if !r.CanMove(next) {
		return fmt.Errorf("cannot move %s to %s", r.State, next)
	}
	return nil
}
func ValidateInspector(id string) error {
	if strings.TrimSpace(id) == "" {
		return fmt.Errorf("inspector is required")
	}
	if len(id) > 64 {
		return fmt.Errorf("inspector is too long")
	}
	return nil
}
func ValidateDecision(v string) error {
	n := domain.NormalizeDecision(v)
	if n == "pending" {
		return fmt.Errorf("decision must be approve or reject")
	}
	return nil
}
func ValidateQueryState(v string) error {
	if v == "" {
		return nil
	}
	for _, s := range domain.AllowedStates() {
		if string(s) == v {
			return nil
		}
	}
	return fmt.Errorf("invalid state filter")
}
func ValidateRequiredReview(n int) error {
	if n < 1 || n > 20 {
		return fmt.Errorf("required reviewers out of range")
	}
	return nil
}

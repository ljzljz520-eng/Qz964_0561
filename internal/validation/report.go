package validation

import (
	"fmt"
	"strings"
)

type Query struct {
	State, Risk, Source string
	Limit               int
}

func (q Query) Clean() Query {
	q.State = strings.TrimSpace(q.State)
	q.Risk = strings.TrimSpace(q.Risk)
	q.Source = strings.TrimSpace(q.Source)
	if q.Limit <= 0 {
		q.Limit = 100
	}
	if q.Limit > 1000 {
		q.Limit = 1000
	}
	return q
}
func ValidateQuery(q Query) error {
	q = q.Clean()
	if err := ValidateQueryState(q.State); err != nil {
		return err
	}
	if q.Risk != "" && q.Risk != "low" && q.Risk != "medium" && q.Risk != "high" && q.Risk != "critical" {
		return fmt.Errorf("invalid risk filter")
	}
	return nil
}
func ValidatePagination(offset, limit int) error {
	if offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	if limit < 1 || limit > 1000 {
		return fmt.Errorf("limit must be 1..1000")
	}
	return nil
}
func MatchText(value, needle string) bool {
	if needle == "" {
		return true
	}
	return strings.Contains(strings.ToLower(value), strings.ToLower(needle))
}

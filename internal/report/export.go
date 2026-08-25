package report

import (
	"encoding/json"
	"timber-safety/internal/domain"
)

func (b Builder) JSON(records []domain.Record) ([]byte, error) {
	return json.Marshal(map[string]any{"title": b.Title(), "records": records})
}
func (b Builder) Filter(records []domain.Record, level domain.RiskLevel) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if level == "" || r.RiskLevel == level {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) CountOpen(records []domain.Record) int {
	n := 0
	for _, r := range records {
		if r.IsOpen() {
			n++
		}
	}
	return n
}
func (b Builder) Latest(records []domain.Record) domain.Record {
	var latest domain.Record
	for _, r := range records {
		if r.UpdatedAt.After(latest.UpdatedAt) {
			latest = r
		}
	}
	return latest
}

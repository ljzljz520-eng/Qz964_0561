package report

import "timber-safety/internal/domain"

type Summary struct {
	Total, Open, HighRisk int
	States                map[domain.RecordState]int
	Risks                 map[domain.RiskLevel]int
}

func (b Builder) BuildSummary(records []domain.Record) Summary {
	s := Summary{States: map[domain.RecordState]int{}, Risks: map[domain.RiskLevel]int{}}
	for _, r := range records {
		s.Total++
		s.States[r.State]++
		s.Risks[r.RiskLevel]++
		if r.IsOpen() {
			s.Open++
		}
		if r.RiskLevel == domain.RiskHigh || r.RiskLevel == domain.RiskCritical {
			s.HighRisk++
		}
	}
	return s
}
func (b Builder) PercentOpen(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	n := 0
	for _, r := range records {
		if r.IsOpen() {
			n++
		}
	}
	return float64(n) * 100 / float64(len(records))
}
func (b Builder) PercentReviewed(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	n := 0
	for _, r := range records {
		if r.State == domain.StateReviewed || r.State == domain.StateArchived {
			n++
		}
	}
	return float64(n) * 100 / float64(len(records))
}
func (b Builder) RiskOrder() []domain.RiskLevel {
	return []domain.RiskLevel{domain.RiskCritical, domain.RiskHigh, domain.RiskMedium, domain.RiskLow}
}
func (b Builder) StateOrder() []domain.RecordState {
	return []domain.RecordState{domain.StateDraft, domain.StateRegistered, domain.StateProcessing, domain.StateReviewed, domain.StateArchived, domain.StateRejected}
}

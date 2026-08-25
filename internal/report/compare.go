package report

import "timber-safety/internal/domain"

type Comparison struct {
	Left, Right                domain.Record
	MoistureDelta, VolumeDelta float64
	SameSource, SameRisk       bool
}

func (b Builder) Compare(left, right domain.Record) Comparison {
	return Comparison{Left: left, Right: right, MoistureDelta: right.Measurements.Moisture - left.Measurements.Moisture, VolumeDelta: right.Measurements.Volume() - left.Measurements.Volume(), SameSource: left.Source == right.Source, SameRisk: left.RiskLevel == right.RiskLevel}
}
func (b Builder) MostMoist(records []domain.Record) (domain.Record, bool) {
	if len(records) == 0 {
		return domain.Record{}, false
	}
	best := records[0]
	for _, r := range records[1:] {
		if r.Measurements.Moisture > best.Measurements.Moisture {
			best = r
		}
	}
	return best, true
}
func (b Builder) Largest(records []domain.Record) (domain.Record, bool) {
	if len(records) == 0 {
		return domain.Record{}, false
	}
	best := records[0]
	for _, r := range records[1:] {
		if r.Measurements.Volume() > best.Measurements.Volume() {
			best = r
		}
	}
	return best, true
}
func (b Builder) MostConfirmed(records []domain.Record) (domain.Record, bool) {
	if len(records) == 0 {
		return domain.Record{}, false
	}
	best := records[0]
	for _, r := range records[1:] {
		if len(r.Confirmations) > len(best.Confirmations) {
			best = r
		}
	}
	return best, true
}
func (b Builder) AverageDefects(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	n := 0
	for _, r := range records {
		n += len(r.Measurements.Defects)
	}
	return float64(n) / float64(len(records))
}
func (b Builder) AverageVolume(records []domain.Record) float64 {
	if len(records) == 0 {
		return 0
	}
	n := 0.0
	for _, r := range records {
		n += r.Measurements.Volume()
	}
	return n / float64(len(records))
}
func (b Builder) RiskDistribution(records []domain.Record) []int {
	out := make([]int, 4)
	for _, r := range records {
		switch r.RiskLevel {
		case domain.RiskLow:
			out[0]++
		case domain.RiskMedium:
			out[1]++
		case domain.RiskHigh:
			out[2]++
		case domain.RiskCritical:
			out[3]++
		}
	}
	return out
}
func (b Builder) StateDistribution(records []domain.Record) []int {
	out := make([]int, 6)
	for _, r := range records {
		switch r.State {
		case domain.StateDraft:
			out[0]++
		case domain.StateRegistered:
			out[1]++
		case domain.StateProcessing:
			out[2]++
		case domain.StateReviewed:
			out[3]++
		case domain.StateArchived:
			out[4]++
		case domain.StateRejected:
			out[5]++
		}
	}
	return out
}
func (b Builder) Stable(records []domain.Record) bool {
	for _, r := range records {
		if r.Version < 1 || r.ID == "" {
			return false
		}
	}
	return true
}
func (b Builder) ValidForExport(records []domain.Record) bool {
	for _, r := range records {
		if r.TimberCode == "" || r.Source == "" || !r.Measurements.IsComplete() {
			return false
		}
	}
	return true
}

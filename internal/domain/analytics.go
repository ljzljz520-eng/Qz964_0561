package domain

import (
	"math"
	"sort"
)

type MeasurementStats struct {
	AverageMoisture, AverageLength float64
	MaxMoisture, MinMoisture       float64
	Samples                        int
}

func BuildMeasurementStats(records []Record) MeasurementStats {
	s := MeasurementStats{}
	if len(records) == 0 {
		return s
	}
	s.MinMoisture = math.MaxFloat64
	for _, r := range records {
		s.Samples++
		s.AverageMoisture += r.Measurements.Moisture
		s.AverageLength += r.Measurements.Length
		if r.Measurements.Moisture > s.MaxMoisture {
			s.MaxMoisture = r.Measurements.Moisture
		}
		if r.Measurements.Moisture < s.MinMoisture {
			s.MinMoisture = r.Measurements.Moisture
		}
	}
	s.AverageMoisture /= float64(s.Samples)
	s.AverageLength /= float64(s.Samples)
	return s
}
func RiskScore(r Record) int {
	score := DefectSeverity(r.Measurements.Defects)
	if r.Measurements.Moisture > 24 {
		score += 3
	} else if r.Measurements.Moisture > 16 {
		score++
	}
	if r.Measurements.Length > 600 {
		score++
	}
	return score
}
func RankByRisk(records []Record) []Record {
	out := append([]Record(nil), records...)
	sort.SliceStable(out, func(i, j int) bool { return RiskScore(out[i]) > RiskScore(out[j]) })
	return out
}
func StateProgression(r Record) []RecordState {
	p := []RecordState{StateDraft}
	if r.State == StateDraft {
		return p
	}
	p = append(p, StateRegistered)
	if r.State == StateRegistered {
		return p
	}
	p = append(p, StateProcessing)
	if r.State == StateProcessing {
		return p
	}
	p = append(p, StateReviewed)
	if r.State == StateReviewed {
		return p
	}
	return append(p, StateArchived)
}
func IsTerminal(s RecordState) bool      { return s == StateArchived || s == StateRejected }
func RequiresSecondReview(r Record) bool { return RiskScore(r) >= 4 || len(r.Confirmations) < 2 }
func ConfirmationNames(r Record) []string {
	out := make([]string, 0, len(r.Confirmations))
	for _, c := range r.Confirmations {
		out = append(out, c.InspectorID)
	}
	sort.Strings(out)
	return out
}
func DistinctSources(records []Record) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, r := range records {
		if !seen[r.Source] {
			seen[r.Source] = true
			out = append(out, r.Source)
		}
	}
	sort.Strings(out)
	return out
}

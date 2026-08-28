package domain

import "time"

func (r Record) CanMove(to RecordState) bool {
	if r.State == StateDraft && to == StateRegistered {
		return true
	}
	if r.State == StateRegistered && (to == StateProcessing || to == StateRejected) {
		return true
	}
	if r.State == StateProcessing && (to == StateReviewed || to == StateRejected) {
		return true
	}
	if r.State == StateReviewed && to == StateArchived {
		return true
	}
	return false
}
func (r *Record) Move(to RecordState, now time.Time) bool {
	if !r.CanMove(to) {
		return false
	}
	r.State = to
	r.UpdatedAt = now
	r.Version++
	return true
}
func (r Record) CalculateRisk() RiskLevel {
	score := 0
	if r.Measurements.Moisture > 24 {
		score += 3
	} else if r.Measurements.Moisture > 16 {
		score += 1
	}
	if len(r.Measurements.Defects) >= 3 {
		score += 3
	} else if len(r.Measurements.Defects) > 0 {
		score++
	}
	if r.Measurements.Length > 600 {
		score++
	}
	if score >= 6 {
		return RiskCritical
	}
	if score >= 4 {
		return RiskHigh
	}
	if score >= 2 {
		return RiskMedium
	}
	return RiskLow
}
func (r *Record) RefreshRisk() { r.RiskLevel = r.CalculateRisk(); r.UpdatedAt = time.Now().UTC() }
func (r Record) ReviewState(required int) RecordState {
	if required <= 0 {
		return StateReviewed
	}
	if len(r.Confirmations) == 0 {
		return StateProcessing
	}
	for _, c := range r.Confirmations {
		if c.Decision == "reject" {
			return StateRejected
		}
	}
	if len(r.Confirmations) >= required {
		return StateReviewed
	}
	return StateProcessing
}
func (r Record) ApprovedBy(id string) bool {
	for _, c := range r.Confirmations {
		if c.InspectorID == id && c.Decision == "approve" {
			return true
		}
	}
	return false
}

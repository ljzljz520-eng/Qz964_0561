package domain

func ShouldNotify(r Record) bool {
	return r.State == StateProcessing || r.State == StateRejected || r.RiskLevel == RiskHigh || r.RiskLevel == RiskCritical
}
func CanArchive(r Record) bool {
	return r.State == StateReviewed && len(r.Confirmations) >= 2 && r.RiskLevel != RiskCritical
}
func ReviewLabel(r Record) string {
	if r.State == StateRejected {
		return "驳回"
	}
	if len(r.Confirmations) >= 2 {
		return "已确认"
	}
	if len(r.Confirmations) > 0 {
		return "部分确认"
	}
	return "待确认"
}
func NextState(r Record) RecordState {
	switch r.State {
	case StateDraft:
		return StateRegistered
	case StateRegistered:
		return StateProcessing
	case StateProcessing:
		return StateReviewed
	case StateReviewed:
		return StateArchived
	}
	return r.State
}
func RiskThreshold(level RiskLevel) int {
	switch level {
	case RiskLow:
		return 0
	case RiskMedium:
		return 2
	case RiskHigh:
		return 4
	case RiskCritical:
		return 6
	}
	return 99
}
func MeetsRisk(level RiskLevel, score int) bool { return score >= RiskThreshold(level) }

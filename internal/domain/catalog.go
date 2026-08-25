package domain

import "sort"

type RiskBand struct {
	Level    RiskLevel
	MinScore int
	Guidance string
}

var RiskBands = []RiskBand{{RiskLow, 0, "常规堆放"}, {RiskMedium, 2, "增加通风"}, {RiskHigh, 4, "隔离并复检"}, {RiskCritical, 6, "停止处理"}}

func RiskGuidance(level RiskLevel) string {
	for _, b := range RiskBands {
		if b.Level == level {
			return b.Guidance
		}
	}
	return "未知"
}
func AllowedStates() []RecordState {
	return []RecordState{StateDraft, StateRegistered, StateProcessing, StateReviewed, StateArchived, StateRejected}
}
func SortRecords(items []Record) []Record {
	out := append([]Record(nil), items...)
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.Before(out[j].UpdatedAt) })
	return out
}
func DefectSeverity(defects []string) int {
	score := 0
	for _, d := range defects {
		if d == "裂纹" || d == "腐朽" {
			score += 2
		} else {
			score++
		}
	}
	return score
}
func NormalizeDecision(v string) string {
	if v == "通过" || v == "approve" {
		return "approve"
	}
	if v == "拒绝" || v == "reject" {
		return "reject"
	}
	return "pending"
}

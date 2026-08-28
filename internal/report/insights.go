package report

import (
	"sort"
	"strings"
	"timber-safety/internal/domain"
)

type Insight struct {
	Code, Title, Detail string
	Severity            int
}

func (b Builder) Insights(records []domain.Record) []Insight {
	out := []Insight{}
	for _, r := range records {
		if r.Measurements.Moisture > 24 {
			out = append(out, Insight{Code: "MOISTURE", Title: r.TimberCode, Detail: "含水率超过安全阈值", Severity: 3})
		}
		if len(r.Measurements.Defects) >= 3 {
			out = append(out, Insight{Code: "DEFECTS", Title: r.TimberCode, Detail: "缺陷数量需要复核", Severity: 2})
		}
		if r.State == domain.StateRejected {
			out = append(out, Insight{Code: "REJECTED", Title: r.TimberCode, Detail: "资料已驳回", Severity: 2})
		}
		if len(r.Confirmations) == 0 && r.IsOpen() {
			out = append(out, Insight{Code: "UNREVIEWED", Title: r.TimberCode, Detail: "仍等待审核确认", Severity: 1})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Severity > out[j].Severity })
	return out
}
func (b Builder) InsightText(records []domain.Record) string {
	items := b.Insights(records)
	if len(items) == 0 {
		return b.EmptyMessage()
	}
	parts := []string{}
	for _, i := range items {
		parts = append(parts, i.Title+": "+i.Detail)
	}
	return strings.Join(parts, "\n")
}
func (b Builder) CountSeverity(records []domain.Record, severity int) int {
	n := 0
	for _, i := range b.Insights(records) {
		if i.Severity >= severity {
			n++
		}
	}
	return n
}
func (b Builder) SourceRisk(records []domain.Record) map[string]domain.RiskLevel {
	out := map[string]domain.RiskLevel{}
	for _, r := range records {
		if current, ok := out[r.Source]; !ok || domain.RiskThreshold(r.RiskLevel) > domain.RiskThreshold(current) {
			out[r.Source] = r.RiskLevel
		}
	}
	return out
}
func (b Builder) SourceVolume(records []domain.Record) map[string]float64 {
	out := map[string]float64{}
	for _, r := range records {
		out[r.Source] += r.Measurements.Volume()
	}
	return out
}
func (b Builder) MoistureAverageBySource(records []domain.Record) map[string]float64 {
	sum := map[string]float64{}
	count := map[string]int{}
	for _, r := range records {
		sum[r.Source] += r.Measurements.Moisture
		count[r.Source]++
	}
	for source, total := range sum {
		if count[source] > 0 {
			sum[source] = total / float64(count[source])
		}
	}
	return sum
}
func (b Builder) Confirmed(records []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if len(r.Confirmations) > 0 {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) Unconfirmed(records []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if len(r.Confirmations) == 0 {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) Rejected(records []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.State == domain.StateRejected {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) HighRisk(records []domain.Record) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.RiskLevel == domain.RiskHigh || r.RiskLevel == domain.RiskCritical {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) ByState(records []domain.Record, state domain.RecordState) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.State == state {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) BySource(records []domain.Record, source string) []domain.Record {
	out := []domain.Record{}
	for _, r := range records {
		if r.Source == source {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) Search(records []domain.Record, text string) []domain.Record {
	needle := strings.ToLower(strings.TrimSpace(text))
	if needle == "" {
		return append([]domain.Record(nil), records...)
	}
	out := []domain.Record{}
	for _, r := range records {
		if strings.Contains(strings.ToLower(r.TimberCode), needle) || strings.Contains(strings.ToLower(r.Species), needle) || strings.Contains(strings.ToLower(r.Source), needle) {
			out = append(out, r)
		}
	}
	return out
}
func (b Builder) SortByCode(records []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), records...)
	sort.Slice(out, func(i, j int) bool { return out[i].TimberCode < out[j].TimberCode })
	return out
}
func (b Builder) SortByUpdated(records []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), records...)
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.After(out[j].UpdatedAt) })
	return out
}
func (b Builder) SortByMoisture(records []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), records...)
	sort.Slice(out, func(i, j int) bool { return out[i].Measurements.Moisture > out[j].Measurements.Moisture })
	return out
}
func (b Builder) LatestBySource(records []domain.Record) map[string]domain.Record {
	out := map[string]domain.Record{}
	for _, r := range records {
		if old, ok := out[r.Source]; !ok || r.UpdatedAt.After(old.UpdatedAt) {
			out[r.Source] = r
		}
	}
	return out
}
func (b Builder) RiskSummaryLine(records []domain.Record) string {
	s := b.BuildSummary(records)
	return strings.Join([]string{"总数 " + itoa(s.Total), "开放 " + itoa(s.Open), "高风险 " + itoa(s.HighRisk)}, " | ")
}
func itoa(v int) string {
	digits := "0123456789"
	if v == 0 {
		return "0"
	}
	out := ""
	for v > 0 {
		out = string(digits[v%10]) + out
		v /= 10
	}
	return out
}

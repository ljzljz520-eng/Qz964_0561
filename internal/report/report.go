package report

import (
	"fmt"
	"strings"
	"timber-safety/internal/domain"
)

type Builder struct{}

func New() Builder            { return Builder{} }
func (Builder) Title() string { return "木材加工安全报告" }
func (Builder) Render(records []domain.Record) string {
	lines := []string{"木材加工安全报告"}
	for _, r := range records {
		lines = append(lines, fmt.Sprintf("%s | %s | %s | %s", r.TimberCode, r.State, r.RiskLevel, domain.RiskGuidance(r.RiskLevel)))
	}
	return strings.Join(lines, "\n")
}
func (Builder) RiskCounts(records []domain.Record) map[domain.RiskLevel]int {
	out := map[domain.RiskLevel]int{}
	for _, r := range records {
		out[r.RiskLevel]++
	}
	return out
}
func (Builder) StateCounts(records []domain.Record) map[domain.RecordState]int {
	out := map[domain.RecordState]int{}
	for _, r := range records {
		out[r.State]++
	}
	return out
}
func (Builder) NeedsEscalation(r domain.Record) bool {
	return r.RiskLevel == domain.RiskHigh || r.RiskLevel == domain.RiskCritical
}

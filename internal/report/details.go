package report

import (
	"fmt"
	"sort"
	"strings"
	"timber-safety/internal/domain"
)

type Detail struct {
	ID, Code, State, Risk, Guidance string
	ConfirmationCount               int
	Open                            bool
}

func (b Builder) Details(records []domain.Record) []Detail {
	out := make([]Detail, 0, len(records))
	for _, r := range records {
		out = append(out, Detail{r.ID, r.TimberCode, string(r.State), string(r.RiskLevel), domain.RiskGuidance(r.RiskLevel), len(r.Confirmations), r.IsOpen()})
	}
	return out
}
func (b Builder) GroupBySource(records []domain.Record) map[string][]domain.Record {
	out := map[string][]domain.Record{}
	for _, r := range records {
		out[r.Source] = append(out[r.Source], r)
	}
	return out
}
func (b Builder) Timeline(records []domain.Record) []domain.Record {
	out := append([]domain.Record(nil), records...)
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.Before(out[j].UpdatedAt) })
	return out
}
func (b Builder) RenderDetail(r domain.Record) string {
	return fmt.Sprintf("%s\n编号: %s\n状态: %s\n风险: %s\n确认: %d", b.Title(), r.TimberCode, r.State, r.RiskLevel, len(r.Confirmations))
}
func (b Builder) RenderSources(records []domain.Record) string {
	groups := b.GroupBySource(records)
	parts := []string{}
	for source, items := range groups {
		parts = append(parts, fmt.Sprintf("%s:%d", source, len(items)))
	}
	sort.Strings(parts)
	return strings.Join(parts, "\n")
}
func (b Builder) SafetyNotice(r domain.Record) string {
	if r.RiskLevel == domain.RiskCritical {
		return "立即停止处理并隔离"
	}
	if r.RiskLevel == domain.RiskHigh {
		return "暂停处理并复检"
	}
	return "按标准流程处理"
}

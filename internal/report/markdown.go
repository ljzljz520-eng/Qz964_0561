package report

import (
	"fmt"
	"strings"
	"timber-safety/internal/domain"
)

func (b Builder) Markdown(records []domain.Record) string {
	lines := []string{"# 木材加工安全报告", "", "| 编号 | 状态 | 风险 | 确认数 |", "|---|---|---|---:|"}
	for _, r := range records {
		lines = append(lines, fmt.Sprintf("| %s | %s | %s | %d |", r.TimberCode, r.State, r.RiskLevel, len(r.Confirmations)))
	}
	return strings.Join(lines, "\n")
}
func (b Builder) CSV(records []domain.Record) string {
	lines := []string{"id,timber_code,state,risk,confirmations"}
	for _, r := range records {
		lines = append(lines, fmt.Sprintf("%s,%s,%s,%s,%d", r.ID, r.TimberCode, r.State, r.RiskLevel, len(r.Confirmations)))
	}
	return strings.Join(lines, "\n")
}
func (b Builder) Section(title string, records []domain.Record) string {
	return fmt.Sprintf("## %s\n\n%s", title, b.Render(records))
}
func (b Builder) RiskSection(records []domain.Record) string {
	parts := []string{}
	for _, r := range records {
		if r.RiskLevel == domain.RiskHigh || r.RiskLevel == domain.RiskCritical {
			parts = append(parts, b.RenderDetail(r))
		}
	}
	return strings.Join(parts, "\n\n")
}
func (b Builder) OpenSection(records []domain.Record) string {
	open := []domain.Record{}
	for _, r := range records {
		if r.IsOpen() {
			open = append(open, r)
		}
	}
	return b.Render(open)
}
func (b Builder) ArchiveSection(records []domain.Record) string {
	arch := []domain.Record{}
	for _, r := range records {
		if r.State == domain.StateArchived {
			arch = append(arch, r)
		}
	}
	return b.Render(arch)
}
func (b Builder) EmptyMessage() string { return "暂无符合条件的木材记录" }

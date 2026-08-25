package notify

import (
	"fmt"
	"timber-safety/internal/domain"
)

func SubjectFor(r domain.Record) string {
	if r.State == domain.StateRejected {
		return "木材资料已驳回"
	}
	if r.RiskLevel == domain.RiskCritical {
		return "木材重大风险警报"
	}
	return "木材状态更新"
}
func BodyFor(r domain.Record) string {
	return fmt.Sprintf("%s 当前%s，风险%s。%s", r.TimberCode, r.State, r.RiskLevel, domain.RiskGuidance(r.RiskLevel))
}
func (n *Notifier) NotifyReview(r domain.Record) error {
	return n.Send("review-desk", SubjectFor(r), BodyFor(r))
}
func (n *Notifier) Latest() (Message, bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if len(n.messages) == 0 {
		return Message{}, false
	}
	return n.messages[len(n.messages)-1], true
}
func (n *Notifier) ByRecipient(recipient string) []Message {
	n.mu.Lock()
	defer n.mu.Unlock()
	out := []Message{}
	for _, m := range n.messages {
		if m.Recipient == recipient {
			out = append(out, m)
		}
	}
	return out
}

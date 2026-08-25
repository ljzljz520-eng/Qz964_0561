package notify

import (
	"fmt"
	"sync"
	"timber-safety/internal/domain"
)

type Message struct{ Recipient, Subject, Body string }
type Notifier struct {
	mu       sync.Mutex
	messages []Message
}

func New() *Notifier { return &Notifier{messages: []Message{}} }
func (n *Notifier) Send(recipient, subject, body string) error {
	if recipient == "" || subject == "" {
		return fmt.Errorf("recipient and subject required")
	}
	n.mu.Lock()
	defer n.mu.Unlock()
	n.messages = append(n.messages, Message{recipient, subject, body})
	return nil
}
func (n *Notifier) NotifyRecord(r domain.Record) error {
	return n.Send("safety-desk", "木材状态更新", r.SummaryLine())
}
func (n *Notifier) Messages() []Message {
	n.mu.Lock()
	defer n.mu.Unlock()
	return append([]Message(nil), n.messages...)
}
func (n *Notifier) Count() int { n.mu.Lock(); defer n.mu.Unlock(); return len(n.messages) }
func (n *Notifier) Clear()     { n.mu.Lock(); defer n.mu.Unlock(); n.messages = nil }

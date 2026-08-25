package domain

import (
	"fmt"
	"strings"
	"time"
)

type RecordState string

const (
	StateDraft      RecordState = "draft"
	StateRegistered RecordState = "registered"
	StateProcessing RecordState = "processing"
	StateReviewed   RecordState = "reviewed"
	StateArchived   RecordState = "archived"
	StateRejected   RecordState = "rejected"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type Measurements struct {
	Length, Width, Height float64
	Moisture              float64
	Defects               []string
}
type Confirmation struct {
	RecordID, InspectorID, Decision, Notes string
	ConfirmedAt                            time.Time
}
type Record struct {
	ID            string         `json:"id"`
	TimberCode    string         `json:"timber_code"`
	Species       string         `json:"species"`
	Source        string         `json:"source"`
	State         RecordState    `json:"state"`
	RiskLevel     RiskLevel      `json:"risk_level"`
	Measurements  Measurements   `json:"measurements"`
	Confirmations []Confirmation `json:"confirmations"`
	Version       int            `json:"version"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}
type User struct {
	ID, Name, Role string
	Active         bool
	CreatedAt      time.Time
}
type Event struct {
	ID, RecordID, Type, ActorID, Details string
	OccurredAt                           time.Time
}
type Audit struct {
	ID, ActorID, Action, TargetID, Message string
	OccurredAt                             time.Time
}

func NewRecord(id, code, species, source string, m Measurements, now time.Time) Record {
	return Record{ID: id, TimberCode: code, Species: species, Source: source, State: StateDraft, RiskLevel: RiskLow, Measurements: m, Version: 1, CreatedAt: now, UpdatedAt: now}
}
func (r Record) IsOpen() bool           { return r.State != StateArchived && r.State != StateRejected }
func (r Record) ConfirmationCount() int { return len(r.Confirmations) }
func (r Record) HasInspector(id string) bool {
	for _, c := range r.Confirmations {
		if c.InspectorID == id {
			return true
		}
	}
	return false
}
func (r Record) SummaryLine() string {
	return fmt.Sprintf("%s %s %s %s", r.TimberCode, r.Species, r.State, r.RiskLevel)
}
func (r *Record) Normalize() {
	r.TimberCode = strings.ToUpper(strings.TrimSpace(r.TimberCode))
	r.Species = strings.TrimSpace(r.Species)
	r.Source = strings.TrimSpace(r.Source)
}

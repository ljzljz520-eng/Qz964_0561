package domain

import (
	"encoding/json"
	"fmt"
	"time"
)

func EncodeRecord(r Record) ([]byte, error) { return json.Marshal(r) }
func DecodeRecord(b []byte) (Record, error) {
	var r Record
	err := json.Unmarshal(b, &r)
	return r, err
}
func EncodeUser(u User) ([]byte, error)   { return json.Marshal(u) }
func DecodeUser(b []byte) (User, error)   { var u User; err := json.Unmarshal(b, &u); return u, err }
func EncodeEvent(e Event) ([]byte, error) { return json.Marshal(e) }
func DecodeEvent(b []byte) (Event, error) { var e Event; err := json.Unmarshal(b, &e); return e, err }
func EncodeAudit(a Audit) ([]byte, error) { return json.Marshal(a) }
func DecodeAudit(b []byte) (Audit, error) { var a Audit; err := json.Unmarshal(b, &a); return a, err }
func NewEvent(id, recordID, kind, actor, details string, t time.Time) Event {
	return Event{ID: id, RecordID: recordID, Type: kind, ActorID: actor, Details: details, OccurredAt: t}
}
func NewAudit(id, actor, action, target, message string, t time.Time) Audit {
	return Audit{ID: id, ActorID: actor, Action: action, TargetID: target, Message: message, OccurredAt: t}
}
func ValidateState(s RecordState) error {
	for _, v := range AllowedStates() {
		if s == v {
			return nil
		}
	}
	return fmt.Errorf("unknown state %q", s)
}

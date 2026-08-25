package domain

import "time"

type Lifecycle struct {
	RecordID                                           string
	Created, Registered, Processed, Reviewed, Archived *time.Time
}

func (l *Lifecycle) Mark(state RecordState, at time.Time) {
	switch state {
	case StateDraft:
		l.Created = &at
	case StateRegistered:
		l.Registered = &at
	case StateProcessing:
		l.Processed = &at
	case StateReviewed:
		l.Reviewed = &at
	case StateArchived:
		l.Archived = &at
	}
}
func (l Lifecycle) Duration(from, to *time.Time) time.Duration {
	if from == nil || to == nil {
		return 0
	}
	return to.Sub(*from)
}
func (l Lifecycle) TotalDuration() time.Duration {
	start := l.Created
	end := l.Archived
	if end == nil {
		end = l.Reviewed
	}
	if end == nil {
		end = l.Processed
	}
	if end == nil {
		end = l.Registered
	}
	return l.Duration(start, end)
}
func (l Lifecycle) Current() RecordState {
	if l.Archived != nil {
		return StateArchived
	}
	if l.Reviewed != nil {
		return StateReviewed
	}
	if l.Processed != nil {
		return StateProcessing
	}
	if l.Registered != nil {
		return StateRegistered
	}
	return StateDraft
}
func (l Lifecycle) Complete() bool { return l.Archived != nil }
func (l Lifecycle) MissingMilestones() []RecordState {
	out := []RecordState{}
	if l.Registered == nil {
		out = append(out, StateRegistered)
	}
	if l.Processed == nil {
		out = append(out, StateProcessing)
	}
	if l.Reviewed == nil {
		out = append(out, StateReviewed)
	}
	if l.Archived == nil {
		out = append(out, StateArchived)
	}
	return out
}

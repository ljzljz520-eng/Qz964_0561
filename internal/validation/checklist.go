package validation

import (
	"fmt"
	"timber-safety/internal/domain"
)

func ValidateChecklist(c domain.Checklist) error {
	if c.RecordID == "" {
		return fmt.Errorf("record id required")
	}
	if len(c.Items) == 0 {
		return fmt.Errorf("checklist empty")
	}
	if c.PassedCount() > len(c.Items) {
		return fmt.Errorf("invalid checklist count")
	}
	return nil
}
func ValidateChecklistCode(code string) error {
	if code == "" {
		return fmt.Errorf("checklist code required")
	}
	return nil
}
func ValidateNote(note string) error {
	if len(note) > 500 {
		return fmt.Errorf("note too long")
	}
	return nil
}
func ChecklistReady(c domain.Checklist) bool       { return ValidateChecklist(c) == nil && c.Complete() }
func MissingChecklist(c domain.Checklist) []string { return c.Missing() }

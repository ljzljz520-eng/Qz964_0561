package validation

import (
	"fmt"
	"regexp"
	"strings"
	"timber-safety/internal/domain"
)

var timberCode = regexp.MustCompile(`^木材[0-9]{1,4}$`)

func ValidateRecord(r domain.Record) error {
	r.Normalize()
	if r.ID == "" {
		return fmt.Errorf("id is required")
	}
	if !timberCode.MatchString(r.TimberCode) {
		return fmt.Errorf("timber code must look like 木材35")
	}
	if strings.TrimSpace(r.Species) == "" {
		return fmt.Errorf("species is required")
	}
	if strings.TrimSpace(r.Source) == "" {
		return fmt.Errorf("source is required")
	}
	if err := ValidateMeasurements(r.Measurements); err != nil {
		return err
	}
	return nil
}
func ValidateMeasurements(m domain.Measurements) error {
	if m.Length <= 0 || m.Width <= 0 || m.Height <= 0 {
		return fmt.Errorf("dimensions must be positive")
	}
	if m.Moisture < 0 || m.Moisture > 100 {
		return fmt.Errorf("moisture must be between 0 and 100")
	}
	for _, d := range m.Defects {
		if strings.TrimSpace(d) == "" {
			return fmt.Errorf("defect cannot be empty")
		}
	}
	return nil
}
func ValidateCode(code string) bool        { return timberCode.MatchString(strings.TrimSpace(code)) }
func NormalizeSource(source string) string { return strings.ToUpper(strings.TrimSpace(source)) }

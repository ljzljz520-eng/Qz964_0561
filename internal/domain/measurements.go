package domain

func (m Measurements) Volume() float64  { return m.Length * m.Width * m.Height }
func (m Measurements) IsComplete() bool { return m.Length > 0 && m.Width > 0 && m.Height > 0 }
func (m Measurements) MoistureBand() string {
	if m.Moisture < 12 {
		return "dry"
	}
	if m.Moisture < 20 {
		return "normal"
	}
	if m.Moisture < 28 {
		return "wet"
	}
	return "very-wet"
}
func (m Measurements) DefectCount() int { return len(m.Defects) }
func (m Measurements) HasDefect(name string) bool {
	for _, d := range m.Defects {
		if d == name {
			return true
		}
	}
	return false
}
func (m *Measurements) AddDefect(name string) bool {
	if name == "" || m.HasDefect(name) {
		return false
	}
	m.Defects = append(m.Defects, name)
	return true
}
func (m *Measurements) RemoveDefect(name string) bool {
	for i, d := range m.Defects {
		if d == name {
			m.Defects = append(m.Defects[:i], m.Defects[i+1:]...)
			return true
		}
	}
	return false
}
func (m Measurements) SortedDefects() []string {
	out := append([]string(nil), m.Defects...)
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}
func (m Measurements) SafetyFlags() []string {
	flags := []string{}
	if m.Moisture > 24 {
		flags = append(flags, "high-moisture")
	}
	if m.Length > 600 {
		flags = append(flags, "oversize")
	}
	if len(m.Defects) > 2 {
		flags = append(flags, "multiple-defects")
	}
	return flags
}

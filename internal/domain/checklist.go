package domain

type ChecklistItem struct {
	Code, Label string
	Required    bool
	Passed      bool
	Note        string
}
type Checklist struct {
	RecordID string
	Items    []ChecklistItem
}

func DefaultChecklist(recordID string) Checklist {
	return Checklist{RecordID: recordID, Items: []ChecklistItem{{"source", "来源可追溯", true, false, ""}, {"measurements", "尺寸完整", true, false, ""}, {"moisture", "含水率已测", true, false, ""}, {"defects", "缺陷已标记", false, false, ""}, {"storage", "堆放区域确认", true, false, ""}, {"ppe", "防护装备确认", true, false, ""}}}
}
func (c *Checklist) Mark(code string, passed bool, note string) bool {
	for i := range c.Items {
		if c.Items[i].Code == code {
			c.Items[i].Passed = passed
			c.Items[i].Note = note
			return true
		}
	}
	return false
}
func (c Checklist) Complete() bool {
	for _, i := range c.Items {
		if i.Required && !i.Passed {
			return false
		}
	}
	return true
}
func (c Checklist) RequiredCount() int {
	n := 0
	for _, i := range c.Items {
		if i.Required {
			n++
		}
	}
	return n
}
func (c Checklist) PassedCount() int {
	n := 0
	for _, i := range c.Items {
		if i.Passed {
			n++
		}
	}
	return n
}
func (c Checklist) Missing() []string {
	out := []string{}
	for _, i := range c.Items {
		if i.Required && !i.Passed {
			out = append(out, i.Label)
		}
	}
	return out
}
func (c Checklist) Progress() float64 {
	if len(c.Items) == 0 {
		return 1
	}
	return float64(c.PassedCount()) / float64(len(c.Items))
}
func (c Checklist) Has(code string) bool {
	for _, i := range c.Items {
		if i.Code == code {
			return true
		}
	}
	return false
}
func (c Checklist) Notes() map[string]string {
	out := map[string]string{}
	for _, i := range c.Items {
		if i.Note != "" {
			out[i.Code] = i.Note
		}
	}
	return out
}
func (c Checklist) Clone() Checklist {
	out := Checklist{RecordID: c.RecordID, Items: make([]ChecklistItem, len(c.Items))}
	copy(out.Items, c.Items)
	return out
}

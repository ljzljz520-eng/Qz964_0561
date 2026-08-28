package domain

import "encoding/json"

func (c Checklist) Marshal() ([]byte, error) { return json.Marshal(c) }
func UnmarshalChecklist(data []byte) (Checklist, error) {
	var c Checklist
	err := json.Unmarshal(data, &c)
	return c, err
}
func (l Lifecycle) Marshal() ([]byte, error) { return json.Marshal(l) }
func UnmarshalLifecycle(data []byte) (Lifecycle, error) {
	var l Lifecycle
	err := json.Unmarshal(data, &l)
	return l, err
}
func CloneRecord(r Record) Record {
	r.Confirmations = append([]Confirmation(nil), r.Confirmations...)
	r.Measurements.Defects = append([]string(nil), r.Measurements.Defects...)
	return r
}
func RecordsEqual(a, b Record) bool {
	ab, _ := json.Marshal(a)
	bb, _ := json.Marshal(b)
	return string(ab) == string(bb)
}

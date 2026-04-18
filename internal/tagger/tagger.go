package tagger

import (
	"encoding/json"
)

// Rule adds a tag value to a field when a condition is met.
type Rule struct {
	Field    string `json:"field"`
	Value    string `json:"value"`
	Tag      string `json:"tag"`
	TagValue string `json:"tag_value"`
}

// Tagger applies tagging rules to JSON log lines.
type Tagger struct {
	rules []Rule
}

// New creates a Tagger with the given rules.
func New(rules []Rule) *Tagger {
	return &Tagger{rules: rules}
}

// Apply processes a JSON line and returns it with tags applied.
func (t *Tagger) Apply(line string) (string, error) {
	if len(t.rules) == 0 {
		return line, nil
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line, err
	}
	for _, r := range t.rules {
		if applyRule(m, r) {
			m[r.Tag] = r.TagValue
		}
	}
	b, err := json.Marshal(m)
	if err != nil {
		return line, err
	}
	return string(b), nil
}

func applyRule(m map[string]any, r Rule) bool {
	v, ok := m[r.Field]
	if !ok {
		return false
	}
	s, ok := v.(string)
	if !ok {
		return false
	}
	return s == r.Value
}

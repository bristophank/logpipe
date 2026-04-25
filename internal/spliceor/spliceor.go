package spliceor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a splice operation: insert values from Fields into Target
// at the given position ("before", "after", or "replace").
type Rule struct {
	Target   string   `json:"target"`
	Fields   []string `json:"fields"`
	Position string   `json:"position"` // before | after | replace
	Sep      string   `json:"sep"`
}

// Spliceor reorders or injects field values relative to a target field.
type Spliceor struct {
	rules []Rule
}

// New creates a Spliceor with the given rules.
func New(rules []Rule) *Spliceor {
	return &Spliceor{rules: rules}
}

// Apply processes a JSON log line and returns the modified line.
func (s *Spliceor) Apply(line string) (string, error) {
	if len(s.rules) == 0 {
		return line, nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}
	for _, r := range s.rules {
		applyRule(obj, r)
	}
	return toString(obj)
}

func applyRule(obj map[string]interface{}, r Rule) {
	sep := r.Sep
	if sep == "" {
		sep = " "
	}

	parts := make([]string, 0, len(r.Fields))
	for _, f := range r.Fields {
		if v, ok := obj[f]; ok {
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	injected := strings.Join(parts, sep)

	targetVal := ""
	if v, ok := obj[r.Target]; ok {
		targetVal = fmt.Sprintf("%v", v)
	}

	switch r.Position {
	case "before":
		obj[r.Target] = injected + sep + targetVal
	case "replace":
		obj[r.Target] = injected
	default: // after
		obj[r.Target] = targetVal + sep + injected
	}
}

func toString(obj map[string]interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

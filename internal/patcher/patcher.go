// Package patcher applies conditional field patches to JSON log lines.
// A patch sets a target field to a static value only when a condition matches.
package patcher

import (
	"encoding/json"
)

// Rule defines a conditional patch operation.
type Rule struct {
	// Field is the JSON field to inspect for the condition.
	Field string `json:"field"`
	// Op is the match operator: "eq", "contains", "exists".
	Op string `json:"op"`
	// Value is the value to match against (unused for "exists").
	Value string `json:"value"`
	// Target is the field to set when the condition matches.
	Target string `json:"target"`
	// Patch is the value to assign to Target.
	Patch string `json:"patch"`
}

// Patcher conditionally patches fields in JSON log lines.
type Patcher struct {
	rules []Rule
}

// New creates a Patcher with the given rules.
func New(rules []Rule) *Patcher {
	return &Patcher{rules: rules}
}

// Apply processes a JSON line, returning it with patches applied.
// Lines that are not valid JSON are returned unchanged.
func (p *Patcher) Apply(line string) string {
	if len(p.rules) == 0 {
		return line
	}
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	modified := false
	for _, r := range p.rules {
		if matches(m, r) {
			m[r.Target] = r.Patch
			modified = true
		}
	}
	if !modified {
		return line
	}
	out, err := json.Marshal(m)
	if err != nil {
		return line
	}
	return string(out)
}

func matches(m map[string]any, r Rule) bool {
	v, ok := m[r.Field]
	switch r.Op {
	case "exists":
		return ok
	case "eq":
		if !ok {
			return false
		}
		s, ok := v.(string)
		return ok && s == r.Value
	case "contains":
		if !ok {
			return false
		}
		s, ok := v.(string)
		if !ok {
			return false
		}
		return len(r.Value) > 0 && contains(s, r.Value)
	}
	return false
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && indexStr(s, sub) >= 0)
}

func indexStr(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

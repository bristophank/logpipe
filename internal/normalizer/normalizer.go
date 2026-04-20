// Package normalizer standardises field names and values in structured log lines.
package normalizer

import (
	"encoding/json"
	"strings"
)

// Rule describes a single normalisation operation.
type Rule struct {
	// Field is the JSON key to normalise.
	Field string `json:"field"`
	// Transform is one of: "lowercase", "uppercase", "trim", "snake_case".
	Transform string `json:"transform"`
	// Rename optionally renames the field after transformation (empty = keep name).
	Rename string `json:"rename,omitempty"`
}

// Normalizer applies a set of Rules to each log line.
type Normalizer struct {
	rules []Rule
}

// New creates a Normalizer with the given rules.
func New(rules []Rule) *Normalizer {
	return &Normalizer{rules: rules}
}

// Apply parses line as JSON, applies all rules, and returns the modified JSON.
// Lines that are not valid JSON are returned unchanged.
func (n *Normalizer) Apply(line string) string {
	if len(n.rules) == 0 {
		return line
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	for _, r := range n.rules {
		applyRule(m, r)
	}
	out, err := json.Marshal(m)
	if err != nil {
		return line
	}
	return string(out)
}

func applyRule(m map[string]interface{}, r Rule) {
	v, ok := m[r.Field]
	if !ok {
		return
	}
	s, ok := v.(string)
	if !ok {
		return
	}
	switch r.Transform {
	case "lowercase":
		s = strings.ToLower(s)
	case "uppercase":
		s = strings.ToUpper(s)
	case "trim":
		s = strings.TrimSpace(s)
	case "snake_case":
		s = toSnakeCase(s)
	}
	delete(m, r.Field)
	dest := r.Field
	if r.Rename != "" {
		dest = r.Rename
	}
	m[dest] = s
}

// toSnakeCase converts a string to snake_case (spaces and hyphens → underscore, lowercased).
func toSnakeCase(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

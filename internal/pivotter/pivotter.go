// Package pivotter rotates a repeated field value into separate top-level keys.
// For example, given a field "metrics" that is an array of {name, value} objects,
// it produces a flat map of name→value pairs on the log line.
package pivotter

import (
	"encoding/json"
	"fmt"
)

// Rule defines how to pivot an array field.
type Rule struct {
	// Source is the array field to pivot (e.g. "metrics").
	Source string `json:"source"`
	// KeyField is the sub-field used as the new key (e.g. "name").
	KeyField string `json:"key_field"`
	// ValueField is the sub-field used as the new value (e.g. "value").
	ValueField string `json:"value_field"`
	// Prefix is optionally prepended to each generated key.
	Prefix string `json:"prefix"`
	// DeleteSource removes the original array field after pivoting.
	DeleteSource bool `json:"delete_source"`
}

// Pivotter applies pivot rules to JSON log lines.
type Pivotter struct {
	rules []Rule
}

// New creates a Pivotter with the given rules.
func New(rules []Rule) *Pivotter {
	return &Pivotter{rules: rules}
}

// Apply pivots array fields in the JSON line according to configured rules.
// Lines that are not valid JSON are returned unchanged.
func (p *Pivotter) Apply(line string) string {
	if len(p.rules) == 0 {
		return line
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}

	for _, r := range p.rules {
		m = applyRule(m, r)
	}

	return toString(m)
}

func applyRule(m map[string]any, r Rule) map[string]any {
	raw, ok := m[r.Source]
	if !ok {
		return m
	}

	items, ok := raw.([]any)
	if !ok {
		return m
	}

	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		keyVal, ok := obj[r.KeyField]
		if !ok {
			continue
		}
		val, ok := obj[r.ValueField]
		if !ok {
			continue
		}
		newKey := fmt.Sprintf("%s%v", r.Prefix, keyVal)
		m[newKey] = val
	}

	if r.DeleteSource {
		delete(m, r.Source)
	}
	return m
}

func toString(m map[string]any) string {
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}

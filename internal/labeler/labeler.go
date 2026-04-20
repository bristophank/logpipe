// Package labeler applies static or dynamic labels to structured log lines.
package labeler

import (
	"encoding/json"
	"fmt"
)

// Rule defines a label to add to every log line.
type Rule struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Labeler applies a fixed set of label rules to JSON log lines.
type Labeler struct {
	rules []Rule
}

// New returns a Labeler configured with the given rules.
// Rules with empty keys are silently skipped.
func New(rules []Rule) *Labeler {
	filtered := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Key != "" {
			filtered = append(filtered, r)
		}
	}
	return &Labeler{rules: filtered}
}

// Apply adds configured labels to the JSON log line.
// If line is empty or not valid JSON the original line is returned unchanged.
func (l *Labeler) Apply(line string) (string, error) {
	if line == "" {
		return line, nil
	}
	if len(l.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, fmt.Errorf("labeler: invalid JSON: %w", err)
	}

	for _, r := range l.rules {
		obj[r.Key] = r.Value
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line, fmt.Errorf("labeler: marshal error: %w", err)
	}
	return string(out), nil
}

// Package coalescer merges multiple JSON log fields into a single field,
// using the first non-empty value found among the candidates.
package coalescer

import (
	"encoding/json"
	"strings"
)

// Rule defines a coalesce operation: pick the first non-empty value from
// Candidates and write it to Target. If none are found, Target is omitted
// unless Default is set.
type Rule struct {
	Target     string   `json:"target"`
	Candidates []string `json:"candidates"`
	Default    string   `json:"default"`
	DeleteSrc  bool     `json:"delete_src"`
}

// Coalescer applies coalesce rules to JSON log lines.
type Coalescer struct {
	rules []Rule
}

// New creates a Coalescer with the given rules.
func New(rules []Rule) *Coalescer {
	return &Coalescer{rules: rules}
}

// Apply processes a single JSON log line and returns the modified line.
// Lines that are empty or cannot be parsed are returned unchanged.
func (c *Coalescer) Apply(line string) string {
	if strings.TrimSpace(line) == "" {
		return line
	}
	if len(c.rules) == 0 {
		return line
	}

	var record map[string]any
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return line
	}

	for _, rule := range c.rules {
		if rule.Target == "" || len(rule.Candidates) == 0 {
			continue
		}

		var chosen any
		var chosenKey string
		for _, key := range rule.Candidates {
			if v, ok := record[key]; ok && !isEmpty(v) {
				chosen = v
				chosenKey = key
				break
			}
		}

		if chosen == nil && rule.Default != "" {
			chosen = rule.Default
		}

		if chosen != nil {
			record[rule.Target] = chosen
			if rule.DeleteSrc && chosenKey != "" {
				for _, key := range rule.Candidates {
					if key != rule.Target {
						delete(record, key)
					}
				}
				_ = chosenKey
			}
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return line
	}
	return string(out)
}

func isEmpty(v any) bool {
	if v == nil {
		return true
	}
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s) == ""
	}
	return false
}

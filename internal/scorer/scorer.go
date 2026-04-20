// Package scorer assigns a numeric score to each log line based on
// configurable field-matching rules. The score is written back into
// the JSON object under a configurable output key.
package scorer

import (
	"encoding/json"
	"strings"
)

// Rule defines a single scoring rule.
type Rule struct {
	// Field is the JSON key to inspect.
	Field string
	// Contains is a substring that must appear in the field value.
	Contains string
	// Score is the value added when the rule matches.
	Score float64
}

// Config holds scorer configuration.
type Config struct {
	// Rules is the ordered list of scoring rules.
	Rules []Rule
	// OutputKey is the JSON key under which the total score is stored.
	// Defaults to "score" when empty.
	OutputKey string
}

// Scorer evaluates rules against each log line and injects a score.
type Scorer struct {
	rules     []Rule
	outputKey string
}

// New returns a new Scorer from cfg.
func New(cfg Config) *Scorer {
	key := cfg.OutputKey
	if key == "" {
		key = "score"
	}
	return &Scorer{rules: cfg.Rules, outputKey: key}
}

// Apply parses line as JSON, evaluates all rules, injects the cumulative
// score under the configured output key, and returns the modified JSON.
// If line is empty or not valid JSON it is returned unchanged.
func (s *Scorer) Apply(line string) string {
	if strings.TrimSpace(line) == "" {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	var total float64
	for _, r := range s.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		str, ok := val.(string)
		if !ok {
			continue
		}
		if r.Contains == "" || strings.Contains(str, r.Contains) {
			total += r.Score
		}
	}

	obj[s.outputKey] = total

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

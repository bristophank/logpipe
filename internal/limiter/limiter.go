// Package limiter drops log lines that exceed a maximum field value length.
package limiter

import (
	"encoding/json"
)

// Rule defines a field and the maximum allowed byte length for its string value.
type Rule struct {
	Field  string `json:"field"`
	MaxLen int    `json:"max_len"`
}

// Limiter filters out log lines where a watched field exceeds the configured
// maximum length. Lines that do not match any rule are always passed through.
type Limiter struct {
	rules []Rule
}

// New creates a Limiter with the given rules. Rules with MaxLen <= 0 are ignored.
func New(rules []Rule) *Limiter {
	valid := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if r.Field != "" && r.MaxLen > 0 {
			valid = append(valid, r)
		}
	}
	return &Limiter{rules: valid}
}

// Allow returns true when the line should be forwarded downstream.
// A line is dropped (false) if any configured field's string value exceeds
// its maximum allowed length. Non-string field values are never dropped.
func (l *Limiter) Allow(line string) bool {
	if len(l.rules) == 0 {
		return true
	}
	if line == "" {
		return true
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return true
	}

	for _, r := range l.rules {
		v, ok := record[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		if len(s) > r.MaxLen {
			return false
		}
	}
	return true
}

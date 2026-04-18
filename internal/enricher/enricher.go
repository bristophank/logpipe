package enricher

import (
	"encoding/json"
	"time"
)

// Rule defines a field to add with a static value or a special source.
type Rule struct {
	Field  string `json:"field"`
	Value  string `json:"value"`   // static value; overrides Source
	Source string `json:"source"` // "timestamp", "hostname"
}

// Enricher adds fields to log lines.
type Enricher struct {
	rules    []Rule
	hostname string
}

// New creates an Enricher with the given rules.
func New(rules []Rule, hostname string) *Enricher {
	return &Enricher{rules: rules, hostname: hostname}
}

// Apply parses line as JSON, adds configured fields, and returns the result.
// Non-JSON lines are returned unchanged.
func (e *Enricher) Apply(line string) string {
	if len(e.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, r := range e.rules {
		if r.Value != "" {
			obj[r.Field] = r.Value
			continue
		}
		switch r.Source {
		case "timestamp":
			obj[r.Field] = time.Now().UTC().Format(time.RFC3339)
		case "hostname":
			obj[r.Field] = e.hostname
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

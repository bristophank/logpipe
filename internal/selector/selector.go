package selector

import (
	"encoding/json"
	"strings"
)

// Rule defines which fields to keep or drop from a log line.
type Rule struct {
	Fields []string `json:"fields"`
	Mode   string   `json:"mode"` // "keep" or "drop"
}

// Selector filters fields from structured log lines.
type Selector struct {
	rules []Rule
}

// New creates a Selector with the given rules.
func New(rules []Rule) *Selector {
	return &Selector{rules: rules}
}

// Apply processes a JSON log line, keeping or dropping fields per rules.
// Lines that are empty or cannot be parsed are returned unchanged.
func (s *Selector) Apply(line string) string {
	if strings.TrimSpace(line) == "" {
		return line
	}
	if len(s.rules) == 0 {
		return line
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return line
	}

	for _, rule := range s.rules {
		switch strings.ToLower(rule.Mode) {
		case "keep":
			record = keepFields(record, rule.Fields)
		case "drop":
			record = dropFields(record, rule.Fields)
		}
	}

	return toString(record)
}

func keepFields(record map[string]interface{}, fields []string) map[string]interface{} {
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	out := make(map[string]interface{}, len(fields))
	for k, v := range record {
		if _, ok := set[k]; ok {
			out[k] = v
		}
	}
	return out
}

func dropFields(record map[string]interface{}, fields []string) map[string]interface{} {
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	for _, f := range fields {
		if _, ok := set[f]; ok {
			delete(record, f)
		}
	}
	return record
}

func toString(record map[string]interface{}) string {
	b, err := json.Marshal(record)
	if err != nil {
		return "{}"
	}
	return string(b)
}

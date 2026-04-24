package unwrapper

import (
	"encoding/json"
)

// Rule defines a field to unwrap and an optional prefix for child keys.
type Rule struct {
	Field  string `json:"field"`
	Prefix string `json:"prefix"`
	Delete bool   `json:"delete"`
}

// Unwrapper promotes nested object fields to the top level.
type Unwrapper struct {
	rules []Rule
}

// New creates an Unwrapper with the given rules.
func New(rules []Rule) *Unwrapper {
	return &Unwrapper{rules: rules}
}

// Apply processes a JSON log line, unwrapping nested objects according to rules.
func (u *Unwrapper) Apply(line string) (string, error) {
	if len(u.rules) == 0 {
		return line, nil
	}

	var record map[string]any
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return line, nil
	}

	for _, rule := range u.rules {
		applyRule(record, rule)
	}

	out, err := json.Marshal(record)
	if err != nil {
		return line, err
	}
	return string(out), nil
}

func applyRule(record map[string]any, rule Rule) {
	val, ok := record[rule.Field]
	if !ok {
		return
	}

	nested, ok := val.(map[string]any)
	if !ok {
		return
	}

	for k, v := range nested {
		key := k
		if rule.Prefix != "" {
			key = rule.Prefix + k
		}
		record[key] = v
	}

	if rule.Delete {
		delete(record, rule.Field)
	}
}

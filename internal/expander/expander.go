package expander

import (
	"encoding/json"
	"strings"
)

// Rule defines a field to expand and the delimiter to split on.
type Rule struct {
	Field     string `json:"field"`
	Delimiter string `json:"delimiter"`
	Target    string `json:"target"` // if empty, replaces field in-place
}

// Expander splits a string field into a JSON array.
type Expander struct {
	rules []Rule
}

// New creates an Expander with the given rules.
func New(rules []Rule) *Expander {
	return &Expander{rules: rules}
}

// Apply processes a JSON line and returns the modified line.
func (e *Expander) Apply(line string) (string, error) {
	if len(e.rules) == 0 {
		return line, nil
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line, err
	}

	for _, r := range e.rules {
		v, ok := m[r.Field]
		if !ok {
			continue
		}
		str, ok := v.(string)
		if !ok {
			continue
		}
		delim := r.Delimiter
		if delim == "" {
			delim = ","
		}
		parts := strings.Split(str, delim)
		trimmed := make([]string, 0, len(parts))
		for _, p := range parts {
			if t := strings.TrimSpace(p); t != "" {
				trimmed = append(trimmed, t)
			}
		}
		target := r.Target
		if target == "" {
			target = r.Field
		}
		m[target] = trimmed
	}

	return toString(m)
}

func toString(m map[string]any) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

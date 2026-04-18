package truncator

import (
	"encoding/json"
)

// Rule defines a field and its max byte length.
type Rule struct {
	Field  string
	MaxLen int
}

// Truncator shortens string fields in JSON log lines that exceed a max length.
type Truncator struct {
	rules []Rule
}

// New creates a Truncator with the given rules.
func New(rules []Rule) *Truncator {
	return &Truncator{rules: rules}
}

// Apply truncates fields in the JSON line according to configured rules.
// Returns the original line if it is not valid JSON.
func (t *Truncator) Apply(line string) string {
	if len(t.rules) == 0 {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	changed := false
	for _, r := range t.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if len(s) > r.MaxLen {
			obj[r.Field] = s[:r.MaxLen]
			changed = true
		}
	}

	if !changed {
		return line
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

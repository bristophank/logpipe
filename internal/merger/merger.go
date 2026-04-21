package merger

import (
	"encoding/json"
)

// Rule defines how two fields should be merged into a target field.
type Rule struct {
	Fields    []string `json:"fields"`     // source fields to merge
	Target    string   `json:"target"`     // destination field name
	Separator string   `json:"separator"`  // separator between values (default: " ")
	DeleteSrc bool     `json:"delete_src"` // remove source fields after merge
}

// Merger combines multiple string fields into a single field.
type Merger struct {
	rules []Rule
}

// New creates a Merger with the given rules.
func New(rules []Rule) *Merger {
	return &Merger{rules: rules}
}

// Apply processes a JSON log line and merges fields according to configured rules.
// Returns the original line unchanged if no rules are defined or the line is not valid JSON.
func (m *Merger) Apply(line string) (string, error) {
	if len(m.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, err
	}

	for _, rule := range m.rules {
		if len(rule.Fields) == 0 || rule.Target == "" {
			continue
		}
		sep := rule.Separator
		if sep == "" {
			sep = " "
		}
		parts := make([]string, 0, len(rule.Fields))
		for _, f := range rule.Fields {
			if v, ok := obj[f]; ok {
				parts = append(parts, toString(v))
			}
		}
		merged := ""
		for i, p := range parts {
			if i > 0 {
				merged += sep
			}
			merged += p
		}
		obj[rule.Target] = merged
		if rule.DeleteSrc {
			for _, f := range rule.Fields {
				if f != rule.Target {
					delete(obj, f)
				}
			}
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return string(out), nil
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}

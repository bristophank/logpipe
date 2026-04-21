package classifier

import (
	"encoding/json"
	"strings"
)

// Rule defines a classification rule: if field matches value, assign category.
type Rule struct {
	Field    string `json:"field"`
	Contains string `json:"contains,omitempty"`
	Equals   string `json:"equals,omitempty"`
	Category string `json:"category"`
}

// Classifier assigns a category label to log lines based on field matching rules.
type Classifier struct {
	rules  []Rule
	outKey string
}

// New creates a Classifier that writes the result into outKey.
// If outKey is empty it defaults to "category".
func New(outKey string, rules []Rule) *Classifier {
	if outKey == "" {
		outKey = "category"
	}
	return &Classifier{rules: rules, outKey: outKey}
}

// Apply classifies a single JSON log line and returns the annotated line.
// Lines that do not match any rule are returned unchanged.
func (c *Classifier) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, err
	}
	for _, r := range c.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		s := toString(val)
		matched := false
		switch {
		case r.Equals != "" && s == r.Equals:
			matched = true
		case r.Contains != "" && strings.Contains(s, r.Contains):
			matched = true
		}
		if matched {
			obj[c.outKey] = r.Category
			out, err := json.Marshal(obj)
			if err != nil {
				return line, err
			}
			return string(out), nil
		}
	}
	return line, nil
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

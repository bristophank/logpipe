// Package dropper removes log lines that match specified field conditions.
package dropper

import (
	"encoding/json"
	"strings"
)

// Rule defines a condition for dropping a log line.
type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"` // eq, contains, exists
	Value    string `json:"value"`
}

// Dropper evaluates rules and drops matching log lines.
type Dropper struct {
	rules []Rule
}

// New creates a new Dropper with the given rules.
func New(rules []Rule) *Dropper {
	return &Dropper{rules: rules}
}

// ShouldDrop returns true if the line matches any drop rule.
func (d *Dropper) ShouldDrop(line string) bool {
	if len(d.rules) == 0 {
		return false
	}
	if strings.TrimSpace(line) == "" {
		return false
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return false
	}

	for _, rule := range d.rules {
		if matchRule(obj, rule) {
			return true
		}
	}
	return false
}

// Apply returns the line unchanged if it should not be dropped, or empty string if dropped.
func (d *Dropper) Apply(line string) string {
	if d.ShouldDrop(line) {
		return ""
	}
	return line
}

func matchRule(obj map[string]interface{}, rule Rule) bool {
	val, ok := obj[rule.Field]
	switch rule.Operator {
	case "exists":
		return ok
	case "eq":
		if !ok {
			return false
		}
		return toString(val) == rule.Value
	case "contains":
		if !ok {
			return false
		}
		return strings.Contains(toString(val), rule.Value)
	}
	return false
}

func toString(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

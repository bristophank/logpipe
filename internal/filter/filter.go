package filter

import (
	"encoding/json"
	"strings"
)

// Rule defines a single filter condition applied to a log entry field.
type Rule struct {
	Field    string
	Operator string // "eq", "contains", "exists", "not_eq"
	Value    string
}

// Filter holds a set of rules and applies them to log entries.
type Filter struct {
	Rules []Rule
}

// New creates a Filter from a slice of rules.
func New(rules []Rule) *Filter {
	return &Filter{Rules: rules}
}

// Match returns true if the log line (JSON) satisfies all filter rules.
func (f *Filter) Match(line string) bool {
	if len(f.Rules) == 0 {
		return true
	}

	var entry map[string]interface{}
	if err := json.Unmarshal([]byte(line), &entry); err != nil {
		return false
	}

	for _, rule := range f.Rules {
		if !applyRule(rule, entry) {
			return false
		}
	}
	return true
}

func applyRule(rule Rule, entry map[string]interface{}) bool {
	val, exists := entry[rule.Field]

	switch rule.Operator {
	case "exists":
		return exists
	case "eq":
		if !exists {
			return false
		}
		return toString(val) == rule.Value
	case "not_eq":
		if !exists {
			return false
		}
		return toString(val) != rule.Value
	case "contains":
		if !exists {
			return false
		}
		return strings.Contains(toString(val), rule.Value)
	default:
		return false
	}
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

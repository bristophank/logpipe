package cutter

import (
	"encoding/json"
	"strings"
)

// Rule defines a field to cut a substring from.
type Rule struct {
	Field string `json:"field"`
	Start int    `json:"start"`
	End   int    `json:"end,omitempty"` // 0 means to end of string
	As    string `json:"as,omitempty"` // if set, write result to this field instead
}

// Cutter slices string fields by index range.
type Cutter struct {
	rules []Rule
}

// New creates a Cutter with the given rules.
func New(rules []Rule) *Cutter {
	return &Cutter{rules: rules}
}

// Apply processes a JSON log line, slicing each configured field.
func (c *Cutter) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, err
	}

	for _, r := range c.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		result := cut(s, r.Start, r.End)
		dest := r.Field
		if r.As != "" {
			dest = r.As
		}
		obj[dest] = result
	}

	return toString(obj)
}

// cut returns the substring of s from start to end (exclusive).
// Negative indices count from the end. end==0 means len(s).
func cut(s string, start, end int) string {
	n := len(s)
	if start < 0 {
		start = n + start
	}
	if end <= 0 {
		end = n + end
		if end == 0 {
			end = n
		}
	}
	if start < 0 {
		start = 0
	}
	if end > n {
		end = n
	}
	if start >= end {
		return ""
	}
	return s[start:end]
}

func toString(obj map[string]interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(b)), nil
}

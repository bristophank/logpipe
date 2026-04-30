// Package compactor removes fields with empty, null, or zero values from log lines.
package compactor

import (
	"encoding/json"
)

// Rule defines which fields to compact and what counts as empty.
type Rule struct {
	Field      string `json:"field"`       // field to inspect; empty means all fields
	DropNull   bool   `json:"drop_null"`   // remove null values
	DropEmpty  bool   `json:"drop_empty"`  // remove empty strings
	DropZero   bool   `json:"drop_zero"`   // remove numeric zero values
	DropFalse  bool   `json:"drop_false"`  // remove boolean false values
}

// Compactor applies compaction rules to JSON log lines.
type Compactor struct {
	rules []Rule
}

// New creates a Compactor with the given rules.
func New(rules []Rule) *Compactor {
	return &Compactor{rules: rules}
}

// Apply compacts a single JSON log line and returns the result.
func (c *Compactor) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}

	for _, rule := range c.rules {
		if rule.Field != "" {
			if v, ok := obj[rule.Field]; ok && isEmpty(v, rule) {
				delete(obj, rule.Field)
			}
		} else {
			for k, v := range obj {
				if isEmpty(v, rule) {
					delete(obj, k)
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

func isEmpty(v interface{}, r Rule) bool {
	if v == nil {
		return r.DropNull
	}
	switch val := v.(type) {
	case string:
		return r.DropEmpty && val == ""
	case float64:
		return r.DropZero && val == 0
	case bool:
		return r.DropFalse && !val
	}
	return false
}

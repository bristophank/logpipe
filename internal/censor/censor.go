package censor

import (
	"encoding/json"
	"strings"
)

// Rule defines a field and a list of values that should be censored.
type Rule struct {
	Field  string   `json:"field"`
	Values []string `json:"values"`
	Mask   string   `json:"mask"`
}

// Censor replaces specific field values with a mask string.
type Censor struct {
	rules []Rule
}

const defaultMask = "[CENSORED]"

// New creates a Censor with the provided rules.
func New(rules []Rule) *Censor {
	return &Censor{rules: rules}
}

// Apply censors matching field values in the JSON line and returns the result.
func (c *Censor) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return line, err
	}

	modified := false
	for _, rule := range c.rules {
		val, ok := record[rule.Field]
		if !ok {
			continue
		}
		strVal, ok := val.(string)
		if !ok {
			continue
		}
		for _, blocked := range rule.Values {
			if strings.EqualFold(strVal, blocked) {
				mask := rule.Mask
				if mask == "" {
					mask = defaultMask
				}
				record[rule.Field] = mask
				modified = true
				break
			}
		}
	}

	if !modified {
		return line, nil
	}

	out, err := json.Marshal(record)
	if err != nil {
		return line, err
	}
	return string(out), nil
}

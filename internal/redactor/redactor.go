package redactor

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Rule describes a single redaction rule.
type Rule struct {
	Field   string `json:"field"`
	Pattern string `json:"pattern"`
	Mask    string `json:"mask"`
}

// Redactor applies redaction rules to JSON log lines.
type Redactor struct {
	rules []compiled
}

type compiled struct {
	Rule
	re   *regexp.Regexp
	mask string
}

// New creates a Redactor from the provided rules.
func New(rules []Rule) (*Redactor, error) {
	cs := make([]compiled, 0, len(rules))
	for _, r := range rules {
		c := compiled{Rule: r, mask: r.Mask}
		if c.mask == "" {
			c.mask = "***"
		}
		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, err
			}
			c.re = re
		}
		cs = append(cs, c)
	}
	return &Redactor{rules: cs}, nil
}

// Apply redacts sensitive fields in a JSON line and returns the result.
func (r *Redactor) Apply(line string) string {
	if len(r.rules) == 0 {
		return line
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	for _, c := range r.rules {
		v, ok := obj[c.Field]
		if !ok {
			continue
		}
		if c.re != nil {
			obj[c.Field] = c.re.ReplaceAllString(toString(v), c.mask)
		} else {
			obj[c.Field] = c.mask
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

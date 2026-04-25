// Package bouncer filters log lines based on allowed or blocked field values.
// Rules specify a field, a list of allowed values (allowlist) or blocked values
// (blocklist). A line is dropped if any blocklist rule matches, or if an
// allowlist rule is present and no value matches.
package bouncer

import (
	"encoding/json"
)

// Rule defines a single bouncer rule.
type Rule struct {
	Field     string   `json:"field"`
	Allow     []string `json:"allow,omitempty"`
	Block     []string `json:"block,omitempty"`
}

// Bouncer evaluates log lines against a set of allow/block rules.
type Bouncer struct {
	rules []Rule
}

// New creates a Bouncer with the given rules.
func New(rules []Rule) *Bouncer {
	return &Bouncer{rules: rules}
}

// Allow returns true if the line passes all bouncer rules.
func (b *Bouncer) Allow(line string) bool {
	if len(b.rules) == 0 {
		return true
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return true
	}
	for _, r := range b.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		s := toString(val)
		for _, blocked := range r.Block {
			if s == blocked {
				return false
			}
		}
		if len(r.Allow) > 0 {
			matched := false
			for _, allowed := range r.Allow {
				if s == allowed {
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}
	}
	return true
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(v)
		return string(b)
	}
}

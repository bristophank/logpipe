// Package condenser merges consecutive log lines that share the same value
// for a configured key field into a single output line, collapsing repeated
// entries and optionally recording a count of how many were merged.
package condenser

import (
	"encoding/json"
	"strings"
)

// Rule describes how to condense lines.
type Rule struct {
	// Field is the JSON key whose value is used to detect consecutive duplicates.
	Field string `json:"field"`
	// CountField, when non-empty, adds the merge count to the output line.
	CountField string `json:"count_field,omitempty"`
}

// Condenser holds the condensing state.
type Condenser struct {
	rules    []Rule
	lastKey  string
	lastLine map[string]any
	count    int
}

// New returns a Condenser configured with the given rules.
func New(rules []Rule) *Condenser {
	return &Condenser{rules: rules}
}

// Add processes a raw JSON line. It returns a flushed line (non-empty) when a
// new key group starts, and an empty string while accumulating duplicates.
// Callers must call Flush after the stream ends to retrieve the final line.
func (c *Condenser) Add(line string) string {
	line = strings.TrimSpace(line)
	if line == "" {
		return ""
	}
	if len(c.rules) == 0 {
		return line
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	// Use first matching rule.
	rule := c.rules[0]
	key, _ := obj[rule.Field].(string)

	if c.lastLine == nil {
		// First line ever.
		c.lastKey = key
		c.lastLine = obj
		c.count = 1
		return ""
	}

	if key == c.lastKey {
		c.count++
		return ""
	}

	// Key changed — flush previous group and start a new one.
	out := c.buildOutput(rule)
	c.lastKey = key
	c.lastLine = obj
	c.count = 1
	return out
}

// Flush returns the final accumulated line and resets state.
func (c *Condenser) Flush() string {
	if c.lastLine == nil {
		return ""
	}
	if len(c.rules) == 0 {
		return ""
	}
	out := c.buildOutput(c.rules[0])
	c.lastLine = nil
	c.lastKey = ""
	c.count = 0
	return out
}

func (c *Condenser) buildOutput(rule Rule) string {
	copy := make(map[string]any, len(c.lastLine))
	for k, v := range c.lastLine {
		copy[k] = v
	}
	if rule.CountField != "" {
		copy[rule.CountField] = c.count
	}
	b, err := json.Marshal(copy)
	if err != nil {
		return ""
	}
	return string(b)
}

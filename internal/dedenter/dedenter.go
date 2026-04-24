// Package dedenter removes or normalises indentation from string fields in
// structured log lines. Each rule targets a specific field and strips leading
// whitespace / tab characters, optionally collapsing interior runs of
// whitespace to a single space.
package dedenter

import (
	"encoding/json"
	"strings"
)

// Rule describes how to dedent a single field.
type Rule struct {
	// Field is the JSON key whose value will be processed.
	Field string `json:"field"`
	// CollapseInner replaces interior whitespace runs with a single space when
	// true.
	CollapseInner bool `json:"collapse_inner"`
}

// Dedenter applies dedent rules to JSON log lines.
type Dedenter struct {
	rules []Rule
}

// New returns a Dedenter configured with the supplied rules. If rules is
// empty every line is passed through unchanged.
func New(rules []Rule) *Dedenter {
	return &Dedenter{rules: rules}
}

// Apply processes a single JSON log line and returns the modified line.
// Lines that cannot be decoded are returned as-is.
func (d *Dedenter) Apply(line string) string {
	if len(d.rules) == 0 {
		return line
	}
	var rec map[string]any
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		return line
	}
	changed := false
	for _, r := range d.rules {
		v, ok := rec[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		dedented := dedent(s, r.CollapseInner)
		if dedented != s {
			rec[r.Field] = dedented
			changed = true
		}
	}
	if !changed {
		return line
	}
	return toString(rec)
}

// dedent strips leading whitespace from every line of s and optionally
// collapses interior whitespace runs.
func dedent(s string, collapseInner bool) string {
	lines := strings.Split(s, "\n")
	for i, l := range lines {
		l = strings.TrimLeft(l, " \t")
		if collapseInner {
			l = strings.Join(strings.Fields(l), " ")
		}
		lines[i] = l
	}
	return strings.Join(lines, "\n")
}

func toString(rec map[string]any) string {
	b, err := json.Marshal(rec)
	if err != nil {
		return ""
	}
	return string(b)
}

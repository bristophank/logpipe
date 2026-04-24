// Package stencil applies a template pattern to JSON log lines,
// rendering a new string field by interpolating existing field values.
package stencil

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"
)

// Rule defines a single stencil operation.
type Rule struct {
	// Target is the field name to write the rendered result into.
	Target string
	// Template is a Go text/template string where {{.field}} references log fields.
	Template string
	// Overwrite controls whether an existing target field is replaced.
	Overwrite bool
}

// Stencil renders template strings against JSON log line fields.
type Stencil struct {
	rules []parsedRule
}

type parsedRule struct {
	Rule
	tmpl *template.Template
}

// New creates a Stencil from the given rules.
// Returns an error if any template fails to parse.
func New(rules []Rule) (*Stencil, error) {
	parsed := make([]parsedRule, 0, len(rules))
	for _, r := range rules {
		t, err := template.New("").Option("missingkey=zero").Parse(r.Template)
		if err != nil {
			return nil, err
		}
		parsed = append(parsed, parsedRule{Rule: r, tmpl: t})
	}
	return &Stencil{rules: parsed}, nil
}

// Apply renders each rule's template against the fields in line and injects
// the result into the JSON object. Non-JSON lines are returned unchanged.
func (s *Stencil) Apply(line string) string {
	if len(s.rules) == 0 {
		return line
	}
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return line
	}
	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(trimmed), &fields); err != nil {
		return line
	}
	for _, r := range s.rules {
		if _, exists := fields[r.Target]; exists && !r.Overwrite {
			continue
		}
		var buf bytes.Buffer
		if err := r.tmpl.Execute(&buf, fields); err != nil {
			continue
		}
		fields[r.Target] = buf.String()
	}
	out, err := json.Marshal(fields)
	if err != nil {
		return line
	}
	return string(out)
}

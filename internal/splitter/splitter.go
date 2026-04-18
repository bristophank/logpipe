// Package splitter fans out a single log line to multiple named outputs
// based on a field value match.
package splitter

import (
	"encoding/json"
	"io"
)

// Rule maps a field value to a sink name.
type Rule struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Sink    string `json:"sink"`
}

// Splitter routes lines to sinks by inspecting a JSON field.
type Splitter struct {
	rules   []Rule
	sinks   map[string]io.Writer
	fallback io.Writer
}

// New creates a Splitter with the given rules, sinks, and optional fallback.
func New(rules []Rule, sinks map[string]io.Writer, fallback io.Writer) *Splitter {
	return &Splitter{rules: rules, sinks: sinks, fallback: fallback}
}

// Write evaluates line against all rules and writes to matching sinks.
// If no rule matches and a fallback is set, the line is written there.
func (s *Splitter) Write(line []byte) error {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return err
	}

	matched := false
	for _, r := range s.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		if toString(v) == r.Value {
			if w, ok := s.sinks[r.Sink]; ok {
				w.Write(append(line, '\n')) //nolint:errcheck
				matched = true
			}
		}
	}

	if !matched && s.fallback != nil {
		s.fallback.Write(append(line, '\n')) //nolint:errcheck
	}
	return nil
}

func toString(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}

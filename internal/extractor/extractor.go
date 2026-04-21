package extractor

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a single extraction: pull a value from Field,
// apply an optional regex-like prefix/suffix trim, and write it to Target.
type Rule struct {
	Field  string `json:"field"`
	Target string `json:"target"`
	Prefix string `json:"prefix,omitempty"`
	Suffix string `json:"suffix,omitempty"`
}

// Extractor copies or derives values from existing fields into new fields.
type Extractor struct {
	rules []Rule
}

// New creates an Extractor with the given rules.
func New(rules []Rule) *Extractor {
	return &Extractor{rules: rules}
}

// Apply processes a single JSON log line, extracting fields per the rules.
// Lines that are empty or cannot be decoded are returned unchanged.
func (e *Extractor) Apply(line string) (string, error) {
	if strings.TrimSpace(line) == "" {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return line, nil
	}

	if len(e.rules) == 0 {
		return line, nil
	}

	for _, r := range e.rules {
		raw, ok := record[r.Field]
		if !ok {
			continue
		}
		val := fmt.Sprintf("%v", raw)
		val = applyTrim(val, r.Prefix, r.Suffix)
		target := r.Target
		if target == "" {
			target = r.Field + "_extracted"
		}
		record[target] = val
	}

	out, err := json.Marshal(record)
	if err != nil {
		return line, err
	}
	return string(out), nil
}

func applyTrim(val, prefix, suffix string) string {
	if prefix != "" {
		val = strings.TrimPrefix(val, prefix)
	}
	if suffix != "" {
		val = strings.TrimSuffix(val, suffix)
	}
	return val
}

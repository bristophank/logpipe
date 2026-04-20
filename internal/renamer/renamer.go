package renamer

import (
	"encoding/json"
	"strings"
)

// Rule defines a single field rename operation.
type Rule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Renamer renames fields in JSON log lines according to a set of rules.
type Renamer struct {
	rules []Rule
}

// New creates a Renamer with the given rules. Rules with empty From or To
// fields are silently ignored.
func New(rules []Rule) *Renamer {
	valid := make([]Rule, 0, len(rules))
	for _, r := range rules {
		if strings.TrimSpace(r.From) != "" && strings.TrimSpace(r.To) != "" {
			)
		}
	}
	return &Renamer{rules: valid}
}

// Apply renames fields in the parsed JSON object according to the configured
// rules and returns the re-encoded JSON line. If there are no rules or the
// input is not valid JSON the original line is returned unchanged.
func (rn *Renamer) Apply(line string) string {
	if len(rn.rules) == 0 {
		return line
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, r := range rn.rules {
		val, ok := obj[r.From]
		if !ok {
			continue
		}
		delete(obj, r.From)
		obj[r.To] = val
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

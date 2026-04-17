package transformer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a single field transformation.
type Rule struct {
	Field  string `json:"field"`
	Op     string `json:"op"`     // set, rename, delete, uppercase, lowercase
	Value  string `json:"value"`  // used by set, rename
}

// Transformer applies field-level mutations to JSON log lines.
type Transformer struct {
	rules []Rule
}

// New creates a Transformer with the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply mutates the JSON line according to the configured rules.
// Returns the original line on parse error.
func (t *Transformer) Apply(line string) string {
	if len(t.rules) == 0 {
		return line
	}
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line
	}
	for _, r := range t.rules {
		applyRule(m, r)
	}
	out, err := json.Marshal(m)
	if err != nil {
		return line
	}
	return string(out)
}

func applyRule(m map[string]interface{}, r Rule) {
	switch r.Op {
	case "set":
		m[r.Field] = r.Value
	case "delete":
		delete(m, r.Field)
	case "rename":
		if v, ok := m[r.Field]; ok {
			m[r.Value] = v
			delete(m, r.Field)
		}
	case "uppercase":
		if v, ok := m[r.Field]; ok {
			m[r.Field] = strings.ToUpper(fmt.Sprintf("%v", v))
		}
	case "lowercase":
		if v, ok := m[r.Field]; ok {
			m[r.Field] = strings.ToLower(fmt.Sprintf("%v", v))
		}
	}
}

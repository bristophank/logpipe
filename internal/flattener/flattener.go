package flattener

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines how a nested JSON object should be flattened.
type Rule struct {
	Prefix    string // optional prefix for flattened keys
	Separator string // separator between key levels (default: ".")
	Fields    []string // specific top-level fields to flatten; empty means all
}

// Flattener collapses nested JSON objects into a single level.
type Flattener struct {
	rules []Rule
}

// New creates a Flattener with the given rules.
func New(rules []Rule) *Flattener {
	for i := range rules {
		if rules[i].Separator == "" {
			rules[i].Separator = "."
		}
	}
	return &Flattener{rules: rules}
}

// Apply flattens the JSON line according to configured rules.
// Non-JSON or empty lines are returned unchanged.
func (f *Flattener) Apply(line string) string {
	if strings.TrimSpace(line) == "" {
		return line
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}
	if len(f.rules) == 0 {
		return line
	}
	for _, rule := range f.rules {
		if len(rule.Fields) == 0 {
			flat := make(map[string]any)
			flattenMap(obj, rule.Prefix, rule.Separator, flat)
			obj = flat
		} else {
			for _, field := range rule.Fields {
				val, ok := obj[field]
				if !ok {
					continue
				}
				nested, ok := val.(map[string]any)
				if !ok {
					continue
				}
				delete(obj, field)
				prefix := field
				if rule.Prefix != "" {
					prefix = rule.Prefix + rule.Separator + field
				}
				flattenMap(nested, prefix, rule.Separator, obj)
			}
		}
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(b)
}

func flattenMap(src map[string]any, prefix, sep string, dst map[string]any) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = fmt.Sprintf("%s%s%s", prefix, sep, k)
		}
		if nested, ok := v.(map[string]any); ok {
			flattenMap(nested, key, sep, dst)
		} else {
			dst[key] = v
		}
	}
}

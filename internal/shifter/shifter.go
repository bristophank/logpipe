package shifter

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Rule defines a shift operation on a numeric JSON field.
type Rule struct {
	Field  string  // field to shift
	By     float64 // amount to add (negative to subtract)
	Scale  float64 // multiplier applied after shift (0 means 1.0)
	Target string  // write result here; defaults to Field
}

// Shifter applies arithmetic shift operations to numeric log fields.
type Shifter struct {
	rules []Rule
}

// New creates a Shifter with the given rules.
func New(rules []Rule) *Shifter {
	for i := range rules {
		if rules[i].Scale == 0 {
			rules[i].Scale = 1.0
		}
		if rules[i].Target == "" {
			rules[i].Target = rules[i].Field
		}
	}
	return &Shifter{rules: rules}
}

// Apply parses line as JSON, applies each shift rule, and returns the result.
func (s *Shifter) Apply(line string) (string, error) {
	if len(s.rules) == 0 {
		return line, nil
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}
	for _, r := range s.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		f, err := toFloat(v)
		if err != nil {
			continue
		}
		result := (f + r.By) * r.Scale
		obj[r.Target] = result
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line, fmt.Errorf("shifter: marshal: %w", err)
	}
	return string(out), nil
}

func toFloat(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case int:
		return float64(x), nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		if err != nil {
			return 0, fmt.Errorf("not a number: %q", x)
		}
		return f, nil
	}
	return 0, fmt.Errorf("unsupported type %T", v)
}

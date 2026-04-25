package clamper

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rule defines a numeric field to clamp between Min and Max.
type Rule struct {
	Field string  `json:"field"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

// Clamper constrains numeric field values to [Min, Max].
type Clamper struct {
	rules []Rule
}

// New returns a Clamper with the given rules.
func New(rules []Rule) *Clamper {
	return &Clamper{rules: rules}
}

// Apply clamps numeric fields in the JSON line according to the configured rules.
// Lines that are not valid JSON are returned unchanged.
func (c *Clamper) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}

	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}

	changed := false
	for _, r := range c.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		f, err := toFloat(v)
		if err != nil {
			continue
		}
		clamped := clamp(f, r.Min, r.Max)
		if clamped != f {
			obj[r.Field] = clamped
			changed = true
		}
	}

	if !changed {
		return line, nil
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return line, fmt.Errorf("clamper: marshal: %w", err)
	}
	return string(b), nil
}

func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

func toFloat(v any) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case int:
		return float64(t), nil
	case string:
		return strconv.ParseFloat(t, 64)
	case json.Number:
		return t.Float64()
	}
	return 0, fmt.Errorf("unsupported type %T", v)
}

package comparator

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rule defines a comparison to apply to a log field.
type Rule struct {
	Field  string  `json:"field"`
	Op     string  `json:"op"`      // "gt", "lt", "gte", "lte", "eq", "neq"
	Value  float64 `json:"value"`
	Target string  `json:"target"`  // field to write result into
}

// Comparator evaluates numeric comparisons on log fields and writes a boolean result.
type Comparator struct {
	rules []Rule
}

// New creates a Comparator with the given rules.
func New(rules []Rule) *Comparator {
	return &Comparator{rules: rules}
}

// Apply processes a JSON log line, evaluating each rule and writing the result.
func (c *Comparator) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}
	var obj map[string]any
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, err
	}
	for _, r := range c.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		num, err := toFloat(v)
		if err != nil {
			continue
		}
		result := compare(num, r.Op, r.Value)
		target := r.Target
		if target == "" {
			target = r.Field + "_" + r.Op
		}
		obj[target] = result
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return string(out), nil
}

func compare(num float64, op string, value float64) bool {
	switch op {
	case "gt":
		return num > value
	case "lt":
		return num < value
	case "gte":
		return num >= value
	case "lte":
		return num <= value
	case "eq":
		return num == value
	case "neq":
		return num != value
	default:
		return false
	}
}

func toFloat(v any) (float64, error) {
	switch t := v.(type) {
	case float64:
		return t, nil
	case int:
		return float64(t), nil
	case string:
		return strconv.ParseFloat(t, 64)
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}

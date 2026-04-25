package caster

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rule defines a single cast operation: convert field to a target format.
type Rule struct {
	Field  string `json:"field"`
	Format string `json:"format"` // "upper", "lower", "string", "int", "float", "bool"
}

// Caster applies format coercions to string representations of JSON values.
type Caster struct {
	rules []Rule
}

// New returns a Caster configured with the given rules.
func New(rules []Rule) *Caster {
	return &Caster{rules: rules}
}

// Apply processes a JSON line, coercing each targeted field to its declared format.
// Lines that are empty or fail to parse are returned unchanged.
func (c *Caster) Apply(line string) string {
	if len(c.rules) == 0 || line == "" {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, r := range c.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		casted, err := castTo(v, r.Format)
		if err != nil {
			continue
		}
		obj[r.Field] = casted
	}

	return toString(obj)
}

func castTo(v interface{}, format string) (interface{}, error) {
	raw := fmt.Sprintf("%v", v)
	switch format {
	case "string":
		return raw, nil
	case "int":
		f, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return nil, err
		}
		return int64(f), nil
	case "float":
		return strconv.ParseFloat(raw, 64)
	case "bool":
		return strconv.ParseBool(raw)
	default:
		return nil, fmt.Errorf("unknown format: %s", format)
	}
}

func toString(obj map[string]interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(b)
}

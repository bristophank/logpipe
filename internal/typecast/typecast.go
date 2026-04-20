package typecast

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rule defines a field and the target type to cast it to.
type Rule struct {
	Field string `json:"field"`
	To    string `json:"to"` // "string", "int", "float", "bool"
}

// Caster applies type-casting rules to JSON log lines.
type Caster struct {
	rules []Rule
}

// New creates a new Caster with the given rules.
func New(rules []Rule) *Caster {
	return &Caster{rules: rules}
}

// Apply parses the JSON line, casts configured fields, and returns the result.
func (c *Caster) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}

	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return line, fmt.Errorf("typecast: invalid json: %w", err)
	}

	for _, r := range c.rules {
		v, ok := m[r.Field]
		if !ok {
			continue
		}
		casted, err := castValue(v, r.To)
		if err != nil {
			continue // skip fields that cannot be cast
		}
		m[r.Field] = casted
	}

	out, err := json.Marshal(m)
	if err != nil {
		return line, fmt.Errorf("typecast: marshal error: %w", err)
	}
	return string(out), nil
}

func castValue(v any, to string) (any, error) {
	raw := fmt.Sprintf("%v", v)
	switch to {
	case "string":
		return raw, nil
	case "int":
		n, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			// try via float first (e.g. "3.0" -> 3)
			f, ferr := strconv.ParseFloat(raw, 64)
			if ferr != nil {
				return nil, err
			}
			return int64(f), nil
		}
		return n, nil
	case "float":
		return strconv.ParseFloat(raw, 64)
	case "bool":
		return strconv.ParseBool(raw)
	default:
		return nil, fmt.Errorf("typecast: unknown target type %q", to)
	}
}

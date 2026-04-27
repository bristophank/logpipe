package scaler

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rule defines a numeric field scaling operation.
type Rule struct {
	Field  string  `json:"field"`
	Min    float64 `json:"min"`
	Max    float64 `json:"max"`
	NewMin float64 `json:"new_min"`
	NewMax float64 `json:"new_max"`
}

// Scaler rescales numeric fields from one range to another.
type Scaler struct {
	rules []Rule
}

// New creates a Scaler with the given rules.
func New(rules []Rule) *Scaler {
	return &Scaler{rules: rules}
}

// Apply rescales numeric fields in the JSON line according to configured rules.
// Returns the original line if no rules match or the input is invalid JSON.
func (s *Scaler) Apply(line string) (string, error) {
	if len(s.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}

	for _, r := range s.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}

		f, err := toFloat(val)
		if err != nil {
			continue
		}

		scaled := rescale(f, r.Min, r.Max, r.NewMin, r.NewMax)
		obj[r.Field] = scaled
	}

	return toString(obj)
}

// rescale maps v from [min, max] to [newMin, newMax].
// If the source range is zero, returns newMin.
func rescale(v, min, max, newMin, newMax float64) float64 {
	span := max - min
	if span == 0 {
		return newMin
	}
	return newMin + (v-min)/(span)*(newMax-newMin)
}

func toFloat(v interface{}) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case string:
		return strconv.ParseFloat(x, 64)
	}
	return 0, fmt.Errorf("not numeric")
}

func toString(obj map[string]interface{}) (string, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

package timestamper

import (
	"encoding/json"
	"time"
)

// Rule defines a field to stamp and the time format to use.
type Rule struct {
	Field     string `json:"field"`
	Format    string `json:"format"`
	Overwrite bool   `json:"overwrite"`
}

// Timestamper injects timestamp fields into JSON log lines.
type Timestamper struct {
	rules []Rule
	now   func() time.Time
}

// New creates a Timestamper with the given rules.
// If now is nil, time.Now is used.
func New(rules []Rule, now func() time.Time) *Timestamper {
	if now == nil {
		now = time.Now
	}
	return &Timestamper{rules: rules, now: now}
}

// Apply stamps timestamp fields onto a JSON log line according to configured rules.
// Lines that are not valid JSON are returned unchanged.
func (t *Timestamper) Apply(line string) (string, error) {
	if len(t.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, nil
	}

	ts := t.now()
	for _, r := range t.rules {
		if r.Field == "" {
			continue
		}
		if _, exists := obj[r.Field]; exists && !r.Overwrite {
			continue
		}
		fmt := r.Format
		if fmt == "" {
			fmt = time.RFC3339
		}
		obj[r.Field] = ts.UTC().Format(fmt)
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return string(b), nil
}

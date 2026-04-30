package windower

import (
	"encoding/json"
	"time"
)

// Rule defines a tumbling window over a numeric field.
type Rule struct {
	Field  string        // field to accumulate
	Window time.Duration // window duration
	Alias  string        // output field name for the window sum
}

// Windower accumulates numeric field values within a fixed tumbling window
// and emits a summary line when the window expires.
type Windower struct {
	rules  []Rule
	bucket map[string]float64
	start  time.Time
	now    func() time.Time
}

// New returns a Windower configured with the given rules.
func New(rules []Rule) *Windower {
	return &Windower{
		rules:  rules,
		bucket: make(map[string]float64),
		now:    time.Now,
	}
}

// Add ingests a JSON log line and accumulates matching fields.
// It returns a flushed summary line if the window has expired, otherwise "".
func (w *Windower) Add(line string) (string, bool) {
	if len(w.rules) == 0 || line == "" {
		return "", false
	}

	var rec map[string]any
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		return "", false
	}

	if w.start.IsZero() {
		w.start = w.now()
	}

	for _, r := range w.rules {
		if v, ok := rec[r.Field]; ok {
			switch n := v.(type) {
			case float64:
				w.bucket[r.Alias] += n
			}
		}
	}

	if w.now().Sub(w.start) >= w.rules[0].Window {
		return w.flush()
	}
	return "", false
}

// Flush emits the current bucket as a JSON summary and resets state.
func (w *Windower) Flush() (string, bool) {
	return w.flush()
}

func (w *Windower) flush() (string, bool) {
	if len(w.bucket) == 0 {
		w.start = time.Time{}
		return "", false
	}
	out := make(map[string]any, len(w.bucket))
	for k, v := range w.bucket {
		out[k] = v
	}
	w.bucket = make(map[string]float64)
	w.start = time.Time{}
	b, err := json.Marshal(out)
	if err != nil {
		return "", false
	}
	return string(b), true
}

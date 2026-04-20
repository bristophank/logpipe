package counter

import (
	"encoding/json"
	"fmt"
	"sync"
)

// Rule defines a field to count distinct values for.
type Rule struct {
	Field string `json:"field"`
	Alias string `json:"alias,omitempty"` // output key; defaults to field
}

// Counter tracks distinct value frequencies for configured fields.
type Counter struct {
	mu    sync.Mutex
	rules []Rule
	counts map[string]map[string]int64 // field -> value -> count
}

// New creates a Counter for the given rules.
func New(rules []Rule) *Counter {
	return &Counter{
		rules:  rules,
		counts: make(map[string]map[string]int64),
	}
}

// Add records field values from a JSON log line.
// Returns an error only if the line is not valid JSON.
func (c *Counter) Add(line string) error {
	if line == "" {
		return nil
	}
	var rec map[string]any
	if err := json.Unmarshal([]byte(line), &rec); err != nil {
		return fmt.Errorf("counter: invalid json: %w", err)
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, r := range c.rules {
		v, ok := rec[r.Field]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", v)
		if c.counts[r.Field] == nil {
			c.counts[r.Field] = make(map[string]int64)
		}
		c.counts[r.Field][key]++
	}
	return nil
}

// Snapshot returns a copy of the current counts keyed by alias (or field).
func (c *Counter) Snapshot() map[string]map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make(map[string]map[string]int64, len(c.rules))
	for _, r := range c.rules {
		alias := r.Alias
		if alias == "" {
			alias = r.Field
		}
		src := c.counts[r.Field]
		copy := make(map[string]int64, len(src))
		for k, v := range src {
			copy[k] = v
		}
		out[alias] = copy
	}
	return out
}

// Reset clears all accumulated counts.
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts = make(map[string]map[string]int64)
}

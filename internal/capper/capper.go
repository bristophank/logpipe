// Package capper limits the number of log lines emitted per field value
// within a sliding time window. Once a value exceeds its cap it is dropped
// until the window resets.
package capper

import (
	"encoding/json"
	"sync"
	"time"
)

// Rule defines a cap for a specific field.
type Rule struct {
	Field  string        // JSON field to inspect
	Max    int           // maximum lines allowed per value per window
	Window time.Duration // sliding window duration
}

type entry struct {
	count     int
	windowEnd time.Time
}

// Capper tracks per-value counts and drops lines that exceed the cap.
type Capper struct {
	mu    sync.Mutex
	rules []Rule
	state map[string]map[string]*entry // rule field -> value -> entry
}

// New creates a Capper with the given rules.
func New(rules []Rule) *Capper {
	state := make(map[string]map[string]*entry, len(rules))
	for _, r := range rules {
		state[r.Field] = make(map[string]*entry)
	}
	return &Capper{rules: rules, state: state}
}

// Allow returns true if the line should be passed through, false if it should
// be dropped. Lines with invalid JSON are always passed through.
func (c *Capper) Allow(line string) bool {
	if len(c.rules) == 0 {
		return true
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return true
	}
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, r := range c.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		key, ok := v.(string)
		if !ok {
			continue
		}
		e, exists := c.state[r.Field][key]
		if !exists || now.After(e.windowEnd) {
			c.state[r.Field][key] = &entry{count: 1, windowEnd: now.Add(r.Window)}
			continue
		}
		e.count++
		if e.count > r.Max {
			return false
		}
	}
	return true
}

// Reset clears all accumulated state.
func (c *Capper) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, r := range c.rules {
		c.state[r.Field] = make(map[string]*entry)
	}
}

package sequencer

import (
	"encoding/json"
	"sync"
)

// Rule defines a field to add a sequence number to.
type Rule struct {
	Field  string // destination field name
	Start  int    // initial value (default 0)
	Step   int    // increment per line (default 1)
}

// Sequencer assigns monotonically increasing sequence numbers to log lines.
type Sequencer struct {
	rules    []Rule
	counters []int
	mu       sync.Mutex
}

// New creates a Sequencer with the given rules.
func New(rules []Rule) *Sequencer {
	counters := make([]int, len(rules))
	for i, r := range rules {
		if r.Step == 0 {
			rules[i].Step = 1
		}
		counters[i] = r.Start
	}
	return &Sequencer{rules: rules, counters: counters}
}

// Apply injects sequence numbers into a JSON log line.
func (s *Sequencer) Apply(line string) (string, error) {
	if len(s.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, err
	}

	s.mu.Lock()
	for i, r := range s.rules {
		obj[r.Field] = s.counters[i]
		s.counters[i] += r.Step
	}
	s.mu.Unlock()

	out, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return string(out), nil
}

// Reset resets all counters to their starting values.
func (s *Sequencer) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, r := range s.rules {
		s.counters[i] = r.Start
	}
}

// Snapshot returns the current counter values keyed by field name.
func (s *Sequencer) Snapshot() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]int, len(s.rules))
	for i, r := range s.rules {
		out[r.Field] = s.counters[i]
	}
	return out
}

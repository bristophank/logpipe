// Package sequencer assigns a monotonically increasing sequence number to each
// JSON log line. This is useful for preserving ordering guarantees when logs
// are fanned out to multiple sinks or processed concurrently.
package sequencer

import (
	"encoding/json"
	"sync/atomic"
)

// Sequencer stamps each log line with a sequence number field.
type Sequencer struct {
	field  string
	start  int64
	counter atomic.Int64
}

// Config holds options for the Sequencer.
type Config struct {
	// Field is the JSON key to write the sequence number into.
	// Defaults to "seq" if empty.
	Field string

	// Start is the initial sequence value. Defaults to 1.
	Start int64
}

// New creates a Sequencer from the given Config.
func New(cfg Config) *Sequencer {
	field := cfg.Field
	if field == "" {
		field = "seq"
	}
	start := cfg.Start
	if start == 0 {
		start = 1
	}
	s := &Sequencer{
		field: field,
		start: start,
	}
	s.counter.Store(start)
	return s
}

// Apply parses line as JSON, injects the next sequence number, and returns
// the re-encoded line. If line is empty or not valid JSON the original bytes
// are returned unchanged.
func (s *Sequencer) Apply(line []byte) ([]byte, error) {
	if len(line) == 0 {
		return line, nil
	}

	var obj map[string]any
	if err := json.Unmarshal(line, &obj); err != nil {
		// Not valid JSON — pass through without stamping.
		return line, nil
	}

	seq := s.counter.Add(1) - 1 // fetch-then-increment; first value == start
	obj[s.field] = seq

	out, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return out, nil
}

// Reset sets the counter back to the configured start value. Safe for
// concurrent use.
func (s *Sequencer) Reset() {
	s.counter.Store(s.start)
}

// Current returns the next sequence number that will be assigned without
// advancing the counter.
func (s *Sequencer) Current() int64 {
	return s.counter.Load()
}

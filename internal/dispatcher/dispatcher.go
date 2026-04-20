// Package dispatcher routes parsed log lines to named sinks based on
// configurable field-match rules, falling back to a default sink when no
// rule matches.
package dispatcher

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Rule describes a single dispatch condition.
type Rule struct {
	Field string // JSON field name to inspect
	Value string // expected value (exact match)
	Sink  string // destination sink name
}

// Config holds the dispatcher configuration.
type Config struct {
	Rules       []Rule
	DefaultSink string // used when no rule matches; empty means drop
}

// Dispatcher evaluates each log line against a set of rules and writes it to
// the appropriate sink.
type Dispatcher struct {
	cfg   Config
	sinks map[string]io.Writer
}

// New creates a Dispatcher with the given config and sink map.
func New(cfg Config, sinks map[string]io.Writer) (*Dispatcher, error) {
	for i, r := range cfg.Rules {
		if r.Field == "" {
			return nil, fmt.Errorf("rule %d: field must not be empty", i)
		}
		if r.Sink == "" {
			return nil, fmt.Errorf("rule %d: sink must not be empty", i)
		}
		if _, ok := sinks[r.Sink]; !ok {
			return nil, fmt.Errorf("rule %d: unknown sink %q", i, r.Sink)
		}
	}
	if cfg.DefaultSink != "" {
		if _, ok := sinks[cfg.DefaultSink]; !ok {
			return nil, fmt.Errorf("default sink %q not found", cfg.DefaultSink)
		}
	}
	return &Dispatcher{cfg: cfg, sinks: sinks}, nil
}

// Dispatch evaluates line against all rules and writes it to the matching sink.
// It returns the name of the sink used, or an empty string if the line was dropped.
func (d *Dispatcher) Dispatch(line string) (string, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return "", nil
	}

	var fields map[string]interface{}
	if err := json.Unmarshal([]byte(line), &fields); err != nil {
		// non-JSON lines go to default sink
		return d.write(d.cfg.DefaultSink, line)
	}

	for _, r := range d.cfg.Rules {
		v, ok := fields[r.Field]
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", v) == r.Value {
			return d.write(r.Sink, line)
		}
	}

	return d.write(d.cfg.DefaultSink, line)
}

func (d *Dispatcher) write(sinkName, line string) (string, error) {
	if sinkName == "" {
		return "", nil
	}
	w := d.sinks[sinkName]
	_, err := fmt.Fprintln(w, line)
	if err != nil {
		return "", err
	}
	return sinkName, nil
}

package aggregator

import (
	"encoding/json"
	"sync"
	"time"
)

// Aggregator groups log lines by a key field and emits count summaries on a flush interval.
type Aggregator struct {
	mu       sync.Mutex
	field    string
	counts   map[string]int
	window   time.Duration
	stop     chan struct{}
	output   func(line string)
}

// New creates an Aggregator that groups by field and flushes every window.
// output is called with each summary JSON line on flush.
func New(field string, window time.Duration, output func(string)) *Aggregator {
	a := &Aggregator{
		field:  field,
		counts: make(map[string]int),
		window: window,
		stop:   make(chan struct{}),
		output: output,
	}
	if window > 0 {
		go a.run()
	}
	return a
}

// Add records a log line into the aggregation bucket.
func (a *Aggregator) Add(line string) {
	var m map[string]any
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		return
	}
	val, ok := m[a.field]
	if !ok {
		return
	}
	key := toString(val)
	a.mu.Lock()
	a.counts[key]++
	a.mu.Unlock()
}

// Flush emits current counts and resets state.
func (a *Aggregator) Flush() {
	a.mu.Lock()
	snap := a.counts
	a.counts = make(map[string]int)
	a.mu.Unlock()
	for k, v := range snap {
		out := map[string]any{"field": a.field, "value": k, "count": v}
		b, _ := json.Marshal(out)
		a.output(string(b))
	}
}

// Stop halts the background flush goroutine.
func (a *Aggregator) Stop() {
	close(a.stop)
}

func (a *Aggregator) run() {
	t := time.NewTicker(a.window)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			a.Flush()
		case <-a.stop:
			return
		}
	}
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

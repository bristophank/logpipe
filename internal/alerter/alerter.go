// Package alerter fires alerts when a log field threshold is exceeded.
package alerter

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"
	"time"
)

// Rule defines a threshold condition on a numeric log field.
type Rule struct {
	Field     string  `json:"field"`
	Threshold float64 `json:"threshold"`
	Window    int     `json:"window_seconds"`
	SinkName  string  `json:"sink"`
}

type bucket struct {
	count int
	reset time.Time
}

// Alerter checks log lines against threshold rules and writes alerts.
type Alerter struct {
	mu      sync.Mutex
	rules   []Rule
	buckets map[string]*bucket
	now     func() time.Time
}

// New creates an Alerter with the given rules.
func New(rules []Rule) *Alerter {
	return &Alerter{
		rules:   rules,
		buckets: make(map[string]*bucket),
		now:     time.Now,
	}
}

// Check evaluates line against all rules, writing alerts to w.
func (a *Alerter) Check(line string, w io.Writer) error {
	if len(a.rules) == 0 {
		return nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return nil
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	for _, r := range a.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		num, ok := val.(float64)
		if !ok {
			continue
		}
		if num <= r.Threshold {
			continue
		}
		key := r.Field + ":" + r.SinkName
		b, exists := a.buckets[key]
		now := a.now()
		if !exists || now.After(b.reset) {
			b = &bucket{reset: now.Add(time.Duration(r.Window) * time.Second)}
			a.buckets[key] = b
		}
		b.count++
		msg := fmt.Sprintf(`{"alert":true,"field":%q,"value":%g,"threshold":%g,"count":%d}\n`,
			r.Field, num, r.Threshold, b.count)
		fmt.Fprint(w, msg)
	}
	return nil
}

// Reset clears all window buckets.
func (a *Alerter) Reset() {
	a.mu.Lock()
	a.buckets = make(map[string]*bucket)
	a.mu.Unlock()
}

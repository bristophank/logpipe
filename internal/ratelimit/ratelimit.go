package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a maximum number of lines per second per sink.
type Limiter struct {
	mu       sync.Mutex
	rate     int
	counters map[string]int
	window   map[string]time.Time
}

// New creates a Limiter allowing up to rate lines/sec per sink.
// A rate <= 0 disables limiting.
func New(rate int) *Limiter {
	return &Limiter{
		rate:     rate,
		counters: make(map[string]int),
		window:   make(map[string]time.Time),
	}
}

// Allow returns true if the line for the given sink is within the rate limit.
func (l *Limiter) Allow(sink string) bool {
	if l.rate <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	start, ok := l.window[sink]
	if !ok || now.Sub(start) >= time.Second {
		l.window[sink] = now
		l.counters[sink] = 1
		return true
	}
	l.counters[sink]++
	return l.counters[sink] <= l.rate
}

// Reset clears counters for all sinks.
func (l *Limiter) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.counters = make(map[string]int)
	l.window = make(map[string]time.Time)
}

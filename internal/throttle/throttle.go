package throttle

import (
	"sync"
	"time"
)

// Throttler suppresses repeated identical log lines within a cooldown window.
type Throttler struct {
	mu       sync.Mutex
	cooldown time.Duration
	seen     map[string]time.Time
	now      func() time.Time
}

// New creates a Throttler with the given cooldown duration.
// If cooldown is zero, all lines are allowed through.
func New(cooldown time.Duration) *Throttler {
	return &Throttler{
		cooldown: cooldown,
		seen:     make(map[string]time.Time),
		now:      time.Now,
	}
}

// Allow returns true if the line should be forwarded.
// A line is suppressed if the same line was seen within the cooldown window.
func (t *Throttler) Allow(line string) bool {
	if t.cooldown == 0 {
		return true
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.now()
	if last, ok := t.seen[line]; ok && now.Sub(last) < t.cooldown {
		return false
	}
	t.seen[line] = now
	return true
}

// Reset clears all tracked lines.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.seen = make(map[string]time.Time)
}

// Len returns the number of currently tracked lines.
func (t *Throttler) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.seen)
}

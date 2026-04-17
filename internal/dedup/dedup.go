package dedup

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"
)

// Deduplicator suppresses repeated log lines within a sliding time window.
type Deduplicator struct {
	mu      sync.Mutex
	seen    map[string]time.Time
	window  time.Duration
	now     func() time.Time
}

// New creates a Deduplicator that suppresses duplicate lines seen within window.
// If window is zero, deduplication is disabled and all lines pass.
func New(window time.Duration) *Deduplicator {
	return &Deduplicator{
		seen:   make(map[string]time.Time),
		window: window,
		now:    time.Now,
	}
}

// Allow returns true if the line should be forwarded (not a recent duplicate).
func (d *Deduplicator) Allow(line string) bool {
	if d.window == 0 {
		return true
	}
	key := hash(line)
	now := d.now()
	d.mu.Lock()
	defer d.mu.Unlock()
	d.evict(now)
	if _, exists := d.seen[key]; exists {
		return false
	}
	d.seen[key] = now
	return true
}

// Reset clears all tracked entries.
func (d *Deduplicator) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]time.Time)
}

// evict removes entries older than the window. Must be called with mu held.
func (d *Deduplicator) evict(now time.Time) {
	for k, t := range d.seen {
		if now.Sub(t) >= d.window {
			delete(d.seen, k)
		}
	}
}

func hash(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

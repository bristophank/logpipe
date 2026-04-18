// Package tee fans out a single log stream to multiple writers.
package tee

import (
	"io"
	"sync"
)

// Tee writes each line to all registered writers.
type Tee struct {
	mu      sync.RWMutex
	writers map[string]io.Writer
}

// New returns an empty Tee.
func New() *Tee {
	return &Tee{writers: make(map[string]io.Writer)}
}

// Add registers a named writer. Replaces any existing writer with the same name.
func (t *Tee) Add(name string, w io.Writer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.writers[name] = w
}

// Remove unregisters a named writer.
func (t *Tee) Remove(name string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.writers, name)
}

// Write sends p to every registered writer.
// Errors from individual writers are collected; the first error is returned.
func (t *Tee) Write(p []byte) (int, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	var firstErr error
	for _, w := range t.writers {
		if _, err := w.Write(p); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return len(p), firstErr
}

// Len returns the number of registered writers.
func (t *Tee) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.writers)
}

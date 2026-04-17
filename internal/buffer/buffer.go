package buffer

import (
	"sync"
)

// Buffer is a thread-safe ring buffer for log lines.
type Buffer struct {
	mu       sync.Mutex
	lines    []string
	cap      int
	head     int
	count    int
	dropped  int64
}

// New creates a Buffer with the given capacity.
func New(capacity int) *Buffer {
	if capacity <= 0 {
		capacity = 1
	}
	return &Buffer{
		lines: make([]string, capacity),
		cap:   capacity,
	}
}

// Write adds a line to the buffer. If full, the oldest entry is overwritten
// and the dropped counter is incremented.
func (b *Buffer) Write(line string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.count == b.cap {
		b.dropped++
		b.lines[b.head] = line
		b.head = (b.head + 1) % b.cap
		return
	}

	idx := (b.head + b.count) % b.cap
	b.lines[idx] = line
	b.count++
}

// Flush returns all buffered lines in order and resets the buffer.
func (b *Buffer) Flush() []string {
	b.mu.Lock()
	defer b.mu.Unlock()

	out := make([]string, b.count)
	for i := 0; i < b.count; i++ {
		out[i] = b.lines[(b.head+i)%b.cap]
	}
	b.head = 0
	b.count = 0
	return out
}

// Dropped returns the number of lines dropped due to overflow.
func (b *Buffer) Dropped() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.dropped
}

// Len returns the current number of buffered lines.
func (b *Buffer) Len() int {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.count
}

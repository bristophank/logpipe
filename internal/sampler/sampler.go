package sampler

import "sync/atomic"

// Sampler drops log lines based on a 1-in-N sampling strategy.
// A rate of 0 or 1 means keep everything.
type Sampler struct {
	rate    uint64
	counter atomic.Uint64
}

// New returns a Sampler that keeps 1 out of every rate lines.
// rate <= 1 disables sampling (all lines pass).
func New(rate uint64) *Sampler {
	if rate == 0 {
		rate = 1
	}
	return &Sampler{rate: rate}
}

// Allow returns true if the current line should be forwarded.
func (s *Sampler) Allow() bool {
	n := s.counter.Add(1)
	return n%s.rate == 1
}

// Rate returns the configured sampling rate.
func (s *Sampler) Rate() uint64 {
	return s.rate
}

// Reset resets the internal counter.
func (s *Sampler) Reset() {
	s.counter.Store(0)
}

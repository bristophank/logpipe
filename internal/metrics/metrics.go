package metrics

import "sync/atomic"

// Counters holds pipeline-level counters.
type Counters struct {
	LinesIn      atomic.Int64
	LinesMatched atomic.Int64
	LinesDropped atomic.Int64
	RouteErrors  atomic.Int64
}

// Collector collects and exposes pipeline metrics.
type Collector struct {
	counters Counters
}

// New returns a new Collector.
func New() *Collector {
	return &Collector{}
}

// IncIn increments the lines-in counter.
func (c *Collector) IncIn() { c.counters.LinesIn.Add(1) }

// IncMatched increments the lines-matched counter.
func (c *Collector) IncMatched() { c.counters.LinesMatched.Add(1) }

// IncDropped increments the lines-dropped counter.
func (c *Collector) IncDropped() { c.counters.LinesDropped.Add(1) }

// IncRouteError increments the route-error counter.
func (c *Collector) IncRouteError() { c.counters.RouteErrors.Add(1) }

// Snapshot returns a point-in-time copy of the counters.
func (c *Collector) Snapshot() map[string]int64 {
	return map[string]int64{
		"lines_in":      c.counters.LinesIn.Load(),
		"lines_matched": c.counters.LinesMatched.Load(),
		"lines_dropped": c.counters.LinesDropped.Load(),
		"route_errors":  c.counters.RouteErrors.Load(),
	}
}

// Reset zeroes all counters.
func (c *Collector) Reset() {
	c.counters.LinesIn.Store(0)
	c.counters.LinesMatched.Store(0)
	c.counters.LinesDropped.Store(0)
	c.counters.RouteErrors.Store(0)
}

package metrics

import (
	"testing"
)

func TestCollector_InitialZero(t *testing.T) {
	c := New()
	snap := c.Snapshot()
	for k, v := range snap {
		if v != 0 {
			t.Errorf("expected %s=0, got %d", k, v)
		}
	}
}

func TestCollector_Increments(t *testing.T) {
	c := New()
	c.IncIn()
	c.IncIn()
	c.IncMatched()
	c.IncDropped()
	c.IncRouteError()

	snap := c.Snapshot()
	if snap["lines_in"] != 2 {
		t.Errorf("lines_in: want 2, got %d", snap["lines_in"])
	}
	if snap["lines_matched"] != 1 {
		t.Errorf("lines_matched: want 1, got %d", snap["lines_matched"])
	}
	if snap["lines_dropped"] != 1 {
		t.Errorf("lines_dropped: want 1, got %d", snap["lines_dropped"])
	}
	if snap["route_errors"] != 1 {
		t.Errorf("route_errors: want 1, got %d", snap["route_errors"])
	}
}

func TestCollector_Reset(t *testing.T) {
	c := New()
	c.IncIn()
	c.IncMatched()
	c.Reset()

	snap := c.Snapshot()
	for k, v := range snap {
		if v != 0 {
			t.Errorf("after reset, expected %s=0, got %d", k, v)
		}
	}
}

func TestCollector_SnapshotIndependent(t *testing.T) {
	c := New()
	c.IncIn()
	snap1 := c.Snapshot()
	c.IncIn()
	snap2 := c.Snapshot()
	if snap1["lines_in"] == snap2["lines_in"] {
		t.Error("snapshots should differ after increment")
	}
}

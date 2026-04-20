package counter

import (
	"testing"
)

func TestAdd_NoRules(t *testing.T) {
	c := New(nil)
	if err := c.Add(`{"level":"info"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	snap := c.Snapshot()
	if len(snap) != 0 {
		t.Errorf("expected empty snapshot, got %v", snap)
	}
}

func TestAdd_InvalidJSON(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	if err := c.Add("not-json"); err == nil {
		t.Fatal("expected error for invalid json")
	}
}

func TestAdd_CountsValues(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	lines := []string{
		`{"level":"info"}`,
		`{"level":"error"}`,
		`{"level":"info"}`,
		`{"level":"info"}`,
	}
	for _, l := range lines {
		if err := c.Add(l); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	snap := c.Snapshot()
	if snap["level"]["info"] != 3 {
		t.Errorf("expected 3 for info, got %d", snap["level"]["info"])
	}
	if snap["level"]["error"] != 1 {
		t.Errorf("expected 1 for error, got %d", snap["level"]["error"])
	}
}

func TestAdd_MissingField(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	_ = c.Add(`{"msg":"hello"}`)
	snap := c.Snapshot()
	if len(snap["level"]) != 0 {
		t.Errorf("expected no counts for missing field")
	}
}

func TestSnapshot_UsesAlias(t *testing.T) {
	c := New([]Rule{{Field: "level", Alias: "severity"}})
	_ = c.Add(`{"level":"warn"}`)
	snap := c.Snapshot()
	if _, ok := snap["severity"]; !ok {
		t.Error("expected alias 'severity' in snapshot")
	}
	if _, ok := snap["level"]; ok {
		t.Error("expected raw field 'level' to be absent when alias set")
	}
}

func TestReset_ClearsCounts(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	_ = c.Add(`{"level":"info"}`)
	c.Reset()
	snap := c.Snapshot()
	if len(snap["level"]) != 0 {
		t.Errorf("expected empty after reset, got %v", snap["level"])
	}
}

func TestSnapshot_Independent(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	_ = c.Add(`{"level":"info"}`)
	s1 := c.Snapshot()
	_ = c.Add(`{"level":"info"}`)
	if s1["level"]["info"] != 1 {
		t.Errorf("snapshot should be independent of later mutations")
	}
}

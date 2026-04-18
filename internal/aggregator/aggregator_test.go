package aggregator

import (
	"encoding/json"
	"testing"
	"time"
)

func collect(t *testing.T, field string, lines []string) []map[string]any {
	t.Helper()
	var out []map[string]any
	a := New(field, 0, func(line string) {
		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err == nil {
			out = append(out, m)
		}
	})
	for _, l := range lines {
		a.Add(l)
	}
	a.Flush()
	return out
}

func TestAdd_InvalidJSON(t *testing.T) {
	results := collect(t, "level", []string{"not-json"})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestAdd_MissingField(t *testing.T) {
	results := collect(t, "level", []string{`{"msg":"hello"}`})
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestFlush_CountsGrouped(t *testing.T) {
	lines := []string{
		`{"level":"info"}`,
		`{"level":"info"}`,
		`{"level":"error"}`,
	}
	results := collect(t, "level", lines)
	counts := map[string]int{}
	for _, r := range results {
		v := r["value"].(string)
		counts[v] = int(r["count"].(float64))
	}
	if counts["info"] != 2 {
		t.Errorf("expected info=2, got %d", counts["info"])
	}
	if counts["error"] != 1 {
		t.Errorf("expected error=1, got %d", counts["error"])
	}
}

func TestFlush_ResetsState(t *testing.T) {
	var calls int
	a := New("level", 0, func(_ string) { calls++ })
	a.Add(`{"level":"info"}`)
	a.Flush()
	a.Flush() // second flush should emit nothing
	if calls != 1 {
		t.Errorf("expected 1 output line, got %d", calls)
	}
}

func TestAutoFlush_WithWindow(t *testing.T) {
	var calls int
	a := New("level", 50*time.Millisecond, func(_ string) { calls++ })
	a.Add(`{"level":"warn"}`)
	time.Sleep(120 * time.Millisecond)
	a.Stop()
	if calls == 0 {
		t.Error("expected at least one auto-flush")
	}
}

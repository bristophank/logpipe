package windower

import (
	"encoding/json"
	"testing"
	"time"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestAdd_NoRules_ReturnsFalse(t *testing.T) {
	w := New(nil)
	_, ok := w.Add(`{"bytes":100}`)
	if ok {
		t.Fatal("expected no output with no rules")
	}
}

func TestAdd_InvalidJSON_ReturnsFalse(t *testing.T) {
	w := New([]Rule{{Field: "bytes", Alias: "total_bytes", Window: time.Second}})
	_, ok := w.Add("not-json")
	if ok {
		t.Fatal("expected false for invalid JSON")
	}
}

func TestAdd_AccumulatesWithinWindow(t *testing.T) {
	w := New([]Rule{{Field: "bytes", Alias: "total_bytes", Window: time.Hour}})
	w.now = func() time.Time { return time.Unix(0, 0) }

	_, ok := w.Add(`{"bytes":100}`)
	if ok {
		t.Fatal("should not flush mid-window")
	}
	_, ok = w.Add(`{"bytes":200}`)
	if ok {
		t.Fatal("should not flush mid-window")
	}
}

func TestAdd_FlushesOnWindowExpiry(t *testing.T) {
	now := time.Unix(0, 0)
	w := New([]Rule{{Field: "bytes", Alias: "total_bytes", Window: time.Second}})
	w.now = func() time.Time { return now }

	w.Add(`{"bytes":100}`)
	w.Add(`{"bytes":150}`)

	// advance past the window
	now = now.Add(2 * time.Second)
	w.now = func() time.Time { return now }

	line, ok := w.Add(`{"bytes":50}`)
	if !ok {
		t.Fatal("expected flush after window expiry")
	}
	m := decode(t, line)
	if m["total_bytes"].(float64) != 300 {
		t.Fatalf("expected 300, got %v", m["total_bytes"])
	}
}

func TestFlush_ResetsBucket(t *testing.T) {
	now := time.Unix(0, 0)
	w := New([]Rule{{Field: "count", Alias: "total_count", Window: time.Hour}})
	w.now = func() time.Time { return now }

	w.Add(`{"count":5}`)
	line, ok := w.Flush()
	if !ok {
		t.Fatal("expected flush to succeed")
	}
	m := decode(t, line)
	if m["total_count"].(float64) != 5 {
		t.Fatalf("expected 5, got %v", m["total_count"])
	}

	// second flush should be empty
	_, ok = w.Flush()
	if ok {
		t.Fatal("expected empty flush after reset")
	}
}

func TestAdd_MissingField_Skipped(t *testing.T) {
	now := time.Unix(0, 0)
	w := New([]Rule{{Field: "bytes", Alias: "total_bytes", Window: time.Hour}})
	w.now = func() time.Time { return now }

	w.Add(`{"other":999}`)
	w.Add(`{"bytes":42}`)

	line, ok := w.Flush()
	if !ok {
		t.Fatal("expected flush")
	}
	m := decode(t, line)
	if m["total_bytes"].(float64) != 42 {
		t.Fatalf("expected 42, got %v", m["total_bytes"])
	}
}

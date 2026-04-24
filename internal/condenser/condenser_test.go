package condenser

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]any {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestAdd_NoRules_Passthrough(t *testing.T) {
	c := New(nil)
	out := c.Add(`{"level":"info","msg":"hello"}`)
	if out == "" {
		t.Fatal("expected passthrough, got empty")
	}
}

func TestAdd_EmptyLine_Skipped(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	out := c.Add("   ")
	if out != "" {
		t.Fatalf("expected empty, got %q", out)
	}
}

func TestAdd_InvalidJSON_Passthrough(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	out := c.Add("not json")
	if out != "not json" {
		t.Fatalf("expected passthrough, got %q", out)
	}
}

func TestAdd_SameKey_Accumulates(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	if out := c.Add(`{"level":"info","msg":"a"}`); out != "" {
		t.Fatalf("first line should buffer, got %q", out)
	}
	if out := c.Add(`{"level":"info","msg":"b"}`); out != "" {
		t.Fatalf("duplicate should buffer, got %q", out)
	}
}

func TestAdd_KeyChange_FlushesGroup(t *testing.T) {
	c := New([]Rule{{Field: "level", CountField: "_count"}})
	c.Add(`{"level":"info","msg":"a"}`)
	c.Add(`{"level":"info","msg":"b"}`)
	out := c.Add(`{"level":"warn","msg":"c"}`)
	if out == "" {
		t.Fatal("expected flushed line on key change")
	}
	m := decode(t, out)
	if m["level"] != "info" {
		t.Fatalf("expected level=info, got %v", m["level"])
	}
	if cnt, ok := m["_count"].(float64); !ok || cnt != 2 {
		t.Fatalf("expected _count=2, got %v", m["_count"])
	}
}

func TestFlush_ReturnsLastGroup(t *testing.T) {
	c := New([]Rule{{Field: "level", CountField: "n"}})
	c.Add(`{"level":"error","msg":"x"}`)
	c.Add(`{"level":"error","msg":"y"}`)
	c.Add(`{"level":"error","msg":"z"}`)
	out := c.Flush()
	if out == "" {
		t.Fatal("expected final flush output")
	}
	m := decode(t, out)
	if cnt, _ := m["n"].(float64); cnt != 3 {
		t.Fatalf("expected n=3, got %v", m["n"])
	}
}

func TestFlush_EmptyState_ReturnsEmpty(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	if out := c.Flush(); out != "" {
		t.Fatalf("expected empty flush on fresh condenser, got %q", out)
	}
}

func TestFlush_ResetsState(t *testing.T) {
	c := New([]Rule{{Field: "level"}})
	c.Add(`{"level":"info"}`)
	c.Flush()
	if out := c.Flush(); out != "" {
		t.Fatalf("second flush should be empty, got %q", out)
	}
}

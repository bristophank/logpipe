package alerter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestCheck_NoRules(t *testing.T) {
	a := New(nil)
	var buf bytes.Buffer
	_ = a.Check(`{"latency":999}`, &buf)
	if buf.Len() != 0 {
		t.Fatalf("expected no output")
	}
}

func TestCheck_BelowThreshold(t *testing.T) {
	a := New([]Rule{{Field: "latency", Threshold: 500, Window: 60}})
	var buf bytes.Buffer
	_ = a.Check(`{"latency":100}`, &buf)
	if buf.Len() != 0 {
		t.Fatalf("expected no alert")
	}
}

func TestCheck_ExceedsThreshold(t *testing.T) {
	a := New([]Rule{{Field: "latency", Threshold: 500, Window: 60}})
	var buf bytes.Buffer
	_ = a.Check(`{"latency":600}`, &buf)
	if buf.Len() == 0 {
		t.Fatal("expected alert output")
	}
	if !strings.Contains(buf.String(), `"alert":true`) {
		t.Fatalf("unexpected output: %s", buf.String())
	}
}

func TestCheck_CountIncrementsInWindow(t *testing.T) {
	a := New([]Rule{{Field: "err", Threshold: 0, Window: 60}})
	var buf bytes.Buffer
	for i := 0; i < 3; i++ {
		_ = a.Check(`{"err":1}`, &buf)
	}
	if !strings.Contains(buf.String(), `"count":3`) {
		t.Fatalf("expected count 3, got: %s", buf.String())
	}
}

func TestCheck_WindowReset(t *testing.T) {
	a := New([]Rule{{Field: "err", Threshold: 0, Window: 1}})
	now := time.Now()
	a.now = func() time.Time { return now }
	var buf bytes.Buffer
	_ = a.Check(`{"err":1}`, &buf)
	// advance past window
	a.now = func() time.Time { return now.Add(2 * time.Second) }
	buf.Reset()
	_ = a.Check(`{"err":1}`, &buf)
	if !strings.Contains(buf.String(), `"count":1`) {
		t.Fatalf("expected reset count, got: %s", buf.String())
	}
}

func TestCheck_InvalidJSON(t *testing.T) {
	a := New([]Rule{{Field: "x", Threshold: 1, Window: 10}})
	var buf bytes.Buffer
	if err := a.Check("not-json", &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatal("expected no output for invalid JSON")
	}
}

package flattener

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestStream_FlattensLines(t *testing.T) {
	f := New([]Rule{{Separator: "."}})
	input := "{\"a\":{\"b\":1}}\n{\"x\":{\"y\":2}}\n"
	var buf bytes.Buffer
	if err := Stream(f, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	var m1 map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &m1); err != nil {
		t.Fatalf("unmarshal line 0: %v", err)
	}
	if m1["a.b"] != float64(1) {
		t.Errorf("expected a.b=1, got %v", m1["a.b"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	f := New([]Rule{{Separator: "."}})
	input := "\n{\"a\":{\"b\":1}}\n\n"
	var buf bytes.Buffer
	if err := Stream(f, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Errorf("expected 1 line, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	f := New([]Rule{{Separator: "."}})
	input := "not json\n"
	var buf bytes.Buffer
	if err := Stream(f, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "not json" {
		t.Errorf("expected passthrough, got %q", got)
	}
}

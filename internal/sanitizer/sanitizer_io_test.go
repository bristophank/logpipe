package sanitizer

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_SanitizesLines(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "trim"}})
	input := `{"msg":"  hello  "}
{"msg":" world "}
`
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, s); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["msg"] != "hello" {
		t.Errorf("line 0: expected 'hello', got %q", m0["msg"])
	}
	m1 := decode(t, lines[1])
	if m1["msg"] != "world" {
		t.Errorf("line 1: expected 'world', got %q", m1["msg"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "trim"}})
	input := `{"msg":"hello"}

{"msg":"world"}
`
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, s); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	s := New([]Rule{{Field: "msg", Mode: "strip_html"}})
	input := "not-json\n"
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, s); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if got := strings.TrimSpace(buf.String()); got != "not-json" {
		t.Errorf("expected passthrough, got %q", got)
	}
}

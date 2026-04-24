package dedenter

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_DedentsLines(t *testing.T) {
	d := New([]Rule{{Field: "msg"}})
	input := `{"msg":"   hello"}` + "\n" + `{"msg":"\tworld"}` + "\n"
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, d); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["msg"] != "hello" {
		t.Errorf("line 0: unexpected msg %v", m0["msg"])
	}
	m1 := decode(t, lines[1])
	if m1["msg"] != "world" {
		t.Errorf("line 1: unexpected msg %v", m1["msg"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	d := New([]Rule{{Field: "msg"}})
	input := "\n" + `{"msg":" hi"}` + "\n\n"
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, d); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	d := New([]Rule{{Field: "msg"}})
	input := "not-json\n"
	var buf bytes.Buffer
	if err := Stream(strings.NewReader(input), &buf, d); err != nil {
		t.Fatalf("Stream: %v", err)
	}
	if !strings.Contains(buf.String(), "not-json") {
		t.Errorf("invalid JSON should pass through unchanged")
	}
}

package shifter

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_ShiftsLines(t *testing.T) {
	sh := New([]Rule{{Field: "n", By: 1}})
	input := `{"n":4}
{"n":9}
`
	var buf bytes.Buffer
	if err := Stream(sh, strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["n"].(float64) != 5 {
		t.Errorf("line 0: expected n=5, got %v", m0["n"])
	}
	m1 := decode(t, lines[1])
	if m1["n"].(float64) != 10 {
		t.Errorf("line 1: expected n=10, got %v", m1["n"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	sh := New([]Rule{{Field: "x", By: 1}})
	input := "\n{\"x\":1}\n\n"
	var buf bytes.Buffer
	if err := Stream(sh, strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %v", len(lines), lines)
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	sh := New([]Rule{{Field: "x", By: 1}})
	input := "not-json\n"
	var buf bytes.Buffer
	if err := Stream(sh, strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "not-json") {
		t.Errorf("expected passthrough of invalid JSON, got %q", buf.String())
	}
}

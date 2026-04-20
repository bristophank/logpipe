package sequencer

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_SequencesLines(t *testing.T) {
	seq := New([]Rule{{Field: "seq", Start: 1, Step: 1}})
	input := `{"msg":"a"}
{"msg":"b"}
{"msg":"c"}
`
	var buf bytes.Buffer
	if err := seq.Stream(strings.NewReader(input), &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	for i, line := range lines {
		m := decode(t, line)
		got := int(m["seq"].(float64))
		if got != i+1 {
			t.Fatalf("line %d: expected seq=%d, got %d", i, i+1, got)
		}
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	seq := New([]Rule{{Field: "n"}})
	input := `{"a":1}

{"a":2}
`
	var buf bytes.Buffer
	seq.Stream(strings.NewReader(input), &buf)
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	seq := New([]Rule{{Field: "n"}})
	input := "not-json\n"
	var buf bytes.Buffer
	seq.Stream(strings.NewReader(input), &buf)
	got := strings.TrimSpace(buf.String())
	if got != "not-json" {
		t.Fatalf("expected passthrough of invalid JSON, got %q", got)
	}
}

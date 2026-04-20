package counter

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_PassesAllLines(t *testing.T) {
	input := "{\"level\":\"info\"}\n{\"level\":\"error\"}\n"
	c := New([]Rule{{Field: "level"}})
	var out bytes.Buffer
	if err := c.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 output lines, got %d", len(lines))
	}
}

func TestStream_CountsAccumulate(t *testing.T) {
	input := "{\"level\":\"info\"}\n{\"level\":\"info\"}\n{\"level\":\"warn\"}\n"
	c := New([]Rule{{Field: "level"}})
	var out bytes.Buffer
	_ = c.Stream(strings.NewReader(input), &out)
	snap := c.Snapshot()
	if snap["level"]["info"] != 2 {
		t.Errorf("expected 2 for info, got %d", snap["level"]["info"])
	}
	if snap["level"]["warn"] != 1 {
		t.Errorf("expected 1 for warn, got %d", snap["level"]["warn"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	input := "{\"level\":\"info\"}\n\n{\"level\":\"error\"}\n"
	c := New([]Rule{{Field: "level"}})
	var out bytes.Buffer
	_ = c.Stream(strings.NewReader(input), &out)
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 non-empty lines, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	input := "not-json\n{\"level\":\"info\"}\n"
	c := New([]Rule{{Field: "level"}})
	var out bytes.Buffer
	if err := c.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected both lines passed through, got %d", len(lines))
	}
}

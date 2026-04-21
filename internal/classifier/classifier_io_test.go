package classifier

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_ClassifiesLines(t *testing.T) {
	c := New("", []Rule{
		{Field: "level", Equals: "error", Category: "critical"},
	})
	input := `{"level":"error","msg":"boom"}
{"level":"info","msg":"ok"}
`
	var buf bytes.Buffer
	if err := c.Stream(strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["category"] != "critical" {
		t.Fatalf("line 0: expected critical, got %v", m0["category"])
	}
	m1 := decode(t, lines[1])
	if _, ok := m1["category"]; ok {
		t.Fatalf("line 1: unexpected category field")
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	c := New("", []Rule{{Field: "x", Equals: "y", Category: "z"}})
	input := "\n\n{\"x\":\"y\"}\n"
	var buf bytes.Buffer
	if err := c.Stream(strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	c := New("", []Rule{{Field: "level", Equals: "error", Category: "bad"}})
	input := "not json\n"
	var buf bytes.Buffer
	if err := c.Stream(strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	got := strings.TrimSpace(buf.String())
	if got != "not json" {
		t.Fatalf("expected passthrough, got %q", got)
	}
}

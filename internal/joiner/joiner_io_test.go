package joiner

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_JoinsLines(t *testing.T) {
	rules := []Rule{{PrimaryKey: "uid", SecondaryKey: "id", Fields: []string{"role"}}}
	j := New(rules)
	_ = j.Index(0, `{"id":"7","role":"admin"}`)

	input := `{"uid":"7","msg":"hello"}
{"uid":"8","msg":"world"}
`
	var buf bytes.Buffer
	if err := Stream(j, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["role"] != "admin" {
		t.Errorf("expected role=admin in first line, got %v", m0["role"])
	}
	m1 := decode(t, lines[1])
	if _, ok := m1["role"]; ok {
		t.Error("expected no role in second line")
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	j := New(nil)
	input := "\n{\"a\":1}\n\n{\"b\":2}\n"
	var buf bytes.Buffer
	if err := Stream(j, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 output lines, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	j := New([]Rule{{PrimaryKey: "id", SecondaryKey: "id", Fields: []string{"x"}}})
	input := "not-json\n"
	var buf bytes.Buffer
	if err := Stream(j, strings.NewReader(input), &buf); err != nil {
		t.Fatalf("Stream error: %v", err)
	}
	got := strings.TrimSpace(buf.String())
	if got != "not-json" {
		t.Errorf("expected passthrough, got %q", got)
	}
}

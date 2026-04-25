package spliceor

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_SplicesLines(t *testing.T) {
	rules := []Rule{{Target: "msg", Fields: []string{"level"}, Position: "after", Sep: " "}}
	s := New(rules)

	input := `{"msg":"hello","level":"info"}
{"msg":"world","level":"warn"}
`
	r := strings.NewReader(input)
	var buf bytes.Buffer
	if err := Stream(s, r, &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["msg"] != "hello info" {
		t.Errorf("line 0 msg: %v", m0["msg"])
	}
	m1 := decode(t, lines[1])
	if m1["msg"] != "world warn" {
		t.Errorf("line 1 msg: %v", m1["msg"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"level"}, Position: "after"}})
	input := "\n\n{\"msg\":\"hi\",\"level\":\"debug\"}\n"
	r := strings.NewReader(input)
	var buf bytes.Buffer
	if err := Stream(s, r, &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d: %v", len(lines), lines)
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	s := New([]Rule{{Target: "msg", Fields: []string{"level"}, Position: "after"}})
	input := "not-json\n"
	r := strings.NewReader(input)
	var buf bytes.Buffer
	if err := Stream(s, r, &buf); err != nil {
		t.Fatal(err)
	}
	out := strings.TrimSpace(buf.String())
	if out != "not-json" {
		t.Errorf("expected passthrough, got %q", out)
	}
}

package expander

import (
	"bytes"
	"strings"
	"testing"
)

func TestStream_ExpandsLines(t *testing.T) {
	e := New([]Rule{{Field: "tags", Delimiter: ","}})
	input := `{"tags":"a,b,c"}` + "\n" + `{"tags":"x,y"}` + "\n"
	var buf bytes.Buffer
	if err := e.Stream(strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	vals := m0["tags"].([]any)
	if len(vals) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(vals))
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	e := New([]Rule{{Field: "tags", Delimiter: ","}})
	input := "\n" + `{"tags":"a,b"}` + "\n\n"
	var buf bytes.Buffer
	if err := e.Stream(strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	e := New([]Rule{{Field: "tags", Delimiter: ","}})
	input := "not-json\n"
	var buf bytes.Buffer
	if err := e.Stream(strings.NewReader(input), &buf); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "not-json") {
		t.Fatal("expected invalid JSON to pass through")
	}
}

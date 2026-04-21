package zipper_test

import (
	"strings"
	"testing"

	"github.com/yourorg/logpipe/internal/zipper"
)

func TestStream_ZipsLines(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"first", "last"}, Target: "full", Sep: " "},
	})
	input := `{"first":"Ada","last":"Lovelace"}` + "\n" +
		`{"first":"Alan","last":"Turing"}` + "\n"

	var out strings.Builder
	if err := z.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["full"] != "Ada Lovelace" {
		t.Errorf("line 0: expected 'Ada Lovelace', got %v", m0["full"])
	}
	m1 := decode(t, lines[1])
	if m1["full"] != "Alan Turing" {
		t.Errorf("line 1: expected 'Alan Turing', got %v", m1["full"])
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"a", "b"}, Target: "ab", Sep: "-"},
	})
	input := "\n" + `{"a":"x","b":"y"}` + "\n\n"
	var out strings.Builder
	if err := z.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	z := zipper.New([]zipper.Rule{
		{Keys: []string{"a", "b"}, Target: "ab", Sep: "-"},
	})
	input := "not-json\n"
	var out strings.Builder
	if err := z.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out.String()) != "not-json" {
		t.Errorf("expected passthrough, got %q", out.String())
	}
}

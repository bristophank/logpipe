package tagger

import (
	"strings"
	"testing"
)

func TestStream_AppliesTags(t *testing.T) {
	tg := New([]Rule{{Field: "level", Value: "error", Tag: "alert", TagValue: "true"}})
	input := `{"level":"error","msg":"bad"}` + "\n" + `{"level":"info","msg":"ok"}` + "\n"
	var out strings.Builder
	if err := tg.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	m0 := decode(t, lines[0])
	if m0["alert"] != "true" {
		t.Fatalf("first line missing alert tag")
	}
	m1 := decode(t, lines[1])
	if _, ok := m1["alert"]; ok {
		t.Fatal("second line should not have alert tag")
	}
}

func TestStream_SkipsEmptyLines(t *testing.T) {
	tg := New(nil)
	input := "\n" + `{"msg":"hi"}` + "\n\n"
	var out strings.Builder
	if err := tg.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 1 {
		t.Fatalf("expected 1 line, got %d", len(lines))
	}
}

func TestStream_InvalidJSONPassthrough(t *testing.T) {
	tg := New([]Rule{{Field: "level", Value: "error", Tag: "alert", TagValue: "true"}})
	input := "not-json\n"
	var out strings.Builder
	if err := tg.Stream(strings.NewReader(input), &out); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out.String()) != "not-json" {
		t.Fatalf("expected passthrough, got %q", out.String())
	}
}
